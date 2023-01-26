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

	// logger, _ := zap.NewProduction()
	// defer logger.Sync() // flushes buffer, if any
	// sugar := logger.Sugar()

	// envHandler := &Handler{logger: sugar, dbConn: db}
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

	e.GET("/hello", HelloWorld)

	e.Logger.Fatal(e.Start(":8000"))

	/*

		r := mux.NewRouter()

		cors := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodPost,
				http.MethodGet,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
		})

		r.HandleFunc("/uploadfile/", envHandler.UploadFileHandler).Methods("POST")
		r.HandleFunc("/images/", envHandler.GetAllFiles).Methods("GET")
		r.HandleFunc("/images/{id}", envHandler.DeleteImage).Methods("DELETE")
		r.HandleFunc("/download_csv/", envHandler.DownloadCSV).Methods("GET")
		r.HandleFunc("/download_image/", envHandler.DownloadImage).Methods("GET")
		r.HandleFunc("/file/", FileHandler)
		r.HandleFunc("/hello", HelloWorld)

		handler := cors.Handler(r)

		log.Fatal(http.ListenAndServe(":8000", handler))
	*/
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

/*
const MAX_UPLOAD_SIZE = 1024 * 1024 // 1 MB

func (handler *Handler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	handler.logger.Infof("Content Length %d ", r.ContentLength)

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

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

	// Create the uploads fodler if it doesn't already exist
	err = os.MkdirAll("./uploads", os.ModePerm)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.logger.Infof("File name from user %s ", r.Form.Get("file_name"))

	// Create a new file in the uploads directory
	filePath := ""
	if r.Form.Get("file_name") != "" {
		filePath = fmt.Sprintf("./uploads/%s%s", r.Form.Get("file_name"), filepath.Ext(fileHeader.Filename))
	} else {
		filePath = fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	}
	// dst, err := os.Create(filePath)
	// file_path := fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create(filePath)
	fmt.Println(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// file save was successsful
	query := "INSERT INTO images (name, file_path) VALUES (?, ?)"

	handler.logger.Infow("Running SQL statement",
		"SQL", query,
	)

	_, err = handler.dbConn.Exec(query, r.Form.Get("file_name"), filePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")

	w.Write([]byte(`{"response":"Upload successful"}`))

	changeImageToBlackAndWhite()

}

// When user successfully uploads image, they can click "convert to black and white". The new image will
// show as a thumbnail, and then they can click download to download the new image
// https://stackoverflow.com/questions/42516203/converting-rgba-image-to-grayscale-golang
func changeImageToBlackAndWhite() {
	fmt.Println("Converting image to black and white...")
}

func (handler *Handler) GetAllFiles(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Retreiving all images....")

	searchParams := r.URL.Query().Get("q")
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

	err = json.NewEncoder(w).Encode(images)
	w.Header().Set("Content-Type", "application/json")

	if len(images) == 0 {
		err = json.NewEncoder(w).Encode(make([]Image, 0))
	}

	if err != nil {
		log.Fatal(err)
	}

}

func (handler *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Deleting image....")

	vars := mux.Vars(r)

	query := "DELETE FROM websites WHERE id = ?"

	handler.logger.Infow("Running SQL statement",
		"id of image delete", vars["id"],
		"SQL", query,
	)

	_, err := handler.dbConn.Exec(vars["id"])

	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")

}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var fileInfo FileInfo

		err := json.NewDecoder(r.Body).Decode(&fileInfo)

		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Printf("FILE NAME IS %s\n", fileInfo.Name)

		// fmt.Println("RECEIVING A POST REQUEST")
		// fmt.Printf("BODY %v\n", r.Body)
		// bodyBytes, _ := ioutil.ReadAll(r.Body)
		// bodyString := string(bodyBytes)
		// fmt.Print(bodyString)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// send csv to client to automatically download
// https://medium.com/wesionary-team/create-csv-file-in-go-server-and-download-from-reactjs-4f22f148290b
// https://stackoverflow.com/questions/68162651/go-how-to-response-csv-file
// https://medium.com/wesionary-team/create-csv-file-in-go-server-and-download-from-reactjs-4f22f148290b
func (handler *Handler) DownloadCSV(w http.ResponseWriter, r *http.Request) {
	// open file
	f, err := os.Open("./data/airtravel.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	w.Header().Add("Content-Disposition", `attachment; filename="test.csv"`)
	http.ServeFile(w, r, "./data/airtravel.csv")

	//io.Copy(w, f)
}

func (handler *Handler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	// open file
	f, err := os.Open("./data/IMG_7015.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	w.Header().Add("Content-Disposition", `attachment; filename="image.jpeg"`)
	http.ServeFile(w, r, "./data/IMG_7015.jpg")

	//io.Copy(w, f)
}

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
