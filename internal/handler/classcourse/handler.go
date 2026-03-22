package classcourse

import (
	"net/http"
	"strings"

	domain "student_service_app/backend/internal/domain/classcourse"
	dto "student_service_app/backend/internal/dto/classcourse"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/classcourse"
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
	start, err := mapper.ParseDate(req.StartDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	end, err := mapper.ParseDate(req.EndDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	model := &domain.ClassCourse{
		CourseCode:        req.CourseCode,
		CourseName:        req.CourseName,
		ClassName:         req.ClassName,
		Category:          req.Category,
		Subject:           req.Subject,
		Level:             req.Level,
		StartDate:         start,
		EndDate:           end,
		ScheduleText:      req.ScheduleText,
		DaysOfWeek:        strings.Join(req.DaysOfWeek, ","),
		TimeStart:         req.TimeStart,
		TimeEnd:           req.TimeEnd,
		Room:              req.Room,
		AssignedTeacherID: req.AssignedTeacherID,
		MaxStudents:       req.MaxStudents,
		Status:            req.Status,
		BaseCourseFee:     req.BaseCourseFee,
		RegistrationFee:   req.RegistrationFee,
		ExamFee:           req.ExamFee,
		CertificateFee:    req.CertificateFee,
		Note:              req.Note,
	}
	if err := h.service.Create(r.Context(), model); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, mapper.ClassCourseToDTO(*model))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.ClassCourseToDTO(item))
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
	response.JSON(w, http.StatusOK, mapper.ClassCourseToDTO(*item))
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
	start, err := mapper.ParseDate(req.StartDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	end, err := mapper.ParseDate(req.EndDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	model := &domain.ClassCourse{
		ID:                id,
		CourseCode:        req.CourseCode,
		CourseName:        req.CourseName,
		ClassName:         req.ClassName,
		Category:          req.Category,
		Subject:           req.Subject,
		Level:             req.Level,
		StartDate:         start,
		EndDate:           end,
		ScheduleText:      req.ScheduleText,
		DaysOfWeek:        strings.Join(req.DaysOfWeek, ","),
		TimeStart:         req.TimeStart,
		TimeEnd:           req.TimeEnd,
		Room:              req.Room,
		AssignedTeacherID: req.AssignedTeacherID,
		MaxStudents:       req.MaxStudents,
		Status:            req.Status,
		BaseCourseFee:     req.BaseCourseFee,
		RegistrationFee:   req.RegistrationFee,
		ExamFee:           req.ExamFee,
		CertificateFee:    req.CertificateFee,
		Note:              req.Note,
	}
	if err := h.service.Update(r.Context(), model); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.ClassCourseToDTO(*model))
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

func (h *Handler) CreateOptionalFee(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	var req dto.OptionalFeeItemRequest
	if err := h.base.Decode(r, &req); err != nil {
		response.Error(w, err)
		return
	}
	item := &domain.OptionalFeeItem{
		ClassCourseID: id,
		ItemName:      req.ItemName,
		DefaultAmount: req.DefaultAmount,
		IsOptional:    req.IsOptional,
		IsActive:      req.IsActive,
	}
	if err := h.service.CreateOptionalFee(r.Context(), item); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, mapper.OptionalFeeToDTO(*item))
}

func (h *Handler) ListOptionalFees(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	items, err := h.service.ListOptionalFees(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.OptionalFeeItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.OptionalFeeToDTO(item))
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) UpdateOptionalFee(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	var req dto.OptionalFeeItemRequest
	if err := h.base.Decode(r, &req); err != nil {
		response.Error(w, err)
		return
	}
	item := &domain.OptionalFeeItem{
		ID:            id,
		ItemName:      req.ItemName,
		DefaultAmount: req.DefaultAmount,
		IsOptional:    req.IsOptional,
		IsActive:      req.IsActive,
	}
	if err := h.service.UpdateOptionalFee(r.Context(), item); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.OptionalFeeToDTO(*item))
}

func (h *Handler) DeleteOptionalFee(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.DeleteOptionalFee(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"deleted": true})
}
