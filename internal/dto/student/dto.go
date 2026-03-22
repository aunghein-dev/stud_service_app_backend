package student

type CreateRequest struct {
	StudentCode   string `json:"student_code" validate:"required,max=50"`
	FullName      string `json:"full_name" validate:"required,max=150"`
	Gender        string `json:"gender" validate:"omitempty,oneof=male female other"`
	DateOfBirth   string `json:"date_of_birth"`
	Phone         string `json:"phone" validate:"required,max=30"`
	GuardianName  string `json:"guardian_name" validate:"max=150"`
	GuardianPhone string `json:"guardian_phone" validate:"max=30"`
	Address       string `json:"address" validate:"max=255"`
	SchoolName    string `json:"school_name" validate:"max=150"`
	GradeLevel    string `json:"grade_level" validate:"max=50"`
	Note          string `json:"note" validate:"max=500"`
	IsActive      *bool  `json:"is_active"`
}

type UpdateRequest struct {
	FullName      string `json:"full_name" validate:"required,max=150"`
	Gender        string `json:"gender" validate:"omitempty,oneof=male female other"`
	DateOfBirth   string `json:"date_of_birth"`
	Phone         string `json:"phone" validate:"required,max=30"`
	GuardianName  string `json:"guardian_name" validate:"max=150"`
	GuardianPhone string `json:"guardian_phone" validate:"max=30"`
	Address       string `json:"address" validate:"max=255"`
	SchoolName    string `json:"school_name" validate:"max=150"`
	GradeLevel    string `json:"grade_level" validate:"max=50"`
	Note          string `json:"note" validate:"max=500"`
	IsActive      bool   `json:"is_active"`
}

type Response struct {
	ID            int64  `json:"id"`
	StudentCode   string `json:"student_code"`
	FullName      string `json:"full_name"`
	Gender        string `json:"gender"`
	DateOfBirth   string `json:"date_of_birth,omitempty"`
	Phone         string `json:"phone"`
	GuardianName  string `json:"guardian_name"`
	GuardianPhone string `json:"guardian_phone"`
	Address       string `json:"address"`
	SchoolName    string `json:"school_name"`
	GradeLevel    string `json:"grade_level"`
	Note          string `json:"note"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
