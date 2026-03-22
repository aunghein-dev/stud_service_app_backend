package teacher

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	domain "student_service_app/backend/internal/domain/teacher"
	"student_service_app/backend/internal/errs"
	"student_service_app/backend/internal/repository/teacher"
)

type Service interface {
	Create(ctx context.Context, t *domain.Teacher) error
	List(ctx context.Context, filter common.ListFilter) ([]domain.Teacher, error)
	GetByID(ctx context.Context, id int64) (*domain.Teacher, error)
	Update(ctx context.Context, t *domain.Teacher) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo teacher.Repository
}

func NewService(repo teacher.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, t *domain.Teacher) error {
	existing, err := s.repo.GetByCode(ctx, t.TeacherCode)
	if err != nil {
		return err
	}
	if existing != nil {
		return errs.Conflict("teacher_code already exists")
	}
	if !t.IsActive {
		t.IsActive = true
	}
	return s.repo.Create(ctx, t)
}

func (s *service) List(ctx context.Context, filter common.ListFilter) ([]domain.Teacher, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.Teacher, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errs.NotFound("teacher not found")
	}
	return t, nil
}

func (s *service) Update(ctx context.Context, t *domain.Teacher) error {
	exists, err := s.repo.ExistsByID(ctx, t.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("teacher not found")
	}
	return s.repo.Update(ctx, t)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("teacher not found")
	}
	return s.repo.Delete(ctx, id)
}
