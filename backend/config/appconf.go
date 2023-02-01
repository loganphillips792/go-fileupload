package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConf struct {
	GorillaSessionsHashKey  string
	GorillaSessionsBlockKey string
}

func Init() (*AppConf, error) {
	// get config
	godotenv.Load(".env")

	config := &AppConf{
		GorillaSessionsHashKey:  os.Getenv("GORILLA_SESSIONS_HASH_KEY"),
		GorillaSessionsBlockKey: os.Getenv("GORILLA_SESSIONS_BLOCK_KEY"),
	}

	return config, nil
}
