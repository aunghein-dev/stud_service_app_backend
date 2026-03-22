package mapper

import (
	domain "student_service_app/backend/internal/domain/report"
	dto "student_service_app/backend/internal/dto/report"
)

func DashboardToDTO(r domain.DashboardSummary) dto.DashboardResponse {
	return dto.DashboardResponse{
		TotalStudents:      r.TotalStudents,
		TotalTeachers:      r.TotalTeachers,
		TotalActiveClasses: r.TotalActiveClasses,
		TodayIncome:        r.TodayIncome,
		TodayExpenses:      r.TodayExpenses,
		TodayGross:         r.TodayGross,
		MonthlyIncome:      r.MonthlyIncome,
		MonthlyExpenses:    r.MonthlyExpenses,
		MonthlyGross:       r.MonthlyGross,
		PendingDuesCount:   r.PendingDuesCount,
	}
}

func GrossToDTO(r domain.GrossReport) dto.GrossResponse {
	rows := make([]dto.GrossRowResponse, 0, len(r.Rows))
	for _, row := range r.Rows {
		rows = append(rows, dto.GrossRowResponse{
			ClassCourseID: row.ClassCourseID,
			ClassName:     row.ClassName,
			Income:        row.Income,
			Expenses:      row.Expenses,
			Gross:         row.Gross,
		})
	}
	return dto.GrossResponse{
		Rows:          rows,
		TotalIncome:   r.TotalIncome,
		TotalExpenses: r.TotalExpenses,
		TotalGross:    r.TotalGross,
	}
}
