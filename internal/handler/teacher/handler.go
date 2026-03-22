package teacher

import (
	"net/http"

	domain "student_service_app/backend/internal/domain/teacher"
	dto "student_service_app/backend/internal/dto/teacher"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/teacher"
)

type Handler struct {
	base    handler.Base
	service servicepkg.Service
}

func NewHandler(base handler.Base, service servicepkg.Service) *Handler {
	return &Handler{base: base, service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest
	if err := h.base.Decode(r, &req); err != nil {
		response.Error(w, err)
		return
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	t := &domain.Teacher{
		TeacherCode:      req.TeacherCode,
		TeacherName:      req.TeacherName,
		Phone:            req.Phone,
		Address:          req.Address,
		SubjectSpecialty: req.SubjectSpecialty,
		SalaryType:       req.SalaryType,
		DefaultFeeAmount: req.DefaultFeeAmount,
		Note:             req.Note,
		IsActive:         active,
	}
	if err := h.service.Create(r.Context(), t); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, mapper.TeacherToDTO(*t))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.TeacherToDTO(item))
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.TeacherToDTO(*item))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	var req dto.UpdateRequest
	if err := h.base.Decode(r, &req); err != nil {
		response.Error(w, err)
		return
	}
	t := &domain.Teacher{
		ID:               id,
		TeacherName:      req.TeacherName,
		Phone:            req.Phone,
		Address:          req.Address,
		SubjectSpecialty: req.SubjectSpecialty,
		SalaryType:       req.SalaryType,
		DefaultFeeAmount: req.DefaultFeeAmount,
		Note:             req.Note,
		IsActive:         req.IsActive,
	}
	if err := h.service.Update(r.Context(), t); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.TeacherToDTO(*t))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"deleted": true})
}
