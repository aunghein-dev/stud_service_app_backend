BEGIN;

CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) NOT NULL UNIQUE,
    school_name VARCHAR(150) NOT NULL,
    school_address VARCHAR(255) NOT NULL DEFAULT '',
    school_phone VARCHAR(30) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tenant_users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(30) NOT NULL DEFAULT 'owner',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, email)
);

INSERT INTO tenants (slug, school_name, school_address, school_phone, is_active)
SELECT 'default-school',
       COALESCE(s.school_name, 'My School'),
       COALESCE(s.school_address, ''),
       COALESCE(s.school_phone, ''),
       TRUE
FROM settings s
WHERE NOT EXISTS (SELECT 1 FROM tenants)
LIMIT 1;

INSERT INTO tenants (slug, school_name, school_address, school_phone, is_active)
SELECT 'default-school', 'My School', '', '', TRUE
WHERE NOT EXISTS (SELECT 1 FROM tenants);

ALTER TABLE students ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE teachers ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE branches ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE class_courses ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE optional_fee_items ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE payment_transactions ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE expense_transactions ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE receipts ADD COLUMN IF NOT EXISTS tenant_id BIGINT;
ALTER TABLE settings ADD COLUMN IF NOT EXISTS tenant_id BIGINT;

UPDATE students
SET tenant_id = (SELECT id FROM tenants ORDER BY id ASC LIMIT 1)
WHERE tenant_id IS NULL;

UPDATE teachers
SET tenant_id = (SELECT id FROM tenants ORDER BY id ASC LIMIT 1)
WHERE tenant_id IS NULL;

UPDATE branches
SET tenant_id = (SELECT id FROM tenants ORDER BY id ASC LIMIT 1)
WHERE tenant_id IS NULL;

UPDATE class_courses
SET tenant_id = COALESCE(class_courses.tenant_id, (SELECT id FROM tenants ORDER BY id ASC LIMIT 1))
WHERE tenant_id IS NULL;

UPDATE optional_fee_items ofi
SET tenant_id = c.tenant_id
FROM class_courses c
WHERE ofi.class_course_id = c.id
  AND ofi.tenant_id IS NULL;

UPDATE enrollments e
SET tenant_id = c.tenant_id
FROM class_courses c
WHERE e.class_course_id = c.id
  AND e.tenant_id IS NULL;

UPDATE payment_transactions p
SET tenant_id = e.tenant_id
FROM enrollments e
WHERE p.enrollment_id = e.id
  AND p.tenant_id IS NULL;

UPDATE expense_transactions
SET tenant_id = COALESCE(
    (SELECT tenant_id FROM class_courses WHERE id = expense_transactions.class_course_id),
    (SELECT tenant_id FROM teachers WHERE id = expense_transactions.teacher_id),
    (SELECT id FROM tenants ORDER BY id ASC LIMIT 1)
)
WHERE tenant_id IS NULL;

UPDATE receipts r
SET tenant_id = e.tenant_id
FROM enrollments e
WHERE r.enrollment_id = e.id
  AND r.tenant_id IS NULL;

UPDATE settings
SET tenant_id = (SELECT id FROM tenants ORDER BY id ASC LIMIT 1)
WHERE tenant_id IS NULL;

ALTER TABLE students ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE teachers ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE branches ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE class_courses ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE optional_fee_items ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE enrollments ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE payment_transactions ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE expense_transactions ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE receipts ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE settings ALTER COLUMN tenant_id SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'students_tenant_id_fkey') THEN
        ALTER TABLE students ADD CONSTRAINT students_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'teachers_tenant_id_fkey') THEN
        ALTER TABLE teachers ADD CONSTRAINT teachers_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'branches_tenant_id_fkey') THEN
        ALTER TABLE branches ADD CONSTRAINT branches_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'class_courses_tenant_id_fkey') THEN
        ALTER TABLE class_courses ADD CONSTRAINT class_courses_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'optional_fee_items_tenant_id_fkey') THEN
        ALTER TABLE optional_fee_items ADD CONSTRAINT optional_fee_items_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'enrollments_tenant_id_fkey') THEN
        ALTER TABLE enrollments ADD CONSTRAINT enrollments_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'payment_transactions_tenant_id_fkey') THEN
        ALTER TABLE payment_transactions ADD CONSTRAINT payment_transactions_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'expense_transactions_tenant_id_fkey') THEN
        ALTER TABLE expense_transactions ADD CONSTRAINT expense_transactions_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'receipts_tenant_id_fkey') THEN
        ALTER TABLE receipts ADD CONSTRAINT receipts_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'settings_tenant_id_fkey') THEN
        ALTER TABLE settings ADD CONSTRAINT settings_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);
    END IF;
END $$;

ALTER TABLE students DROP CONSTRAINT IF EXISTS students_student_code_key;
ALTER TABLE teachers DROP CONSTRAINT IF EXISTS teachers_teacher_code_key;
ALTER TABLE branches DROP CONSTRAINT IF EXISTS branches_branch_code_key;
ALTER TABLE class_courses DROP CONSTRAINT IF EXISTS class_courses_course_code_class_name_key;
ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS enrollments_enrollment_code_key;
ALTER TABLE receipts DROP CONSTRAINT IF EXISTS receipts_receipt_no_key;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'students_tenant_student_code_key') THEN
        ALTER TABLE students ADD CONSTRAINT students_tenant_student_code_key UNIQUE (tenant_id, student_code);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'teachers_tenant_teacher_code_key') THEN
        ALTER TABLE teachers ADD CONSTRAINT teachers_tenant_teacher_code_key UNIQUE (tenant_id, teacher_code);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'branches_tenant_branch_code_key') THEN
        ALTER TABLE branches ADD CONSTRAINT branches_tenant_branch_code_key UNIQUE (tenant_id, branch_code);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'class_courses_tenant_course_class_key') THEN
        ALTER TABLE class_courses ADD CONSTRAINT class_courses_tenant_course_class_key UNIQUE (tenant_id, course_code, class_name);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'enrollments_tenant_enrollment_code_key') THEN
        ALTER TABLE enrollments ADD CONSTRAINT enrollments_tenant_enrollment_code_key UNIQUE (tenant_id, enrollment_code);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'receipts_tenant_receipt_no_key') THEN
        ALTER TABLE receipts ADD CONSTRAINT receipts_tenant_receipt_no_key UNIQUE (tenant_id, receipt_no);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'settings_tenant_id_key') THEN
        ALTER TABLE settings ADD CONSTRAINT settings_tenant_id_key UNIQUE (tenant_id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_students_tenant_id ON students(tenant_id);
CREATE INDEX IF NOT EXISTS idx_teachers_tenant_id ON teachers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_class_courses_tenant_id ON class_courses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_optional_fee_items_tenant_id ON optional_fee_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_tenant_id ON enrollments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_tenant_id ON payment_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_expense_transactions_tenant_id ON expense_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_receipts_tenant_id ON receipts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_settings_tenant_id ON settings(tenant_id);

COMMIT;
