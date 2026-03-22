package classcourse

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"student_service_app/backend/internal/domain/classcourse"
	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, c *classcourse.ClassCourse) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO class_courses
	(tenant_id, course_code, course_name, class_name, category, subject, level, start_date, end_date, schedule_text, days_of_week, time_start, time_end, room,
	assigned_teacher_id, max_students, status, base_course_fee, registration_fee, exam_fee, certificate_fee, note)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
	RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		tenantID, c.CourseCode, c.CourseName, c.ClassName, c.Category, c.Subject, c.Level, c.StartDate, c.EndDate, c.ScheduleText,
		c.DaysOfWeek, c.TimeStart, c.TimeEnd, c.Room, c.AssignedTeacherID, c.MaxStudents, c.Status,
		c.BaseCourseFee, c.RegistrationFee, c.ExamFee, c.CertificateFee, c.Note,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *postgresRepository) List(ctx context.Context, filter common.ListFilter) ([]classcourse.ClassCourse, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	args := []any{tenantID}
	where := []string{"tenant_id = $1"}
	if filter.Query != "" {
		args = append(args, "%"+filter.Query+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(class_name ILIKE $%d OR course_name ILIKE $%d OR course_code ILIKE $%d)", idx, idx, idx))
	}
	if filter.ClassName != "" {
		args = append(args, "%"+filter.ClassName+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("class_name ILIKE $%d", idx))
	}
	if filter.CourseCategory != "" {
		args = append(args, filter.CourseCategory)
		idx := len(args)
		where = append(where, fmt.Sprintf("category=$%d", idx))
	}
	if filter.ClassStatus != "" {
		args = append(args, filter.ClassStatus)
		idx := len(args)
		where = append(where, fmt.Sprintf("status=$%d", idx))
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	args = append(args, filter.Limit, filter.Offset)
	lIdx := len(args) - 1
	oIdx := len(args)
	query := fmt.Sprintf(`SELECT id, course_code, course_name, class_name, category, subject, level, start_date, end_date, schedule_text, days_of_week, time_start,
	time_end, room, assigned_teacher_id, max_students, status, base_course_fee, registration_fee, exam_fee, certificate_fee, note, created_at, updated_at
	FROM class_courses WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), lIdx, oIdx)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]classcourse.ClassCourse, 0)
	for rows.Next() {
		var c classcourse.ClassCourse
		if err := rows.Scan(&c.ID, &c.CourseCode, &c.CourseName, &c.ClassName, &c.Category, &c.Subject, &c.Level,
			&c.StartDate, &c.EndDate, &c.ScheduleText, &c.DaysOfWeek, &c.TimeStart, &c.TimeEnd, &c.Room,
			&c.AssignedTeacherID, &c.MaxStudents, &c.Status, &c.BaseCourseFee, &c.RegistrationFee,
			&c.ExamFee, &c.CertificateFee, &c.Note, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*classcourse.ClassCourse, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, course_code, course_name, class_name, category, subject, level, start_date, end_date, schedule_text, days_of_week, time_start,
	time_end, room, assigned_teacher_id, max_students, status, base_course_fee, registration_fee, exam_fee, certificate_fee, note, created_at, updated_at
	FROM class_courses WHERE id=$1 AND tenant_id=$2`
	var c classcourse.ClassCourse
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&c.ID, &c.CourseCode, &c.CourseName, &c.ClassName, &c.Category, &c.Subject, &c.Level,
		&c.StartDate, &c.EndDate, &c.ScheduleText, &c.DaysOfWeek, &c.TimeStart, &c.TimeEnd, &c.Room,
		&c.AssignedTeacherID, &c.MaxStudents, &c.Status, &c.BaseCourseFee, &c.RegistrationFee,
		&c.ExamFee, &c.CertificateFee, &c.Note, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresRepository) Update(ctx context.Context, c *classcourse.ClassCourse) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE class_courses SET course_code=$2, course_name=$3, class_name=$4, category=$5, subject=$6, level=$7, start_date=$8, end_date=$9,
	schedule_text=$10, days_of_week=$11, time_start=$12, time_end=$13, room=$14, assigned_teacher_id=$15, max_students=$16, status=$17,
	base_course_fee=$18, registration_fee=$19, exam_fee=$20, certificate_fee=$21, note=$22, updated_at=NOW() WHERE id=$1 AND tenant_id=$23 RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query, c.ID, c.CourseCode, c.CourseName, c.ClassName, c.Category, c.Subject, c.Level,
		c.StartDate, c.EndDate, c.ScheduleText, c.DaysOfWeek, c.TimeStart, c.TimeEnd, c.Room, c.AssignedTeacherID,
		c.MaxStudents, c.Status, c.BaseCourseFee, c.RegistrationFee, c.ExamFee, c.CertificateFee, c.Note, tenantID).Scan(&c.UpdatedAt)
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `UPDATE class_courses SET status='closed', updated_at=NOW() WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

func (r *postgresRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return false, err
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM class_courses WHERE id=$1 AND tenant_id=$2)`, id, tenantID).Scan(&exists)
	return exists, err
}

func (r *postgresRepository) CreateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO optional_fee_items (tenant_id, class_course_id, item_name, default_amount, is_optional, is_active)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		tenantID, item.ClassCourseID, item.ItemName, item.DefaultAmount, item.IsOptional, item.IsActive,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *postgresRepository) ListOptionalFees(ctx context.Context, classCourseID int64) ([]classcourse.OptionalFeeItem, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, class_course_id, item_name, default_amount, is_optional, is_active, created_at, updated_at FROM optional_fee_items WHERE class_course_id=$1 AND tenant_id=$2 ORDER BY id DESC`, classCourseID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]classcourse.OptionalFeeItem, 0)
	for rows.Next() {
		var item classcourse.OptionalFeeItem
		if err := rows.Scan(&item.ID, &item.ClassCourseID, &item.ItemName, &item.DefaultAmount, &item.IsOptional, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *postgresRepository) GetOptionalFeeByID(ctx context.Context, id int64) (*classcourse.OptionalFeeItem, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, class_course_id, item_name, default_amount, is_optional, is_active, created_at, updated_at FROM optional_fee_items WHERE id=$1 AND tenant_id=$2`
	var item classcourse.OptionalFeeItem
	if err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(&item.ID, &item.ClassCourseID, &item.ItemName, &item.DefaultAmount, &item.IsOptional, &item.IsActive, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *postgresRepository) UpdateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE optional_fee_items SET item_name=$2, default_amount=$3, is_optional=$4, is_active=$5, updated_at=NOW() WHERE id=$1 AND tenant_id=$6 RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query, item.ID, item.ItemName, item.DefaultAmount, item.IsOptional, item.IsActive, tenantID).Scan(&item.UpdatedAt)
}

func (r *postgresRepository) DeleteOptionalFee(ctx context.Context, id int64) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM optional_fee_items WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}
