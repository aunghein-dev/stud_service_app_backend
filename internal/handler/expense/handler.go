package expense

import (
	"net/http"

	domain "student_service_app/backend/internal/domain/expense"
	dto "student_service_app/backend/internal/dto/expense"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/expense"
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
	expenseDate, err := mapper.ParseDateOrNow(req.ExpenseDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	exp := &domain.Expense{
		ExpenseDate:   expenseDate,
		ExpenseType:   req.ExpenseType,
		TeacherID:     req.TeacherID,
		ClassCourseID: req.ClassCourseID,
		Amount:        req.Amount,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
		ReferenceNo:   req.ReferenceNo,
	}
	if err := h.service.Create(r.Context(), exp); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, mapper.ExpenseToDTO(*exp))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.ExpenseToDTO(item))
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
	response.JSON(w, http.StatusOK, mapper.ExpenseToDTO(*item))
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
	expenseDate, err := mapper.ParseDateOrNow(req.ExpenseDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	exp := &domain.Expense{
		ID:            id,
		ExpenseDate:   expenseDate,
		ExpenseType:   req.ExpenseType,
		TeacherID:     req.TeacherID,
		ClassCourseID: req.ClassCourseID,
		Amount:        req.Amount,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
		ReferenceNo:   req.ReferenceNo,
	}
	if err := h.service.Update(r.Context(), exp); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.ExpenseToDTO(*exp))
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
