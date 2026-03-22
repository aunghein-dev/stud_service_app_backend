package student

import (
	"net/http"

	domain "student_service_app/backend/internal/domain/student"
	dto "student_service_app/backend/internal/dto/student"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/student"
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
	dob, err := mapper.ParseDate(req.DateOfBirth)
	if err != nil {
		response.Error(w, err)
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	st := &domain.Student{
		StudentCode:   req.StudentCode,
		FullName:      req.FullName,
		Gender:        req.Gender,
		DateOfBirth:   dob,
		Phone:         req.Phone,
		GuardianName:  req.GuardianName,
		GuardianPhone: req.GuardianPhone,
		Address:       req.Address,
		SchoolName:    req.SchoolName,
		GradeLevel:    req.GradeLevel,
		Note:          req.Note,
		IsActive:      isActive,
	}
	if err := h.service.Create(r.Context(), st); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, mapper.StudentToDTO(*st))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.StudentToDTO(item))
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
	response.JSON(w, http.StatusOK, mapper.StudentToDTO(*item))
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
	dob, err := mapper.ParseDate(req.DateOfBirth)
	if err != nil {
		response.Error(w, err)
		return
	}
	st := &domain.Student{
		ID:            id,
		FullName:      req.FullName,
		Gender:        req.Gender,
		DateOfBirth:   dob,
		Phone:         req.Phone,
		GuardianName:  req.GuardianName,
		GuardianPhone: req.GuardianPhone,
		Address:       req.Address,
		SchoolName:    req.SchoolName,
		GradeLevel:    req.GradeLevel,
		Note:          req.Note,
		IsActive:      req.IsActive,
	}
	if err := h.service.Update(r.Context(), st); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.StudentToDTO(*st))
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
