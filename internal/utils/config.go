package utils

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

// Configuration file structure
type Config struct {
	Server struct {
		Port uint `mapstructure:"port"`
	} `mapstructure:"server"`
	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     uint   `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
		Options  string `mapstructure:"options"`
	} `mapstructure:"postgres"`
	Migrations string `mapstructure:"migrations"`
	Rabbit     struct {
		Host     string `mapstructure:"host"`
		Port     uint   `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"rabbit"`
	Auth0 struct {
		Issuer  string `mapstructure:"issuer"`
		Backend struct {
			ClientID string `mapstructure:"clientid"`
		} `mapstructure:"backend"`
		Management struct {
			ClientID     string `mapstructure:"clientid"`
			ClientSecret string `mapstructure:"clientsecret"`
			Audience     string `mapstructure:"audience"`
		} `mapstructure:"management"`
	} `mapstructure:"auth0"`
	Keycloak struct {
		BasePath     string `mapstructure:"basepath"`
		realm        string `mapstructure:"realm"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
		ClientID     string `mapstructure:"clientid"`
		ClientSecret string `mapstructure:"clientsecret"`
	} `mapstructure:"keycloak"`
}

// Set default values
func setDefaultValues() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("postgres.host", "postgres")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.database", "mydb")
	viper.SetDefault("postgres.username", "myuser")
	viper.SetDefault("postgres.password", "mypassword")
	viper.SetDefault("postgres.options", "sslmode=disable")
	viper.SetDefault("migrations", "file://migrations")
}

// Singleton config
var config *Config

// Eager configuration Initialisation
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("~/.config/api")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set default values
	setDefaultValues()

	// Read in the config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config.yaml found")
		} else {
			log.Panic(err)
		}
	}

	// Unpack the config attributes
	err := viper.Unmarshal(&config)

	// Error check
	if err != nil {
		panic("Failed to unmarshal configuration")
	}

	// Print configuration
	// fmt.Println(PrettyJSON(config))
}

// Get the config singleton
func GetConfig() *Config {
	return config
}
