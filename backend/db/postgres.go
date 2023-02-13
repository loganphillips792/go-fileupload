package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreatePostgresConnection() (*sqlx.DB, error) {
	connectionStr, err := createPostgresConnectionUrl()

	if err != nil {
		return nil, err
	}

	// db, err := sql.Open("postgres", connectionStr)
	db, err := sqlx.Connect("postgres", connectionStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createPostgresConnectionUrl() (string, error) {
	host, exists := os.LookupEnv("DB_HOST")

	if !exists {
		return "", errors.New("DB_HOST NOT SET")
	}

	username, exists := os.LookupEnv("DB_USERNAME")

	if !exists {
		return "", errors.New("DB_USERNAME NOT SET")
	}

	password, exists := os.LookupEnv("DB_PASSWORD")

	if !exists {
		return "", errors.New("DB_PASSWORD NOT SET")
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		return "", errors.New("DB_NAME NOT SET")
	}

	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		return "", errors.New("DB_PORT NOT SET")
	}

	fmt.Println("URL")
	fmt.Println(fmt.Printf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, dbPort, dbName))
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, dbPort, dbName), nil
}
