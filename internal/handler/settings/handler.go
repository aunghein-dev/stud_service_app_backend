package settings

import (
	"net/http"

	dto "student_service_app/backend/internal/dto/settings"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicepkg "student_service_app/backend/internal/service/settings"
)

type Handler struct {
	base    handler.Base
	service servicepkg.Service
}

func NewHandler(base handler.Base, service servicepkg.Service) *Handler {
	return &Handler{base: base, service: service}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	set, err := h.service.Get(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.SettingsToDTO(*set))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateRequest
	if err := h.base.Decode(r, &req); err != nil {
		response.Error(w, err)
		return
	}
	updated, err := h.service.Update(r.Context(), servicepkg.UpdateInput{
		SchoolName:           req.SchoolName,
		SchoolAddress:        req.SchoolAddress,
		SchoolPhone:          req.SchoolPhone,
		DefaultCurrency:      req.DefaultCurrency,
		ReceiptPrefix:        req.ReceiptPrefix,
		PaymentMethods:       req.PaymentMethods,
		OptionalItemDefaults: req.OptionalItemDefaults,
		PrintPreferences:     req.PrintPreferences,
	})
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.SettingsToDTO(*updated))
}
