package student

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"student_service_app/backend/internal/domain/common"
	domain "student_service_app/backend/internal/domain/student"
	basehandler "student_service_app/backend/internal/handler"
	servicepkg "student_service_app/backend/internal/service/student"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type fakeService struct{}

func (f *fakeService) Create(ctx context.Context, s *domain.Student) error { return nil }
func (f *fakeService) List(ctx context.Context, filter common.ListFilter) ([]domain.Student, error) {
	return []domain.Student{{ID: 1, StudentCode: "S-001", FullName: "Alice", Phone: "091", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}
func (f *fakeService) GetByID(ctx context.Context, id int64) (*domain.Student, error) {
	return nil, nil
}
func (f *fakeService) Update(ctx context.Context, s *domain.Student) error { return nil }
func (f *fakeService) Delete(ctx context.Context, id int64) error          { return nil }

var _ servicepkg.Service = (*fakeService)(nil)

func TestListStudents(t *testing.T) {
	h := NewHandler(basehandler.NewBase(validator.New()), &fakeService{})
	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var payload map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	require.NoError(t, err)
	require.Equal(t, true, payload["success"])
}
