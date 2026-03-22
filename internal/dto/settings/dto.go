package settings

type Response struct {
	SchoolName           string         `json:"school_name"`
	SchoolAddress        string         `json:"school_address"`
	SchoolPhone          string         `json:"school_phone"`
	DefaultCurrency      string         `json:"default_currency"`
	ReceiptPrefix        string         `json:"receipt_prefix"`
	ReceiptLastNumber    int64          `json:"receipt_last_number"`
	PaymentMethods       []string       `json:"payment_methods"`
	OptionalItemDefaults []string       `json:"optional_item_defaults"`
	PrintPreferences     map[string]any `json:"print_preferences"`
}

type UpdateRequest struct {
	SchoolName           string         `json:"school_name" validate:"required,max=150"`
	SchoolAddress        string         `json:"school_address" validate:"max=255"`
	SchoolPhone          string         `json:"school_phone" validate:"max=30"`
	DefaultCurrency      string         `json:"default_currency" validate:"required,max=10"`
	ReceiptPrefix        string         `json:"receipt_prefix" validate:"required,max=20"`
	PaymentMethods       []string       `json:"payment_methods"`
	OptionalItemDefaults []string       `json:"optional_item_defaults"`
	PrintPreferences     map[string]any `json:"print_preferences"`
}
