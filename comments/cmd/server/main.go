package main

import (
	"context"
	"errors"
	"example/comments/internal/app"
	"example/comments/internal/logger"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	logger.Infow(ctx, "App starting")

	appCartService, err := app.NewApp(ctx, os.Getenv("CONFIG_FILE"))
	if err != nil {
		logger.Errorw(ctx, "app start failed", "error", err.Error())
		panic(err)
	}

	if err := appCartService.ListenAndServe(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Infow(ctx, "Server closed")
			return
		}

		logger.Errorw(ctx, "server failed", "error", err.Error())
		panic(err)
	}
}
