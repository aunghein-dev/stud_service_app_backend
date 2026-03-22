BEGIN;

INSERT INTO tenants (slug, school_name, school_address, school_phone, is_active)
VALUES (
    'demo-school',
    'Bright Future Academy',
    'No. 21, Main Road, Yangon',
    '+95-900000000',
    TRUE
)
ON CONFLICT (slug) DO UPDATE SET
    school_name = EXCLUDED.school_name,
    school_address = EXCLUDED.school_address,
    school_phone = EXCLUDED.school_phone,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

INSERT INTO settings (
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
)
SELECT
    t.id,
    'Bright Future Academy',
    'No. 21, Main Road, Yangon',
    '+95-900000000',
    'MMK',
    'BFA',
    0,
    '["cash","bank_transfer","mobile_wallet","other"]'::jsonb,
    '["books","uniform","shoes","stationery","registration_fee","exam_fee","certificate_fee"]'::jsonb,
    '{"show_logo": true, "show_signature": true, "theme": "classic"}'::jsonb
FROM tenants t
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id) DO UPDATE SET
    school_name = EXCLUDED.school_name,
    school_address = EXCLUDED.school_address,
    school_phone = EXCLUDED.school_phone,
    default_currency = EXCLUDED.default_currency,
    receipt_prefix = EXCLUDED.receipt_prefix,
    payment_methods_json = EXCLUDED.payment_methods_json,
    optional_defaults_json = EXCLUDED.optional_defaults_json,
    print_preferences_json = EXCLUDED.print_preferences_json,
    updated_at = NOW();

INSERT INTO branches (tenant_id, branch_code, branch_name, address, phone)
SELECT
    t.id,
    'BR-MAIN',
    'Main Branch',
    'No. 21, Main Road, Yangon',
    '+95-900000000'
FROM tenants t
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id, branch_code) DO UPDATE SET
    branch_name = EXCLUDED.branch_name,
    address = EXCLUDED.address,
    phone = EXCLUDED.phone,
    is_active = TRUE,
    updated_at = NOW();

INSERT INTO teachers (tenant_id, teacher_code, teacher_name, phone, subject_specialty, salary_type, default_fee_amount)
SELECT t.id, 'T-001', 'Emma Johnson', '+1-100-000-0001', 'English Speaking', 'fixed_per_class', 120000
FROM tenants t
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id, teacher_code) DO UPDATE SET
    teacher_name = EXCLUDED.teacher_name,
    phone = EXCLUDED.phone,
    subject_specialty = EXCLUDED.subject_specialty,
    salary_type = EXCLUDED.salary_type,
    default_fee_amount = EXCLUDED.default_fee_amount,
    is_active = TRUE,
    updated_at = NOW();

INSERT INTO teachers (tenant_id, teacher_code, teacher_name, phone, subject_specialty, salary_type, default_fee_amount)
SELECT t.id, 'T-002', 'Michael Chen', '+1-100-000-0002', 'Mathematics', 'fixed_monthly', 800000
FROM tenants t
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id, teacher_code) DO UPDATE SET
    teacher_name = EXCLUDED.teacher_name,
    phone = EXCLUDED.phone,
    subject_specialty = EXCLUDED.subject_specialty,
    salary_type = EXCLUDED.salary_type,
    default_fee_amount = EXCLUDED.default_fee_amount,
    is_active = TRUE,
    updated_at = NOW();

INSERT INTO students (tenant_id, student_code, full_name, gender, phone, guardian_name, guardian_phone, school_name, grade_level)
SELECT t.id, 'S-001', 'Alice Brown', 'female', '+1-200-000-0001', 'David Brown', '+1-300-000-0001', 'City School', 'Grade 8'
FROM tenants t
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id, student_code) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    gender = EXCLUDED.gender,
    phone = EXCLUDED.phone,
    guardian_name = EXCLUDED.guardian_name,
    guardian_phone = EXCLUDED.guardian_phone,
    school_name = EXCLUDED.school_name,
    grade_level = EXCLUDED.grade_level,
    is_active = TRUE,
    updated_at = NOW();

INSERT INTO students (tenant_id, student_code, full_name, gender, phone, guardian_name, guardian_phone, school_name, grade_level)
SELECT t.id, 'S-002', 'Noah Smith', 'male', '+1-200-000-0002', 'Sophia Smith', '+1-300-000-0002', 'Town School', 'Grade 9'
FROM tenants t
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id, student_code) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    gender = EXCLUDED.gender,
    phone = EXCLUDED.phone,
    guardian_name = EXCLUDED.guardian_name,
    guardian_phone = EXCLUDED.guardian_phone,
    school_name = EXCLUDED.school_name,
    grade_level = EXCLUDED.grade_level,
    is_active = TRUE,
    updated_at = NOW();

INSERT INTO class_courses (
    tenant_id,
    course_code,
    course_name,
    class_name,
    category,
    subject,
    level,
    schedule_text,
    days_of_week,
    time_start,
    time_end,
    room,
    branch_id,
    assigned_teacher_id,
    max_students,
    status,
    base_course_fee,
    registration_fee,
    exam_fee,
    certificate_fee
)
SELECT
    t.id,
    'ENG-GEN',
    'General English',
    'General English Batch 1 Morning',
    'english_speaking',
    'English',
    'Intermediate',
    'Mon/Wed/Fri',
    'Mon,Wed,Fri',
    '09:00',
    '10:30',
    'Room A',
    b.id,
    teach.id,
    35,
    'open',
    450000,
    30000,
    20000,
    15000
FROM tenants t
JOIN branches b ON b.tenant_id = t.id AND b.branch_code = 'BR-MAIN'
JOIN teachers teach ON teach.tenant_id = t.id AND teach.teacher_code = 'T-001'
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id, course_code, class_name) DO UPDATE SET
    course_name = EXCLUDED.course_name,
    category = EXCLUDED.category,
    subject = EXCLUDED.subject,
    level = EXCLUDED.level,
    schedule_text = EXCLUDED.schedule_text,
    days_of_week = EXCLUDED.days_of_week,
    time_start = EXCLUDED.time_start,
    time_end = EXCLUDED.time_end,
    room = EXCLUDED.room,
    branch_id = EXCLUDED.branch_id,
    assigned_teacher_id = EXCLUDED.assigned_teacher_id,
    max_students = EXCLUDED.max_students,
    status = EXCLUDED.status,
    base_course_fee = EXCLUDED.base_course_fee,
    registration_fee = EXCLUDED.registration_fee,
    exam_fee = EXCLUDED.exam_fee,
    certificate_fee = EXCLUDED.certificate_fee,
    updated_at = NOW();

INSERT INTO class_courses (
    tenant_id,
    course_code,
    course_name,
    class_name,
    category,
    subject,
    level,
    schedule_text,
    days_of_week,
    time_start,
    time_end,
    room,
    branch_id,
    assigned_teacher_id,
    max_students,
    status,
    base_course_fee,
    registration_fee,
    exam_fee,
    certificate_fee
)
SELECT
    t.id,
    'MATH-ACA',
    'Academic Math',
    'Academic Math Batch 1 Evening',
    'academic',
    'Mathematics',
    'Grade 9',
    'Tue/Thu/Sat',
    'Tue,Thu,Sat',
    '17:00',
    '18:30',
    'Room B',
    b.id,
    teach.id,
    30,
    'open',
    500000,
    30000,
    25000,
    15000
FROM tenants t
JOIN branches b ON b.tenant_id = t.id AND b.branch_code = 'BR-MAIN'
JOIN teachers teach ON teach.tenant_id = t.id AND teach.teacher_code = 'T-002'
WHERE t.slug = 'demo-school'
ON CONFLICT (tenant_id, course_code, class_name) DO UPDATE SET
    course_name = EXCLUDED.course_name,
    category = EXCLUDED.category,
    subject = EXCLUDED.subject,
    level = EXCLUDED.level,
    schedule_text = EXCLUDED.schedule_text,
    days_of_week = EXCLUDED.days_of_week,
    time_start = EXCLUDED.time_start,
    time_end = EXCLUDED.time_end,
    room = EXCLUDED.room,
    branch_id = EXCLUDED.branch_id,
    assigned_teacher_id = EXCLUDED.assigned_teacher_id,
    max_students = EXCLUDED.max_students,
    status = EXCLUDED.status,
    base_course_fee = EXCLUDED.base_course_fee,
    registration_fee = EXCLUDED.registration_fee,
    exam_fee = EXCLUDED.exam_fee,
    certificate_fee = EXCLUDED.certificate_fee,
    updated_at = NOW();

INSERT INTO optional_fee_items (tenant_id, class_course_id, item_name, default_amount, is_optional, is_active)
SELECT t.id, cc.id, 'Books', 30000, TRUE, TRUE
FROM tenants t
JOIN class_courses cc ON cc.tenant_id = t.id AND cc.course_code = 'ENG-GEN' AND cc.class_name = 'General English Batch 1 Morning'
WHERE t.slug = 'demo-school'
  AND NOT EXISTS (
      SELECT 1
      FROM optional_fee_items ofi
      WHERE ofi.tenant_id = t.id
        AND ofi.class_course_id = cc.id
        AND LOWER(ofi.item_name) = 'books'
  );

INSERT INTO optional_fee_items (tenant_id, class_course_id, item_name, default_amount, is_optional, is_active)
SELECT t.id, cc.id, 'Uniform', 50000, TRUE, TRUE
FROM tenants t
JOIN class_courses cc ON cc.tenant_id = t.id AND cc.course_code = 'ENG-GEN' AND cc.class_name = 'General English Batch 1 Morning'
WHERE t.slug = 'demo-school'
  AND NOT EXISTS (
      SELECT 1
      FROM optional_fee_items ofi
      WHERE ofi.tenant_id = t.id
        AND ofi.class_course_id = cc.id
        AND LOWER(ofi.item_name) = 'uniform'
  );

INSERT INTO optional_fee_items (tenant_id, class_course_id, item_name, default_amount, is_optional, is_active)
SELECT t.id, cc.id, 'Stationery', 15000, TRUE, TRUE
FROM tenants t
JOIN class_courses cc ON cc.tenant_id = t.id AND cc.course_code = 'MATH-ACA' AND cc.class_name = 'Academic Math Batch 1 Evening'
WHERE t.slug = 'demo-school'
  AND NOT EXISTS (
      SELECT 1
      FROM optional_fee_items ofi
      WHERE ofi.tenant_id = t.id
        AND ofi.class_course_id = cc.id
        AND LOWER(ofi.item_name) = 'stationery'
  );

COMMIT;
