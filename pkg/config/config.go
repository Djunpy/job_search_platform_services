package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	UsersMrcUrl string

	Environment string `mapstructure:"NODE_ENV"`
	// db
	PostgresDriver string `mapstructure:"POSTGRES_DRIVER"`
	PostgresSource string
	DbName         string
	DbUser         string
	DbPassword     string
	DbHost         string
	DbPort         string
	DbSSLMode      string
	MigrationURL   string

	Port string `mapstructure:"PORT"`

	Origin string

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	SessionDuration        int           `mapstructure:"SESSION_DURATION"`

	HTTPServerAddress string
	HTTPClientAddress string `mapstructure:"HTTP_CLIENT_ADDRESS"`

	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`

	// SMTP
	SMTPAuthAddress     string `mapstructure:"SMTP_AUTH_ADDRESS"`
	SMTPServerAddress   string `mapstructure:"SMTP_SERVER_ADDRESS"`
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`

	// Redis
	RedisAddress string
}

func LoadConfig(path, serviceName string) (config Config, err error) {
	serviceNameToUpper := strings.ToUpper(serviceName)
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	nodeEnv := config.Environment

	// services
	//originUsersMrcAddress := viper.GetString(fmt.Sprintf("%s_%s_HTTP_SERVER_ADDRESS", nodeEnv, "USERS_MRC"))
	//config.HttpUsersMrcAddress = addHttpPrefix(originUsersMrcAddress)
	config.UsersMrcUrl = viper.GetString(fmt.Sprintf("%s_%s_URL", nodeEnv, "USERS_MRC"))

	// redis
	config.RedisAddress = viper.GetString(fmt.Sprintf("%s_%s_REDIS_ADDRESS", nodeEnv, serviceNameToUpper))

	// postgres
	config.DbName = viper.GetString(fmt.Sprintf("%s_%s_POSTGRES_DB", nodeEnv, serviceNameToUpper))
	config.DbUser = viper.GetString(fmt.Sprintf("%s_%s_POSTGRES_USER", nodeEnv, serviceNameToUpper))
	config.DbPassword = viper.GetString(fmt.Sprintf("%s_%s_POSTGRES_PASSWORD", nodeEnv, serviceNameToUpper))
	config.DbHost = viper.GetString(fmt.Sprintf("%s_%s_POSTGRES_HOST", nodeEnv, serviceNameToUpper))
	config.DbPort = viper.GetString(fmt.Sprintf("%s_%s_POSTGRES_PORT", nodeEnv, serviceNameToUpper))
	config.DbSSLMode = viper.GetString(fmt.Sprintf("%s_SSL_MODE", nodeEnv))

	config.MigrationURL = viper.GetString(fmt.Sprintf("%s_MIGRATION_URL", serviceNameToUpper))

	// server
	config.HTTPServerAddress = viper.GetString(fmt.Sprintf("%s_%s_HTTP_SERVER_ADDRESS", nodeEnv, serviceNameToUpper))
	config.Origin = viper.GetString(fmt.Sprintf("%s_ORIGIN", nodeEnv))
	config.PostgresSource = fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s", "postgres", config.DbUser, config.DbPassword,
		config.DbHost, config.DbPort, config.DbName, config.DbSSLMode)

	return
}

func addHttpPrefix(address string) string {
	return "http://" + address
}
