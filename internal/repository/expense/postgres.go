package expense

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/expense"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, e *expense.Expense) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO expense_transactions
	(tenant_id, expense_date, expense_type, teacher_id, class_course_id, amount, description, payment_method, reference_no)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		tenantID, e.ExpenseDate, e.ExpenseType, e.TeacherID, e.ClassCourseID, e.Amount, e.Description, e.PaymentMethod, e.ReferenceNo,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *postgresRepository) List(ctx context.Context, filter common.ListFilter) ([]expense.Expense, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	where := []string{"tenant_id=$1"}
	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(expense_type::text ILIKE $%d OR description ILIKE $%d OR reference_no ILIKE $%d)", idx, idx, idx))
	}
	if filter.ExpenseType != "" {
		args = append(args, filter.ExpenseType)
		idx := len(args)
		where = append(where, fmt.Sprintf("expense_type=$%d", idx))
	}
	if filter.DateFrom != "" {
		args = append(args, filter.DateFrom)
		idx := len(args)
		where = append(where, fmt.Sprintf("expense_date >= $%d", idx))
	}
	if filter.DateTo != "" {
		args = append(args, filter.DateTo)
		idx := len(args)
		where = append(where, fmt.Sprintf("expense_date <= $%d", idx))
	}
	if filter.TeacherName != "" {
		args = append(args, "%"+filter.TeacherName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`teacher_id IN (SELECT id FROM teachers WHERE tenant_id=$1 AND teacher_name ILIKE $%d)`, idx))
	}
	if filter.ClassName != "" {
		args = append(args, "%"+filter.ClassName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`class_course_id IN (SELECT id FROM class_courses WHERE tenant_id=$1 AND class_name ILIKE $%d)`, idx))
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	args = append(args, filter.Limit, filter.Offset)
	lIdx := len(args) - 1
	oIdx := len(args)

	query := fmt.Sprintf(`SELECT id, expense_date, expense_type, teacher_id, class_course_id, amount, description, payment_method, reference_no, created_at, updated_at
	FROM expense_transactions WHERE %s ORDER BY expense_date DESC, created_at DESC LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), lIdx, oIdx)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]expense.Expense, 0)
	for rows.Next() {
		var e expense.Expense
		if err := rows.Scan(&e.ID, &e.ExpenseDate, &e.ExpenseType, &e.TeacherID, &e.ClassCourseID, &e.Amount, &e.Description,
			&e.PaymentMethod, &e.ReferenceNo, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*expense.Expense, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, expense_date, expense_type, teacher_id, class_course_id, amount, description, payment_method, reference_no, created_at, updated_at
	FROM expense_transactions WHERE id=$1 AND tenant_id=$2`
	var e expense.Expense
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&e.ID, &e.ExpenseDate, &e.ExpenseType, &e.TeacherID, &e.ClassCourseID,
		&e.Amount, &e.Description, &e.PaymentMethod, &e.ReferenceNo, &e.CreatedAt, &e.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *postgresRepository) Update(ctx context.Context, e *expense.Expense) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE expense_transactions SET expense_date=$2, expense_type=$3, teacher_id=$4, class_course_id=$5, amount=$6, description=$7,
	payment_method=$8, reference_no=$9, updated_at=NOW() WHERE id=$1 AND tenant_id=$10 RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		e.ID, e.ExpenseDate, e.ExpenseType, e.TeacherID, e.ClassCourseID, e.Amount, e.Description, e.PaymentMethod, e.ReferenceNo, tenantID,
	).Scan(&e.UpdatedAt)
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM expense_transactions WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}
