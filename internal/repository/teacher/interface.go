package teacher

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/teacher"
)

type Repository interface {
	Create(ctx context.Context, t *teacher.Teacher) error
	List(ctx context.Context, filter common.ListFilter) ([]teacher.Teacher, error)
	GetByID(ctx context.Context, id int64) (*teacher.Teacher, error)
	GetByCode(ctx context.Context, code string) (*teacher.Teacher, error)
	Update(ctx context.Context, t *teacher.Teacher) error
	Delete(ctx context.Context, id int64) error
	ExistsByID(ctx context.Context, id int64) (bool, error)
}
