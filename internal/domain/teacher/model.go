package teacher

import "time"

type Teacher struct {
	ID               int64     `db:"id"`
	TeacherCode      string    `db:"teacher_code"`
	TeacherName      string    `db:"teacher_name"`
	Phone            string    `db:"phone"`
	Address          string    `db:"address"`
	SubjectSpecialty string    `db:"subject_specialty"`
	SalaryType       string    `db:"salary_type"`
	DefaultFeeAmount float64   `db:"default_fee_amount"`
	Note             string    `db:"note"`
	IsActive         bool      `db:"is_active"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
