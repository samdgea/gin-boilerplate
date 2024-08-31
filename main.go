package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/samdgea/gin-boilerplate/db"
	"github.com/samdgea/gin-boilerplate/middlewares"
	"github.com/samdgea/gin-boilerplate/routes"
	"os"
)

func main() {
	db.InitPostgres()

	isProd := os.Getenv("APP_IS_PROD")
	if isProd == "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	if isProd != "true" {
		r.ForwardedByClientIP = true
		_ = r.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
	}

	r.Use(middlewares.CORS())

	// API routing
	api := r.Group("/api")

	routes.SetupAuthRoutes(api.Group("/auth"))

	// Start
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
