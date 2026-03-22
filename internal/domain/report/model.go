package report

type DashboardSummary struct {
	TotalStudents      int64   `json:"total_students"`
	TotalTeachers      int64   `json:"total_teachers"`
	TotalActiveClasses int64   `json:"total_active_classes"`
	TodayIncome        float64 `json:"today_income"`
	TodayExpenses      float64 `json:"today_expenses"`
	TodayGross         float64 `json:"today_gross"`
	MonthlyIncome      float64 `json:"monthly_income"`
	MonthlyExpenses    float64 `json:"monthly_expenses"`
	MonthlyGross       float64 `json:"monthly_gross"`
	PendingDuesCount   int64   `json:"pending_dues_count"`
}

type GrossRow struct {
	ClassCourseID int64   `json:"class_course_id"`
	ClassName     string  `json:"class_name"`
	Income        float64 `json:"income"`
	Expenses      float64 `json:"expenses"`
	Gross         float64 `json:"gross"`
}

type GrossReport struct {
	Rows          []GrossRow `json:"rows"`
	TotalIncome   float64    `json:"total_income"`
	TotalExpenses float64    `json:"total_expenses"`
	TotalGross    float64    `json:"total_gross"`
}
