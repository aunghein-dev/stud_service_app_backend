package mapper

import (
	"encoding/json"

	domain "student_service_app/backend/internal/domain/receipt"
	dto "student_service_app/backend/internal/dto/receipt"
)

func ReceiptToDTO(r domain.Receipt) dto.Response {
	payload := map[string]any{}
	if len(r.PayloadJSON) > 0 {
		_ = json.Unmarshal(r.PayloadJSON, &payload)
	}
	return dto.Response{
		ID:              r.ID,
		ReceiptNo:       r.ReceiptNo,
		ReceiptType:     r.ReceiptType,
		StudentID:       r.StudentID,
		EnrollmentID:    r.EnrollmentID,
		PaymentID:       r.PaymentID,
		ClassCourseID:   r.ClassCourseID,
		TotalAmount:     r.TotalAmount,
		PaidAmount:      r.PaidAmount,
		RemainingAmount: r.RemainingAmount,
		Payload:         payload,
		IssuedAt:        DateTimeString(r.IssuedAt),
		CreatedAt:       DateTimeString(r.CreatedAt),
	}
}
