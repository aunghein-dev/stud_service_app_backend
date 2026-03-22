package report

import (
	"context"
	"sort"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/dto/report"
	"student_service_app/backend/internal/mapper"
	reportrepo "student_service_app/backend/internal/repository/report"
)

type Service interface {
	Dashboard(ctx context.Context) (*report.DashboardResponse, error)
	Gross(ctx context.Context, filter common.ListFilter) (*report.GrossResponse, error)
	StudentReport(ctx context.Context, filter common.ListFilter) (map[string]any, error)
	TeacherReport(ctx context.Context, filter common.ListFilter) (map[string]any, error)
	ClassCourseReport(ctx context.Context, filter common.ListFilter) (map[string]any, error)
	TransactionReport(ctx context.Context, filter common.ListFilter) (map[string]any, error)
	Performance(ctx context.Context, filter common.ListFilter) (*report.PerformanceResponse, error)
}

type service struct {
	repo reportrepo.Repository
}

func NewService(repo reportrepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Dashboard(ctx context.Context) (*report.DashboardResponse, error) {
	result, err := s.repo.Dashboard(ctx)
	if err != nil {
		return nil, err
	}
	mapped := mapper.DashboardToDTO(*result)
	return &mapped, nil
}

func (s *service) Gross(ctx context.Context, filter common.ListFilter) (*report.GrossResponse, error) {
	result, err := s.repo.Gross(ctx, filter)
	if err != nil {
		return nil, err
	}
	mapped := mapper.GrossToDTO(*result)
	return &mapped, nil
}

func (s *service) StudentReport(ctx context.Context, filter common.ListFilter) (map[string]any, error) {
	gross, err := s.repo.Gross(ctx, filter)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"summary": map[string]any{
			"total_income":          gross.TotalIncome,
			"total_unpaid_estimate": 0,
		},
		"filters": filter,
	}, nil
}

func (s *service) TeacherReport(ctx context.Context, filter common.ListFilter) (map[string]any, error) {
	gross, err := s.repo.Gross(ctx, filter)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"summary": map[string]any{
			"linked_income":   gross.TotalIncome,
			"linked_expenses": gross.TotalExpenses,
		},
		"filters": filter,
	}, nil
}

func (s *service) ClassCourseReport(ctx context.Context, filter common.ListFilter) (map[string]any, error) {
	gross, err := s.repo.Gross(ctx, filter)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"rows": gross.Rows,
		"totals": map[string]any{
			"income":   gross.TotalIncome,
			"expenses": gross.TotalExpenses,
			"gross":    gross.TotalGross,
		},
	}, nil
}

func (s *service) TransactionReport(ctx context.Context, filter common.ListFilter) (map[string]any, error) {
	dashboard, err := s.repo.Dashboard(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"today_income":       dashboard.TodayIncome,
		"today_expenses":     dashboard.TodayExpenses,
		"pending_dues_count": dashboard.PendingDuesCount,
		"filters":            filter,
	}, nil
}

func (s *service) Performance(ctx context.Context, filter common.ListFilter) (*report.PerformanceResponse, error) {
	bestIncomeClass, err := s.repo.BestClassByMetric(ctx, "income", filter)
	if err != nil {
		return nil, err
	}
	bestGrossClass, err := s.repo.BestClassByMetric(ctx, "gross", filter)
	if err != nil {
		return nil, err
	}
	incomeTrendMap, err := s.repo.MonthlyTrend(ctx, "income", filter)
	if err != nil {
		return nil, err
	}
	expenseTrendMap, err := s.repo.MonthlyTrend(ctx, "expense", filter)
	if err != nil {
		return nil, err
	}

	months := make([]string, 0)
	for m := range incomeTrendMap {
		months = append(months, m)
	}
	for m := range expenseTrendMap {
		if _, ok := incomeTrendMap[m]; !ok {
			months = append(months, m)
		}
	}
	sort.Strings(months)

	incomeTrend := make([]report.TrendPoint, 0, len(months))
	expenseTrend := make([]report.TrendPoint, 0, len(months))
	grossTrend := make([]report.TrendPoint, 0, len(months))
	for _, m := range months {
		income := incomeTrendMap[m]
		expense := expenseTrendMap[m]
		incomeTrend = append(incomeTrend, report.TrendPoint{Period: m, Value: income})
		expenseTrend = append(expenseTrend, report.TrendPoint{Period: m, Value: expense})
		grossTrend = append(grossTrend, report.TrendPoint{Period: m, Value: income - expense})
	}

	return &report.PerformanceResponse{
		BestByIncomeClass: bestIncomeClass,
		BestByGrossClass:  bestGrossClass,
		IncomeTrend:       incomeTrend,
		ExpenseTrend:      expenseTrend,
		GrossTrend:        grossTrend,
	}, nil
}
