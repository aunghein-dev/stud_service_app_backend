package settings

import (
	"context"
	"database/sql"
	"testing"

	"student_service_app/backend/internal/domain/settings"
	settingsrepo "student_service_app/backend/internal/repository/settings"

	"github.com/stretchr/testify/require"
)

type fakeSettingsRepo struct {
	setting    *settings.Setting
	updateCall *settings.Setting
}

func (f *fakeSettingsRepo) Get(ctx context.Context) (*settings.Setting, error) {
	return f.setting, nil
}

func (f *fakeSettingsRepo) Update(ctx context.Context, s *settings.Setting) error {
	clone := *s
	f.updateCall = &clone
	return nil
}

func (f *fakeSettingsRepo) EnsureForTenant(ctx context.Context, tx *sql.Tx, defaults settingsrepo.TenantDefaults) error {
	return nil
}

func (f *fakeSettingsRepo) AllocateReceiptNo(ctx context.Context, tx *sql.Tx) (string, error) {
	return "", nil
}

func TestUpdateSanitizesSchoolBranding(t *testing.T) {
	repo := &fakeSettingsRepo{
		setting: &settings.Setting{
			ID:                   1,
			SchoolName:           "Old",
			SchoolAddress:        "Old address",
			SchoolPhone:          "Old phone",
			DefaultCurrency:      "MMK",
			ReceiptPrefix:        "RC",
			PaymentMethodsJSON:   []byte("[]"),
			OptionalDefaultsJSON: []byte("[]"),
			PrintPreferencesJSON: []byte("{}"),
		},
	}

	svc := NewService(repo)
	updated, err := svc.Update(context.Background(), UpdateInput{
		SchoolName:           "  Bright Future Academy  ",
		SchoolAddress:        "  Main Road  ",
		SchoolPhone:          "  +95-900000000  ",
		DefaultCurrency:      " mmk ",
		ReceiptPrefix:        " bfa ",
		PaymentMethods:       []string{"cash"},
		OptionalItemDefaults: []string{"books"},
		PrintPreferences:     map[string]any{"show_logo": true},
	})

	require.NoError(t, err)
	require.NotNil(t, repo.updateCall)
	require.Equal(t, "Bright Future Academy", updated.SchoolName)
	require.Equal(t, "Main Road", updated.SchoolAddress)
	require.Equal(t, "+95-900000000", updated.SchoolPhone)
	require.Equal(t, "MMK", updated.DefaultCurrency)
	require.Equal(t, "BFA", updated.ReceiptPrefix)
}

func TestUpdateRequiresSchoolName(t *testing.T) {
	repo := &fakeSettingsRepo{
		setting: &settings.Setting{
			ID:                   1,
			PaymentMethodsJSON:   []byte("[]"),
			OptionalDefaultsJSON: []byte("[]"),
			PrintPreferencesJSON: []byte("{}"),
		},
	}

	svc := NewService(repo)
	_, err := svc.Update(context.Background(), UpdateInput{
		SchoolName:       "   ",
		DefaultCurrency:  "MMK",
		ReceiptPrefix:    "RC",
		PrintPreferences: map[string]any{},
	})

	require.Error(t, err)
	require.Nil(t, repo.updateCall)
}
