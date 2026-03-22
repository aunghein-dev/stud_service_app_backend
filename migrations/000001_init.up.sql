CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    student_code VARCHAR(50) NOT NULL UNIQUE,
    full_name VARCHAR(150) NOT NULL,
    gender VARCHAR(20) NOT NULL DEFAULT 'other',
    date_of_birth DATE NULL,
    phone VARCHAR(30) NOT NULL,
    guardian_name VARCHAR(150) NOT NULL DEFAULT '',
    guardian_phone VARCHAR(30) NOT NULL DEFAULT '',
    address VARCHAR(255) NOT NULL DEFAULT '',
    school_name VARCHAR(150) NOT NULL DEFAULT '',
    grade_level VARCHAR(50) NOT NULL DEFAULT '',
    note VARCHAR(500) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS branches (
    id BIGSERIAL PRIMARY KEY,
    branch_code VARCHAR(50) NOT NULL UNIQUE,
    branch_name VARCHAR(150) NOT NULL,
    address VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(30) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS teachers (
    id BIGSERIAL PRIMARY KEY,
    teacher_code VARCHAR(50) NOT NULL UNIQUE,
    teacher_name VARCHAR(150) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    address VARCHAR(255) NOT NULL DEFAULT '',
    subject_specialty VARCHAR(150) NOT NULL DEFAULT '',
    salary_type VARCHAR(50) NOT NULL CHECK (salary_type IN ('fixed_monthly','fixed_per_class','future_percentage_based')),
    default_fee_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    note VARCHAR(500) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS class_courses (
    id BIGSERIAL PRIMARY KEY,
    course_code VARCHAR(50) NOT NULL,
    course_name VARCHAR(150) NOT NULL,
    class_name VARCHAR(150) NOT NULL,
    category VARCHAR(30) NOT NULL CHECK (category IN ('english_speaking','academic','exam_prep','other')),
    subject VARCHAR(150) NOT NULL DEFAULT '',
    level VARCHAR(100) NOT NULL DEFAULT '',
    start_date DATE NULL,
    end_date DATE NULL,
    schedule_text VARCHAR(255) NOT NULL DEFAULT '',
    days_of_week VARCHAR(100) NOT NULL DEFAULT '',
    time_start VARCHAR(10) NOT NULL DEFAULT '',
    time_end VARCHAR(10) NOT NULL DEFAULT '',
    room VARCHAR(50) NOT NULL DEFAULT '',
    branch_id BIGINT NULL REFERENCES branches(id),
    assigned_teacher_id BIGINT NULL REFERENCES teachers(id),
    max_students INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL CHECK (status IN ('planned','open','running','completed','closed')),
    base_course_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    registration_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    exam_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    certificate_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    note VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(course_code, class_name)
);

CREATE TABLE IF NOT EXISTS optional_fee_items (
    id BIGSERIAL PRIMARY KEY,
    class_course_id BIGINT NOT NULL REFERENCES class_courses(id) ON DELETE CASCADE,
    item_name VARCHAR(100) NOT NULL,
    default_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    is_optional BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS enrollments (
    id BIGSERIAL PRIMARY KEY,
    enrollment_code VARCHAR(80) NOT NULL UNIQUE,
    student_id BIGINT NOT NULL REFERENCES students(id),
    class_course_id BIGINT NOT NULL REFERENCES class_courses(id),
    branch_id BIGINT NULL REFERENCES branches(id),
    enrollment_date DATE NOT NULL,
    sub_total NUMERIC(14,2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    final_fee NUMERIC(14,2) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    remaining_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('unpaid','partial','paid')),
    note VARCHAR(500) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS enrollment_optional_items (
    id BIGSERIAL PRIMARY KEY,
    enrollment_id BIGINT NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    optional_fee_item_id BIGINT NULL REFERENCES optional_fee_items(id),
    item_name_snapshot VARCHAR(100) NOT NULL,
    amount_snapshot NUMERIC(14,2) NOT NULL DEFAULT 0,
    quantity INTEGER NOT NULL DEFAULT 1,
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payment_transactions (
    id BIGSERIAL PRIMARY KEY,
    receipt_no VARCHAR(50) NOT NULL,
    student_id BIGINT NOT NULL REFERENCES students(id),
    enrollment_id BIGINT NOT NULL REFERENCES enrollments(id),
    class_course_id BIGINT NOT NULL REFERENCES class_courses(id),
    branch_id BIGINT NULL REFERENCES branches(id),
    payment_date DATE NOT NULL,
    payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN ('cash','bank_transfer','mobile_wallet','other')),
    amount NUMERIC(14,2) NOT NULL,
    note VARCHAR(500) NOT NULL DEFAULT '',
    received_by VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS expense_transactions (
    id BIGSERIAL PRIMARY KEY,
    expense_date DATE NOT NULL,
    expense_type VARCHAR(30) NOT NULL CHECK (expense_type IN ('teacher_fee','books','uniform','shoes','stationery','rent','utilities','marketing','misc')),
    teacher_id BIGINT NULL REFERENCES teachers(id),
    class_course_id BIGINT NULL REFERENCES class_courses(id),
    branch_id BIGINT NULL REFERENCES branches(id),
    amount NUMERIC(14,2) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    payment_method VARCHAR(50) NOT NULL DEFAULT '',
    reference_no VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS receipts (
    id BIGSERIAL PRIMARY KEY,
    receipt_no VARCHAR(50) NOT NULL UNIQUE,
    receipt_type VARCHAR(50) NOT NULL,
    student_id BIGINT NOT NULL REFERENCES students(id),
    enrollment_id BIGINT NOT NULL REFERENCES enrollments(id),
    payment_id BIGINT NULL REFERENCES payment_transactions(id),
    class_course_id BIGINT NOT NULL REFERENCES class_courses(id),
    branch_id BIGINT NULL REFERENCES branches(id),
    total_amount NUMERIC(14,2) NOT NULL,
    paid_amount NUMERIC(14,2) NOT NULL,
    remaining_amount NUMERIC(14,2) NOT NULL,
    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    issued_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS settings (
    id BIGSERIAL PRIMARY KEY,
    school_name VARCHAR(150) NOT NULL,
    school_address VARCHAR(255) NOT NULL DEFAULT '',
    school_phone VARCHAR(30) NOT NULL DEFAULT '',
    default_currency VARCHAR(10) NOT NULL DEFAULT 'MMK',
    receipt_prefix VARCHAR(20) NOT NULL DEFAULT 'RC',
    receipt_last_number BIGINT NOT NULL DEFAULT 0,
    payment_methods_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    optional_defaults_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    print_preferences_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
