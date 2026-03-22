package mapper

import (
	"encoding/json"

	domain "student_service_app/backend/internal/domain/settings"
	dto "student_service_app/backend/internal/dto/settings"
)

func SettingsToDTO(s domain.Setting) dto.Response {
	methods := []string{}
	defaults := []string{}
	prefs := map[string]any{}
	_ = json.Unmarshal(s.PaymentMethodsJSON, &methods)
	_ = json.Unmarshal(s.OptionalDefaultsJSON, &defaults)
	_ = json.Unmarshal(s.PrintPreferencesJSON, &prefs)

	return dto.Response{
		SchoolName:           s.SchoolName,
		SchoolAddress:        s.SchoolAddress,
		SchoolPhone:          s.SchoolPhone,
		DefaultCurrency:      s.DefaultCurrency,
		ReceiptPrefix:        s.ReceiptPrefix,
		ReceiptLastNumber:    s.ReceiptLastNumber,
		PaymentMethods:       methods,
		OptionalItemDefaults: defaults,
		PrintPreferences:     prefs,
	}
}
