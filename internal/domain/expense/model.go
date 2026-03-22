package expense

import "time"

type Expense struct {
	ID            int64     `db:"id"`
	ExpenseDate   time.Time `db:"expense_date"`
	ExpenseType   string    `db:"expense_type"`
	TeacherID     *int64    `db:"teacher_id"`
	ClassCourseID *int64    `db:"class_course_id"`
	Amount        float64   `db:"amount"`
	Description   string    `db:"description"`
	PaymentMethod string    `db:"payment_method"`
	ReferenceNo   string    `db:"reference_no"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
