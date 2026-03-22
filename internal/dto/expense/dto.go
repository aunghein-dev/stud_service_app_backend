package expense

type CreateRequest struct {
	ExpenseDate   string  `json:"expense_date"`
	ExpenseType   string  `json:"expense_type" validate:"required,oneof=teacher_fee books uniform shoes stationery rent utilities marketing misc"`
	TeacherID     *int64  `json:"teacher_id"`
	ClassCourseID *int64  `json:"class_course_id"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Description   string  `json:"description" validate:"max=500"`
	PaymentMethod string  `json:"payment_method" validate:"max=50"`
	ReferenceNo   string  `json:"reference_no" validate:"max=100"`
}

type UpdateRequest = CreateRequest

type Response struct {
	ID            int64   `json:"id"`
	ExpenseDate   string  `json:"expense_date"`
	ExpenseType   string  `json:"expense_type"`
	TeacherID     *int64  `json:"teacher_id,omitempty"`
	ClassCourseID *int64  `json:"class_course_id,omitempty"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	PaymentMethod string  `json:"payment_method"`
	ReferenceNo   string  `json:"reference_no"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
