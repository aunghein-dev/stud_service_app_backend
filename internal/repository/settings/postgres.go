package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"student_service_app/backend/internal/domain/settings"
	"student_service_app/backend/internal/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Get(ctx context.Context) (*settings.Setting, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, school_name, school_address, school_phone, default_currency, receipt_prefix, receipt_last_number,
	payment_methods_json, optional_defaults_json, print_preferences_json, created_at, updated_at
	FROM settings WHERE tenant_id=$1 ORDER BY id ASC LIMIT 1`
	var s settings.Setting
	if err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&s.ID, &s.SchoolName, &s.SchoolAddress, &s.SchoolPhone, &s.DefaultCurrency,
		&s.ReceiptPrefix, &s.ReceiptLastNumber, &s.PaymentMethodsJSON, &s.OptionalDefaultsJSON, &s.PrintPreferencesJSON, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresRepository) Update(ctx context.Context, s *settings.Setting) error {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return err
	}

	query := `WITH updated_settings AS (
		UPDATE settings
		SET school_name=$2,
			school_address=$3,
			school_phone=$4,
			default_currency=$5,
			receipt_prefix=$6,
			payment_methods_json=$7,
			optional_defaults_json=$8,
			print_preferences_json=$9,
			updated_at=NOW()
		WHERE id=$1 AND tenant_id=$10
		RETURNING tenant_id, updated_at
	)
	UPDATE tenants
	SET school_name=$2,
		school_address=$3,
		school_phone=$4,
		updated_at=NOW()
	WHERE id = (SELECT tenant_id FROM updated_settings)
	RETURNING (SELECT updated_at FROM updated_settings)`
	return r.db.QueryRowContext(ctx, query,
		s.ID, s.SchoolName, s.SchoolAddress, s.SchoolPhone, s.DefaultCurrency, s.ReceiptPrefix,
		s.PaymentMethodsJSON, s.OptionalDefaultsJSON, s.PrintPreferencesJSON, tenantID,
	).Scan(&s.UpdatedAt)
}

func (r *postgresRepository) EnsureForTenant(ctx context.Context, tx *sql.Tx, defaults TenantDefaults) error {
	methods, _ := json.Marshal([]string{"cash", "bank_transfer", "mobile_wallet", "other"})
	optionalDefaults, _ := json.Marshal([]string{"books", "uniform", "stationery", "registration_fee", "exam_fee"})
	prefs, _ := json.Marshal(map[string]any{"show_logo": true, "show_signature": true, "theme": "classic"})

	receiptPrefix := defaults.ReceiptPrefix
	if receiptPrefix == "" {
		receiptPrefix = "RC"
	}

	query := `INSERT INTO settings (
		tenant_id,
		school_name,
		school_address,
		school_phone,
		default_currency,
		receipt_prefix,
		receipt_last_number,
		payment_methods_json,
		optional_defaults_json,
		print_preferences_json
	) VALUES ($1,$2,$3,$4,'MMK',$5,0,$6,$7,$8)
	ON CONFLICT (tenant_id) DO UPDATE SET
		school_name=EXCLUDED.school_name,
		school_address=EXCLUDED.school_address,
		school_phone=EXCLUDED.school_phone,
		receipt_prefix=EXCLUDED.receipt_prefix,
		updated_at=NOW()`
	_, err := tx.ExecContext(ctx, query,
		defaults.TenantID,
		defaults.SchoolName,
		defaults.SchoolAddress,
		defaults.SchoolPhone,
		receiptPrefix,
		methods,
		optionalDefaults,
		prefs,
	)
	return err
}

func (r *postgresRepository) AllocateReceiptNo(ctx context.Context, tx *sql.Tx) (string, error) {
	tenantID, err := repository.TenantID(ctx)
	if err != nil {
		return "", err
	}

	var id int64
	var prefix string
	var lastNo int64
	query := `SELECT id, receipt_prefix, receipt_last_number FROM settings WHERE tenant_id=$1 ORDER BY id ASC LIMIT 1 FOR UPDATE`
	if err := tx.QueryRowContext(ctx, query, tenantID).Scan(&id, &prefix, &lastNo); err != nil {
		if err == sql.ErrNoRows {
			if err := r.EnsureForTenant(ctx, tx, TenantDefaults{
				TenantID:      tenantID,
				SchoolName:    "My School",
				SchoolAddress: "",
				SchoolPhone:   "",
				ReceiptPrefix: "RC",
			}); err != nil {
				return "", err
			}
			if err := tx.QueryRowContext(ctx, query, tenantID).Scan(&id, &prefix, &lastNo); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	nextNo := lastNo + 1
	if _, err := tx.ExecContext(ctx, `UPDATE settings SET receipt_last_number=$2, updated_at=NOW() WHERE id=$1 AND tenant_id=$3`, id, nextNo, tenantID); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%06d", prefix, nextNo), nil
}
