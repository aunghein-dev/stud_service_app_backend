package enrollment

type OptionalItemInput struct {
	OptionalFeeItemID *int64  `json:"optional_fee_item_id"`
	ItemName          string  `json:"item_name" validate:"required,max=100"`
	Amount            float64 `json:"amount" validate:"gte=0"`
	Quantity          int     `json:"quantity" validate:"gte=1"`
}

type CreateRequest struct {
	StudentID      int64               `json:"student_id" validate:"required,gt=0"`
	ClassCourseID  int64               `json:"class_course_id" validate:"required,gt=0"`
	EnrollmentDate string              `json:"enrollment_date"`
	DiscountAmount float64             `json:"discount_amount" validate:"gte=0"`
	OptionalItems  []OptionalItemInput `json:"optional_items"`
	InitialPayment float64             `json:"initial_payment" validate:"gte=0"`
	PaymentMethod  string              `json:"payment_method" validate:"omitempty,oneof=cash bank_transfer mobile_wallet other"`
	ReceivedBy     string              `json:"received_by" validate:"max=100"`
	Note           string              `json:"note" validate:"max=500"`
	AllowDuplicate bool                `json:"allow_duplicate"`
}

type UpdateRequest struct {
	DiscountAmount float64 `json:"discount_amount" validate:"gte=0"`
	Note           string  `json:"note" validate:"max=500"`
}

type OptionalItemResponse struct {
	ID                int64   `json:"id"`
	OptionalFeeItemID *int64  `json:"optional_fee_item_id,omitempty"`
	ItemNameSnapshot  string  `json:"item_name_snapshot"`
	AmountSnapshot    float64 `json:"amount_snapshot"`
	Quantity          int     `json:"quantity"`
	TotalAmount       float64 `json:"total_amount"`
}

type Response struct {
	ID              int64                  `json:"id"`
	EnrollmentCode  string                 `json:"enrollment_code"`
	StudentID       int64                  `json:"student_id"`
	StudentName     string                 `json:"student_name,omitempty"`
	GuardianName    string                 `json:"guardian_name,omitempty"`
	ClassCourseID   int64                  `json:"class_course_id"`
	ClassName       string                 `json:"class_name,omitempty"`
	CourseName      string                 `json:"course_name,omitempty"`
	EnrollmentDate  string                 `json:"enrollment_date"`
	SubTotal        float64                `json:"sub_total"`
	DiscountAmount  float64                `json:"discount_amount"`
	FinalFee        float64                `json:"final_fee"`
	PaidAmount      float64                `json:"paid_amount"`
	RemainingAmount float64                `json:"remaining_amount"`
	PaymentStatus   string                 `json:"payment_status"`
	Note            string                 `json:"note"`
	OptionalItems   []OptionalItemResponse `json:"optional_items,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}
