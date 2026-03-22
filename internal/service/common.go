package service

import (
	"net/http"
	"strconv"

	"student_service_app/backend/internal/domain/common"
)

func BuildFilter(r *http.Request) common.ListFilter {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	return common.ListFilter{
		Query:           q.Get("q"),
		DateFrom:        q.Get("date_from"),
		DateTo:          q.Get("date_to"),
		TeacherName:     q.Get("teacher_name"),
		StudentName:     q.Get("student_name"),
		ClassName:       q.Get("class_course_name"),
		ClassStatus:     q.Get("class_status"),
		CourseCategory:  q.Get("course_category"),
		PaymentStatus:   q.Get("payment_status"),
		ReceiptNo:       q.Get("receipt_no"),
		ExpenseType:     q.Get("expense_type"),
		TransactionType: q.Get("transaction_type"),
		Limit:           limit,
		Offset:          offset,
	}
}
