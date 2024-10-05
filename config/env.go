package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Env holds all the environment variables
type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	URL                    string `mapstructure:"URL"`
	API                    string `mapstructure:"API"`
}

// NewEnv creates and loads environment variables from .env file
func NewEnv() (*Env, error) {
	env := Env{}

	// Configure viper to read .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("../")

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the config into the Env struct
	err = viper.Unmarshal(&env)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	// Check if the application is in development mode
	if env.AppEnv == "development" {
		log.Println("App is running in development mode")
	}

	return &env, nil
}
