BEGIN;

DROP INDEX IF EXISTS idx_settings_tenant_id;
DROP INDEX IF EXISTS idx_receipts_tenant_id;
DROP INDEX IF EXISTS idx_expense_transactions_tenant_id;
DROP INDEX IF EXISTS idx_payment_transactions_tenant_id;
DROP INDEX IF EXISTS idx_enrollments_tenant_id;
DROP INDEX IF EXISTS idx_optional_fee_items_tenant_id;
DROP INDEX IF EXISTS idx_class_courses_tenant_id;
DROP INDEX IF EXISTS idx_teachers_tenant_id;
DROP INDEX IF EXISTS idx_students_tenant_id;

ALTER TABLE settings DROP CONSTRAINT IF EXISTS settings_tenant_id_key;
ALTER TABLE receipts DROP CONSTRAINT IF EXISTS receipts_tenant_receipt_no_key;
ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS enrollments_tenant_enrollment_code_key;
ALTER TABLE class_courses DROP CONSTRAINT IF EXISTS class_courses_tenant_course_class_key;
ALTER TABLE branches DROP CONSTRAINT IF EXISTS branches_tenant_branch_code_key;
ALTER TABLE teachers DROP CONSTRAINT IF EXISTS teachers_tenant_teacher_code_key;
ALTER TABLE students DROP CONSTRAINT IF EXISTS students_tenant_student_code_key;

ALTER TABLE settings DROP CONSTRAINT IF EXISTS settings_tenant_id_fkey;
ALTER TABLE receipts DROP CONSTRAINT IF EXISTS receipts_tenant_id_fkey;
ALTER TABLE expense_transactions DROP CONSTRAINT IF EXISTS expense_transactions_tenant_id_fkey;
ALTER TABLE payment_transactions DROP CONSTRAINT IF EXISTS payment_transactions_tenant_id_fkey;
ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS enrollments_tenant_id_fkey;
ALTER TABLE optional_fee_items DROP CONSTRAINT IF EXISTS optional_fee_items_tenant_id_fkey;
ALTER TABLE class_courses DROP CONSTRAINT IF EXISTS class_courses_tenant_id_fkey;
ALTER TABLE branches DROP CONSTRAINT IF EXISTS branches_tenant_id_fkey;
ALTER TABLE teachers DROP CONSTRAINT IF EXISTS teachers_tenant_id_fkey;
ALTER TABLE students DROP CONSTRAINT IF EXISTS students_tenant_id_fkey;

ALTER TABLE settings DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE receipts DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE expense_transactions DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE payment_transactions DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE enrollments DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE optional_fee_items DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE class_courses DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE branches DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE teachers DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE students DROP COLUMN IF EXISTS tenant_id;

DROP TABLE IF EXISTS tenant_users;
DROP TABLE IF EXISTS tenants;

ALTER TABLE students ADD CONSTRAINT students_student_code_key UNIQUE (student_code);
ALTER TABLE teachers ADD CONSTRAINT teachers_teacher_code_key UNIQUE (teacher_code);
ALTER TABLE branches ADD CONSTRAINT branches_branch_code_key UNIQUE (branch_code);
ALTER TABLE class_courses ADD CONSTRAINT class_courses_course_code_class_name_key UNIQUE (course_code, class_name);
ALTER TABLE enrollments ADD CONSTRAINT enrollments_enrollment_code_key UNIQUE (enrollment_code);
ALTER TABLE receipts ADD CONSTRAINT receipts_receipt_no_key UNIQUE (receipt_no);

COMMIT;
