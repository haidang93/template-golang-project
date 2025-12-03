package server

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDatabasePool(connStr string) *pgxpool.Pool {
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Failed to create database config:", err)
	}

	Pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database!")

	return Pool
}
