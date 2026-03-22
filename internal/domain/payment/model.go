package payment

import "time"

type Payment struct {
	ID            int64     `db:"id"`
	ReceiptNo     string    `db:"receipt_no"`
	StudentID     int64     `db:"student_id"`
	EnrollmentID  int64     `db:"enrollment_id"`
	ClassCourseID int64     `db:"class_course_id"`
	PaymentDate   time.Time `db:"payment_date"`
	PaymentMethod string    `db:"payment_method"`
	Amount        float64   `db:"amount"`
	Note          string    `db:"note"`
	ReceivedBy    string    `db:"received_by"`
	CreatedAt     time.Time `db:"created_at"`
}
