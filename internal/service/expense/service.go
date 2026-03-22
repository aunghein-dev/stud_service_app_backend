package expense

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/expense"
	"student_service_app/backend/internal/errs"
	classcourserepo "student_service_app/backend/internal/repository/classcourse"
	expenserepo "student_service_app/backend/internal/repository/expense"
	teacherrepo "student_service_app/backend/internal/repository/teacher"
)

type Service interface {
	Create(ctx context.Context, e *expense.Expense) error
	List(ctx context.Context, filter common.ListFilter) ([]expense.Expense, error)
	GetByID(ctx context.Context, id int64) (*expense.Expense, error)
	Update(ctx context.Context, e *expense.Expense) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo        expenserepo.Repository
	teacherRepo teacherrepo.Repository
	classRepo   classcourserepo.Repository
}

func NewService(repo expenserepo.Repository, teacherRepo teacherrepo.Repository, classRepo classcourserepo.Repository) Service {
	return &service{repo: repo, teacherRepo: teacherRepo, classRepo: classRepo}
}

func (s *service) Create(ctx context.Context, e *expense.Expense) error {
	if e.TeacherID != nil {
		exists, err := s.teacherRepo.ExistsByID(ctx, *e.TeacherID)
		if err != nil {
			return err
		}
		if !exists {
			return errs.BadRequest("teacher_id is invalid")
		}
	}
	if e.ClassCourseID != nil {
		exists, err := s.classRepo.ExistsByID(ctx, *e.ClassCourseID)
		if err != nil {
			return err
		}
		if !exists {
			return errs.BadRequest("class_course_id is invalid")
		}
	}
	return s.repo.Create(ctx, e)
}

func (s *service) List(ctx context.Context, filter common.ListFilter) ([]expense.Expense, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id int64) (*expense.Expense, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errs.NotFound("expense not found")
	}
	return e, nil
}

func (s *service) Update(ctx context.Context, e *expense.Expense) error {
	existing, err := s.repo.GetByID(ctx, e.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errs.NotFound("expense not found")
	}
	return s.repo.Update(ctx, e)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errs.NotFound("expense not found")
	}
	return s.repo.Delete(ctx, id)
}
