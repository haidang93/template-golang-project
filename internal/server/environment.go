package server

import (
	"log"

	"github.com/example/internal/config"
)

func LoadEnv() *config.Environment {
	env := config.Env
	if err := env.Init(); err != nil {
		log.Fatal("Failed to create database config:", err)
	}

	return &env
}
