package enrollment

import (
	"net/http"

	domain "student_service_app/backend/internal/domain/enrollment"
	dto "student_service_app/backend/internal/dto/enrollment"
	"student_service_app/backend/internal/handler"
	"student_service_app/backend/internal/mapper"
	"student_service_app/backend/internal/response"
	servicecommon "student_service_app/backend/internal/service"
	servicepkg "student_service_app/backend/internal/service/enrollment"
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
	enrollmentDate, err := mapper.ParseDateOrNow(req.EnrollmentDate)
	if err != nil {
		response.Error(w, err)
		return
	}
	inputs := make([]servicepkg.OptionalItemInput, 0, len(req.OptionalItems))
	for _, item := range req.OptionalItems {
		inputs = append(inputs, servicepkg.OptionalItemInput{
			OptionalFeeItemID: item.OptionalFeeItemID,
			ItemName:          item.ItemName,
			Amount:            item.Amount,
			Quantity:          item.Quantity,
		})
	}
	created, createdItems, paymentTx, receiptTx, err := h.service.Create(r.Context(), servicepkg.CreateInput{
		StudentID:      req.StudentID,
		ClassCourseID:  req.ClassCourseID,
		EnrollmentDate: enrollmentDate,
		DiscountAmount: req.DiscountAmount,
		OptionalItems:  inputs,
		InitialPayment: req.InitialPayment,
		PaymentMethod:  req.PaymentMethod,
		ReceivedBy:     req.ReceivedBy,
		Note:           req.Note,
		AllowDuplicate: req.AllowDuplicate,
	})
	if err != nil {
		response.Error(w, err)
		return
	}
	resp := mapper.EnrollmentToDTO(*created, createdItems)
	out := map[string]any{"enrollment": resp}
	if paymentTx != nil {
		payment := mapper.PaymentToDTO(*paymentTx)
		out["initial_payment"] = payment
	}
	if receiptTx != nil {
		receipt := mapper.ReceiptToDTO(*receiptTx)
		out["receipt"] = receipt
	}
	response.JSON(w, http.StatusCreated, out)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), servicecommon.BuildFilter(r))
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.EnrollmentToDTO(item, nil))
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) ListByStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	items, err := h.service.ListByStudent(r.Context(), studentID)
	if err != nil {
		response.Error(w, err)
		return
	}
	result := make([]dto.Response, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.EnrollmentToDTO(item, nil))
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := handler.ParseID(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	item, optionalItems, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.EnrollmentToDTO(*item, optionalItems))
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
	enroll := &domain.Enrollment{
		ID:             id,
		DiscountAmount: req.DiscountAmount,
		Note:           req.Note,
	}
	if err := h.service.Update(r.Context(), enroll); err != nil {
		response.Error(w, err)
		return
	}
	updated, items, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, mapper.EnrollmentToDTO(*updated, items))
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
