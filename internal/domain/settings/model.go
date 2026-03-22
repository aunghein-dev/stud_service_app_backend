package settings

import "time"

type Setting struct {
	ID                   int64     `db:"id"`
	SchoolName           string    `db:"school_name"`
	SchoolAddress        string    `db:"school_address"`
	SchoolPhone          string    `db:"school_phone"`
	DefaultCurrency      string    `db:"default_currency"`
	ReceiptPrefix        string    `db:"receipt_prefix"`
	ReceiptLastNumber    int64     `db:"receipt_last_number"`
	PaymentMethodsJSON   []byte    `db:"payment_methods_json"`
	OptionalDefaultsJSON []byte    `db:"optional_defaults_json"`
	PrintPreferencesJSON []byte    `db:"print_preferences_json"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}
