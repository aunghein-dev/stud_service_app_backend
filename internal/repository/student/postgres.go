package student

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/student"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, s *student.Student) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO students
		(tenant_id, student_code, full_name, gender, date_of_birth, phone, guardian_name, guardian_phone, address, school_name, grade_level, note, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		tenantID, s.StudentCode, s.FullName, s.Gender, s.DateOfBirth, s.Phone, s.GuardianName, s.GuardianPhone,
		s.Address, s.SchoolName, s.GradeLevel, s.Note, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *postgresRepository) List(ctx context.Context, filter common.ListFilter) ([]student.Student, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	where := []string{"tenant_id = $1"}

	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(full_name ILIKE $%d OR student_code ILIKE $%d OR phone ILIKE $%d)", idx, idx, idx))
	}
	if filter.StudentName != "" {
		args = append(args, "%"+filter.StudentName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("full_name ILIKE $%d", idx))
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	args = append(args, filter.Limit, filter.Offset)
	limitIdx := len(args) - 1
	offsetIdx := len(args)

	query := fmt.Sprintf(`
		SELECT id, student_code, full_name, gender, date_of_birth, phone, guardian_name, guardian_phone, address, school_name, grade_level, note, is_active, created_at, updated_at
		FROM students
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), limitIdx, offsetIdx)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]student.Student, 0)
	for rows.Next() {
		var s student.Student
		if err := rows.Scan(
			&s.ID, &s.StudentCode, &s.FullName, &s.Gender, &s.DateOfBirth, &s.Phone, &s.GuardianName,
			&s.GuardianPhone, &s.Address, &s.SchoolName, &s.GradeLevel, &s.Note, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*student.Student, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, student_code, full_name, gender, date_of_birth, phone, guardian_name, guardian_phone, address, school_name, grade_level, note, is_active, created_at, updated_at
		FROM students WHERE id = $1 AND tenant_id = $2`
	var s student.Student
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&s.ID, &s.StudentCode, &s.FullName, &s.Gender, &s.DateOfBirth, &s.Phone, &s.GuardianName,
		&s.GuardianPhone, &s.Address, &s.SchoolName, &s.GradeLevel, &s.Note, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresRepository) GetByCode(ctx context.Context, code string) (*student.Student, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, student_code, full_name, gender, date_of_birth, phone, guardian_name, guardian_phone, address, school_name, grade_level, note, is_active, created_at, updated_at FROM students WHERE student_code = $1 AND tenant_id = $2`
	var s student.Student
	if err := r.db.QueryRowContext(ctx, query, code, tenantID).Scan(
		&s.ID, &s.StudentCode, &s.FullName, &s.Gender, &s.DateOfBirth, &s.Phone, &s.GuardianName,
		&s.GuardianPhone, &s.Address, &s.SchoolName, &s.GradeLevel, &s.Note, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresRepository) Update(ctx context.Context, s *student.Student) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE students
		SET full_name=$2, gender=$3, date_of_birth=$4, phone=$5, guardian_name=$6, guardian_phone=$7, address=$8,
			school_name=$9, grade_level=$10, note=$11, is_active=$12, updated_at=NOW()
		WHERE id=$1 AND tenant_id=$13
		RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		s.ID, s.FullName, s.Gender, s.DateOfBirth, s.Phone, s.GuardianName, s.GuardianPhone,
		s.Address, s.SchoolName, s.GradeLevel, s.Note, s.IsActive, tenantID,
	).Scan(&s.UpdatedAt)
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `UPDATE students SET is_active=false, updated_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

func (r *postgresRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return false, err
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM students WHERE id=$1 AND tenant_id=$2)`, id, tenantID).Scan(&exists)
	return exists, err
}
