package config

import "os"

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	//начало
	JWTSecret    string
	JWTAccessTTL string
	//конец
}

func Load() Config {
	return Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable"),
		//начало
		JWTSecret:    getEnv("JWT_SECRET", "change_me_please"),
		JWTAccessTTL: getEnv("JWT_ACCESS_TTL", "1h"),
		//конец
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
