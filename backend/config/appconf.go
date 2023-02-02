package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type AppConf struct {
	GorillaSessionsHashKey  string
	GorillaSessionsBlockKey string
}

func Init() (*AppConf, error) {
	// get config
	err := godotenv.Load(".env")

	if err != nil {
		log.Error("failed when reading .env file")
	}

	config := &AppConf{
		GorillaSessionsHashKey:  os.Getenv("GORILLA_SESSIONS_HASH_KEY"),
		GorillaSessionsBlockKey: os.Getenv("GORILLA_SESSIONS_BLOCK_KEY"),
	}

	return config, nil
}
