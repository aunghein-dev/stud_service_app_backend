package common

type ListQuery struct {
	Q               string `json:"q"`
	DateFrom        string `json:"date_from"`
	DateTo          string `json:"date_to"`
	TeacherName     string `json:"teacher_name"`
	StudentName     string `json:"student_name"`
	ClassCourseName string `json:"class_course_name"`
	CourseCategory  string `json:"course_category"`
	PaymentStatus   string `json:"payment_status"`
	ReceiptNo       string `json:"receipt_no"`
	ExpenseType     string `json:"expense_type"`
	TransactionType string `json:"transaction_type"`
	Limit           int    `json:"limit"`
	Offset          int    `json:"offset"`
}
