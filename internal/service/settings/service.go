package settings

import (
	"context"
	"encoding/json"
	"strings"

	"student_service_app/backend/internal/domain/settings"
	"student_service_app/backend/internal/errs"
	settingsrepo "student_service_app/backend/internal/repository/settings"
)

type Service interface {
	Get(ctx context.Context) (*settings.Setting, error)
	Update(ctx context.Context, req UpdateInput) (*settings.Setting, error)
}

type UpdateInput struct {
	SchoolName           string
	SchoolAddress        string
	SchoolPhone          string
	DefaultCurrency      string
	ReceiptPrefix        string
	PaymentMethods       []string
	OptionalItemDefaults []string
	PrintPreferences     map[string]any
}

type service struct {
	repo settingsrepo.Repository
}

func NewService(repo settingsrepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context) (*settings.Setting, error) {
	set, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	if set == nil {
		return nil, errs.NotFound("settings not found")
	}
	return set, nil
}

func (s *service) Update(ctx context.Context, req UpdateInput) (*settings.Setting, error) {
	existing, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errs.NotFound("settings not found")
	}

	schoolName := strings.TrimSpace(req.SchoolName)
	if schoolName == "" {
		return nil, errs.BadRequest("school_name is required")
	}

	defaultCurrency := strings.ToUpper(strings.TrimSpace(req.DefaultCurrency))
	if defaultCurrency == "" {
		return nil, errs.BadRequest("default_currency is required")
	}

	receiptPrefix := strings.ToUpper(strings.TrimSpace(req.ReceiptPrefix))
	if receiptPrefix == "" {
		return nil, errs.BadRequest("receipt_prefix is required")
	}

	methods, err := json.Marshal(req.PaymentMethods)
	if err != nil {
		return nil, errs.BadRequest("invalid payment_methods")
	}
	defaults, err := json.Marshal(req.OptionalItemDefaults)
	if err != nil {
		return nil, errs.BadRequest("invalid optional_item_defaults")
	}
	prefs, err := json.Marshal(req.PrintPreferences)
	if err != nil {
		return nil, errs.BadRequest("invalid print_preferences")
	}
	existing.SchoolName = schoolName
	existing.SchoolAddress = strings.TrimSpace(req.SchoolAddress)
	existing.SchoolPhone = strings.TrimSpace(req.SchoolPhone)
	existing.DefaultCurrency = defaultCurrency
	existing.ReceiptPrefix = receiptPrefix
	existing.PaymentMethodsJSON = methods
	existing.OptionalDefaultsJSON = defaults
	existing.PrintPreferencesJSON = prefs

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
