package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    ClientPassword string
    AdminPassword  string
    JWTSecret      string
    DBPath         string
    Port           string
}

func LoadConfig() Config {
    // Load .env file if it exists
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    return Config{
        ClientPassword: os.Getenv("CLIENT_PASSWORD"),
        AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
        JWTSecret:      os.Getenv("JWT_SECRET"),
        DBPath:         os.Getenv("DB_PATH"),
        Port:           os.Getenv("PORT"),
    }
}