package api

type Image struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}
