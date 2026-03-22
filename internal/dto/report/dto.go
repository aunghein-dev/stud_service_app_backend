package report

type DashboardResponse struct {
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

type GrossRowResponse struct {
	ClassCourseID int64   `json:"class_course_id"`
	ClassName     string  `json:"class_name"`
	Income        float64 `json:"income"`
	Expenses      float64 `json:"expenses"`
	Gross         float64 `json:"gross"`
}

type GrossResponse struct {
	Rows          []GrossRowResponse `json:"rows"`
	TotalIncome   float64            `json:"total_income"`
	TotalExpenses float64            `json:"total_expenses"`
	TotalGross    float64            `json:"total_gross"`
}

type TrendPoint struct {
	Period string  `json:"period"`
	Value  float64 `json:"value"`
}

type PerformanceResponse struct {
	BestByIncomeClass string       `json:"best_by_income_class"`
	BestByGrossClass  string       `json:"best_by_gross_class"`
	IncomeTrend       []TrendPoint `json:"income_trend"`
	ExpenseTrend      []TrendPoint `json:"expense_trend"`
	GrossTrend        []TrendPoint `json:"gross_trend"`
}
