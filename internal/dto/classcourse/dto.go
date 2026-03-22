package classcourse

type CreateRequest struct {
	CourseCode        string   `json:"course_code" validate:"required,max=50"`
	CourseName        string   `json:"course_name" validate:"required,max=150"`
	ClassName         string   `json:"class_name" validate:"required,max=150"`
	Category          string   `json:"category" validate:"required,oneof=english_speaking academic exam_prep other"`
	Subject           string   `json:"subject" validate:"max=150"`
	Level             string   `json:"level" validate:"max=100"`
	StartDate         string   `json:"start_date"`
	EndDate           string   `json:"end_date"`
	ScheduleText      string   `json:"schedule_text" validate:"max=255"`
	DaysOfWeek        []string `json:"days_of_week"`
	TimeStart         string   `json:"time_start"`
	TimeEnd           string   `json:"time_end"`
	Room              string   `json:"room" validate:"max=50"`
	AssignedTeacherID *int64   `json:"assigned_teacher_id"`
	MaxStudents       int      `json:"max_students" validate:"gte=0"`
	Status            string   `json:"status" validate:"required,oneof=planned open running completed closed"`
	BaseCourseFee     float64  `json:"base_course_fee" validate:"gte=0"`
	RegistrationFee   float64  `json:"registration_fee" validate:"gte=0"`
	ExamFee           float64  `json:"exam_fee" validate:"gte=0"`
	CertificateFee    float64  `json:"certificate_fee" validate:"gte=0"`
	Note              string   `json:"note" validate:"max=500"`
}

type UpdateRequest = CreateRequest

type OptionalFeeItemRequest struct {
	ItemName      string  `json:"item_name" validate:"required,max=100"`
	DefaultAmount float64 `json:"default_amount" validate:"gte=0"`
	IsOptional    bool    `json:"is_optional"`
	IsActive      bool    `json:"is_active"`
}

type OptionalFeeItemResponse struct {
	ID            int64   `json:"id"`
	ClassCourseID int64   `json:"class_course_id"`
	ItemName      string  `json:"item_name"`
	DefaultAmount float64 `json:"default_amount"`
	IsOptional    bool    `json:"is_optional"`
	IsActive      bool    `json:"is_active"`
}

type Response struct {
	ID                int64    `json:"id"`
	CourseCode        string   `json:"course_code"`
	CourseName        string   `json:"course_name"`
	ClassName         string   `json:"class_name"`
	Category          string   `json:"category"`
	Subject           string   `json:"subject"`
	Level             string   `json:"level"`
	StartDate         string   `json:"start_date,omitempty"`
	EndDate           string   `json:"end_date,omitempty"`
	ScheduleText      string   `json:"schedule_text"`
	DaysOfWeek        []string `json:"days_of_week"`
	TimeStart         string   `json:"time_start"`
	TimeEnd           string   `json:"time_end"`
	Room              string   `json:"room"`
	AssignedTeacherID *int64   `json:"assigned_teacher_id"`
	MaxStudents       int      `json:"max_students"`
	Status            string   `json:"status"`
	BaseCourseFee     float64  `json:"base_course_fee"`
	RegistrationFee   float64  `json:"registration_fee"`
	ExamFee           float64  `json:"exam_fee"`
	CertificateFee    float64  `json:"certificate_fee"`
	Note              string   `json:"note"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}
