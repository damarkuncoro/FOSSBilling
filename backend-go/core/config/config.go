package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv          string
	Port            string
	DatabaseURL     string
	RedisURL        string
	JWTSecret       string
	AppURL          string
	DefaultCurrency string
	AllowedOrigins  string
	WorkerConcurrency int

	// Mailer
	MailDriver   string
	MailHost     string
	MailPort     string
	MailUser     string
	MailPass     string
	MailFromAddr string
	MailFromName string

	// Payment Gateways
	MidtransServerKey string
	MidtransClientKey string
	StripeSecretKey   string
	StripePublicKey   string

	TelegramBotToken string
	TelegramChatID   string
}

func Load() *Config {
	return &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/fossbilling?sslmode=disable"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:       getEnv("JWT_SECRET", "super-secret-default-key-change-me"),
		AppURL:          getEnv("APP_URL", "http://localhost:8080"),
		DefaultCurrency: getEnv("DEFAULT_CURRENCY", "USD"),
		AllowedOrigins:  getEnv("ALLOWED_ORIGINS", "*"),
		WorkerConcurrency: getEnvInt("WORKER_CONCURRENCY", 20),

		MailDriver:   getEnv("MAIL_DRIVER", "mock"),
		MailHost:     getEnv("MAIL_HOST", "smtp.mailtrap.io"),
		MailPort:     getEnv("MAIL_PORT", "2525"),
		MailUser:     getEnv("MAIL_USER", ""),
		MailPass:     getEnv("MAIL_PASS", ""),
		MailFromAddr: getEnv("MAIL_FROM_ADDRESS", "noreply@fossbilling.org"),
		MailFromName: getEnv("MAIL_FROM_NAME", "FOSSBilling"),

		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
		StripeSecretKey:   getEnv("STRIPE_SECRET_KEY", ""),
		StripePublicKey:   getEnv("STRIPE_PUBLIC_KEY", ""),

		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:   getEnv("TELEGRAM_CHAT_ID", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
