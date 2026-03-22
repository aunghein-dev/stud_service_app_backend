package payment

type CreateRequest struct {
	EnrollmentID  int64   `json:"enrollment_id" validate:"required,gt=0"`
	PaymentDate   string  `json:"payment_date"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash bank_transfer mobile_wallet other"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Note          string  `json:"note" validate:"max=500"`
	ReceivedBy    string  `json:"received_by" validate:"max=100"`
}

type UpdateRequest struct {
	PaymentDate   string  `json:"payment_date"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash bank_transfer mobile_wallet other"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Note          string  `json:"note" validate:"max=500"`
	ReceivedBy    string  `json:"received_by" validate:"max=100"`
}

type Response struct {
	ID            int64   `json:"id"`
	ReceiptNo     string  `json:"receipt_no"`
	StudentID     int64   `json:"student_id"`
	EnrollmentID  int64   `json:"enrollment_id"`
	ClassCourseID int64   `json:"class_course_id"`
	PaymentDate   string  `json:"payment_date"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Note          string  `json:"note"`
	ReceivedBy    string  `json:"received_by"`
	CreatedAt     string  `json:"created_at"`
}
