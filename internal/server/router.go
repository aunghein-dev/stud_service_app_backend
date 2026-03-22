package server

import (
	"net/http"

	"student_service_app/backend/internal/apidocs"
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
	"student_service_app/backend/internal/middleware"
	authservice "student_service_app/backend/internal/service/auth"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth        *authhandler.Handler
	Students    *studenthandler.Handler
	Teachers    *teacherhandler.Handler
	ClassCourse *classcoursehandler.Handler
	Enrollments *enrollmenthandler.Handler
	Payments    *paymenthandler.Handler
	Expenses    *expensehandler.Handler
	Receipts    *receipthandler.Handler
	Reports     *reporthandler.Handler
	Settings    *settingshandler.Handler
}

func NewRouter(log *zap.Logger, tokenManager *authservice.TokenManager, handlers *Handlers) http.Handler {
	r := chi.NewRouter()
	docs := apidocs.NewService()
	r.Use(middleware.RequestID())
	r.Use(middleware.Recoverer())
	r.Use(middleware.Logger(log))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/docs", docs.ServeScalar)
	r.Get("/docs/", docs.ServeScalar)
	r.Get("/docs/openapi.json", docs.ServeOpenAPI)

	r.Route("/api/v1", func(api chi.Router) {
		api.Route("/auth", func(ar chi.Router) {
			ar.Post("/signup", handlers.Auth.SignUp)
			ar.Post("/login", handlers.Auth.Login)
			ar.With(middleware.RequireAuth(tokenManager)).Get("/me", handlers.Auth.Me)
		})

		api.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(tokenManager))

			protected.Route("/students", func(sr chi.Router) {
				sr.Post("/", handlers.Students.Create)
				sr.Get("/", handlers.Students.List)
				sr.Get("/{id}", handlers.Students.GetByID)
				sr.Put("/{id}", handlers.Students.Update)
				sr.Delete("/{id}", handlers.Students.Delete)
				sr.Get("/{id}/enrollments", handlers.Enrollments.ListByStudent)
			})

			protected.Route("/teachers", func(tr chi.Router) {
				tr.Post("/", handlers.Teachers.Create)
				tr.Get("/", handlers.Teachers.List)
				tr.Get("/{id}", handlers.Teachers.GetByID)
				tr.Put("/{id}", handlers.Teachers.Update)
				tr.Delete("/{id}", handlers.Teachers.Delete)
			})

			protected.Route("/class-courses", func(cr chi.Router) {
				cr.Post("/", handlers.ClassCourse.Create)
				cr.Get("/", handlers.ClassCourse.List)
				cr.Get("/{id}", handlers.ClassCourse.GetByID)
				cr.Put("/{id}", handlers.ClassCourse.Update)
				cr.Delete("/{id}", handlers.ClassCourse.Delete)
				cr.Post("/{id}/optional-fees", handlers.ClassCourse.CreateOptionalFee)
				cr.Get("/{id}/optional-fees", handlers.ClassCourse.ListOptionalFees)
			})

			protected.Route("/optional-fees", func(ofr chi.Router) {
				ofr.Put("/{id}", handlers.ClassCourse.UpdateOptionalFee)
				ofr.Delete("/{id}", handlers.ClassCourse.DeleteOptionalFee)
			})

			protected.Route("/enrollments", func(er chi.Router) {
				er.Post("/", handlers.Enrollments.Create)
				er.Get("/", handlers.Enrollments.List)
				er.Get("/{id}", handlers.Enrollments.GetByID)
				er.Put("/{id}", handlers.Enrollments.Update)
				er.Delete("/{id}", handlers.Enrollments.Delete)
				er.Get("/{id}/payments", handlers.Payments.ListByEnrollment)
			})

			protected.Route("/payments", func(pr chi.Router) {
				pr.Post("/", handlers.Payments.Create)
				pr.Get("/", handlers.Payments.List)
				pr.Get("/{id}", handlers.Payments.GetByID)
				pr.Put("/{id}", handlers.Payments.Update)
				pr.Delete("/{id}", handlers.Payments.Delete)
			})

			protected.Route("/expenses", func(er chi.Router) {
				er.Post("/", handlers.Expenses.Create)
				er.Get("/", handlers.Expenses.List)
				er.Get("/{id}", handlers.Expenses.GetByID)
				er.Put("/{id}", handlers.Expenses.Update)
				er.Delete("/{id}", handlers.Expenses.Delete)
			})

			protected.Route("/receipts", func(rr chi.Router) {
				rr.Get("/", handlers.Receipts.List)
				rr.Get("/{key}", handlers.Receipts.GetByKey)
			})

			protected.Route("/reports", func(rr chi.Router) {
				rr.Get("/dashboard", handlers.Reports.Dashboard)
				rr.Get("/students", handlers.Reports.Students)
				rr.Get("/teachers", handlers.Reports.Teachers)
				rr.Get("/class-courses", handlers.Reports.ClassCourses)
				rr.Get("/gross", handlers.Reports.Gross)
				rr.Get("/transactions", handlers.Reports.Transactions)
				rr.Get("/performance", handlers.Reports.Performance)
			})

			protected.Route("/settings", func(sr chi.Router) {
				sr.Get("/", handlers.Settings.Get)
				sr.Put("/", handlers.Settings.Update)
			})
		})
	})

	return r
}
