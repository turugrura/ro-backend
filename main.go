package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log"

	"ro-backend/configuration"
	"ro-backend/handler"
	_productHandler "ro-backend/handler/product"
	_storeHandler "ro-backend/handler/store"
	"ro-backend/repository"
	"ro-backend/service"

	"github.com/gorilla/handlers"
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
var roPresetForSummaryCollection *mongo.Collection
var roTagCollection *mongo.Collection
var storeCollection *mongo.Collection
var productCollection *mongo.Collection

var appConfig configuration.AppConfig

var limiter = rate.NewLimiter(1, 30)

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
	var storeRepo = repository.NewStoreRepository(storeCollection)
	var productRepo = repository.NewProductRepository(productCollection)

	var userService = service.NewUserService(userRepo, roPresetRepo)
	var tokenService = service.NewTokenService(refreshTokenRepo)
	var authDataService = service.NewAuthenticationDataService(authDataRepo)
	var roPresetService = service.NewRoPresetService(roPresetRepo, roTagRepo)
	var roTagService = service.NewPresetTagService(roTagRepo, roPresetRepo, userRepo)
	var storeService = service.NewStoreService(storeRepo)
	var productService = service.NewProductService(productRepo, storeRepo)

	var roPresetSummaryRepo = repository.NewRoPresetRepository(roPresetForSummaryCollection)
	var presetSummaryService = service.NewSummaryPresetService(roPresetSummaryRepo)

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
	var presetSummaryHandler = handler.NewPresetSummaryHandler(presetSummaryService)
	var storeHandler = _storeHandler.NewStoreHandler(storeService)
	var productHandler = _productHandler.NewProductHandler(productService)

	var helpCheckHandler = handler.NewHelpCheckHandler()

	r := newAppRouter(mux.NewRouter())
	r.use(jsonResponseMiddleware)
	r.use(rateLimitMiddleware)

	r.get("/ping", helpCheckHandler.Ping)

	r.get("/auth/{provider}/callback", authHandler.AuthenticationCallback)
	r.get("/auth/{provider}", gothic.BeginAuthHandler)

	r.post("/login", authHandler.Login)
	r.post("/refresh_token", authHandler.RefreshToken)

	// ------
	if appConfig.Environment == "dev" {
		admin := r.subRouter("/admin")
		admin.use(adminGuard)
		admin.post("/preset_summary", presetSummaryHandler.GenerateSummary)
	}

	// ------
	me := r.subRouter("/me")
	me.use(userGuard)
	me.get("", userHandler.GetMyProfile)
	me.post("", userHandler.PatchMyProfile)
	me.post("/logout", authHandler.Logout)
	me.post("/bulk_ro_presets", roPresetHandler.BulkCreatePresets)
	me.get("/ro_entire_presets", roPresetHandler.GetMyEntirePresets)
	me.get("/ro_presets", roPresetHandler.GetMyPresets)
	me.post("/ro_presets", roPresetHandler.CreatePreset)

	me.get("/ro_presets/{presetId}", roPresetHandler.GetMyPresetById)
	me.post("/ro_presets/{presetId}", roPresetHandler.UpdateMyPreset)
	me.delete("/ro_presets/{presetId}", roPresetHandler.DeleteById)
	me.post("/ro_presets/{presetId}/publish", roPresetHandler.PublishMyPreset)
	me.delete("/ro_presets/{presetId}/publish", roPresetHandler.UnPublishMyPreset)

	me.post("/ro_presets/{presetId}/tags", roPresetHandler.BulkOperationTags)
	me.delete("/ro_presets/{presetId}/tags/{tagId}", roPresetHandler.RemoveTags)

	me.post("/store", storeHandler.UpdateStore)
	me.get("/store", storeHandler.FindMyStore)
	me.post("/products/search", productHandler.GetMyProductList)
	me.post("/products/bulk_create", productHandler.CreateProductList)
	me.post("/products/bulk_update", productHandler.UpdateProductList)
	me.post("/products/bulk_patch", productHandler.PatchProductList)
	me.post("/products/bulk_renew_exp_date", productHandler.RenewExpDateProductList)
	me.post("/products/bulk_delete", productHandler.DeleteProductList)

	// ------
	ro := r.subRouter("/ro_presets")
	ro.use(userGuard)
	ro.get("/class_by_tags/{classId}/{tag}", roPresetHandler.SearchPresetTags)

	// ------ store
	store := r.subRouter("/store")
	store.use(userGuard)
	store.post("", storeHandler.CreateStore)
	store.get("/{storeId}", storeHandler.FindStoreById)
	store.post("/{storeId}/review", storeHandler.ReviewStore)

	// ------ product
	product := r.subRouter("/product")
	// product.use(userGuard)
	product.post("/search", productHandler.SearchProductList)

	// ------
	tag := r.subRouter("/preset_tags")
	tag.use(userGuard)
	tag.post("/{tagId}/like", roPresetHandler.LikeTag)
	tag.delete("/{tagId}/like", roPresetHandler.UnLikeTag)

	headersOk := handlers.AllowedHeaders([]string{"authorization", "Content-Type"})
	origins := handlers.AllowedOrigins(appConfig.Security.AllowedOrigins)
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodDelete})
	maxAge := handlers.MaxAge(86400)
	h := handlers.CORS(headersOk, origins, methods, maxAge)(r.router)
	h = handlers.CompressHandler(h)

	appPort := appConfig.Port
	log.Printf("listening on localhost:%v\n", appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", appPort), h))
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
		{
			Keys: bson.M{
				"name": 1,
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
	roPresetForSummaryCollection = mongoDb.Collection("authorization_codes")
	_, err = roPresetCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"id": 1,
			},
		},
		{
			Keys: bson.M{
				"user_id": 1,
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

	storeCollection = mongoDb.Collection("store")
	_, err = storeCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"owner_id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"name": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		panic(fmt.Errorf("index store: %w", err))
	}

	productCollection = mongoDb.Collection("product")
	_, err = productCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "exp_date", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "baht", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "m", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "m", Value: 1},
				{Key: "baht", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "m", Value: 1},
				{Key: "baht", Value: 1},
				{Key: "exp_date", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "m", Value: 1},
				{Key: "baht", Value: 1},
				{Key: "exp_date", Value: 1},
				{Key: "is_published", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "m", Value: 1},
				{Key: "baht", Value: 1},
				{Key: "exp_date", Value: 1},
				{Key: "is_published", Value: -1},
				{Key: "name", Value: 1},
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("index product: %w", err))
	}

	return
}
