package expense

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/expense"
)

type Repository interface {
	Create(ctx context.Context, e *expense.Expense) error
	List(ctx context.Context, filter common.ListFilter) ([]expense.Expense, error)
	GetByID(ctx context.Context, id int64) (*expense.Expense, error)
	Update(ctx context.Context, e *expense.Expense) error
	Delete(ctx context.Context, id int64) error
}
