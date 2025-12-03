package main

import (
	"github.com/example/internal/pkg/validate"
	"github.com/example/internal/server"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	env := server.LoadEnv()
	httpClient := server.CreateHttpClient()
	DB := server.CreateDatabasePool(env.POSTGRES_CONNSTR)
	defer DB.Close()

	// Create Echo instance
	e := echo.New()

	// setup system middleware
	if !env.IsDev() {
		// Recover from panic to stop the system crash
		e.Use(middleware.Recover())
	}
	// Log incomming request
	e.Use(middleware.Logger())
	e.Validator = &validate.CustomValidator{Validator: validator.New()}

	// initialize all modules and register paths
	server.CreateModule(e, DB, env, httpClient)

	if env.IsLocal() {
		e.Logger.Fatal(e.Start("localhost:" + env.PORT))
	} else {
		e.Logger.Fatal(e.Start(":" + env.PORT))
	}
}
