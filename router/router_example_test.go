package router

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func ExampleLogRoutes() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mux := Setup()

	LogRoutes(mux, logger)

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()

	// Outputs:
	// {"level":"debug","timestamp":"2020-05-21T22:07:12.091+0900","logger":"server/server.go:109","message":"Registering route","method":"GET","route":"/health"}
}
