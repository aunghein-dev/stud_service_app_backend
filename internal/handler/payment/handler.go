package payment

import (
	"net/http"
	"time"

	dto "student_service_app/backend/internal/dto/payment"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/payment"
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
	paymentDate, err := mapper.ParseDateOrNow(req.PaymentDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	paymentTx, receiptTx, err := h.service.Create(r.Context(), servicepkg.CreateInput{
		EnrollmentID:  req.EnrollmentID,
		PaymentDate:   paymentDate,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		Note:          req.Note,
		ReceivedBy:    req.ReceivedBy,
	})
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{
		"payment": mapper.PaymentToDTO(*paymentTx),
		"receipt": mapper.ReceiptToDTO(*receiptTx),
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.PaymentToDTO(item))
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
	response.JSON(w, http.StatusOK, mapper.PaymentToDTO(*item))
}

func (h *Handler) ListByEnrollment(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	items, err := h.service.ListByEnrollment(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.PaymentToDTO(item))
	}
	response.JSON(w, http.StatusOK, result)
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

	paymentDate := time.Time{}
	if req.PaymentDate != "" {
		paymentDate, err = mapper.ParseDateOrNow(req.PaymentDate)
		if err != nil {
			response.Error(w, err)
			return
		}
	}

	updated, err := h.service.Update(r.Context(), id, servicepkg.UpdateInput{
		PaymentDate:   paymentDate,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		Note:          req.Note,
		ReceivedBy:    req.ReceivedBy,
	})
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.PaymentToDTO(*updated))
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
