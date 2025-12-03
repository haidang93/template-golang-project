package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrations struct {
	Migration *migrate.Migrate
}

func newMigration(connStr string) *migrations {
	m, err := migrate.New(
		"file://cmd/db/migrations",
		connStr,
	)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}
	Migration := migrations{
		Migration: m,
	}

	return &Migration
}

func (m *migrations) Up() {
	if err := m.Migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")
}

func (m *migrations) Down() {
	if err := m.Migration.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")
}

func (m *migrations) Steps(step int) {
	if err := m.Migration.Steps(step); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")
}
