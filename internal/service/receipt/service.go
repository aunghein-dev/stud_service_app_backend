package receipt

import (
	"context"

	"student_service_app/backend/internal/domain/common"
	"student_service_app/backend/internal/domain/receipt"
	"student_service_app/backend/internal/errs"
	receiptrepo "student_service_app/backend/internal/repository/receipt"
)

type Service interface {
	List(ctx context.Context, filter common.ListFilter) ([]receipt.Receipt, error)
	GetByID(ctx context.Context, id int64) (*receipt.Receipt, error)
	GetByReceiptNo(ctx context.Context, receiptNo string) (*receipt.Receipt, error)
}

type service struct {
	repo receiptrepo.Repository
}

func NewService(repo receiptrepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, filter common.ListFilter) ([]receipt.Receipt, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id int64) (*receipt.Receipt, error) {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errs.NotFound("receipt not found")
	}
	return r, nil
}

func (s *service) GetByReceiptNo(ctx context.Context, receiptNo string) (*receipt.Receipt, error) {
	r, err := s.repo.GetByReceiptNo(ctx, receiptNo)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errs.NotFound("receipt not found")
	}
	return r, nil
}
