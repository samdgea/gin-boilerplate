package db

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/samdgea/gin-boilerplate/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
	var err error

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// URL format avoids lib/pq parsing bugs with empty password in key=value format
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(user), url.QueryEscape(password), host, port, dbname)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Postgres connection error: ", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Postgres Error ping: ", err)
	}

	err = models.MigrateModels(DB)
	if err != nil {
		log.Fatal("Postgres Error migrate: ", err)
	}

	fmt.Println("Successfully connected to Postgres Database")
}
