package mapper

import (
	domain "student_service_app/backend/internal/domain/enrollment"
	dto "student_service_app/backend/internal/dto/enrollment"
)

func EnrollmentToDTO(e domain.Enrollment, items []domain.EnrollmentOptionalItem) dto.Response {
	itemDTOs := make([]dto.OptionalItemResponse, 0, len(items))
	for _, item := range items {
		itemDTOs = append(itemDTOs, dto.OptionalItemResponse{
			ID:                item.ID,
			OptionalFeeItemID: item.OptionalFeeItemID,
			ItemNameSnapshot:  item.ItemNameSnapshot,
			AmountSnapshot:    item.AmountSnapshot,
			Quantity:          item.Quantity,
			TotalAmount:       item.TotalAmount,
		})
	}
	return dto.Response{
		ID:              e.ID,
		EnrollmentCode:  e.EnrollmentCode,
		StudentID:       e.StudentID,
		StudentName:     e.StudentName,
		GuardianName:    e.GuardianName,
		ClassCourseID:   e.ClassCourseID,
		ClassName:       e.ClassName,
		CourseName:      e.CourseName,
		EnrollmentDate:  e.EnrollmentDate.Format("2006-01-02"),
		SubTotal:        e.SubTotal,
		DiscountAmount:  e.DiscountAmount,
		FinalFee:        e.FinalFee,
		PaidAmount:      e.PaidAmount,
		RemainingAmount: e.RemainingAmount,
		PaymentStatus:   e.PaymentStatus,
		Note:            e.Note,
		OptionalItems:   itemDTOs,
		CreatedAt:       DateTimeString(e.CreatedAt),
		UpdatedAt:       DateTimeString(e.UpdatedAt),
	}
}
