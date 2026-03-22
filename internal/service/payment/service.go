package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"student_service_app/backend/internal/domain/common"
	enrollmentDomain "student_service_app/backend/internal/domain/enrollment"
	"student_service_app/backend/internal/domain/payment"
	receiptDomain "student_service_app/backend/internal/domain/receipt"
	"student_service_app/backend/internal/errs"
	enrollmentrepo "student_service_app/backend/internal/repository/enrollment"
	paymentrepo "student_service_app/backend/internal/repository/payment"
	receiptrepo "student_service_app/backend/internal/repository/receipt"
	settingsrepo "student_service_app/backend/internal/repository/settings"
)

type CreateInput struct {
	EnrollmentID  int64
	PaymentDate   time.Time
	PaymentMethod string
	Amount        float64
	Note          string
	ReceivedBy    string
}

type UpdateInput struct {
	PaymentDate   time.Time
	PaymentMethod string
	Amount        float64
	Note          string
	ReceivedBy    string
}

type Service interface {
	Create(ctx context.Context, input CreateInput) (*payment.Payment, *receiptDomain.Receipt, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*payment.Payment, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter common.ListFilter) ([]payment.Payment, error)
	GetByID(ctx context.Context, id int64) (*payment.Payment, error)
	ListByEnrollment(ctx context.Context, enrollmentID int64) ([]payment.Payment, error)
}

type service struct {
	db          *sql.DB
	paymentRepo paymentrepo.Repository
	enrollRepo  enrollmentrepo.Repository
	receiptRepo receiptrepo.Repository
	settingRepo settingsrepo.Repository
}

func NewService(db *sql.DB, paymentRepo paymentrepo.Repository, enrollRepo enrollmentrepo.Repository, receiptRepo receiptrepo.Repository, settingRepo settingsrepo.Repository) Service {
	return &service{db: db, paymentRepo: paymentRepo, enrollRepo: enrollRepo, receiptRepo: receiptRepo, settingRepo: settingRepo}
}

func (s *service) Create(ctx context.Context, input CreateInput) (*payment.Payment, *receiptDomain.Receipt, error) {
	enroll, err := s.enrollRepo.GetByID(ctx, input.EnrollmentID)
	if err != nil {
		return nil, nil, err
	}
	if enroll == nil {
		return nil, nil, errs.BadRequest("enrollment_id is invalid")
	}
	if input.Amount <= 0 {
		return nil, nil, errs.BadRequest("amount must be > 0")
	}
	if input.Amount > enroll.RemainingAmount {
		return nil, nil, errs.BadRequest("amount exceeds remaining balance")
	}
	if input.PaymentMethod == "" {
		input.PaymentMethod = "cash"
	}
	if input.PaymentDate.IsZero() {
		input.PaymentDate = time.Now().UTC()
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = tx.Rollback() }()

	receiptNo, err := s.settingRepo.AllocateReceiptNo(ctx, tx)
	if err != nil {
		return nil, nil, err
	}

	paymentTx := &payment.Payment{
		ReceiptNo:     receiptNo,
		StudentID:     enroll.StudentID,
		EnrollmentID:  enroll.ID,
		ClassCourseID: enroll.ClassCourseID,
		PaymentDate:   input.PaymentDate,
		PaymentMethod: input.PaymentMethod,
		Amount:        input.Amount,
		Note:          input.Note,
		ReceivedBy:    input.ReceivedBy,
	}
	if err := s.paymentRepo.Create(ctx, tx, paymentTx); err != nil {
		return nil, nil, err
	}

	newPaid := enroll.PaidAmount + input.Amount
	newRemaining := enroll.FinalFee - newPaid
	status := string(common.PaymentStatusUnpaid)
	if newPaid > 0 && newRemaining > 0 {
		status = string(common.PaymentStatusPartial)
	}
	if newRemaining <= 0 {
		newRemaining = 0
		status = string(common.PaymentStatusPaid)
	}
	if err := s.enrollRepo.UpdatePaymentState(ctx, tx, enroll.ID, newPaid, newRemaining, status); err != nil {
		return nil, nil, err
	}

	rcpt, err := buildPaymentReceipt(receiptNo, enroll, paymentTx, newRemaining)
	if err != nil {
		return nil, nil, err
	}
	if err := s.receiptRepo.Create(ctx, tx, rcpt); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	return paymentTx, rcpt, nil
}

func (s *service) List(ctx context.Context, filter common.ListFilter) ([]payment.Payment, error) {
	return s.paymentRepo.List(ctx, filter)
}

func (s *service) Update(ctx context.Context, id int64, input UpdateInput) (*payment.Payment, error) {
	current, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errs.NotFound("payment not found")
	}
	enroll, err := s.enrollRepo.GetByID(ctx, current.EnrollmentID)
	if err != nil {
		return nil, err
	}
	if enroll == nil {
		return nil, errs.NotFound("enrollment not found")
	}
	if input.Amount <= 0 {
		return nil, errs.BadRequest("amount must be > 0")
	}

	maxAllowed := enroll.FinalFee - (enroll.PaidAmount - current.Amount)
	if input.Amount > maxAllowed {
		return nil, errs.BadRequest("amount exceeds remaining balance")
	}
	if input.PaymentMethod == "" {
		input.PaymentMethod = current.PaymentMethod
	}
	if input.PaymentDate.IsZero() {
		input.PaymentDate = current.PaymentDate
	}

	updated := *current
	updated.PaymentDate = input.PaymentDate
	updated.PaymentMethod = input.PaymentMethod
	updated.Amount = input.Amount
	updated.Note = input.Note
	updated.ReceivedBy = input.ReceivedBy

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.paymentRepo.Update(ctx, tx, &updated); err != nil {
		return nil, err
	}

	newPaid := enroll.PaidAmount - current.Amount + updated.Amount
	if newPaid < 0 {
		newPaid = 0
	}
	newRemaining := enroll.FinalFee - newPaid
	status := string(common.PaymentStatusUnpaid)
	if newPaid > 0 && newRemaining > 0 {
		status = string(common.PaymentStatusPartial)
	}
	if newRemaining <= 0 {
		newRemaining = 0
		status = string(common.PaymentStatusPaid)
	}

	if err := s.enrollRepo.UpdatePaymentState(ctx, tx, enroll.ID, newPaid, newRemaining, status); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE receipts SET paid_amount=$2, remaining_amount=$3 WHERE payment_id=$1`, id, updated.Amount, newRemaining); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	current, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return errs.NotFound("payment not found")
	}
	enroll, err := s.enrollRepo.GetByID(ctx, current.EnrollmentID)
	if err != nil {
		return err
	}
	if enroll == nil {
		return errs.NotFound("enrollment not found")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM receipts WHERE payment_id=$1`, id); err != nil {
		return err
	}
	if err := s.paymentRepo.Delete(ctx, tx, id); err != nil {
		return err
	}

	newPaid := enroll.PaidAmount - current.Amount
	if newPaid < 0 {
		newPaid = 0
	}
	newRemaining := enroll.FinalFee - newPaid
	status := string(common.PaymentStatusUnpaid)
	if newPaid > 0 && newRemaining > 0 {
		status = string(common.PaymentStatusPartial)
	}
	if newRemaining <= 0 {
		newRemaining = 0
		status = string(common.PaymentStatusPaid)
	}
	if err := s.enrollRepo.UpdatePaymentState(ctx, tx, enroll.ID, newPaid, newRemaining, status); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) GetByID(ctx context.Context, id int64) (*payment.Payment, error) {
	p, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errs.NotFound("payment not found")
	}
	return p, nil
}

func (s *service) ListByEnrollment(ctx context.Context, enrollmentID int64) ([]payment.Payment, error) {
	return s.paymentRepo.ListByEnrollment(ctx, enrollmentID)
}

func buildPaymentReceipt(receiptNo string, enroll *enrollmentDomain.Enrollment, paymentTx *payment.Payment, remaining float64) (*receiptDomain.Receipt, error) {
	payload := map[string]any{
		"receipt_no": receiptNo,
		"payment": map[string]any{
			"transaction_id": paymentTx.ID,
			"amount":         paymentTx.Amount,
			"payment_method": paymentTx.PaymentMethod,
			"payment_date":   paymentTx.PaymentDate,
			"received_by":    paymentTx.ReceivedBy,
			"note":           paymentTx.Note,
		},
		"enrollment": map[string]any{
			"enrollment_id":    enroll.ID,
			"enrollment_code":  enroll.EnrollmentCode,
			"final_fee":        enroll.FinalFee,
			"previous_paid":    enroll.PaidAmount,
			"new_total_paid":   enroll.PaidAmount + paymentTx.Amount,
			"remaining_amount": remaining,
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	paymentID := paymentTx.ID
	return &receiptDomain.Receipt{
		ReceiptNo:       receiptNo,
		ReceiptType:     "payment",
		StudentID:       enroll.StudentID,
		EnrollmentID:    enroll.ID,
		PaymentID:       &paymentID,
		ClassCourseID:   enroll.ClassCourseID,
		TotalAmount:     enroll.FinalFee,
		PaidAmount:      paymentTx.Amount,
		RemainingAmount: remaining,
		PayloadJSON:     jsonData,
		IssuedAt:        time.Now().UTC(),
	}, nil
}
