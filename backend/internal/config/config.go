package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost   string
	ServerPort   string
	DatabasePath string
	GroqAPIKey   string
	GroqModel    string
	OllamaURL    string
	ChromaURL    string
	LicenseKey   string
	Environment  string
}

func Load() *Config {
	godotenv.Load() // Try to load .env file

	cfg := &Config{
		ServerHost:   getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "./azeru.db"),
		GroqAPIKey:   getEnv("GROQ_API_KEY", ""),
		GroqModel:    getEnv("GROQ_MODEL", "llama3-70b-8192"),
		OllamaURL:    getEnv("OLLAMA_URL", "http://localhost:11434"),
		ChromaURL:    getEnv("CHROMA_URL", "http://localhost:8000"),
		LicenseKey:   getEnv("LICENSE_KEY", ""),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}

	if cfg.GroqAPIKey == "" {
		log.Println("WARNING: GROQ_API_KEY not set")
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
