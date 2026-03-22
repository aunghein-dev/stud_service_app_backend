package enrollment

import (
	"context"
	"database/sql"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/enrollment"
)

type Repository interface {
	ExistsDuplicate(ctx context.Context, studentID, classCourseID int64) (bool, error)
	Create(ctx context.Context, tx *sql.Tx, e *enrollment.Enrollment) error
	AddOptionalItems(ctx context.Context, tx *sql.Tx, items []enrollment.EnrollmentOptionalItem) error
	List(ctx context.Context, filter common.ListFilter) ([]enrollment.Enrollment, error)
	ListByStudent(ctx context.Context, studentID int64) ([]enrollment.Enrollment, error)
	GetByID(ctx context.Context, id int64) (*enrollment.Enrollment, error)
	ListOptionalItems(ctx context.Context, enrollmentID int64) ([]enrollment.EnrollmentOptionalItem, error)
	Update(ctx context.Context, e *enrollment.Enrollment) error
	UpdatePaymentState(ctx context.Context, tx *sql.Tx, id int64, paidAmount, remainingAmount float64, paymentStatus string) error
}
