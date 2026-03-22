package receipt

import (
	"context"
	"database/sql"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/receipt"
)

type Repository interface {
	Create(ctx context.Context, tx *sql.Tx, r *receipt.Receipt) error
	List(ctx context.Context, filter common.ListFilter) ([]receipt.Receipt, error)
	GetByID(ctx context.Context, id int64) (*receipt.Receipt, error)
	GetByReceiptNo(ctx context.Context, receiptNo string) (*receipt.Receipt, error)
}
