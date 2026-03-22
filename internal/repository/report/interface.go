package report

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/report"
)

type Repository interface {
	Dashboard(ctx context.Context) (*report.DashboardSummary, error)
	Gross(ctx context.Context, filter common.ListFilter) (*report.GrossReport, error)
	MonthlyTrend(ctx context.Context, source string, filter common.ListFilter) (map[string]float64, error)
	BestClassByMetric(ctx context.Context, metric string, filter common.ListFilter) (string, error)
}
