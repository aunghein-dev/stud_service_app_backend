package classcourse

import "time"

type ClassCourse struct {
	ID                int64      `db:"id"`
	CourseCode        string     `db:"course_code"`
	CourseName        string     `db:"course_name"`
	ClassName         string     `db:"class_name"`
	Category          string     `db:"category"`
	Subject           string     `db:"subject"`
	Level             string     `db:"level"`
	StartDate         *time.Time `db:"start_date"`
	EndDate           *time.Time `db:"end_date"`
	ScheduleText      string     `db:"schedule_text"`
	DaysOfWeek        string     `db:"days_of_week"`
	TimeStart         string     `db:"time_start"`
	TimeEnd           string     `db:"time_end"`
	Room              string     `db:"room"`
	AssignedTeacherID *int64     `db:"assigned_teacher_id"`
	MaxStudents       int        `db:"max_students"`
	Status            string     `db:"status"`
	BaseCourseFee     float64    `db:"base_course_fee"`
	RegistrationFee   float64    `db:"registration_fee"`
	ExamFee           float64    `db:"exam_fee"`
	CertificateFee    float64    `db:"certificate_fee"`
	Note              string     `db:"note"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

type OptionalFeeItem struct {
	ID            int64     `db:"id"`
	ClassCourseID int64     `db:"class_course_id"`
	ItemName      string    `db:"item_name"`
	DefaultAmount float64   `db:"default_amount"`
	IsOptional    bool      `db:"is_optional"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
