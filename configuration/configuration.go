package configuration

import (
	"fmt"

	"github.com/spf13/viper"
)

type MongoDbConfig struct {
	ConnectionStr string
	DbName        string
}

type AuthProviderConfig struct {
	ClientId     string
	ClientSecret string
	CallbackUrl  string
}

type AuthConfig struct {
	AuthenticationSessionSecret   string
	PostAuthenticationRedirectUrl string
}

type JwtConfig struct {
	Secret                     string
	AccessTokenPeriodInMinutes int
	RefreshTokenPeriodInDays   int
}

type RoConfig struct {
	PresetLimit int
}

type SecurityConfig struct {
	AllowedOrigin string
}

type AppConfig struct {
	Port int
	// dev, prod
	Environment string
	Security    SecurityConfig
	Mongodb     MongoDbConfig
	Auth        AuthConfig
	GoogleAuth  AuthProviderConfig
	Jwt         JwtConfig
	Ro          RoConfig
}

var Config *AppConfig

func getAppConfig() AppConfig {
	if Config == nil {
		Config = &AppConfig{
			Port:        viper.GetInt("app.port"),
			Environment: viper.GetString("app.env"),
			Security: SecurityConfig{
				AllowedOrigin: viper.GetString("security.allowOrigin"),
			},
			Mongodb: MongoDbConfig{
				ConnectionStr: viper.GetString("mongodb.connectionString"),
				DbName:        viper.GetString("mongodb.dbName"),
			},
			Auth: AuthConfig{
				AuthenticationSessionSecret:   viper.GetString("auth.authenticationSessionSecret"),
				PostAuthenticationRedirectUrl: viper.GetString("auth.postAuthenticationRedirectUrl"),
			},
			GoogleAuth: AuthProviderConfig{
				ClientId:     viper.GetString("authProvider.google.clientId"),
				ClientSecret: viper.GetString("authProvider.google.clientSecret"),
				CallbackUrl:  viper.GetString("authProvider.google.callbackUrl"),
			},
			Jwt: JwtConfig{
				Secret:                     viper.GetString("jwt.secret"),
				AccessTokenPeriodInMinutes: viper.GetInt("jwt.accessTokenPeriodInMinutes"),
				RefreshTokenPeriodInDays:   viper.GetInt("jwt.refreshTokenPeriodInDays"),
			},
			Ro: RoConfig{
				PresetLimit: viper.GetInt("ro.preset.limitPerUser"),
			},
		}
	}

	return *Config
}

func InitAppConfig() AppConfig {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	viper.AddConfigPath(".") // optionally look for config in the working directory
	// viper.AutomaticEnv()
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return getAppConfig()
}
