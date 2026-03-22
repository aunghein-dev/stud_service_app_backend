package common

type ListFilter struct {
	Query           string `json:"q,omitempty"`
	DateFrom        string `json:"date_from,omitempty"`
	DateTo          string `json:"date_to,omitempty"`
	TeacherName     string `json:"teacher_name,omitempty"`
	StudentName     string `json:"student_name,omitempty"`
	ClassName       string `json:"class_course_name,omitempty"`
	ClassStatus     string `json:"class_status,omitempty"`
	CourseCategory  string `json:"course_category,omitempty"`
	PaymentStatus   string `json:"payment_status,omitempty"`
	ReceiptNo       string `json:"receipt_no,omitempty"`
	ExpenseType     string `json:"expense_type,omitempty"`
	TransactionType string `json:"transaction_type,omitempty"`
	Limit           int    `json:"limit,omitempty"`
	Offset          int    `json:"offset,omitempty"`
}
