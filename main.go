package main

import (
	"fmt"
	"net/http"
	"time"

	"log"

	"ro-backend/api_router"
	"ro-backend/configuration"
	"ro-backend/handler"
	"ro-backend/repository"
	"ro-backend/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/time/rate"
)

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
	// var storeRepo = repository.NewStoreRepository(storeCollection)
	// var productRepo = repository.NewProductRepository(productCollection)

	var userService = service.NewUserService(userRepo, roPresetRepo)
	var tokenService = service.NewTokenService(refreshTokenRepo)
	var authDataService = service.NewAuthenticationDataService(authDataRepo)
	var roPresetService = service.NewRoPresetService(roPresetRepo, roTagRepo)
	var roTagService = service.NewPresetTagService(roTagRepo, roPresetRepo, userRepo)
	// var storeService = service.NewStoreService(storeRepo)
	// var productService = service.NewProductService(productRepo, storeRepo)

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
	// var storeHandler = _storeHandler.NewStoreHandler(storeService)
	// var productHandler = _productHandler.NewProductHandler(productService)

	var helpCheckHandler = handler.NewHelpCheckHandler()

	r := api_router.NewAppRouter(mux.NewRouter())
	r.Use(jsonResponseMiddleware)
	r.Use(rateLimitMiddleware)

	r.Get("/ping", helpCheckHandler.Ping)

	r.Get("/auth/{provider}/callback", authHandler.AuthenticationCallback)
	r.Get("/auth/{provider}", gothic.BeginAuthHandler)

	r.Post("/login", authHandler.Login)
	r.Post("/refresh_token", authHandler.RefreshToken)

	// ------
	admin := r.SubRouter("/admin")
	admin.Use(adminGuard)
	if appConfig.Environment == "dev" {
		admin.Post("/preset_summary", presetSummaryHandler.GenerateSummary)
	}
	api_router.SetupRouterFriend(friendTranslatorCollection, admin)

	// ------
	me := r.SubRouter("/me")
	me.Use(userGuard)
	me.Get("", userHandler.GetMyProfile)
	me.Post("", userHandler.PatchMyProfile)
	me.Post("/logout", authHandler.Logout)
	me.Post("/bulk_ro_presets", roPresetHandler.BulkCreatePresets)
	me.Get("/ro_entire_presets", roPresetHandler.GetMyEntirePresets)
	me.Get("/ro_presets", roPresetHandler.GetMyPresets)
	me.Post("/ro_presets", roPresetHandler.CreatePreset)

	me.Get("/ro_presets/{presetId}", roPresetHandler.GetMyPresetById)
	me.Post("/ro_presets/{presetId}", roPresetHandler.UpdateMyPreset)
	me.Delete("/ro_presets/{presetId}", roPresetHandler.DeleteById)
	me.Post("/ro_presets/{presetId}/publish", roPresetHandler.PublishMyPreset)
	me.Delete("/ro_presets/{presetId}/publish", roPresetHandler.UnPublishMyPreset)

	me.Post("/ro_presets/{presetId}/tags", roPresetHandler.BulkOperationTags)
	me.Delete("/ro_presets/{presetId}/tags/{tagId}", roPresetHandler.RemoveTags)

	// me.Post("/store", storeHandler.UpdateStore)
	// me.Get("/store", storeHandler.FindMyStore)
	// me.Post("/products/search", productHandler.GetMyProductList)
	// me.Post("/products/bulk_create", productHandler.CreateProductList)
	// me.Post("/products/bulk_update", productHandler.UpdateProductList)
	// me.Post("/products/bulk_patch", productHandler.PatchProductList)
	// me.Post("/products/bulk_renew_exp_date", productHandler.RenewExpDateProductList)
	// me.Post("/products/bulk_delete", productHandler.DeleteProductList)

	// ------
	ro := r.SubRouter("/ro_presets")
	ro.Use(userGuard)
	ro.Get("/class_by_tags/{classId}/{tag}", roPresetHandler.SearchPresetTags)

	// ------ store
	// store := r.SubRouter("/store")
	// store.use(userGuard)
	// store.Post("", storeHandler.CreateStore)
	// store.Get("/{storeId}", storeHandler.FindStoreById)
	// store.Post("/{storeId}/review", storeHandler.ReviewStore)

	// ------ product
	// product := r.SubRouter("/product")
	// product.use(userGuard)
	// product.Post("/search", productHandler.SearchProductList)

	// ------
	tag := r.SubRouter("/preset_tags")
	tag.Use(userGuard)
	tag.Post("/{tagId}/like", roPresetHandler.LikeTag)
	tag.Delete("/{tagId}/like", roPresetHandler.UnLikeTag)

	headersOk := handlers.AllowedHeaders([]string{"authorization", "Content-Type"})
	origins := handlers.AllowedOrigins(appConfig.Security.AllowedOrigins)
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodDelete})
	maxAge := handlers.MaxAge(86400)
	h := handlers.CORS(headersOk, origins, methods, maxAge)(r.Router)
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
