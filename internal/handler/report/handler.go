package report

import (
	"net/http"

	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/report"
)

type Handler struct {
	service servicepkg.Service
}

func NewHandler(service servicepkg.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.Dashboard(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Students(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.StudentReport(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Teachers(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.TeacherReport(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) ClassCourses(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.ClassCourseReport(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Gross(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.Gross(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.TransactionReport(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Performance(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.Performance(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}
