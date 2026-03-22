package student

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	domain "student_service_app/backend/internal/domain/student"
	"student_service_app/backend/internal/errs"
	"student_service_app/backend/internal/repository/student"
)

type Service interface {
	Create(ctx context.Context, s *domain.Student) error
	List(ctx context.Context, filter common.ListFilter) ([]domain.Student, error)
	GetByID(ctx context.Context, id int64) (*domain.Student, error)
	Update(ctx context.Context, s *domain.Student) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo student.Repository
}

func NewService(repo student.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, st *domain.Student) error {
	existing, err := s.repo.GetByCode(ctx, st.StudentCode)
	if err != nil {
		return err
	}
	if existing != nil {
		return errs.Conflict("student_code already exists")
	}
	if st.IsActive == false {
		st.IsActive = true
	}
	return s.repo.Create(ctx, st)
}

func (s *service) List(ctx context.Context, filter common.ListFilter) ([]domain.Student, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.Student, error) {
	st, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, errs.NotFound("student not found")
	}
	return st, nil
}

func (s *service) Update(ctx context.Context, st *domain.Student) error {
	exists, err := s.repo.ExistsByID(ctx, st.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("student not found")
	}
	return s.repo.Update(ctx, st)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("student not found")
	}
	return s.repo.Delete(ctx, id)
}
