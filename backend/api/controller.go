package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	utils "github.com/loganphillips792/fileupload"
	"github.com/loganphillips792/fileupload/config"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.SugaredLogger
	DbConn *sqlx.DB
	Cfg    *config.AppConf
}

func BuildHandler(log *zap.SugaredLogger, db *sqlx.DB, cfg *config.AppConf) *Handler {
	return &Handler{
		Logger: log,
		DbConn: db,
		Cfg:    cfg,
	}
}

func (handler *Handler) UploadFileHandler(c echo.Context) error {
	handler.Logger.Infof("Content Length %d ", c.Request().ContentLength)

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// check if file type is supported
	fileToCheck, _ := file.Open()
	_, err = handler.CheckIfFileTypeIsSupported(fileToCheck)
	if err != nil {
		handler.Logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	defer fileToCheck.Close()

	// open the file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	// // Create the uploads fodler if it doesn't already exist
	// err = os.MkdirAll("uploads", os.ModePerm)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// }

	handler.Logger.Infof("File name from user %s ", c.FormValue("file_name"))

	// Create a new file in the uploads directory
	filePath := ""
	fileNameToInsertIntoDatabase := ""
	if c.FormValue("file_name") != "" {
		filePath = fmt.Sprintf("uploads/%s%s", c.FormValue("file_name"), filepath.Ext(file.Filename))
		fileNameToInsertIntoDatabase = c.FormValue("file_name")
	} else {
		currentUnixTime := time.Now().UnixNano()
		fileNameToInsertIntoDatabase = strconv.FormatInt(currentUnixTime, 10)
		filePath = fmt.Sprintf("uploads/%d%s", currentUnixTime, filepath.Ext(file.Filename))
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// file save was successsful
	query := "INSERT INTO images (name, file_path) VALUES ($1, $2)"

	handler.Logger.Infow("Running SQL statement",
		"SQL", query,
		"name:", fileNameToInsertIntoDatabase,
		"path:", filePath,
	)

	_, err = handler.DbConn.Exec(query, fileNameToInsertIntoDatabase, filePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	// have this run in the background
	go handler.ChangeImageToBlackAndWhite(filePath)

	return c.Blob(http.StatusOK, "application/json", []byte(`{"response":"Upload Successful!!"}`))

}

func (handler *Handler) GetAllFiles(c echo.Context) error {
	handler.Logger.Info("Retreiving all images....")

	searchParams := c.QueryParam("q")
	handler.Logger.Infow("Search parameters passed in",
		"PARAMS", searchParams,
	)

	query := "SELECT * FROM images"

	if searchParams != "" {
		query += " WHERE name LIKE "
	}

	handler.Logger.Infow("Running SQL statement",
		"SQL", query,
	)

	rows, err := handler.DbConn.Query(query)
	if err != nil {
		handler.Logger.Info(err)
		log.Fatal(err)
	}

	defer rows.Close()

	fmt.Println("Getting images")
	var images []Image
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.Id, &image.Name, &image.FilePath, &image.BlackAndWhiteFilePath)

		if err != nil {
			handler.Logger.Error(err)
			// log.Fatal(err)
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
	handler.Logger.Info("Deleting image....")

	id := c.Param("id")

	query := "DELETE FROM images WHERE id = $1"

	handler.Logger.Infow("Running SQL statement",
		"id of image delete", id,
		"SQL", query,
	)

	resp, err := handler.DbConn.Exec(query, id)

	if err != nil {
		log.Fatal(err.Error())
	}

	rowsDeleted, _ := resp.RowsAffected()

	if rowsDeleted == 0 {
		err = c.Blob(http.StatusNotFound, "application/json", []byte(`{"response":"image not found"}`))

		if err != nil {
			log.Fatal(err.Error())
		}

	}

	return c.NoContent(http.StatusNoContent)
}

func (handler *Handler) GetImageById(c echo.Context) error {
	id := c.Param("id")
	handler.Logger.Infof("Getting image of id %s", id)

	query := "SELECT file_path FROM images where id = $1"

	handler.Logger.Infow("Running SQL statement",
		"id of image to find", id,
		"SQL", query,
	)

	row := handler.DbConn.QueryRow(query, id)

	var imagePath string
	err := row.Scan(&imagePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	handler.Logger.Infof("Image path is %s", imagePath)

	return c.File(imagePath)
}

// https://github.com/labstack/echo/blob/v3.3.10/context.go#L542
func (handler *Handler) DownloadImage(c echo.Context) error {
	return c.Attachment("data/IMG_7015.jpg", "download.jpg")
}

// send csv to client to automatically download
func (handler *Handler) DownloadCSV(c echo.Context) error {
	return c.Attachment("data/airtravel.csv", "download.csv")
}

/*

// https://www.reddit.com/r/reactjs/comments/5xgdzh/how_to_correctly_store_user_information/
// https://www.reddit.com/r/reactjs/comments/gek8as/recommended_approach_to_check_if_user_is/
// https://stackoverflow.com/questions/49819183/react-what-is-the-best-way-to-handle-login-and-authentication
// https://stackoverflow.com/questions/70686434/how-to-save-the-users-authorization-in-session-in-react
// https://www.reddit.com/r/programming/comments/nag1cu/jwt_should_not_be_your_default_for_sessions/
// https://www.sohamkamani.com/golang/session-cookie-authentication/

*/

// to create session ID: https://github.com/astaxie/session/blob/master/session.go
func (handler *Handler) Login(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	handler.Logger.Infow("/login",
		"Username", user.Username,
		"Password", user.Password,
	)

	query := "SELECT * FROM users where username = $1"

	// Check if username and password exist
	var userFromDatabase User
	errFromScan := handler.DbConn.QueryRow(query, user.Username).Scan(&userFromDatabase.Id, &userFromDatabase.Username, &userFromDatabase.Email, &userFromDatabase.Password)

	if errFromScan != nil {
		log.Print(errFromScan)
	}

	match := utils.CompareHashAndPassword(userFromDatabase.Password, user.Password)

	if match {
		var hashKey = []byte("very-secret")       // encode value
		var blockKey = []byte("a-lot-secret1111") // encrypt value
		var s = securecookie.New(hashKey, blockKey)

		// create session ID
		sessionId := uuid.New().String()
		expiresAt := time.Now().Add(120 * time.Second).Unix()

		query := "INSERT INTO sessions (session_id, expires_at) VALUES ($1, $2)"

		handler.Logger.Infow("Running SQL statement for session",
			"SQL", query,
		)

		_, err = handler.DbConn.Exec(query, sessionId, expiresAt)

		if err != nil {
			handler.Logger.Errorw("Error inserting into sessions table", "error", err)
		}

		value := map[string]string{
			"sessionId": sessionId,
		}

		if encoded, errFromCookie := s.Encode("user_session", value); errFromCookie == nil {
			cookie := new(http.Cookie)
			cookie.Name = "user_session"
			cookie.Value = encoded
			cookie.Expires = time.Now().Add(24 * time.Hour)
			cookie.HttpOnly = true
			// cookie.Secure = true
			c.SetCookie(cookie)
		}
		return c.String(http.StatusOK, "login successful")
	} else {
		return c.String(http.StatusUnauthorized, "Login failed")
	}
}

func (handler *Handler) Register(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	handler.Logger.Infow("/register",
		"Username", user.Username,
		"Password", user.Password,
	)

	hashedPassword, _ := utils.HashPassword(user.Password)

	query := "INSERT INTO users (username, password) VALUES ($1, $2)"

	handler.Logger.Infow("Running SQL statement",
		"SQL", query,
	)

	_, err = handler.DbConn.Exec(query, user.Username, hashedPassword)

	if err != nil {
		log.Fatal(err.Error())
	}
	return c.String(http.StatusCreated, "User successfully created")
}

func HelloWorld(c echo.Context) error {
	data := []byte(`{"status":"OK!"}`)
	return c.Blob(http.StatusOK, "application/json", data)
}
