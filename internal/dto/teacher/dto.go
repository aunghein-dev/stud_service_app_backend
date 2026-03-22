package teacher

type CreateRequest struct {
	TeacherCode      string  `json:"teacher_code" validate:"required,max=50"`
	TeacherName      string  `json:"teacher_name" validate:"required,max=150"`
	Phone            string  `json:"phone" validate:"required,max=30"`
	Address          string  `json:"address" validate:"max=255"`
	SubjectSpecialty string  `json:"subject_specialty" validate:"max=150"`
	SalaryType       string  `json:"salary_type" validate:"required,oneof=fixed_monthly fixed_per_class future_percentage_based"`
	DefaultFeeAmount float64 `json:"default_fee_amount" validate:"gte=0"`
	Note             string  `json:"note" validate:"max=500"`
	IsActive         *bool   `json:"is_active"`
}

type UpdateRequest struct {
	TeacherName      string  `json:"teacher_name" validate:"required,max=150"`
	Phone            string  `json:"phone" validate:"required,max=30"`
	Address          string  `json:"address" validate:"max=255"`
	SubjectSpecialty string  `json:"subject_specialty" validate:"max=150"`
	SalaryType       string  `json:"salary_type" validate:"required,oneof=fixed_monthly fixed_per_class future_percentage_based"`
	DefaultFeeAmount float64 `json:"default_fee_amount" validate:"gte=0"`
	Note             string  `json:"note" validate:"max=500"`
	IsActive         bool    `json:"is_active"`
}

type Response struct {
	ID               int64   `json:"id"`
	TeacherCode      string  `json:"teacher_code"`
	TeacherName      string  `json:"teacher_name"`
	Phone            string  `json:"phone"`
	Address          string  `json:"address"`
	SubjectSpecialty string  `json:"subject_specialty"`
	SalaryType       string  `json:"salary_type"`
	DefaultFeeAmount float64 `json:"default_fee_amount"`
	Note             string  `json:"note"`
	IsActive         bool    `json:"is_active"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}
