package main

import (
	"context"
	"fmt"
	"net/http"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection
var authDataCollection *mongo.Collection
var refreshTokenCollection *mongo.Collection
var roPresetCollection *mongo.Collection

var appConfig configuration.AppConfig

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

	var userService = service.NewUserService(userRepo)
	var tokenService = service.NewTokenService(refreshTokenRepo)
	var authDataService = service.NewAuthenticationDataService(authDataRepo)
	var roPresetService = service.NewRoPresetService(roPresetRepo)

	var authHandler = handler.NewAuthHandler(handler.AuthHandlerParam{
		UserService:               userService,
		AuthenticationDataService: authDataService,
		TokenService:              tokenService,
	})
	var userHandler = handler.NewUserHandler(userService)
	var roPresetHandler = handler.NewRoPresetHandler(handler.RoPresetHandlerParam{
		RoPresetService: roPresetService,
		UserService:     userService,
	})

	r := newAppRouter(mux.NewRouter())
	r.use(jsonResponseMiddleware)

	r.get("/auth/{provider}/callback", authHandler.AuthenticationCallback)
	r.get("/auth/{provider}", gothic.BeginAuthHandler)

	r.post("/login", authHandler.Login)
	r.post("/refresh_token", authHandler.RefreshToken)

	me := r.subRouter("/me")
	me.use(userGuard)
	me.get("", userHandler.GetMyProfile)
	me.get("/ro_presets", roPresetHandler.GetMyPresets)

	ro := r.subRouter("/ro_presets")
	ro.use(userGuard)
	ro.post("", roPresetHandler.CreatePreset)
	ro.post("/bulk", roPresetHandler.BulkCreatePresets)
	ro.get("/{presetId}", roPresetHandler.GetPresetById)
	ro.post("/{presetId}", roPresetHandler.UpdatePreset)
	ro.delete("/{presetId}", roPresetHandler.DeleteById)

	appPort := appConfig.Port
	log.Printf("listening on localhost:%v\n", appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", appPort), r.router))
}

func initSessionStore() {
	key := appConfig.Auth.AuthenticationSessionSecret // Replace with your SESSION_SECRET or similar
	maxAge := int(time.Hour * 6)                      //
	isProd := false                                   // Set to true when serving over https

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

func connectMongoDB() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.Mongodb.ConnectionStr))
	if err != nil {
		panic(fmt.Errorf("fatal error connect DB: %w", err))
	}

	mongoDb := client.Database(appConfig.Mongodb.DbName)
	userCollection = mongoDb.Collection("users")
	authDataCollection = mongoDb.Collection("authorization_codes")
	refreshTokenCollection = mongoDb.Collection("refresh_tokens")
	roPresetCollection = mongoDb.Collection("ro_presets")

	return
}
