package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"log"

	"ro-backend/configuration"
	"ro-backend/handler"
	"ro-backend/repository"
	"ro-backend/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
)

var userCollection *mongo.Collection
var authDataCollection *mongo.Collection
var refreshTokenCollection *mongo.Collection
var roPresetCollection *mongo.Collection
var roTagCollection *mongo.Collection

var appConfig configuration.AppConfig

var limiter = rate.NewLimiter(1, 3)

func main() {
	initTimeZone()
	appConfig = configuration.InitAppConfig()

	if err := connectMongoDB(); err != nil {
		panic(err)
	}

	initSessionStore()
	initAuthIdentityProviders()

	var userRepo = repository.NewUserRepo(userCollection)
	var authDataRepo = repository.NewAuthenticationDataRepo(authDataCollection)
	var refreshTokenRepo = repository.NewRefreshTokenRepo(refreshTokenCollection)
	var roPresetRepo = repository.NewRoPresetRepository(roPresetCollection)
	var roTagRepo = repository.NewPresetTagRepository(roTagCollection)

	var userService = service.NewUserService(userRepo)
	var tokenService = service.NewTokenService(refreshTokenRepo)
	var authDataService = service.NewAuthenticationDataService(authDataRepo)
	var roPresetService = service.NewRoPresetService(roPresetRepo, roTagRepo)
	var roTagService = service.NewPresetTagService(roTagRepo, roPresetRepo)

	var authHandler = handler.NewAuthHandler(handler.AuthHandlerParam{
		UserService:               userService,
		AuthenticationDataService: authDataService,
		TokenService:              tokenService,
	})
	var userHandler = handler.NewUserHandler(userService)
	var roPresetHandler = handler.NewRoPresetHandler(handler.RoPresetHandlerParam{
		RoPresetService:  roPresetService,
		UserService:      userService,
		PresetTagService: roTagService,
	})

	r := newAppRouter(mux.NewRouter())
	r.use(jsonResponseMiddleware)
	r.use(setSecurityHeadersMiddleware)
	r.use(rateLimitMiddleware)

	r.get("/auth/{provider}/callback", authHandler.AuthenticationCallback)
	r.get("/auth/{provider}", gothic.BeginAuthHandler)

	r.post("/login", authHandler.Login)
	r.post("/refresh_token", authHandler.RefreshToken)

	// ------
	me := r.subRouter("/me")
	me.use(userGuard)
	me.get("", userHandler.GetMyProfile)
	me.post("", userHandler.PatchMyProfile)
	me.post("/bulk_ro_presets", roPresetHandler.BulkCreatePresets)
	me.get("/ro_presets", roPresetHandler.GetMyPresets)
	me.post("/ro_presets", roPresetHandler.CreatePreset)

	me.get("/ro_presets/{presetId}", roPresetHandler.GetMyPresetById)
	me.post("/ro_presets/{presetId}", roPresetHandler.UpdateMyPreset)
	me.delete("/ro_presets/{presetId}", roPresetHandler.DeleteById)
	me.post("/ro_presets/{presetId}/publish", roPresetHandler.PublishMyPreset)
	me.delete("/ro_presets/{presetId}/publish", roPresetHandler.UnPublishMyPreset)

	me.post("/ro_presets/{presetId}/tags", roPresetHandler.AddTags)
	me.delete("/ro_presets/{presetId}/tags/{tagId}", roPresetHandler.RemoveTags)

	// ------
	ro := r.subRouter("/ro_presets")
	ro.use(userGuard)
	ro.get("/class_by_tags/{classId}/{tag}", roPresetHandler.SearchPresetTags)

	// ------
	tag := r.subRouter("/preset_tags")
	tag.use(userGuard)
	tag.post("/{tagId}/like", roPresetHandler.LikeTag)
	tag.delete("/{tagId}/like", roPresetHandler.UnLikeTag)

	appPort := appConfig.Port
	log.Printf("listening on localhost:%v\n", appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", appPort), r.router))
}

func initSessionStore() {
	key := appConfig.Auth.AuthenticationSessionSecret // Replace with your SESSION_SECRET or similar
	maxAge := int(time.Hour * 6)                      //
	isProd := appConfig.Environment == "prod"         // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
}

func initAuthIdentityProviders() {
	ggClientId := appConfig.GoogleAuth.ClientId
	ggClientSecret := appConfig.GoogleAuth.ClientSecret
	goth.UseProviders(
		google.New(ggClientId, ggClientSecret, appConfig.GoogleAuth.CallbackUrl, "email"),
	)
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(fmt.Errorf("fatal error init time zone: %w", err))
	}

	time.Local = ict
}

func jsonResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func setSecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Max-Age", "86400")
		w.Header().Add("Access-Control-Allow-Origin", appConfig.Security.AllowedOrigin)
		w.Header().Add("Access-Control-Allow-Methods", strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, ","))
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func connectMongoDB() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.Mongodb.ConnectionStr))
	if err != nil {
		panic(fmt.Errorf("fatal error connect DB: %w", err))
	}

	mongoDb := client.Database(appConfig.Mongodb.DbName)
	userCollection = mongoDb.Collection("users")
	_, err = userCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		panic(fmt.Errorf("index users: %w", err))
	}

	authDataCollection = mongoDb.Collection("authorization_codes")
	_, err = authDataCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"code": 1,
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("index authorization_codes: %w", err))
	}

	refreshTokenCollection = mongoDb.Collection("refresh_tokens")

	roPresetCollection = mongoDb.Collection("ro_presets")
	_, err = roPresetCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"id": 1,
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("index ro_presets: %w", err))
	}

	roTagCollection = mongoDb.Collection("preset_tags")
	_, err = roTagCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tag", Value: 1},
				{Key: "class_id", Value: 1},
			},
		},
		{
			Keys: bson.M{
				"preset_id": 1,
			},
		},
		{
			Keys: bson.D{
				{Key: "preset_id", Value: 1},
				{Key: "tag", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "total_like", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("index preset_tags: %w", err))
	}

	return
}
