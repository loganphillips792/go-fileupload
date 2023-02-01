package main

// https://gist.github.com/subfuzion/08c5d85437d5d4f00e58

// https://freshman.tech/file-upload-golang/

// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

// https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server

// https://stackoverflow.com/questions/21948243/how-can-i-post-files-and-json-data-together-with-curl

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/loganphillips792/fileupload/api"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {

	// get config
	godotenv.Load(".env")

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

	g := e.Group("/api")
	g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:user_session",
		Validator: func(key string, c echo.Context) (bool, error) {
			sugar.Info("Validating in Middlware...")
			var hashKey = []byte(os.Getenv("GORILLA_SESSIONS_HASH_KEY"))   // encode value
			var blockKey = []byte(os.Getenv("GORILLA_SESSIONS_BLOCK_KEY")) // encrypt value
			var s = securecookie.New(hashKey, blockKey)
			value := make(map[string]string)

			err = s.Decode("user_session", key, &value)

			if err != nil {
				sugar.Errorw("Error when decoding cookie value", err)
				return false, errors.New("authentication failed. Please login again")
			}
			sugar.Infow("Decryption", "The decrypted value is", value["sessionId"])

			// we now have the decrypted session id. We will now look it up in the sessions table
			query := "SELECT * FROM sessions where session_id = ?"

			// Check if username and password exist
			var sessionId string
			var sessionData string
			errFromScan := db.QueryRow(query, value["sessionId"]).Scan(&sessionId, &sessionData)

			if errFromScan != nil {
				log.Print(errFromScan)
			}

			if sessionId == value["sessionId"] {
				sugar.Info("Middlware: Session successfully validated. Coninuting processing")
				return true, nil
			} else {
				return false, errors.New("authentication failed. Please login again")
			}
		},
	}))
	// g.Use(ProcessRequest)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
	}))

	g.GET("/hello", api.HelloWorld)

	e.GET("/images/", envHandler.GetAllFiles)
	e.POST("/uploadfile/", envHandler.UploadFileHandler, middleware.BodyLimit("1M")) // Body limit middleware sets the maximum allowed size for a request body, if the size exceeds the configured limit, it sends “413 - Request Entity Too Large” response. The body limit is determined based on both Content-Length request header and actual content read, which makes it super secure
	e.DELETE("/images/:id", envHandler.DeleteImage)
	e.GET("/test", envHandler.GetImageByPath)
	e.GET("/download_image/", envHandler.DownloadImage)
	e.GET("/download_csv/", envHandler.DownloadCSV)
	e.POST("/register/", envHandler.Register)
	e.POST("/login/", envHandler.Login)

	e.Logger.Fatal(e.Start(":8000"))
}

// midddlware function
// func ProcessRequest(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		fmt.Println("PROCESSING REQUEST MIDDLEWARE")
// 		if err := next(c); err != nil {
// 			c.Error(err)
// 		}

// 		return nil
// 	}
// }

func initializeDatabase() *sql.DB {
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

	return db
}
