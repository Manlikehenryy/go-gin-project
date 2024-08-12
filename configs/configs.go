package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	MONGO_DB_URI string
	MODE         string
	JWT_SECRET   string
	APP_URL      string
}

var Env *Config

func Init() {
	
	Env = &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	Env.PORT = os.Getenv("PORT")
	Env.MONGO_DB_URI = os.Getenv("MONGO_DB_URI")
	Env.MODE = os.Getenv("MODE")
	Env.JWT_SECRET = os.Getenv("JWT_SECRET")
	Env.APP_URL = os.Getenv("APP_URL")
}
