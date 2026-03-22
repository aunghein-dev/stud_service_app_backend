package payment

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/payment"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, tx *sql.Tx, p *payment.Payment) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO payment_transactions
	(tenant_id, receipt_no, student_id, enrollment_id, class_course_id, payment_date, payment_method, amount, note, received_by)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING id, created_at`
	return tx.QueryRowContext(ctx, query,
		tenantID, p.ReceiptNo, p.StudentID, p.EnrollmentID, p.ClassCourseID, p.PaymentDate, p.PaymentMethod,
		p.Amount, p.Note, p.ReceivedBy,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *postgresRepository) Update(ctx context.Context, tx *sql.Tx, p *payment.Payment) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE payment_transactions
		SET payment_date=$2, payment_method=$3, amount=$4, note=$5, received_by=$6
		WHERE id=$1 AND tenant_id=$7`
	_, err = tx.ExecContext(ctx, query, p.ID, p.PaymentDate, p.PaymentMethod, p.Amount, p.Note, p.ReceivedBy, tenantID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM payment_transactions WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

func (r *postgresRepository) List(ctx context.Context, filter common.ListFilter) ([]payment.Payment, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	where := []string{"tenant_id=$1"}
	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`(
			receipt_no ILIKE $%d OR
			student_id IN (
				SELECT id FROM students
				WHERE tenant_id=$1 AND (full_name ILIKE $%d OR guardian_name ILIKE $%d OR student_code ILIKE $%d OR phone ILIKE $%d OR guardian_phone ILIKE $%d)
			)
		)`, idx, idx, idx, idx, idx, idx))
	}
	if filter.ReceiptNo != "" {
		args = append(args, "%"+filter.ReceiptNo+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("receipt_no ILIKE $%d", idx))
	}
	if filter.DateFrom != "" {
		args = append(args, filter.DateFrom)
		idx := len(args)
		where = append(where, fmt.Sprintf("payment_date >= $%d", idx))
	}
	if filter.DateTo != "" {
		args = append(args, filter.DateTo)
		idx := len(args)
		where = append(where, fmt.Sprintf("payment_date <= $%d", idx))
	}
	if filter.StudentName != "" {
		args = append(args, "%"+filter.StudentName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`student_id IN (SELECT id FROM students WHERE tenant_id=$1 AND full_name ILIKE $%d)`, idx))
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
	query := fmt.Sprintf(`SELECT id, receipt_no, student_id, enrollment_id, class_course_id, payment_date, payment_method, amount, note, received_by, created_at
	FROM payment_transactions WHERE %s ORDER BY payment_date DESC, created_at DESC LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), lIdx, oIdx)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]payment.Payment, 0)
	for rows.Next() {
		var p payment.Payment
		if err := rows.Scan(&p.ID, &p.ReceiptNo, &p.StudentID, &p.EnrollmentID, &p.ClassCourseID, &p.PaymentDate,
			&p.PaymentMethod, &p.Amount, &p.Note, &p.ReceivedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*payment.Payment, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	var p payment.Payment
	query := `SELECT id, receipt_no, student_id, enrollment_id, class_course_id, payment_date, payment_method, amount, note, received_by, created_at
	FROM payment_transactions WHERE id=$1 AND tenant_id=$2`
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&p.ID, &p.ReceiptNo, &p.StudentID, &p.EnrollmentID, &p.ClassCourseID,
		&p.PaymentDate, &p.PaymentMethod, &p.Amount, &p.Note, &p.ReceivedBy, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresRepository) ListByEnrollment(ctx context.Context, enrollmentID int64) ([]payment.Payment, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, receipt_no, student_id, enrollment_id, class_course_id, payment_date, payment_method, amount, note, received_by, created_at
	FROM payment_transactions WHERE enrollment_id=$1 AND tenant_id=$2 ORDER BY payment_date DESC, created_at DESC`, enrollmentID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]payment.Payment, 0)
	for rows.Next() {
		var p payment.Payment
		if err := rows.Scan(&p.ID, &p.ReceiptNo, &p.StudentID, &p.EnrollmentID, &p.ClassCourseID, &p.PaymentDate,
			&p.PaymentMethod, &p.Amount, &p.Note, &p.ReceivedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}
