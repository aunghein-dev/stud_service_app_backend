package mapper

import (
	domain "student_service_app/backend/internal/domain/payment"
	dto "student_service_app/backend/internal/dto/payment"
)

func PaymentToDTO(p domain.Payment) dto.Response {
	return dto.Response{
		ID:            p.ID,
		ReceiptNo:     p.ReceiptNo,
		StudentID:     p.StudentID,
		EnrollmentID:  p.EnrollmentID,
		ClassCourseID: p.ClassCourseID,
		PaymentDate:   p.PaymentDate.Format("2006-01-02"),
		PaymentMethod: p.PaymentMethod,
		Amount:        p.Amount,
		Note:          p.Note,
		ReceivedBy:    p.ReceivedBy,
		CreatedAt:     DateTimeString(p.CreatedAt),
	}
}
