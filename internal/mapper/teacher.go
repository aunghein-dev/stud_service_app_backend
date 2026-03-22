package mapper

import (
	domain "student_service_app/backend/internal/domain/teacher"
	dtoteacher "student_service_app/backend/internal/dto/teacher"
)

func TeacherToDTO(t domain.Teacher) dtoteacher.Response {
	return dtoteacher.Response{
		ID:               t.ID,
		TeacherCode:      t.TeacherCode,
		TeacherName:      t.TeacherName,
		Phone:            t.Phone,
		Address:          t.Address,
		SubjectSpecialty: t.SubjectSpecialty,
		SalaryType:       t.SalaryType,
		DefaultFeeAmount: t.DefaultFeeAmount,
		Note:             t.Note,
		IsActive:         t.IsActive,
		CreatedAt:        DateTimeString(t.CreatedAt),
		UpdatedAt:        DateTimeString(t.UpdatedAt),
	}
}
