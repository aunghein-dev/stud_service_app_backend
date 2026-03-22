package di

import (
	"net/http"

	"student_service_app/backend/internal/config"

	"go.uber.org/zap"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	Server *http.Server
}
