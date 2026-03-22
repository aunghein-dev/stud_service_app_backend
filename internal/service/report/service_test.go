package report

import (
	"context"
	"testing"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/report"

	"github.com/stretchr/testify/require"
)

type fakeRepo struct{}

func (f *fakeRepo) Dashboard(ctx context.Context) (*report.DashboardSummary, error) {
	return &report.DashboardSummary{TodayIncome: 100, TodayExpenses: 40}, nil
}

func (f *fakeRepo) Gross(ctx context.Context, filter common.ListFilter) (*report.GrossReport, error) {
	return &report.GrossReport{TotalIncome: 1000, TotalExpenses: 300, TotalGross: 700}, nil
}

func (f *fakeRepo) MonthlyTrend(ctx context.Context, source string, filter common.ListFilter) (map[string]float64, error) {
	if source == "income" {
		return map[string]float64{"2026-01": 300, "2026-02": 500}, nil
	}
	return map[string]float64{"2026-01": 100, "2026-02": 150}, nil
}

func (f *fakeRepo) BestClassByMetric(ctx context.Context, metric string, filter common.ListFilter) (string, error) {
	if metric == "income" {
		return "General English", nil
	}
	return "Academic Math", nil
}

func TestPerformance(t *testing.T) {
	svc := NewService(&fakeRepo{})
	res, err := svc.Performance(context.Background(), common.ListFilter{})
	require.NoError(t, err)
	require.Equal(t, "General English", res.BestByIncomeClass)
	require.Equal(t, "Academic Math", res.BestByGrossClass)
	require.Len(t, res.GrossTrend, 2)
	require.Equal(t, 200.0, res.GrossTrend[0].Value)
}
