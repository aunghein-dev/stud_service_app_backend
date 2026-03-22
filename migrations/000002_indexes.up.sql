CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_students_full_name ON students USING gin (full_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_students_student_code ON students(student_code);
CREATE INDEX IF NOT EXISTS idx_students_phone ON students(phone);

CREATE INDEX IF NOT EXISTS idx_teachers_teacher_name ON teachers USING gin (teacher_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_teachers_teacher_code ON teachers(teacher_code);

CREATE INDEX IF NOT EXISTS idx_class_courses_class_name ON class_courses USING gin (class_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_class_courses_course_name ON class_courses USING gin (course_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_class_courses_category ON class_courses(category);
CREATE INDEX IF NOT EXISTS idx_class_courses_assigned_teacher_id ON class_courses(assigned_teacher_id);
CREATE INDEX IF NOT EXISTS idx_class_courses_branch_id ON class_courses(branch_id);

CREATE INDEX IF NOT EXISTS idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_class_course_id ON enrollments(class_course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_branch_id ON enrollments(branch_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_enrollment_date ON enrollments(enrollment_date);
CREATE INDEX IF NOT EXISTS idx_enrollments_payment_status ON enrollments(payment_status);

CREATE INDEX IF NOT EXISTS idx_payments_receipt_no ON payment_transactions(receipt_no);
CREATE INDEX IF NOT EXISTS idx_payments_student_id ON payment_transactions(student_id);
CREATE INDEX IF NOT EXISTS idx_payments_class_course_id ON payment_transactions(class_course_id);
CREATE INDEX IF NOT EXISTS idx_payments_branch_id ON payment_transactions(branch_id);
CREATE INDEX IF NOT EXISTS idx_payments_payment_date ON payment_transactions(payment_date);

CREATE INDEX IF NOT EXISTS idx_expenses_type ON expense_transactions(expense_type);
CREATE INDEX IF NOT EXISTS idx_expenses_date ON expense_transactions(expense_date);
CREATE INDEX IF NOT EXISTS idx_expenses_teacher_id ON expense_transactions(teacher_id);
CREATE INDEX IF NOT EXISTS idx_expenses_class_course_id ON expense_transactions(class_course_id);
CREATE INDEX IF NOT EXISTS idx_expenses_branch_id ON expense_transactions(branch_id);

CREATE INDEX IF NOT EXISTS idx_receipts_receipt_no ON receipts(receipt_no);
CREATE INDEX IF NOT EXISTS idx_receipts_issued_at ON receipts(issued_at);
CREATE INDEX IF NOT EXISTS idx_receipts_student_id ON receipts(student_id);
CREATE INDEX IF NOT EXISTS idx_receipts_branch_id ON receipts(branch_id);
