package auth

import (
	"net/http"

	dto "student_service_app/backend/internal/dto/auth"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/response"
	servicepkg "student_service_app/backend/internal/service/auth"
)

type Handler struct {
	base    handler.Base
	service servicepkg.Service
}

func NewHandler(base handler.Base, service servicepkg.Service) *Handler {
	return &Handler{base: base, service: service}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest
	if err := h.base.Decode(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	session, err := h.service.SignUp(r.Context(), servicepkg.SignUpInput{
		SchoolName:    req.SchoolName,
		TenantSlug:    req.TenantSlug,
		AdminName:     req.AdminName,
		Email:         req.Email,
		Password:      req.Password,
		SchoolPhone:   req.SchoolPhone,
		SchoolAddress: req.SchoolAddress,
	})
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.FromSession(*session))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := h.base.Decode(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	session, err := h.service.Login(r.Context(), servicepkg.LoginInput{
		TenantSlug: req.TenantSlug,
		Email:      req.Email,
		Password:   req.Password,
	})
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.FromSession(*session))
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	session, err := h.service.Me(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.FromSession(*session))
}
