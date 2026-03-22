package mapper

import (
	domain "student_service_app/backend/internal/domain/student"
	dtostudent "student_service_app/backend/internal/dto/student"
)

func StudentToDTO(s domain.Student) dtostudent.Response {
	return dtostudent.Response{
		ID:            s.ID,
		StudentCode:   s.StudentCode,
		FullName:      s.FullName,
		Gender:        s.Gender,
		DateOfBirth:   DateStringPtr(s.DateOfBirth),
		Phone:         s.Phone,
		GuardianName:  s.GuardianName,
		GuardianPhone: s.GuardianPhone,
		Address:       s.Address,
		SchoolName:    s.SchoolName,
		GradeLevel:    s.GradeLevel,
		Note:          s.Note,
		IsActive:      s.IsActive,
		CreatedAt:     DateTimeString(s.CreatedAt),
		UpdatedAt:     DateTimeString(s.UpdatedAt),
	}
}
