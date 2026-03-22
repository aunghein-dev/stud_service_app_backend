package report

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/report"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Dashboard(ctx context.Context) (*report.DashboardSummary, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	res := &report.DashboardSummary{}
	q := `
	SELECT
		(SELECT COUNT(*) FROM students WHERE tenant_id=$1) AS total_students,
		(SELECT COUNT(*) FROM teachers WHERE tenant_id=$1) AS total_teachers,
		(SELECT COUNT(*) FROM class_courses WHERE tenant_id=$1 AND status IN ('open','running')) AS total_active_classes,
		COALESCE((SELECT SUM(amount) FROM payment_transactions WHERE tenant_id=$1 AND payment_date = CURRENT_DATE),0) AS today_income,
		COALESCE((SELECT SUM(amount) FROM expense_transactions WHERE tenant_id=$1 AND expense_date = CURRENT_DATE),0) AS today_expenses,
		COALESCE((SELECT COUNT(*) FROM enrollments WHERE tenant_id=$1 AND payment_status IN ('unpaid','partial')),0) AS pending_dues_count,
		COALESCE((SELECT SUM(amount) FROM payment_transactions WHERE tenant_id=$1 AND DATE_TRUNC('month', payment_date)=DATE_TRUNC('month', CURRENT_DATE)),0) AS monthly_income,
		COALESCE((SELECT SUM(amount) FROM expense_transactions WHERE tenant_id=$1 AND DATE_TRUNC('month', expense_date)=DATE_TRUNC('month', CURRENT_DATE)),0) AS monthly_expenses`
	if err := r.db.QueryRowContext(ctx, q, tenantID).Scan(
		&res.TotalStudents, &res.TotalTeachers, &res.TotalActiveClasses,
		&res.TodayIncome, &res.TodayExpenses, &res.PendingDuesCount,
		&res.MonthlyIncome, &res.MonthlyExpenses,
	); err != nil {
		return nil, err
	}
	res.TodayGross = res.TodayIncome - res.TodayExpenses
	res.MonthlyGross = res.MonthlyIncome - res.MonthlyExpenses
	return res, nil
}

func (r *postgresRepository) Gross(ctx context.Context, filter common.ListFilter) (*report.GrossReport, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	whereIncome := []string{"p.tenant_id = $1"}
	whereExpense := []string{"e.tenant_id = $1"}
	whereBase := []string{"c.tenant_id = $1"}
	if filter.DateFrom != "" {
		args = append(args, filter.DateFrom)
		idx := len(args)
		whereIncome = append(whereIncome, fmt.Sprintf("p.payment_date >= $%d", idx))
		whereExpense = append(whereExpense, fmt.Sprintf("e.expense_date >= $%d", idx))
	}
	if filter.DateTo != "" {
		args = append(args, filter.DateTo)
		idx := len(args)
		whereIncome = append(whereIncome, fmt.Sprintf("p.payment_date <= $%d", idx))
		whereExpense = append(whereExpense, fmt.Sprintf("e.expense_date <= $%d", idx))
	}
	if filter.ClassName != "" {
		args = append(args, "%"+filter.ClassName+"%")
		idx := len(args)
		whereIncome = append(whereIncome, fmt.Sprintf("c.class_name ILIKE $%d", idx))
		whereExpense = append(whereExpense, fmt.Sprintf("c.class_name ILIKE $%d", idx))
		whereBase = append(whereBase, fmt.Sprintf("c.class_name ILIKE $%d", idx))
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	args = append(args, filter.Limit, filter.Offset)
	lIdx := len(args) - 1
	oIdx := len(args)

	query := fmt.Sprintf(`
	WITH income AS (
		SELECT p.class_course_id, SUM(p.amount) AS income
		FROM payment_transactions p
		JOIN class_courses c ON c.id = p.class_course_id
		WHERE %s
		GROUP BY p.class_course_id
	),
	expense AS (
		SELECT e.class_course_id, SUM(e.amount) AS expense
		FROM expense_transactions e
		JOIN class_courses c ON c.id = e.class_course_id
		WHERE e.class_course_id IS NOT NULL AND %s
		GROUP BY e.class_course_id
	),
	base AS (
		SELECT
			c.id AS class_course_id,
			c.class_name,
			COALESCE(i.income,0) AS income,
			COALESCE(ex.expense,0) AS expenses,
			(COALESCE(i.income,0) - COALESCE(ex.expense,0)) AS gross
		FROM class_courses c
		LEFT JOIN income i ON i.class_course_id = c.id
		LEFT JOIN expense ex ON ex.class_course_id = c.id
		WHERE %s
	)
	SELECT
		class_course_id,
		class_name,
		income,
		expenses,
		gross,
		COALESCE(SUM(income) OVER (), 0) AS total_income,
		COALESCE(SUM(expenses) OVER (), 0) AS total_expenses,
		COALESCE(SUM(gross) OVER (), 0) AS total_gross
	FROM base
	ORDER BY class_name ASC
	LIMIT $%d OFFSET $%d`, strings.Join(whereIncome, " AND "), strings.Join(whereExpense, " AND "), strings.Join(whereBase, " AND "), lIdx, oIdx)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &report.GrossReport{Rows: make([]report.GrossRow, 0)}
	for rows.Next() {
		var row report.GrossRow
		var totalIncome float64
		var totalExpenses float64
		var totalGross float64
		if err := rows.Scan(&row.ClassCourseID, &row.ClassName, &row.Income, &row.Expenses, &row.Gross, &totalIncome, &totalExpenses, &totalGross); err != nil {
			return nil, err
		}
		resp.Rows = append(resp.Rows, row)
		resp.TotalIncome = totalIncome
		resp.TotalExpenses = totalExpenses
		resp.TotalGross = totalGross
	}
	return resp, rows.Err()
}

func (r *postgresRepository) MonthlyTrend(ctx context.Context, source string, filter common.ListFilter) (map[string]float64, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := ""
	switch source {
	case "income":
		query = `SELECT TO_CHAR(DATE_TRUNC('month', payment_date), 'YYYY-MM') AS ym, SUM(amount)
		FROM payment_transactions
		WHERE tenant_id=$1
		GROUP BY DATE_TRUNC('month', payment_date) ORDER BY ym ASC`
	case "expense":
		query = `SELECT TO_CHAR(DATE_TRUNC('month', expense_date), 'YYYY-MM') AS ym, SUM(amount)
		FROM expense_transactions
		WHERE tenant_id=$1
		GROUP BY DATE_TRUNC('month', expense_date) ORDER BY ym ASC`
	default:
		return map[string]float64{}, nil
	}

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]float64{}
	for rows.Next() {
		var month string
		var value float64
		if err := rows.Scan(&month, &value); err != nil {
			return nil, err
		}
		out[month] = value
	}
	return out, rows.Err()
}

func (r *postgresRepository) BestClassByMetric(ctx context.Context, metric string, filter common.ListFilter) (string, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return "", err
	}

	if metric == "income" {
		query := `SELECT c.class_name
		FROM payment_transactions p
		JOIN class_courses c ON c.id=p.class_course_id
		WHERE p.tenant_id=$1 AND c.tenant_id=$1
		GROUP BY c.id, c.class_name
		ORDER BY SUM(p.amount) DESC
		LIMIT 1`
		var className string
		if err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&className); err != nil {
			if err == sql.ErrNoRows {
				return "", nil
			}
			return "", err
		}
		return className, nil
	}

	if metric == "gross" {
		query := `
		WITH income AS (
			SELECT class_course_id, SUM(amount) amount FROM payment_transactions WHERE tenant_id=$1 GROUP BY class_course_id
		), exp AS (
			SELECT class_course_id, SUM(amount) amount FROM expense_transactions WHERE tenant_id=$1 AND class_course_id IS NOT NULL GROUP BY class_course_id
		)
		SELECT c.class_name
		FROM class_courses c
		LEFT JOIN income i ON i.class_course_id = c.id
		LEFT JOIN exp e ON e.class_course_id = c.id
		WHERE c.tenant_id=$1
		ORDER BY (COALESCE(i.amount,0)-COALESCE(e.amount,0)) DESC
		LIMIT 1`
		var className string
		if err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&className); err != nil {
			if err == sql.ErrNoRows {
				return "", nil
			}
			return "", err
		}
		return className, nil
	}

	return "", nil
}
