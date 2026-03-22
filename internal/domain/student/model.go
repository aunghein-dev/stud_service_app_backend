package student

import "time"

type Student struct {
	ID            int64      `db:"id"`
	StudentCode   string     `db:"student_code"`
	FullName      string     `db:"full_name"`
	Gender        string     `db:"gender"`
	DateOfBirth   *time.Time `db:"date_of_birth"`
	Phone         string     `db:"phone"`
	GuardianName  string     `db:"guardian_name"`
	GuardianPhone string     `db:"guardian_phone"`
	Address       string     `db:"address"`
	SchoolName    string     `db:"school_name"`
	GradeLevel    string     `db:"grade_level"`
	Note          string     `db:"note"`
	IsActive      bool       `db:"is_active"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}
