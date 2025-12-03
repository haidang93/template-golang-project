package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var Env Environment = Environment{}

type Environment struct {
	PORT                              string
	ENV                               string
	VERSION                           string
	POSTGRES_CONNSTR                  string
	REDIS_HOST                        string
	REDIS_USERNAME                    string
	REDIS_PORT                        string
	REDIS_PASSWORD                    string
	JWT_REFRESH_TOKEN_EXPIRATION_TIME string
	JWT_SECRET                        string
	SUPABASE_BUCKET_NAME              string
	SUPABASE_BUCKET_URL               string
	SUPABASE_BUCKET_ACCESS_TOKEN      string
	SUPABASE_BUCKET_ACCESS_KEY        string
	SUPABASE_BUCKET_SECRET_ACCESS_KEY string
	MAILTRAP_API_TOKEN                string
	MAILTRAP_SEND_URL                 string
}

func (e *Environment) IsDev() bool {
	return e.ENV != "PRODUCTION"
}

func (e *Environment) IsLocal() bool {
	return e.ENV != "PRODUCTION"
}

func (env *Environment) Init() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	env.ENV = os.Getenv("ENV")
	env.PORT = os.Getenv("PORT")
	env.VERSION = os.Getenv("VERSION")
	env.POSTGRES_CONNSTR = os.Getenv("POSTGRES_CONNSTR")
	env.REDIS_HOST = os.Getenv("REDIS_HOST")
	env.REDIS_USERNAME = os.Getenv("REDIS_USERNAME")
	env.REDIS_PORT = os.Getenv("REDIS_PORT")
	env.REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	env.JWT_REFRESH_TOKEN_EXPIRATION_TIME = os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION_TIME")
	env.JWT_SECRET = os.Getenv("JWT_SECRET")
	env.SUPABASE_BUCKET_NAME = os.Getenv("SUPABASE_BUCKET_NAME")
	env.SUPABASE_BUCKET_URL = os.Getenv("SUPABASE_BUCKET_URL")
	env.SUPABASE_BUCKET_ACCESS_TOKEN = os.Getenv("SUPABASE_BUCKET_ACCESS_TOKEN")
	env.SUPABASE_BUCKET_ACCESS_KEY = os.Getenv("SUPABASE_BUCKET_ACCESS_KEY")
	env.SUPABASE_BUCKET_SECRET_ACCESS_KEY = os.Getenv("SUPABASE_BUCKET_SECRET_ACCESS_KEY")
	env.MAILTRAP_API_TOKEN = os.Getenv("MAILTRAP_API_TOKEN")
	env.MAILTRAP_SEND_URL = os.Getenv("MAILTRAP_SEND_URL")

	return nil
}
