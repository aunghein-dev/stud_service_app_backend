package di

import (
	"student_service_app/backend/internal/handler"
	authhandler "student_service_app/backend/internal/handler/auth"
	classcoursehandler "student_service_app/backend/internal/handler/classcourse"
	enrollmenthandler "student_service_app/backend/internal/handler/enrollment"
	expensehandler "student_service_app/backend/internal/handler/expense"
	paymenthandler "student_service_app/backend/internal/handler/payment"
	receipthandler "student_service_app/backend/internal/handler/receipt"
	reporthandler "student_service_app/backend/internal/handler/report"
	settingshandler "student_service_app/backend/internal/handler/settings"
	studenthandler "student_service_app/backend/internal/handler/student"
	teacherhandler "student_service_app/backend/internal/handler/teacher"
	"student_service_app/backend/internal/server"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func NewBase(v *validator.Validate) handler.Base {
	return handler.NewBase(v)
}

func NewHandlers(
	auth *authhandler.Handler,
	students *studenthandler.Handler,
	teachers *teacherhandler.Handler,
	classCourses *classcoursehandler.Handler,
	enrollments *enrollmenthandler.Handler,
	payments *paymenthandler.Handler,
	expenses *expensehandler.Handler,
	receipts *receipthandler.Handler,
	reports *reporthandler.Handler,
	settings *settingshandler.Handler,
) *server.Handlers {
	return &server.Handlers{
		Auth:        auth,
		Students:    students,
		Teachers:    teachers,
		ClassCourse: classCourses,
		Enrollments: enrollments,
		Payments:    payments,
		Expenses:    expenses,
		Receipts:    receipts,
		Reports:     reports,
		Settings:    settings,
	}
}
