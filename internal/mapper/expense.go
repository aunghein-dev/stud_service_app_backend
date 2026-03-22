package mapper

import (
	domain "student_service_app/backend/internal/domain/expense"
	dto "student_service_app/backend/internal/dto/expense"
)

func ExpenseToDTO(e domain.Expense) dto.Response {
	return dto.Response{
		ID:            e.ID,
		ExpenseDate:   e.ExpenseDate.Format("2006-01-02"),
		ExpenseType:   e.ExpenseType,
		TeacherID:     e.TeacherID,
		ClassCourseID: e.ClassCourseID,
		Amount:        e.Amount,
		Description:   e.Description,
		PaymentMethod: e.PaymentMethod,
		ReferenceNo:   e.ReferenceNo,
		CreatedAt:     DateTimeString(e.CreatedAt),
		UpdatedAt:     DateTimeString(e.UpdatedAt),
	}
}
