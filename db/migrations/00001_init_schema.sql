-- +goose Up
-- SQL schema for Amrita CSE On-Duty (OD) & Attendance Management Workflow System

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    roll_no VARCHAR(64) UNIQUE,
    department VARCHAR(64) NOT NULL DEFAULT 'CSE',
    role VARCHAR(32) NOT NULL CHECK (role IN ('student', 'faculty_advisor', 'club_coordinator', 'hod', 'dept_admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS od_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) UNIQUE NOT NULL,
    description TEXT,
    max_days_per_sem INT NOT NULL DEFAULT 15,
    requires_coordinator BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO od_categories (name, description, max_days_per_sem, requires_coordinator) VALUES
('Hackathon', 'National / International competitive hackathons and coding summits', 8, true),
('Conference / Paper Presentation', 'IEEE/ACM conferences and research presentations', 10, true),
('Techfest / Anokha', 'Official university technical fest organizing and participation', 7, true),
('Sports & Athletics', 'Inter-collegiate and university athletic representations', 12, true),
('Club & Outreach Activity', 'Authorized departmental club initiatives and community sessions', 5, true),
('Medical Leave', 'Hospitalization or health-related certified duty leaves', 15, false)
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS od_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES od_categories(id),
    event_name VARCHAR(255) NOT NULL,
    event_venue VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INT NOT NULL CHECK (total_days > 0),
    academic_year VARCHAR(16) NOT NULL DEFAULT '2026-2027',
    semester INT NOT NULL CHECK (semester BETWEEN 1 AND 8),
    reason TEXT NOT NULL,
    attachment_url VARCHAR(512),
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_ADVISOR' CHECK (
        status IN ('DRAFT', 'PENDING_ADVISOR', 'PENDING_COORDINATOR', 'PENDING_HOD', 'APPROVED', 'REJECTED')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS od_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES od_requests(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL REFERENCES users(id),
    step VARCHAR(32) NOT NULL CHECK (step IN ('ADVISOR', 'COORDINATOR', 'HOD')),
    action VARCHAR(16) NOT NULL CHECK (action IN ('APPROVE', 'REJECT')),
    comments TEXT,
    decided_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action VARCHAR(64) NOT NULL,
    entity_type VARCHAR(64) NOT NULL,
    entity_id VARCHAR(64) NOT NULL,
    performed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    details_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_od_requests_student_id ON od_requests(student_id);
CREATE INDEX IF NOT EXISTS idx_od_requests_status ON od_requests(status);
CREATE INDEX IF NOT EXISTS idx_od_requests_dates ON od_requests(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_od_approvals_request_id ON od_approvals(request_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_type, entity_id);

-- +goose Down
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS od_approvals;
DROP TABLE IF EXISTS od_requests;
DROP TABLE IF EXISTS od_categories;
DROP TABLE IF EXISTS users;
