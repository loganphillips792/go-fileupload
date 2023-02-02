package api

import (
	"database/sql"
)

type Image struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

type User struct {
	Id       int
	Username string
	Email    sql.NullString
	Password string
}
