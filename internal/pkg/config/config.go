package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds configuration values required to connect to MinIO.
// Each field is populated from environment variables.
type Config struct {
	HTTPServerPort          int    `mapstructure:"HTTP_SERVER_PORT"`
	ConnectRPCServerAddress int    `mapstructure:"CONNECT_RPC_SERVER_PORT"`
	FileDirectoryName       string `mapstructure:"FILE_DIRECTORY_NAME"`
	MinioHost               string `mapstructure:"MINIO_HOST"`
	MinioAccessKeyID        string `mapstructure:"MINIO_ACCESS_KEY_ID"`
	MinioAccessKey          string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioUseSSL             bool   `mapstructure:"MINIO_USE_SSL"`
	BucketName              string `mapstructure:"BUCKET_NAME"`
}

// NewConfig loads configuration from environment variables and optionally
// a .env file (if present in the working directory).
//
// It returns a Config struct populated with the retrieved values.
// If any error occurs during unmarshaling, the error is returned.
func NewConfig() (*Config, error) {
	// This checks if the application is running in a local development environment
	// by looking for a `.env` file in the repo root (current working directory).
	// If the file exists, it loads environment variables from it using godotenv.
	// This allows local development without requiring environment variables to be set globally.
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting current directory: %w", err)
	}
	// Assume repo root contains .env; adjust relative path as needed
	envPath := filepath.Join(cwd, ".env")
	fmt.Println(">>>> path", envPath)
	// Check if .env exists
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("MINIO_USE_SSL", false)
	viper.SetDefault("HTTP_SERVER_PORT", 3333)
	viper.SetDefault("CONNECT_RPC_SERVER_PORT", 5555)
	viper.SetDefault("BUCKET_NAME", "files")
	viper.AutomaticEnv()

	viper.BindEnv("HTTP_SERVER_PORT")
	viper.BindEnv("CONNECT_RPC_SERVER_PORT")
	viper.BindEnv("FILE_DIRECTORY_NAME")
	viper.BindEnv("MINIO_HOST")
	viper.BindEnv("MINIO_ACCESS_KEY_ID")
	viper.BindEnv("MINIO_ACCESS_KEY")
	viper.BindEnv("MINIO_USE_SSL")

	viper.BindEnv("BUCKET_NAME")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal error: %w", err)
	}
	return &cfg, nil
}
