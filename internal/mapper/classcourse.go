package mapper

import (
	"strings"

	domain "student_service_app/backend/internal/domain/classcourse"
	dto "student_service_app/backend/internal/dto/classcourse"
)

func ClassCourseToDTO(c domain.ClassCourse) dto.Response {
	days := []string{}
	if c.DaysOfWeek != "" {
		days = strings.Split(c.DaysOfWeek, ",")
	}
	return dto.Response{
		ID:                c.ID,
		CourseCode:        c.CourseCode,
		CourseName:        c.CourseName,
		ClassName:         c.ClassName,
		Category:          c.Category,
		Subject:           c.Subject,
		Level:             c.Level,
		StartDate:         DateStringPtr(c.StartDate),
		EndDate:           DateStringPtr(c.EndDate),
		ScheduleText:      c.ScheduleText,
		DaysOfWeek:        days,
		TimeStart:         c.TimeStart,
		TimeEnd:           c.TimeEnd,
		Room:              c.Room,
		AssignedTeacherID: c.AssignedTeacherID,
		MaxStudents:       c.MaxStudents,
		Status:            c.Status,
		BaseCourseFee:     c.BaseCourseFee,
		RegistrationFee:   c.RegistrationFee,
		ExamFee:           c.ExamFee,
		CertificateFee:    c.CertificateFee,
		Note:              c.Note,
		CreatedAt:         DateTimeString(c.CreatedAt),
		UpdatedAt:         DateTimeString(c.UpdatedAt),
	}
}

func OptionalFeeToDTO(item domain.OptionalFeeItem) dto.OptionalFeeItemResponse {
	return dto.OptionalFeeItemResponse{
		ID:            item.ID,
		ClassCourseID: item.ClassCourseID,
		ItemName:      item.ItemName,
		DefaultAmount: item.DefaultAmount,
		IsOptional:    item.IsOptional,
		IsActive:      item.IsActive,
	}
}
