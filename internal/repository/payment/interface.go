package payment

import (
	"context"
	"database/sql"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/payment"
)

type Repository interface {
	Create(ctx context.Context, tx *sql.Tx, p *payment.Payment) error
	Update(ctx context.Context, tx *sql.Tx, p *payment.Payment) error
	Delete(ctx context.Context, tx *sql.Tx, id int64) error
	List(ctx context.Context, filter common.ListFilter) ([]payment.Payment, error)
	GetByID(ctx context.Context, id int64) (*payment.Payment, error)
	ListByEnrollment(ctx context.Context, enrollmentID int64) ([]payment.Payment, error)
}
