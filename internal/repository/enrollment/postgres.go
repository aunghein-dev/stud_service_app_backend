package enrollment

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/enrollment"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) ExistsDuplicate(ctx context.Context, studentID, classCourseID int64) (bool, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return false, err
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM enrollments WHERE student_id=$1 AND class_course_id=$2 AND tenant_id=$3)`, studentID, classCourseID, tenantID).Scan(&exists)
	return exists, err
}

func (r *postgresRepository) Create(ctx context.Context, tx *sql.Tx, e *enrollment.Enrollment) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO enrollments
	(tenant_id, enrollment_code, student_id, class_course_id, enrollment_date, sub_total, discount_amount, final_fee, paid_amount, remaining_amount, payment_status, note)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	RETURNING id, created_at, updated_at`
	return tx.QueryRowContext(ctx, query,
		tenantID, e.EnrollmentCode, e.StudentID, e.ClassCourseID, e.EnrollmentDate, e.SubTotal,
		e.DiscountAmount, e.FinalFee, e.PaidAmount, e.RemainingAmount, e.PaymentStatus, e.Note,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *postgresRepository) AddOptionalItems(ctx context.Context, tx *sql.Tx, items []enrollment.EnrollmentOptionalItem) error {
	if len(items) == 0 {
		return nil
	}
	query := `INSERT INTO enrollment_optional_items
	(enrollment_id, optional_fee_item_id, item_name_snapshot, amount_snapshot, quantity, total_amount)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`
	for i := range items {
		if err := tx.QueryRowContext(ctx, query,
			items[i].EnrollmentID, items[i].OptionalFeeItemID, items[i].ItemNameSnapshot,
			items[i].AmountSnapshot, items[i].Quantity, items[i].TotalAmount,
		).Scan(&items[i].ID, &items[i].CreatedAt); err != nil {
			return err
		}
	}
	return nil
}

func (r *postgresRepository) List(ctx context.Context, filter common.ListFilter) ([]enrollment.Enrollment, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	where := []string{"e.tenant_id = $1"}
	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`(
			e.enrollment_code ILIKE $%d OR
			e.student_id IN (
				SELECT id FROM students
				WHERE tenant_id = $1 AND (full_name ILIKE $%d OR guardian_name ILIKE $%d OR student_code ILIKE $%d OR phone ILIKE $%d OR guardian_phone ILIKE $%d)
			) OR
			e.class_course_id IN (
				SELECT id FROM class_courses
				WHERE tenant_id = $1 AND (class_name ILIKE $%d OR course_name ILIKE $%d OR course_code ILIKE $%d)
			)
		)`, idx, idx, idx, idx, idx, idx, idx, idx, idx))
	}
	if filter.PaymentStatus != "" {
		args = append(args, filter.PaymentStatus)
		idx := len(args)
		where = append(where, fmt.Sprintf("e.payment_status=$%d", idx))
	}
	if filter.StudentName != "" {
		args = append(args, "%"+filter.StudentName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`e.student_id IN (SELECT id FROM students WHERE tenant_id = $1 AND full_name ILIKE $%d)`, idx))
	}
	if filter.ClassName != "" {
		args = append(args, "%"+filter.ClassName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`e.class_course_id IN (SELECT id FROM class_courses WHERE tenant_id = $1 AND class_name ILIKE $%d)`, idx))
	}
	if filter.DateFrom != "" {
		args = append(args, filter.DateFrom)
		idx := len(args)
		where = append(where, fmt.Sprintf("e.enrollment_date >= $%d", idx))
	}
	if filter.DateTo != "" {
		args = append(args, filter.DateTo)
		idx := len(args)
		where = append(where, fmt.Sprintf("e.enrollment_date <= $%d", idx))
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	args = append(args, filter.Limit, filter.Offset)
	lIdx := len(args) - 1
	oIdx := len(args)
	query := fmt.Sprintf(`SELECT
	e.id, e.enrollment_code, e.student_id, COALESCE(s.full_name, ''), COALESCE(s.guardian_name, ''),
	e.class_course_id, COALESCE(c.class_name, ''), COALESCE(c.course_name, ''),
	e.enrollment_date, e.sub_total, e.discount_amount, e.final_fee, e.paid_amount,
	e.remaining_amount, e.payment_status, e.note, e.created_at, e.updated_at
	FROM enrollments e
	LEFT JOIN students s ON s.id = e.student_id AND s.tenant_id = e.tenant_id
	LEFT JOIN class_courses c ON c.id = e.class_course_id AND c.tenant_id = e.tenant_id
	WHERE %s
	ORDER BY e.enrollment_date DESC, e.created_at DESC
	LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), lIdx, oIdx)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]enrollment.Enrollment, 0)
	for rows.Next() {
		var e enrollment.Enrollment
		if err := rows.Scan(
			&e.ID, &e.EnrollmentCode, &e.StudentID, &e.StudentName, &e.GuardianName,
			&e.ClassCourseID, &e.ClassName, &e.CourseName,
			&e.EnrollmentDate, &e.SubTotal, &e.DiscountAmount, &e.FinalFee, &e.PaidAmount,
			&e.RemainingAmount, &e.PaymentStatus, &e.Note, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, rows.Err()
}

func (r *postgresRepository) ListByStudent(ctx context.Context, studentID int64) ([]enrollment.Enrollment, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT
	e.id, e.enrollment_code, e.student_id, COALESCE(s.full_name, ''), COALESCE(s.guardian_name, ''),
	e.class_course_id, COALESCE(c.class_name, ''), COALESCE(c.course_name, ''),
	e.enrollment_date, e.sub_total, e.discount_amount, e.final_fee, e.paid_amount,
	e.remaining_amount, e.payment_status, e.note, e.created_at, e.updated_at
	FROM enrollments e
	LEFT JOIN students s ON s.id = e.student_id AND s.tenant_id = e.tenant_id
	LEFT JOIN class_courses c ON c.id = e.class_course_id AND c.tenant_id = e.tenant_id
	WHERE e.student_id=$1 AND e.tenant_id=$2
	ORDER BY e.enrollment_date DESC, e.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, studentID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]enrollment.Enrollment, 0)
	for rows.Next() {
		var e enrollment.Enrollment
		if err := rows.Scan(
			&e.ID, &e.EnrollmentCode, &e.StudentID, &e.StudentName, &e.GuardianName,
			&e.ClassCourseID, &e.ClassName, &e.CourseName,
			&e.EnrollmentDate, &e.SubTotal, &e.DiscountAmount, &e.FinalFee, &e.PaidAmount,
			&e.RemainingAmount, &e.PaymentStatus, &e.Note, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*enrollment.Enrollment, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	var e enrollment.Enrollment
	query := `SELECT
	e.id, e.enrollment_code, e.student_id, COALESCE(s.full_name, ''), COALESCE(s.guardian_name, ''),
	e.class_course_id, COALESCE(c.class_name, ''), COALESCE(c.course_name, ''),
	e.enrollment_date, e.sub_total, e.discount_amount, e.final_fee, e.paid_amount,
	e.remaining_amount, e.payment_status, e.note, e.created_at, e.updated_at
	FROM enrollments e
	LEFT JOIN students s ON s.id = e.student_id AND s.tenant_id = e.tenant_id
	LEFT JOIN class_courses c ON c.id = e.class_course_id AND c.tenant_id = e.tenant_id
	WHERE e.id=$1 AND e.tenant_id=$2`
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&e.ID, &e.EnrollmentCode, &e.StudentID, &e.StudentName, &e.GuardianName,
		&e.ClassCourseID, &e.ClassName, &e.CourseName,
		&e.EnrollmentDate, &e.SubTotal, &e.DiscountAmount, &e.FinalFee, &e.PaidAmount,
		&e.RemainingAmount, &e.PaymentStatus, &e.Note, &e.CreatedAt, &e.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *postgresRepository) ListOptionalItems(ctx context.Context, enrollmentID int64) ([]enrollment.EnrollmentOptionalItem, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, enrollment_id, optional_fee_item_id, item_name_snapshot, amount_snapshot, quantity, total_amount, created_at
	FROM enrollment_optional_items
	WHERE enrollment_id=$1 AND enrollment_id IN (SELECT id FROM enrollments WHERE id=$1 AND tenant_id=$2)
	ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, enrollmentID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]enrollment.EnrollmentOptionalItem, 0)
	for rows.Next() {
		var i enrollment.EnrollmentOptionalItem
		if err := rows.Scan(&i.ID, &i.EnrollmentID, &i.OptionalFeeItemID, &i.ItemNameSnapshot, &i.AmountSnapshot, &i.Quantity, &i.TotalAmount, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

func (r *postgresRepository) Update(ctx context.Context, e *enrollment.Enrollment) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE enrollments SET discount_amount=$2, final_fee=$3, remaining_amount=$4, payment_status=$5, note=$6, updated_at=NOW() WHERE id=$1 AND tenant_id=$7 RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query, e.ID, e.DiscountAmount, e.FinalFee, e.RemainingAmount, e.PaymentStatus, e.Note, tenantID).Scan(&e.UpdatedAt)
}

func (r *postgresRepository) UpdatePaymentState(ctx context.Context, tx *sql.Tx, id int64, paidAmount, remainingAmount float64, paymentStatus string) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE enrollments SET paid_amount=$2, remaining_amount=$3, payment_status=$4, updated_at=NOW() WHERE id=$1 AND tenant_id=$5`, id, paidAmount, remainingAmount, paymentStatus, tenantID)
	return err
}
