package api

import (
	"database/sql"

	"go.uber.org/zap"
)

type Image struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

type Handler struct {
	Logger *zap.SugaredLogger
	DbConn *sql.DB
}
