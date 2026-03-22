package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authdomain "student_service_app/backend/internal/domain/auth"
	basehandler "student_service_app/backend/internal/handler"
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
	authservice "student_service_app/backend/internal/service/auth"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fakeAuthService struct {
	session authdomain.Session
}

func (f *fakeAuthService) SignUp(ctx context.Context, input authservice.SignUpInput) (*authdomain.Session, error) {
	return &f.session, nil
}

func (f *fakeAuthService) Login(ctx context.Context, input authservice.LoginInput) (*authdomain.Session, error) {
	return &f.session, nil
}

func (f *fakeAuthService) Me(ctx context.Context) (*authdomain.Session, error) {
	return &f.session, nil
}

func TestRouterMountedEndpoints(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	base := basehandler.NewBase(validate)
	tokenManager := authservice.NewTokenManagerWithValues("test-secret", 24*time.Hour)

	session := authdomain.Session{
		User: authdomain.SessionUser{
			ID:       11,
			FullName: "School Owner",
			Email:    "owner@example.com",
			Role:     "owner",
		},
		Tenant: authdomain.SessionTenant{
			ID:            22,
			Slug:          "bright-future",
			SchoolName:    "Bright Future Academy",
			SchoolAddress: "Main Road",
			SchoolPhone:   "099999999",
		},
	}
	accessToken, err := tokenManager.Sign(session)
	require.NoError(t, err)
	session.AccessToken = accessToken

	router := NewRouter(zap.NewNop(), tokenManager, &Handlers{
		Auth:        authhandler.NewHandler(base, &fakeAuthService{session: session}),
		Students:    studenthandler.NewHandler(base, nil),
		Teachers:    teacherhandler.NewHandler(base, nil),
		ClassCourse: classcoursehandler.NewHandler(base, nil),
		Enrollments: enrollmenthandler.NewHandler(base, nil),
		Payments:    paymenthandler.NewHandler(base, nil),
		Expenses:    expensehandler.NewHandler(base, nil),
		Receipts:    receipthandler.NewHandler(nil),
		Reports:     reporthandler.NewHandler(nil),
		Settings:    settingshandler.NewHandler(base, nil),
	})

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		token      string
		wantStatus int
	}{
		{name: "health", method: http.MethodGet, path: "/healthz", wantStatus: http.StatusOK},
		{name: "scalar docs", method: http.MethodGet, path: "/docs", wantStatus: http.StatusOK},
		{name: "openapi docs", method: http.MethodGet, path: "/docs/openapi.json", wantStatus: http.StatusOK},
		{
			name:       "auth signup",
			method:     http.MethodPost,
			path:       "/api/v1/auth/signup",
			body:       `{"school_name":"Bright Future Academy","tenant_slug":"bright-future","admin_name":"Owner","email":"owner@example.com","password":"password123","school_phone":"099999999","school_address":"Main Road"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "auth login",
			method:     http.MethodPost,
			path:       "/api/v1/auth/login",
			body:       `{"tenant_slug":"bright-future","email":"owner@example.com","password":"password123"}`,
			wantStatus: http.StatusOK,
		},
		{name: "auth me", method: http.MethodGet, path: "/api/v1/auth/me", token: accessToken, wantStatus: http.StatusOK},
		{name: "students create unauthorized", method: http.MethodPost, path: "/api/v1/students", wantStatus: http.StatusUnauthorized},
		{name: "students list unauthorized", method: http.MethodGet, path: "/api/v1/students", wantStatus: http.StatusUnauthorized},
		{name: "students get unauthorized", method: http.MethodGet, path: "/api/v1/students/1", wantStatus: http.StatusUnauthorized},
		{name: "students update unauthorized", method: http.MethodPut, path: "/api/v1/students/1", wantStatus: http.StatusUnauthorized},
		{name: "students delete unauthorized", method: http.MethodDelete, path: "/api/v1/students/1", wantStatus: http.StatusUnauthorized},
		{name: "students enrollments unauthorized", method: http.MethodGet, path: "/api/v1/students/1/enrollments", wantStatus: http.StatusUnauthorized},
		{name: "teachers create unauthorized", method: http.MethodPost, path: "/api/v1/teachers", wantStatus: http.StatusUnauthorized},
		{name: "teachers list unauthorized", method: http.MethodGet, path: "/api/v1/teachers", wantStatus: http.StatusUnauthorized},
		{name: "teachers get unauthorized", method: http.MethodGet, path: "/api/v1/teachers/1", wantStatus: http.StatusUnauthorized},
		{name: "teachers update unauthorized", method: http.MethodPut, path: "/api/v1/teachers/1", wantStatus: http.StatusUnauthorized},
		{name: "teachers delete unauthorized", method: http.MethodDelete, path: "/api/v1/teachers/1", wantStatus: http.StatusUnauthorized},
		{name: "classes create unauthorized", method: http.MethodPost, path: "/api/v1/class-courses", wantStatus: http.StatusUnauthorized},
		{name: "classes list unauthorized", method: http.MethodGet, path: "/api/v1/class-courses", wantStatus: http.StatusUnauthorized},
		{name: "classes get unauthorized", method: http.MethodGet, path: "/api/v1/class-courses/1", wantStatus: http.StatusUnauthorized},
		{name: "classes update unauthorized", method: http.MethodPut, path: "/api/v1/class-courses/1", wantStatus: http.StatusUnauthorized},
		{name: "classes delete unauthorized", method: http.MethodDelete, path: "/api/v1/class-courses/1", wantStatus: http.StatusUnauthorized},
		{name: "optional fee create unauthorized", method: http.MethodPost, path: "/api/v1/class-courses/1/optional-fees", wantStatus: http.StatusUnauthorized},
		{name: "optional fee list unauthorized", method: http.MethodGet, path: "/api/v1/class-courses/1/optional-fees", wantStatus: http.StatusUnauthorized},
		{name: "optional fee update unauthorized", method: http.MethodPut, path: "/api/v1/optional-fees/1", wantStatus: http.StatusUnauthorized},
		{name: "optional fee delete unauthorized", method: http.MethodDelete, path: "/api/v1/optional-fees/1", wantStatus: http.StatusUnauthorized},
		{name: "enrollments create unauthorized", method: http.MethodPost, path: "/api/v1/enrollments", wantStatus: http.StatusUnauthorized},
		{name: "enrollments list unauthorized", method: http.MethodGet, path: "/api/v1/enrollments", wantStatus: http.StatusUnauthorized},
		{name: "enrollments get unauthorized", method: http.MethodGet, path: "/api/v1/enrollments/1", wantStatus: http.StatusUnauthorized},
		{name: "enrollments update unauthorized", method: http.MethodPut, path: "/api/v1/enrollments/1", wantStatus: http.StatusUnauthorized},
		{name: "enrollments delete unauthorized", method: http.MethodDelete, path: "/api/v1/enrollments/1", wantStatus: http.StatusUnauthorized},
		{name: "enrollment payments unauthorized", method: http.MethodGet, path: "/api/v1/enrollments/1/payments", wantStatus: http.StatusUnauthorized},
		{name: "payments create unauthorized", method: http.MethodPost, path: "/api/v1/payments", wantStatus: http.StatusUnauthorized},
		{name: "payments list unauthorized", method: http.MethodGet, path: "/api/v1/payments", wantStatus: http.StatusUnauthorized},
		{name: "payments get unauthorized", method: http.MethodGet, path: "/api/v1/payments/1", wantStatus: http.StatusUnauthorized},
		{name: "payments update unauthorized", method: http.MethodPut, path: "/api/v1/payments/1", wantStatus: http.StatusUnauthorized},
		{name: "payments delete unauthorized", method: http.MethodDelete, path: "/api/v1/payments/1", wantStatus: http.StatusUnauthorized},
		{name: "expenses create unauthorized", method: http.MethodPost, path: "/api/v1/expenses", wantStatus: http.StatusUnauthorized},
		{name: "expenses list unauthorized", method: http.MethodGet, path: "/api/v1/expenses", wantStatus: http.StatusUnauthorized},
		{name: "expenses get unauthorized", method: http.MethodGet, path: "/api/v1/expenses/1", wantStatus: http.StatusUnauthorized},
		{name: "expenses update unauthorized", method: http.MethodPut, path: "/api/v1/expenses/1", wantStatus: http.StatusUnauthorized},
		{name: "expenses delete unauthorized", method: http.MethodDelete, path: "/api/v1/expenses/1", wantStatus: http.StatusUnauthorized},
		{name: "receipts list unauthorized", method: http.MethodGet, path: "/api/v1/receipts", wantStatus: http.StatusUnauthorized},
		{name: "receipts get unauthorized", method: http.MethodGet, path: "/api/v1/receipts/ABC-001", wantStatus: http.StatusUnauthorized},
		{name: "reports dashboard unauthorized", method: http.MethodGet, path: "/api/v1/reports/dashboard", wantStatus: http.StatusUnauthorized},
		{name: "reports students unauthorized", method: http.MethodGet, path: "/api/v1/reports/students", wantStatus: http.StatusUnauthorized},
		{name: "reports teachers unauthorized", method: http.MethodGet, path: "/api/v1/reports/teachers", wantStatus: http.StatusUnauthorized},
		{name: "reports classes unauthorized", method: http.MethodGet, path: "/api/v1/reports/class-courses", wantStatus: http.StatusUnauthorized},
		{name: "reports gross unauthorized", method: http.MethodGet, path: "/api/v1/reports/gross", wantStatus: http.StatusUnauthorized},
		{name: "reports transactions unauthorized", method: http.MethodGet, path: "/api/v1/reports/transactions", wantStatus: http.StatusUnauthorized},
		{name: "reports performance unauthorized", method: http.MethodGet, path: "/api/v1/reports/performance", wantStatus: http.StatusUnauthorized},
		{name: "settings get unauthorized", method: http.MethodGet, path: "/api/v1/settings", wantStatus: http.StatusUnauthorized},
		{name: "settings update unauthorized", method: http.MethodPut, path: "/api/v1/settings", wantStatus: http.StatusUnauthorized},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			require.Equal(t, tc.wantStatus, rr.Code)
		})
	}
}
