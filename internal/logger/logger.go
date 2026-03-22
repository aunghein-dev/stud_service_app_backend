package logger

import (
	"student_service_app/backend/internal/config"

	"go.uber.org/zap"
)

func New(cfg *config.Config) (*zap.Logger, error) {
	if cfg.App.Env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
