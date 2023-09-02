package main

// https://stackoverflow.com/questions/21948243/how-can-i-post-files-and-json-data-together-with-curl

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/loganphillips792/fileupload/api"
	"github.com/loganphillips792/fileupload/config"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	// configuration
	cfg, configError := config.Init()

	if configError != nil {
		log.Fatal("config error")
	}

	// Set up logging
	logger, _ := zap.NewProduction()
	err := logger.Sync() // flushes buffer, if any
	if err != nil {      // for linting
		log.Print("Error when encoding json")
	}

	// db, postgresErr := db.CreatePostgresConnection()

	// if postgresErr != nil {
	// 	logger.Error(postgresErr.Error())
	// }

	db := initializeDatabase()
	defer db.Close()

	sugar := logger.Sugar()
	handler := api.BuildHandler(sugar, db, cfg)

	e := echo.New()
	setupRouter(e, db, handler, sugar, cfg)
	e.Logger.Fatal(e.Start(":8000"))
}

func setupRouter(e *echo.Echo, db *sqlx.DB, handler *api.Handler, sugar *zap.SugaredLogger, cfg *config.AppConf) {

	g := e.Group("/api")
	g.Use(api.ApiMiddleware(db, handler, sugar, cfg))

	// https://echo.labstack.com/middleware/logger/
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
	}))

	g.GET("/hello", api.HelloWorld)

	e.GET("/images/", handler.GetAllFiles)
	e.POST("/uploadfile/", handler.UploadFileHandler, middleware.BodyLimit("1M")) // Body limit middleware sets the maximum allowed size for a request body, if the size exceeds the configured limit, it sends “413 - Request Entity Too Large” response. The body limit is determined based on both Content-Length request header and actual content read, which makes it super secure
	e.DELETE("/images/:id", handler.DeleteImage)
	e.GET("/images/:id", handler.GetImageById)
	e.GET("/download_image/", handler.DownloadImage)
	e.GET("/download_csv/", handler.DownloadCSV)
	e.POST("/register/", handler.Register)
	e.POST("/login/", handler.Login)

}

func initializeDatabase() *sqlx.DB {
	log.Print("Initializing SQL Lite database...")

	// TODO: only create and seed database if it doesn't exist
	file, openFileErr := os.Open("data.db")

	if openFileErr != nil {
		log.Print(openFileErr.Error())
	}

	if errors.Is(openFileErr, os.ErrNotExist) {
		file, _ = os.Create("data.db")
	}

	file.Close()

	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		log.Fatal(err.Error())
	}

	sqlxDb := sqlx.NewDb(db, "sqlite3")

	// run sql files

	if errors.Is(openFileErr, os.ErrNotExist) {
		c, err := ioutil.ReadFile("script.sql")

		if err != nil {
			log.Fatal(err.Error())
		}

		sql := string(c)

		_, err = db.Exec(sql)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	// db, err := sqlx.Open("sqlite3", "data.db")

	// if err != nil {
	// 	return nil
	// }

	// err = db.Ping()

	// if err != nil {
	// 	db.Close()
	// 	return nil
	// }

	// }

	return sqlxDb
}
