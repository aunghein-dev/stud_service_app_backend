package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"student_service_app/backend/internal/di"

	"go.uber.org/zap"
)

func main() {
	app, err := di.InitializeApp()
	if err != nil {
		panic(err)
	}

	go func() {
		app.Logger.Info("starting api server")
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Fatal("server crashed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = app.Server.Shutdown(ctx)
	app.Logger.Info("server stopped")
}
