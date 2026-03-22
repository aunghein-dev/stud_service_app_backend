package student

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/student"
)

type Repository interface {
	Create(ctx context.Context, s *student.Student) error
	List(ctx context.Context, filter common.ListFilter) ([]student.Student, error)
	GetByID(ctx context.Context, id int64) (*student.Student, error)
	GetByCode(ctx context.Context, code string) (*student.Student, error)
	Update(ctx context.Context, s *student.Student) error
	Delete(ctx context.Context, id int64) error
	ExistsByID(ctx context.Context, id int64) (bool, error)
}
