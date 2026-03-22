package enrollment

import "time"

type Enrollment struct {
	ID              int64     `db:"id"`
	EnrollmentCode  string    `db:"enrollment_code"`
	StudentID       int64     `db:"student_id"`
	StudentName     string    `db:"student_name"`
	GuardianName    string    `db:"guardian_name"`
	ClassCourseID   int64     `db:"class_course_id"`
	ClassName       string    `db:"class_name"`
	CourseName      string    `db:"course_name"`
	EnrollmentDate  time.Time `db:"enrollment_date"`
	SubTotal        float64   `db:"sub_total"`
	DiscountAmount  float64   `db:"discount_amount"`
	FinalFee        float64   `db:"final_fee"`
	PaidAmount      float64   `db:"paid_amount"`
	RemainingAmount float64   `db:"remaining_amount"`
	PaymentStatus   string    `db:"payment_status"`
	Note            string    `db:"note"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type EnrollmentOptionalItem struct {
	ID                int64     `db:"id"`
	EnrollmentID      int64     `db:"enrollment_id"`
	OptionalFeeItemID *int64    `db:"optional_fee_item_id"`
	ItemNameSnapshot  string    `db:"item_name_snapshot"`
	AmountSnapshot    float64   `db:"amount_snapshot"`
	Quantity          int       `db:"quantity"`
	TotalAmount       float64   `db:"total_amount"`
	CreatedAt         time.Time `db:"created_at"`
}
