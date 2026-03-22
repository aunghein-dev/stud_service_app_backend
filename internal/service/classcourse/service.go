package classcourse

import (
	"context"

	"student_service_app/backend/internal/domain/classcourse"
	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/errs"
	classcourserepo "student_service_app/backend/internal/repository/classcourse"
	teacherrepo "student_service_app/backend/internal/repository/teacher"
)

type Service interface {
	Create(ctx context.Context, c *classcourse.ClassCourse) error
	List(ctx context.Context, filter common.ListFilter) ([]classcourse.ClassCourse, error)
	GetByID(ctx context.Context, id int64) (*classcourse.ClassCourse, error)
	Update(ctx context.Context, c *classcourse.ClassCourse) error
	Delete(ctx context.Context, id int64) error
	CreateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error
	ListOptionalFees(ctx context.Context, classCourseID int64) ([]classcourse.OptionalFeeItem, error)
	UpdateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error
	DeleteOptionalFee(ctx context.Context, id int64) error
}

type service struct {
	repo        classcourserepo.Repository
	teacherRepo teacherrepo.Repository
}

func NewService(repo classcourserepo.Repository, teacherRepo teacherrepo.Repository) Service {
	return &service{repo: repo, teacherRepo: teacherRepo}
}

func (s *service) Create(ctx context.Context, c *classcourse.ClassCourse) error {
	if c.AssignedTeacherID != nil {
		exists, err := s.teacherRepo.ExistsByID(ctx, *c.AssignedTeacherID)
		if err != nil {
			return err
		}
		if !exists {
			return errs.BadRequest("assigned_teacher_id is invalid")
		}
	}
	return s.repo.Create(ctx, c)
}

func (s *service) List(ctx context.Context, filter common.ListFilter) ([]classcourse.ClassCourse, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id int64) (*classcourse.ClassCourse, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errs.NotFound("class course not found")
	}
	return c, nil
}

func (s *service) Update(ctx context.Context, c *classcourse.ClassCourse) error {
	exists, err := s.repo.ExistsByID(ctx, c.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("class course not found")
	}
	if c.AssignedTeacherID != nil {
		teacherExists, err := s.teacherRepo.ExistsByID(ctx, *c.AssignedTeacherID)
		if err != nil {
			return err
		}
		if !teacherExists {
			return errs.BadRequest("assigned_teacher_id is invalid")
		}
	}
	return s.repo.Update(ctx, c)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("class course not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) CreateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error {
	exists, err := s.repo.ExistsByID(ctx, item.ClassCourseID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("class course not found")
	}
	return s.repo.CreateOptionalFee(ctx, item)
}

func (s *service) ListOptionalFees(ctx context.Context, classCourseID int64) ([]classcourse.OptionalFeeItem, error) {
	return s.repo.ListOptionalFees(ctx, classCourseID)
}

func (s *service) UpdateOptionalFee(ctx context.Context, item *classcourse.OptionalFeeItem) error {
	existing, err := s.repo.GetOptionalFeeByID(ctx, item.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errs.NotFound("optional fee not found")
	}
	item.ClassCourseID = existing.ClassCourseID
	return s.repo.UpdateOptionalFee(ctx, item)
}

func (s *service) DeleteOptionalFee(ctx context.Context, id int64) error {
	existing, err := s.repo.GetOptionalFeeByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errs.NotFound("optional fee not found")
	}
	return s.repo.DeleteOptionalFee(ctx, id)
}
