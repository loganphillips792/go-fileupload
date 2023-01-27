package main

// https://gist.github.com/subfuzion/08c5d85437d5d4f00e58

// https://freshman.tech/file-upload-golang/

// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

// https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server

// https://stackoverflow.com/questions/21948243/how-can-i-post-files-and-json-data-together-with-curl

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.SugaredLogger
	dbConn *sql.DB
}

type FileInfo struct {
	Name string
}

type Image struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

func main() {

	db := initializeDatabase()
	defer db.Close()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	envHandler := &Handler{logger: sugar, dbConn: db}
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

	e.GET("/hello", HelloWorld)
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
	file, err := os.Create("data.db")

	if err != nil {
		log.Fatal(err.Error())
	}

	file.Close()

	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal(err.Error())
	}

	// run sql files

	c, err := ioutil.ReadFile("./script.sql")

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

func (handler *Handler) UploadFileHandler(c echo.Context) error {
	handler.logger.Infof("Content Length %d ", c.Request().ContentLength)

	file, err := c.FormFile("file")

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	// Create the uploads fodler if it doesn't already exist
	err = os.MkdirAll("./uploads", os.ModePerm)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	handler.logger.Infof("File name from user %s ", c.FormValue("file_name"))

	// Create a new file in the uploads directory
	filePath := ""
	fileNameToInsertIntoDatabase := ""
	if c.FormValue("file_name") != "" {
		filePath = fmt.Sprintf("./uploads/%s%s", c.FormValue("file_name"), filepath.Ext(file.Filename))
		fileNameToInsertIntoDatabase = c.FormValue("file_name")
	} else {
		currentUnixTime := time.Now().UnixNano()
		fileNameToInsertIntoDatabase = strconv.FormatInt(currentUnixTime, 10)
		filePath = fmt.Sprintf("./uploads/%d%s", currentUnixTime, filepath.Ext(file.Filename))
	}

	dst, err := os.Create(filePath)
	fmt.Println(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// file save was successsful
	query := "INSERT INTO images (name, file_path) VALUES (?, ?)"

	handler.logger.Infow("Running SQL statement",
		"SQL", query,
	)

	_, err = handler.dbConn.Exec(query, fileNameToInsertIntoDatabase, filePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	changeImageToBlackAndWhite()

	return c.Blob(http.StatusOK, "application/json", []byte(`{"response":"Upload Successful!!"}`))

	// check that the file is only image file
	// https://freshman.tech/file-upload-golang/#restrict-the-type-of-the-uploaded-file

	// buff := make([]byte, 512)
	// _, err = file.Read(buff)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// filetype := http.DetectContentType(buff)

	// if filetype != "image/jpeg" && filetype != "image/png" {
	// 	http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
	// 	return
	// }

	// handler.logger.Infof("File type is %s ", filetype)

}

// When user successfully uploads image, they can click "convert to black and white". The new image will
// show as a thumbnail, and then they can click download to download the new image
// https://stackoverflow.com/questions/42516203/converting-rgba-image-to-grayscale-golang
func changeImageToBlackAndWhite() {
	fmt.Println("Converting image to black and white...")
}

func (handler *Handler) GetAllFiles(c echo.Context) error {
	handler.logger.Info("Retreiving all images....")

	searchParams := c.QueryParam("q")
	handler.logger.Infow("Search parameters passed in",
		"PARAMS", searchParams,
	)

	query := "SELECT * FROM images"

	if searchParams != "" {
		query += " LIKE name "
	}

	handler.logger.Infow("Running SQL statement",
		"SQL", query,
	)

	rows, err := handler.dbConn.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var images []Image
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.Id, &image.Name, &image.FilePath)

		if err != nil {
			log.Fatal(err)
		}

		images = append(images, image)
	}

	if len(images) == 0 {
		// Is this needed? Or will 0 automatically be handled?
		return c.JSON(http.StatusOK, make([]Image, 0))
	}

	return c.JSON(http.StatusOK, images)
}

func (handler *Handler) DeleteImage(c echo.Context) error {
	handler.logger.Info("Deleting image....")

	id := c.Param("id")

	query := "DELETE FROM images WHERE id = ?"

	handler.logger.Infow("Running SQL statement",
		"id of image delete", id,
		"SQL", query,
	)

	resp, err := handler.dbConn.Exec(query, id)

	if err != nil {
		log.Fatal(err.Error())
	}

	rowsDeleted, _ := resp.RowsAffected()

	if rowsDeleted == 0 {
		c.Blob(http.StatusNotFound, "application/json", []byte(`{"response":"image not found"}`))
	}

	return c.NoContent(http.StatusNoContent)
}

// https://github.com/labstack/echo/blob/v3.3.10/context.go#L542
func (handler *Handler) DownloadImage(c echo.Context) error {
	return c.Attachment("./data/IMG_7015.jpg", "download.jpg")
}

// send csv to client to automatically download
// https://medium.com/wesionary-team/create-csv-file-in-go-server-and-download-from-reactjs-4f22f148290b
// https://stackoverflow.com/questions/68162651/go-how-to-response-csv-file
// https://medium.com/wesionary-team/create-csv-file-in-go-server-and-download-from-reactjs-4f22f148290b
func (handler *Handler) DownloadCSV(c echo.Context) error {
	return c.Attachment("./data/airtravel.csv", "download.csv")
}

/*

// https://www.reddit.com/r/reactjs/comments/5xgdzh/how_to_correctly_store_user_information/
// https://www.reddit.com/r/reactjs/comments/gek8as/recommended_approach_to_check_if_user_is/
// https://stackoverflow.com/questions/49819183/react-what-is-the-best-way-to-handle-login-and-authentication
// https://stackoverflow.com/questions/70686434/how-to-save-the-users-authorization-in-session-in-react
// https://www.reddit.com/r/programming/comments/nag1cu/jwt_should_not_be_your_default_for_sessions/
// https://www.sohamkamani.com/golang/session-cookie-authentication/
func Login() {

}
*/
func HelloWorld(c echo.Context) error {
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"status":"OK"}`))
	data := []byte(`{"status":"OK!"}`)
	return c.Blob(http.StatusOK, "application/json", data)
}
