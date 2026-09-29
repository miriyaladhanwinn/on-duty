-- name: GetUserByEmail :one
SELECT id, email, password_hash, name, roll_no, department, role, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, password_hash, name, roll_no, department, role, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, roll_no, department, role)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, name, roll_no, department, role, created_at;

-- name: ListCategories :many
SELECT id, name, description, max_days_per_sem, requires_coordinator
FROM od_categories
ORDER BY id ASC;

-- name: CreateODRequest :one
INSERT INTO od_requests (
    student_id, category_id, event_name, event_venue,
    start_date, end_date, total_days, academic_year, semester, reason, attachment_url, status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, student_id, category_id, event_name, event_venue, start_date, end_date, total_days, academic_year, semester, reason, attachment_url, status, created_at;

-- name: GetODRequestByID :one
SELECT r.id, r.student_id, r.category_id, r.event_name, r.event_venue,
       r.start_date, r.end_date, r.total_days, r.academic_year, r.semester,
       r.reason, r.attachment_url, r.status, r.created_at, r.updated_at,
       u.name AS student_name, u.roll_no AS student_roll_no, c.name AS category_name
FROM od_requests r
JOIN users u ON r.student_id = u.id
JOIN od_categories c ON r.category_id = c.id
WHERE r.id = $1 LIMIT 1;

-- name: ListODRequestsByStudent :many
SELECT r.id, r.student_id, r.category_id, r.event_name, r.event_venue,
       r.start_date, r.end_date, r.total_days, r.academic_year, r.semester,
       r.reason, r.attachment_url, r.status, r.created_at, c.name AS category_name
FROM od_requests r
JOIN od_categories c ON r.category_id = c.id
WHERE r.student_id = $1
ORDER BY r.created_at DESC;

-- name: ListPendingODRequests :many
SELECT r.id, r.student_id, r.category_id, r.event_name, r.event_venue,
       r.start_date, r.end_date, r.total_days, r.academic_year, r.semester,
       r.reason, r.attachment_url, r.status, r.created_at,
       u.name AS student_name, u.roll_no AS student_roll_no, c.name AS category_name
FROM od_requests r
JOIN users u ON r.student_id = u.id
JOIN od_categories c ON r.category_id = c.id
WHERE r.status = $1
ORDER BY r.created_at ASC;

-- name: UpdateODRequestStatus :exec
UPDATE od_requests
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: CreateApprovalAction :one
INSERT INTO od_approvals (request_id, approver_id, step, action, comments)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, request_id, approver_id, step, action, comments, decided_at;

-- name: ListApprovalsByRequestID :many
SELECT a.id, a.request_id, a.approver_id, a.step, a.action, a.comments, a.decided_at,
       u.name AS approver_name, u.role AS approver_role
FROM od_approvals a
JOIN users u ON a.approver_id = u.id
WHERE a.request_id = $1
ORDER BY a.decided_at ASC;

-- name: CreateAuditLog :exec
INSERT INTO audit_logs (action, entity_type, entity_id, performed_by, details_json)
VALUES ($1, $2, $3, $4, $5);

-- name: GetStudentSemesterODTotalDays :one
SELECT COALESCE(SUM(total_days), 0)::INT AS total_approved_days
FROM od_requests
WHERE student_id = $1 AND semester = $2 AND academic_year = $3 AND status = 'APPROVED';
