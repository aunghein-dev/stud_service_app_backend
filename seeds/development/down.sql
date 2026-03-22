BEGIN;

DELETE FROM optional_fee_items
WHERE tenant_id = (SELECT id FROM tenants WHERE slug = 'demo-school')
  AND (
      (class_course_id IN (
          SELECT id
          FROM class_courses
          WHERE tenant_id = (SELECT id FROM tenants WHERE slug = 'demo-school')
            AND course_code = 'ENG-GEN'
            AND class_name = 'General English Batch 1 Morning'
      ) AND LOWER(item_name) IN ('books', 'uniform'))
      OR
      (class_course_id IN (
          SELECT id
          FROM class_courses
          WHERE tenant_id = (SELECT id FROM tenants WHERE slug = 'demo-school')
            AND course_code = 'MATH-ACA'
            AND class_name = 'Academic Math Batch 1 Evening'
      ) AND LOWER(item_name) = 'stationery')
  );

DELETE FROM class_courses
WHERE tenant_id = (SELECT id FROM tenants WHERE slug = 'demo-school')
  AND (course_code, class_name) IN (
      ('ENG-GEN', 'General English Batch 1 Morning'),
      ('MATH-ACA', 'Academic Math Batch 1 Evening')
  );

DELETE FROM students
WHERE tenant_id = (SELECT id FROM tenants WHERE slug = 'demo-school')
  AND student_code IN ('S-001', 'S-002');

DELETE FROM teachers
WHERE tenant_id = (SELECT id FROM tenants WHERE slug = 'demo-school')
  AND teacher_code IN ('T-001', 'T-002');

DELETE FROM branches
WHERE tenant_id = (SELECT id FROM tenants WHERE slug = 'demo-school')
  AND branch_code = 'BR-MAIN';

DELETE FROM settings s
USING tenants t
WHERE s.tenant_id = t.id
  AND t.slug = 'demo-school'
  AND NOT EXISTS (SELECT 1 FROM tenant_users u WHERE u.tenant_id = t.id);

DELETE FROM tenants t
WHERE t.slug = 'demo-school'
  AND NOT EXISTS (SELECT 1 FROM tenant_users u WHERE u.tenant_id = t.id)
  AND NOT EXISTS (SELECT 1 FROM branches b WHERE b.tenant_id = t.id)
  AND NOT EXISTS (SELECT 1 FROM teachers teach WHERE teach.tenant_id = t.id)
  AND NOT EXISTS (SELECT 1 FROM students stu WHERE stu.tenant_id = t.id)
  AND NOT EXISTS (SELECT 1 FROM class_courses cc WHERE cc.tenant_id = t.id);

COMMIT;
