package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ClientKey    string
	ClientSecret string
	RedirectURI  string
	ServerPort   string
	Debug        bool
	AuthURL      = "https://www.tiktok.com/v2/auth/authorize/"
	TokenURL     = "https://open.tiktokapis.com/v2/oauth/token/"
)

func LoadConfig() {
	// .env dosyasını yükle (varsa)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ClientKey = getEnv("TIKTOK_CLIENT_KEY", "")
	ClientSecret = getEnv("TIKTOK_CLIENT_SECRET", "")
	RedirectURI = getEnv("TIKTOK_REDIRECT_URI", "http://localhost:8080/callback")
	ServerPort = getEnv("SERVER_PORT", "8080")
	Debug = getEnv("DEBUG", "false") == "true"

	if ClientKey == "" || ClientSecret == "" {
		log.Fatal("TIKTOK_CLIENT_KEY and TIKTOK_CLIENT_SECRET must be set")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DebugLog prints debug message if DEBUG is enabled
func DebugLog(format string, args ...interface{}) {
	if Debug {
		log.Printf("[DEBUG] "+format, args...)
	}
}
