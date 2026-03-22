package common

type PaymentStatus string

const (
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPartial PaymentStatus = "partial"
	PaymentStatusPaid    PaymentStatus = "paid"
)

type SalaryType string

const (
	SalaryTypeFixedMonthly    SalaryType = "fixed_monthly"
	SalaryTypeFixedPerClass   SalaryType = "fixed_per_class"
	SalaryTypePercentageBased SalaryType = "future_percentage_based"
)

type ClassStatus string

const (
	ClassStatusPlanned   ClassStatus = "planned"
	ClassStatusOpen      ClassStatus = "open"
	ClassStatusRunning   ClassStatus = "running"
	ClassStatusCompleted ClassStatus = "completed"
	ClassStatusClosed    ClassStatus = "closed"
)
