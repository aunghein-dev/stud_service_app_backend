package receipt

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/receipt"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, tx *sql.Tx, rc *receipt.Receipt) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO receipts
	(tenant_id, receipt_no, receipt_type, student_id, enrollment_id, payment_id, class_course_id, total_amount, paid_amount, remaining_amount, payload_json, issued_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	RETURNING id, created_at`
	return tx.QueryRowContext(ctx, query,
		tenantID, rc.ReceiptNo, rc.ReceiptType, rc.StudentID, rc.EnrollmentID, rc.PaymentID, rc.ClassCourseID,
		rc.TotalAmount, rc.PaidAmount, rc.RemainingAmount, rc.PayloadJSON, rc.IssuedAt,
	).Scan(&rc.ID, &rc.CreatedAt)
}

func (r *postgresRepository) List(ctx context.Context, filter common.ListFilter) ([]receipt.Receipt, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	where := []string{"tenant_id=$1"}
	if filter.ReceiptNo != "" {
		args = append(args, "%"+filter.ReceiptNo+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("receipt_no ILIKE $%d", idx))
	}
	if filter.DateFrom != "" {
		args = append(args, filter.DateFrom)
		idx := len(args)
		where = append(where, fmt.Sprintf("issued_at >= $%d", idx))
	}
	if filter.DateTo != "" {
		args = append(args, filter.DateTo)
		idx := len(args)
		where = append(where, fmt.Sprintf("issued_at <= $%d", idx))
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	args = append(args, filter.Limit, filter.Offset)
	lIdx := len(args) - 1
	oIdx := len(args)

	query := fmt.Sprintf(`SELECT id, receipt_no, receipt_type, student_id, enrollment_id, payment_id, class_course_id, total_amount, paid_amount, remaining_amount, payload_json, issued_at, created_at
	FROM receipts WHERE %s ORDER BY issued_at DESC, created_at DESC LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), lIdx, oIdx)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]receipt.Receipt, 0)
	for rows.Next() {
		var rc receipt.Receipt
		if err := rows.Scan(&rc.ID, &rc.ReceiptNo, &rc.ReceiptType, &rc.StudentID, &rc.EnrollmentID, &rc.PaymentID, &rc.ClassCourseID,
			&rc.TotalAmount, &rc.PaidAmount, &rc.RemainingAmount, &rc.PayloadJSON, &rc.IssuedAt, &rc.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, rc)
	}
	return result, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*receipt.Receipt, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, receipt_no, receipt_type, student_id, enrollment_id, payment_id, class_course_id, total_amount, paid_amount, remaining_amount, payload_json, issued_at, created_at
	FROM receipts WHERE id=$1 AND tenant_id=$2`
	var rc receipt.Receipt
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&rc.ID, &rc.ReceiptNo, &rc.ReceiptType, &rc.StudentID, &rc.EnrollmentID, &rc.PaymentID, &rc.ClassCourseID,
		&rc.TotalAmount, &rc.PaidAmount, &rc.RemainingAmount, &rc.PayloadJSON, &rc.IssuedAt, &rc.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rc, nil
}

func (r *postgresRepository) GetByReceiptNo(ctx context.Context, receiptNo string) (*receipt.Receipt, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, receipt_no, receipt_type, student_id, enrollment_id, payment_id, class_course_id, total_amount, paid_amount, remaining_amount, payload_json, issued_at, created_at
	FROM receipts WHERE receipt_no=$1 AND tenant_id=$2`
	var rc receipt.Receipt
	if err := r.db.QueryRowContext(ctx, query, receiptNo, tenantID).Scan(&rc.ID, &rc.ReceiptNo, &rc.ReceiptType, &rc.StudentID, &rc.EnrollmentID, &rc.PaymentID, &rc.ClassCourseID,
		&rc.TotalAmount, &rc.PaidAmount, &rc.RemainingAmount, &rc.PayloadJSON, &rc.IssuedAt, &rc.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rc, nil
}
