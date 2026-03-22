//go:build wireinject
// +build wireinject

package di

import (
	"student_service_app/backend/internal/config"
	"student_service_app/backend/internal/db"
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
	"student_service_app/backend/internal/logger"
	authrepo "student_service_app/backend/internal/repository/auth"
	classcourserepo "student_service_app/backend/internal/repository/classcourse"
	enrollmentrepo "student_service_app/backend/internal/repository/enrollment"
	expenserepo "student_service_app/backend/internal/repository/expense"
	paymentrepo "student_service_app/backend/internal/repository/payment"
	receiptrepo "student_service_app/backend/internal/repository/receipt"
	reportrepo "student_service_app/backend/internal/repository/report"
	settingsrepo "student_service_app/backend/internal/repository/settings"
	studentrepo "student_service_app/backend/internal/repository/student"
	teacherrepo "student_service_app/backend/internal/repository/teacher"
	"student_service_app/backend/internal/server"
	authsvc "student_service_app/backend/internal/service/auth"
	classcoursesvc "student_service_app/backend/internal/service/classcourse"
	enrollmentsvc "student_service_app/backend/internal/service/enrollment"
	expensesvc "student_service_app/backend/internal/service/expense"
	paymentsvc "student_service_app/backend/internal/service/payment"
	receiptsvc "student_service_app/backend/internal/service/receipt"
	reportsvc "student_service_app/backend/internal/service/report"
	settingssvc "student_service_app/backend/internal/service/settings"
	studentsvc "student_service_app/backend/internal/service/student"
	teachersvc "student_service_app/backend/internal/service/teacher"

	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		config.Load,
		logger.New,
		db.NewPostgres,
		NewValidator,
		NewBase,

		authrepo.NewPostgresRepository,
		studentrepo.NewPostgresRepository,
		teacherrepo.NewPostgresRepository,
		classcourserepo.NewPostgresRepository,
		enrollmentrepo.NewPostgresRepository,
		paymentrepo.NewPostgresRepository,
		expenserepo.NewPostgresRepository,
		receiptrepo.NewPostgresRepository,
		settingsrepo.NewPostgresRepository,
		reportrepo.NewPostgresRepository,

		authsvc.NewTokenManager,
		authsvc.NewService,
		studentsvc.NewService,
		teachersvc.NewService,
		classcoursesvc.NewService,
		enrollmentsvc.NewService,
		paymentsvc.NewService,
		expensesvc.NewService,
		receiptsvc.NewService,
		reportsvc.NewService,
		settingssvc.NewService,

		authhandler.NewHandler,
		studenthandler.NewHandler,
		teacherhandler.NewHandler,
		classcoursehandler.NewHandler,
		enrollmenthandler.NewHandler,
		paymenthandler.NewHandler,
		expensehandler.NewHandler,
		receipthandler.NewHandler,
		reporthandler.NewHandler,
		settingshandler.NewHandler,

		NewHandlers,
		server.NewRouter,
		server.NewHTTPServer,
		wire.Struct(new(App), "Config", "Logger", "Server"),
	)
	return nil, nil
}
