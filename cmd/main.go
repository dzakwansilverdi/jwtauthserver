package main

import (
	"jwtauthserver/auth"
	"jwtauthserver/auth/api"
	"jwtauthserver/auth/database"
	"jwtauthserver/pkg/uuid"
	"net/http"
	"os"

	"github.com/go-kit/log"

	"github.com/jmoiron/sqlx"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	db, err := database.Connect()
	if err != nil {
		logger.Log("Error Connecting DB", err)
		os.Exit(1)
	}
	defer db.Close()

	svc := NewService(db)
	var h http.Handler
	{
		h = api.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		logger.Log("transport", "HTTP", "addr", ":8080")
		errs <- http.ListenAndServe(":8080", h)
	}()

	logger.Log("exit", <-errs)
}

func NewService(db *sqlx.DB) auth.Service {
	User := database.NewUsersRepository(db)
	IDprov := uuid.New()
	return auth.NewService(User, IDprov)
}
