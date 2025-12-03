package main

import (
	"log"
	"os"
	"strconv"

	"github.com/example/internal/server"
)

func main() {
	env := server.LoadEnv()
	Migration := newMigration(env.POSTGRES_CONNSTR)
	defer Migration.Migration.Close()

	parameter := os.Args[2]

	switch parameter {
	case "up":
		Migration.Up()
	case "down":
		Migration.Down()
	case "step":
		step, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		Migration.Steps(step)
	default:
	}

}
