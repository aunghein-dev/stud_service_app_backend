package classcourse

import (
	"context"

	"student_service_app/backend/internal/domain/classcourse"
	"student_service_app/backend/internal/domain/common"
)

type Repository interface {
	Create(ctx context.Context, c *classcourse.ClassCourse) error
	List(ctx context.Context, filter common.ListFilter) ([]classcourse.ClassCourse, error)
	GetByID(ctx context.Context, id int64) (*classcourse.ClassCourse, error)
	Update(ctx context.Context, c *classcourse.ClassCourse) error
	Delete(ctx context.Context, id int64) error
	ExistsByID(ctx context.Context, id int64) (bool, error)

	CreateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error
	ListOptionalFees(ctx context.Context, classCourseID int64) ([]classcourse.OptionalFeeItem, error)
	GetOptionalFeeByID(ctx context.Context, id int64) (*classcourse.OptionalFeeItem, error)
	UpdateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error
	DeleteOptionalFee(ctx context.Context, id int64) error
}
