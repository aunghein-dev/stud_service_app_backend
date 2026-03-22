package receipt

type Response struct {
	ID              int64          `json:"id"`
	ReceiptNo       string         `json:"receipt_no"`
	ReceiptType     string         `json:"receipt_type"`
	StudentID       int64          `json:"student_id"`
	EnrollmentID    int64          `json:"enrollment_id"`
	PaymentID       *int64         `json:"payment_id,omitempty"`
	ClassCourseID   int64          `json:"class_course_id"`
	TotalAmount     float64        `json:"total_amount"`
	PaidAmount      float64        `json:"paid_amount"`
	RemainingAmount float64        `json:"remaining_amount"`
	Payload         map[string]any `json:"payload"`
	IssuedAt        string         `json:"issued_at"`
	CreatedAt       string         `json:"created_at"`
}
