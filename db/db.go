package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"os"

	"github.com/samdgea/gin-boilerplate/models"
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

	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	//fmt.Println(dsn)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port),
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

	err = models.MigrateUserModel(DB)
	if err != nil {
		log.Fatal("Postgres Error migrate: ", err)
	}

	fmt.Println("Successfully connected to Postgres Database")
}
