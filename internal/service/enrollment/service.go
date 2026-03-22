package enrollment

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"time"

	classcourseDomain "student_service_app/backend/internal/domain/classcourse"
	"student_service_app/backend/internal/domain/common"
	domain "student_service_app/backend/internal/domain/enrollment"
	paymentDomain "student_service_app/backend/internal/domain/payment"
	receiptDomain "student_service_app/backend/internal/domain/receipt"
	"student_service_app/backend/internal/errs"
	classcourserepo "student_service_app/backend/internal/repository/classcourse"
	enrollmentrepo "student_service_app/backend/internal/repository/enrollment"
	paymentrepo "student_service_app/backend/internal/repository/payment"
	receiptrepo "student_service_app/backend/internal/repository/receipt"
	settingsrepo "student_service_app/backend/internal/repository/settings"
	studentrepo "student_service_app/backend/internal/repository/student"
)

type OptionalItemInput struct {
	OptionalFeeItemID *int64
	ItemName          string
	Amount            float64
	Quantity          int
}

type CreateInput struct {
	StudentID      int64
	ClassCourseID  int64
	EnrollmentDate time.Time
	DiscountAmount float64
	OptionalItems  []OptionalItemInput
	InitialPayment float64
	PaymentMethod  string
	ReceivedBy     string
	Note           string
	AllowDuplicate bool
}

type Service interface {
	Create(ctx context.Context, input CreateInput) (*domain.Enrollment, []domain.EnrollmentOptionalItem, *paymentDomain.Payment, *receiptDomain.Receipt, error)
	List(ctx context.Context, filter common.ListFilter) ([]domain.Enrollment, error)
	ListByStudent(ctx context.Context, studentID int64) ([]domain.Enrollment, error)
	GetByID(ctx context.Context, id int64) (*domain.Enrollment, []domain.EnrollmentOptionalItem, error)
	Update(ctx context.Context, e *domain.Enrollment) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	db          *sql.DB
	repo        enrollmentrepo.Repository
	studentRepo studentrepo.Repository
	classRepo   classcourserepo.Repository
	paymentRepo paymentrepo.Repository
	receiptRepo receiptrepo.Repository
	settingRepo settingsrepo.Repository
}

func NewService(
	db *sql.DB,
	repo enrollmentrepo.Repository,
	studentRepo studentrepo.Repository,
	classRepo classcourserepo.Repository,
	paymentRepo paymentrepo.Repository,
	receiptRepo receiptrepo.Repository,
	settingRepo settingsrepo.Repository,
) Service {
	return &service{
		db:          db,
		repo:        repo,
		studentRepo: studentRepo,
		classRepo:   classRepo,
		paymentRepo: paymentRepo,
		receiptRepo: receiptRepo,
		settingRepo: settingRepo,
	}
}

func (s *service) Create(ctx context.Context, input CreateInput) (*domain.Enrollment, []domain.EnrollmentOptionalItem, *paymentDomain.Payment, *receiptDomain.Receipt, error) {
	exists, err := s.studentRepo.ExistsByID(ctx, input.StudentID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if !exists {
		return nil, nil, nil, nil, errs.BadRequest("student_id is invalid")
	}

	classCourse, err := s.classRepo.GetByID(ctx, input.ClassCourseID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if classCourse == nil {
		return nil, nil, nil, nil, errs.BadRequest("class_course_id is invalid")
	}

	if !input.AllowDuplicate {
		dup, err := s.repo.ExistsDuplicate(ctx, input.StudentID, input.ClassCourseID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if dup {
			return nil, nil, nil, nil, errs.Conflict("student already enrolled in this class")
		}
	}

	subTotal := classCourse.BaseCourseFee + classCourse.RegistrationFee + classCourse.ExamFee + classCourse.CertificateFee
	optionalRows := make([]domain.EnrollmentOptionalItem, 0, len(input.OptionalItems))
	for _, item := range input.OptionalItems {
		if item.Quantity <= 0 {
			return nil, nil, nil, nil, errs.BadRequest("optional item quantity must be >= 1")
		}
		if item.Amount < 0 {
			return nil, nil, nil, nil, errs.BadRequest("optional item amount must be >= 0")
		}
		total := item.Amount * float64(item.Quantity)
		subTotal += total
		optionalRows = append(optionalRows, domain.EnrollmentOptionalItem{
			OptionalFeeItemID: item.OptionalFeeItemID,
			ItemNameSnapshot:  item.ItemName,
			AmountSnapshot:    item.Amount,
			Quantity:          item.Quantity,
			TotalAmount:       total,
		})
	}

	if input.DiscountAmount < 0 {
		return nil, nil, nil, nil, errs.BadRequest("discount_amount cannot be negative")
	}
	finalFee := subTotal - input.DiscountAmount
	if finalFee < 0 {
		finalFee = 0
	}
	if input.InitialPayment < 0 {
		return nil, nil, nil, nil, errs.BadRequest("initial_payment cannot be negative")
	}
	if input.InitialPayment > finalFee {
		return nil, nil, nil, nil, errs.BadRequest("initial_payment cannot exceed final_fee")
	}

	paid := input.InitialPayment
	remaining, status := computePaymentState(finalFee, paid)

	enrollmentDate := input.EnrollmentDate
	if enrollmentDate.IsZero() {
		enrollmentDate = time.Now().UTC()
	}
	enrollment := &domain.Enrollment{
		EnrollmentCode:  generateEnrollmentCode(),
		StudentID:       input.StudentID,
		ClassCourseID:   input.ClassCourseID,
		EnrollmentDate:  enrollmentDate,
		SubTotal:        subTotal,
		DiscountAmount:  input.DiscountAmount,
		FinalFee:        finalFee,
		PaidAmount:      paid,
		RemainingAmount: remaining,
		PaymentStatus:   status,
		Note:            input.Note,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.repo.Create(ctx, tx, enrollment); err != nil {
		return nil, nil, nil, nil, err
	}
	for i := range optionalRows {
		optionalRows[i].EnrollmentID = enrollment.ID
	}
	if err := s.repo.AddOptionalItems(ctx, tx, optionalRows); err != nil {
		return nil, nil, nil, nil, err
	}

	receiptNo, err := s.settingRepo.AllocateReceiptNo(ctx, tx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var paymentTx *paymentDomain.Payment
	if paid > 0 {
		method := input.PaymentMethod
		if method == "" {
			method = "cash"
		}
		paymentDate := enrollmentDate
		paymentTx = &paymentDomain.Payment{
			ReceiptNo:     receiptNo,
			StudentID:     enrollment.StudentID,
			EnrollmentID:  enrollment.ID,
			ClassCourseID: enrollment.ClassCourseID,
			PaymentDate:   paymentDate,
			PaymentMethod: method,
			Amount:        paid,
			Note:          input.Note,
			ReceivedBy:    input.ReceivedBy,
		}
		if err := s.paymentRepo.Create(ctx, tx, paymentTx); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	rcpt, err := buildEnrollmentReceipt(receiptNo, enrollment, classCourse, optionalRows, paymentTx, input)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := s.receiptRepo.Create(ctx, tx, rcpt); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, nil, err
	}

	return enrollment, optionalRows, paymentTx, rcpt, nil
}

func (s *service) List(ctx context.Context, filter common.ListFilter) ([]domain.Enrollment, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) ListByStudent(ctx context.Context, studentID int64) ([]domain.Enrollment, error) {
	exists, err := s.studentRepo.ExistsByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errs.NotFound("student not found")
	}
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.Enrollment, []domain.EnrollmentOptionalItem, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if e == nil {
		return nil, nil, errs.NotFound("enrollment not found")
	}
	items, err := s.repo.ListOptionalItems(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return e, items, nil
}

func (s *service) Update(ctx context.Context, e *domain.Enrollment) error {
	existing, err := s.repo.GetByID(ctx, e.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errs.NotFound("enrollment not found")
	}
	if e.DiscountAmount < 0 {
		return errs.BadRequest("discount_amount cannot be negative")
	}
	e.FinalFee = existing.SubTotal - e.DiscountAmount
	if e.FinalFee < 0 {
		e.FinalFee = 0
	}
	e.PaidAmount = existing.PaidAmount
	e.RemainingAmount, e.PaymentStatus = computePaymentState(e.FinalFee, e.PaidAmount)
	return s.repo.Update(ctx, e)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errs.NotFound("enrollment not found")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Remove dependent records first to satisfy FK constraints.
	queries := []string{
		`DELETE FROM receipts WHERE enrollment_id=$1`,
		`DELETE FROM payment_transactions WHERE enrollment_id=$1`,
		`DELETE FROM enrollment_optional_items WHERE enrollment_id=$1`,
		`DELETE FROM enrollments WHERE id=$1`,
	}

	for _, q := range queries {
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func generateEnrollmentCode() string {
	return fmt.Sprintf("ENR-%d", time.Now().UTC().UnixNano())
}

func computePaymentState(finalFee, paidAmount float64) (float64, string) {
	remaining := math.Max(finalFee-paidAmount, 0)
	status := string(common.PaymentStatusUnpaid)
	if paidAmount > 0 && remaining > 0 {
		status = string(common.PaymentStatusPartial)
	}
	if remaining == 0 {
		status = string(common.PaymentStatusPaid)
	}
	return remaining, status
}

func buildEnrollmentReceipt(
	receiptNo string,
	enrollment *domain.Enrollment,
	classCourse *classcourseDomain.ClassCourse,
	items []domain.EnrollmentOptionalItem,
	paymentTx *paymentDomain.Payment,
	input CreateInput,
) (*receiptDomain.Receipt, error) {
	payload := map[string]any{
		"receipt_no":      receiptNo,
		"enrollment_code": enrollment.EnrollmentCode,
		"class_course_id": enrollment.ClassCourseID,
		"class_name":      classCourse.ClassName,
		"course_name":     classCourse.CourseName,
		"fee_breakdown": map[string]any{
			"base_course_fee":  classCourse.BaseCourseFee,
			"registration_fee": classCourse.RegistrationFee,
			"exam_fee":         classCourse.ExamFee,
			"certificate_fee":  classCourse.CertificateFee,
			"sub_total":        enrollment.SubTotal,
			"discount_amount":  enrollment.DiscountAmount,
			"final_fee":        enrollment.FinalFee,
			"paid_amount":      enrollment.PaidAmount,
			"remaining_amount": enrollment.RemainingAmount,
			"payment_status":   enrollment.PaymentStatus,
		},
		"optional_items": items,
		"payment": map[string]any{
			"initial_payment": input.InitialPayment,
			"payment_method":  input.PaymentMethod,
			"received_by":     input.ReceivedBy,
			"note":            input.Note,
		},
	}
	if paymentTx != nil {
		payload["payment_transaction_id"] = paymentTx.ID
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	typeValue := "registration"
	if paymentTx != nil {
		typeValue = "registration_payment"
	}

	var paymentID *int64
	if paymentTx != nil {
		paymentID = &paymentTx.ID
	}

	return &receiptDomain.Receipt{
		ReceiptNo:       receiptNo,
		ReceiptType:     typeValue,
		StudentID:       enrollment.StudentID,
		EnrollmentID:    enrollment.ID,
		PaymentID:       paymentID,
		ClassCourseID:   enrollment.ClassCourseID,
		TotalAmount:     enrollment.FinalFee,
		PaidAmount:      enrollment.PaidAmount,
		RemainingAmount: enrollment.RemainingAmount,
		PayloadJSON:     jsonData,
		IssuedAt:        time.Now().UTC(),
	}, nil
}
