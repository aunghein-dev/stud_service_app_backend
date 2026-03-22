package receipt

import "time"

type Receipt struct {
	ID              int64     `db:"id"`
	ReceiptNo       string    `db:"receipt_no"`
	ReceiptType     string    `db:"receipt_type"`
	StudentID       int64     `db:"student_id"`
	EnrollmentID    int64     `db:"enrollment_id"`
	PaymentID       *int64    `db:"payment_id"`
	ClassCourseID   int64     `db:"class_course_id"`
	TotalAmount     float64   `db:"total_amount"`
	PaidAmount      float64   `db:"paid_amount"`
	RemainingAmount float64   `db:"remaining_amount"`
	PayloadJSON     []byte    `db:"payload_json"`
	IssuedAt        time.Time `db:"issued_at"`
	CreatedAt       time.Time `db:"created_at"`
}
