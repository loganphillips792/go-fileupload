package api

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
	utils "github.com/loganphillips792/fileupload"
	"github.com/loganphillips792/fileupload/config"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.SugaredLogger
	DbConn *sql.DB
	Cfg    *config.AppConf
}

func BuildHandler(log *zap.SugaredLogger, db *sql.DB, cfg *config.AppConf) *Handler {
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

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	// Create the uploads fodler if it doesn't already exist
	err = os.MkdirAll("uploads", os.ModePerm)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

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

	handler.Logger.Infow("Running SQL statement",
		"SQL", query,
	)

	_, err = handler.DbConn.Exec(query, fileNameToInsertIntoDatabase, filePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	// have this run in the background
	go handler.changeImageToBlackAndWhite(filePath)

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
func (handler *Handler) changeImageToBlackAndWhite(filePath string) {
	fmt.Println("Converting image to black and white...")
	fmt.Println("File path is", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		handler.Logger.Error("Error ")
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		handler.Logger.Error("Error ")
	}

	// Create a new image with the same dimensions as the original
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// https://stackoverflow.com/a/42518487
	// Iterate through each pixel in the original image
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			oldPixel := img.At(x, y)
			pixel := color.GrayModel.Convert(oldPixel)
			newImg.Set(x, y, pixel)
		}
	}

	// Create a new file for the black and white image
	newFile, err := os.Create("uploads/bw.jpeg")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Encode the new image as a JPEG
	imageEncodeError := jpeg.Encode(newFile, newImg, nil)

	if imageEncodeError != nil {
		handler.Logger.Error("Error when encoding image: ", imageEncodeError)
	}

	time.Sleep(10 * time.Second)

	handler.Logger.Info("changeImageToBlackAndWhite go routine finished")
}

func (handler *Handler) GetAllFiles(c echo.Context) error {
	handler.Logger.Info("Retreiving all images....")

	searchParams := c.QueryParam("q")
	handler.Logger.Infow("Search parameters passed in",
		"PARAMS", searchParams,
	)

	query := "SELECT * FROM images"

	if searchParams != "" {
		query += " LIKE name "
	}

	handler.Logger.Infow("Running SQL statement",
		"SQL", query,
	)

	rows, err := handler.DbConn.Query(query)

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
	handler.Logger.Info("Deleting image....")

	id := c.Param("id")

	query := "DELETE FROM images WHERE id = ?"

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

func (handler *Handler) GetImageByPath(c echo.Context) error {
	return c.File(("uploads/YO.jpeg"))
}

// https://github.com/labstack/echo/blob/v3.3.10/context.go#L542
func (handler *Handler) DownloadImage(c echo.Context) error {
	return c.Attachment("data/IMG_7015.jpg", "download.jpg")
}

// send csv to client to automatically download
// https://medium.com/wesionary-team/create-csv-file-in-go-server-and-download-from-reactjs-4f22f148290b
// https://stackoverflow.com/questions/68162651/go-how-to-response-csv-file
// https://medium.com/wesionary-team/create-csv-file-in-go-server-and-download-from-reactjs-4f22f148290b
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

	query := "SELECT * FROM users where username = ?"

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

		query := "INSERT INTO sessions (session_id, expires_at) VALUES (?, ?)"

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

	query := "INSERT INTO users (username, password) VALUES (?, ?)"

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
