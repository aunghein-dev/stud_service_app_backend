package teacher

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/teacher"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, t *teacher.Teacher) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO teachers (tenant_id, teacher_code, teacher_name, phone, address, subject_specialty, salary_type, default_fee_amount, note, is_active)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		tenantID, t.TeacherCode, t.TeacherName, t.Phone, t.Address, t.SubjectSpecialty,
		t.SalaryType, t.DefaultFeeAmount, t.Note, t.IsActive,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *postgresRepository) List(ctx context.Context, filter common.ListFilter) ([]teacher.Teacher, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	where := []string{"tenant_id = $1"}
	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(teacher_name ILIKE $%d OR teacher_code ILIKE $%d)", idx, idx))
	}
	if filter.TeacherName != "" {
		args = append(args, "%"+filter.TeacherName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("teacher_name ILIKE $%d", idx))
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	args = append(args, filter.Limit, filter.Offset)
	lIdx := len(args) - 1
	oIdx := len(args)
	query := fmt.Sprintf(`SELECT id, teacher_code, teacher_name, phone, address, subject_specialty, salary_type, default_fee_amount, note, is_active, created_at, updated_at
	FROM teachers WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), lIdx, oIdx)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]teacher.Teacher, 0)
	for rows.Next() {
		var t teacher.Teacher
		if err := rows.Scan(&t.ID, &t.TeacherCode, &t.TeacherName, &t.Phone, &t.Address, &t.SubjectSpecialty, &t.SalaryType,
			&t.DefaultFeeAmount, &t.Note, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*teacher.Teacher, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, teacher_code, teacher_name, phone, address, subject_specialty, salary_type, default_fee_amount, note, is_active, created_at, updated_at
	FROM teachers WHERE id=$1 AND tenant_id=$2`
	var t teacher.Teacher
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&t.ID, &t.TeacherCode, &t.TeacherName, &t.Phone, &t.Address, &t.SubjectSpecialty, &t.SalaryType,
		&t.DefaultFeeAmount, &t.Note, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *postgresRepository) GetByCode(ctx context.Context, code string) (*teacher.Teacher, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, teacher_code, teacher_name, phone, address, subject_specialty, salary_type, default_fee_amount, note, is_active, created_at, updated_at
	FROM teachers WHERE teacher_code=$1 AND tenant_id=$2`
	var t teacher.Teacher
	if err := r.db.QueryRowContext(ctx, query, code, tenantID).Scan(&t.ID, &t.TeacherCode, &t.TeacherName, &t.Phone, &t.Address, &t.SubjectSpecialty, &t.SalaryType,
		&t.DefaultFeeAmount, &t.Note, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *postgresRepository) Update(ctx context.Context, t *teacher.Teacher) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE teachers SET teacher_name=$2, phone=$3, address=$4, subject_specialty=$5, salary_type=$6, default_fee_amount=$7,
	note=$8, is_active=$9, updated_at=NOW() WHERE id=$1 AND tenant_id=$10 RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query, t.ID, t.TeacherName, t.Phone, t.Address, t.SubjectSpecialty, t.SalaryType, t.DefaultFeeAmount,
		t.Note, t.IsActive, tenantID).Scan(&t.UpdatedAt)
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `UPDATE teachers SET is_active=false, updated_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

func (r *postgresRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return false, err
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM teachers WHERE id=$1 AND tenant_id=$2)`, id, tenantID).Scan(&exists)
	return exists, err
}
