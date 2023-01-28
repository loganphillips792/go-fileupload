package main

// https://gist.github.com/subfuzion/08c5d85437d5d4f00e58

// https://freshman.tech/file-upload-golang/

// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

// https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server

// https://stackoverflow.com/questions/21948243/how-can-i-post-files-and-json-data-together-with-curl

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/loganphillips792/fileupload/api"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {

	db := initializeDatabase()
	defer db.Close()

	logger, _ := zap.NewProduction()
	//logger.Sync() // flushes buffer, if any
	err := logger.Sync() // flushes buffer, if any
	// for linting
	if err != nil {
		log.Print("Error when encoding json")
	}

	sugar := logger.Sugar()

	envHandler := &api.Handler{Logger: sugar, DbConn: db}
	/*
		sugar.Infow("failed to fetch URL",
			// Structured context as loosely typed key-value pairs.
			"url", url,
			"attempt", 3,
			"backoff", time.Second,
		)
		sugar.Infof("Failed to fetch URL: %s", url)
	*/

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
	}))

	e.GET("/hello", api.HelloWorld)
	e.GET("/images/", envHandler.GetAllFiles)
	e.POST("/uploadfile/", envHandler.UploadFileHandler, middleware.BodyLimit("1M")) // Body limit middleware sets the maximum allowed size for a request body, if the size exceeds the configured limit, it sends “413 - Request Entity Too Large” response. The body limit is determined based on both Content-Length request header and actual content read, which makes it super secure
	e.DELETE("/images/:id", envHandler.DeleteImage)
	e.GET("/download_image/", envHandler.DownloadImage)
	e.GET("/download_csv/", envHandler.DownloadCSV)
	e.Logger.Fatal(e.Start(":8000"))
}

func initializeDatabase() *sql.DB {
	log.Print("Initializing SQL Lite database...")
	// TODO: only create if it doesn't exist
	file, err := os.Create("../data.db")

	if err != nil {
		log.Fatal(err.Error())
	}

	file.Close()

	db, err := sql.Open("sqlite3", "../data.db")

	if err != nil {
		log.Fatal(err.Error())
	}

	// run sql files

	c, err := ioutil.ReadFile("../script.sql")

	if err != nil {
		log.Fatal(err.Error())
	}

	sql := string(c)

	_, err = db.Exec(sql)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return db
}
