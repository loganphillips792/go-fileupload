package api

import (
	"errors"
	"log"

	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/loganphillips792/fileupload/config"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func ApiMiddleware(db *sqlx.DB, handler *Handler, sugar *zap.SugaredLogger, cfg *config.AppConf) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:user_session",
		Validator: func(key string, c echo.Context) (bool, error) {
			sugar.Info("Validating in Middlware...")
			var hashKey = []byte(cfg.GorillaSessionsHashKey)   // encode value
			var blockKey = []byte(cfg.GorillaSessionsBlockKey) // encrypt value
			var s = securecookie.New(hashKey, blockKey)
			value := make(map[string]string)

			err := s.Decode("user_session", key, &value)

			if err != nil {
				sugar.Errorw("Error when decoding cookie value", err)
				return false, errors.New("authentication failed. Please login again")
			}
			sugar.Infow("Decryption", "The decrypted value is", value["sessionId"])

			// we now have the decrypted session id. We will now look it up in the sessions table
			query := "SELECT * FROM sessions where session_id = $1"

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
	})
}
