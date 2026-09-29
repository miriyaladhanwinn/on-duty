package db

import (
	"time"
)

type UserRole string

const (
	RoleStudent          UserRole = "student"
	RoleFacultyAdvisor   UserRole = "faculty_advisor"
	RoleClubCoordinator  UserRole = "club_coordinator"
	RoleHOD              UserRole = "hod"
	RoleDeptAdmin        UserRole = "dept_admin"
)

type ODStatus string

const (
	StatusDraft              ODStatus = "DRAFT"
	StatusPendingAdvisor     ODStatus = "PENDING_ADVISOR"
	StatusPendingCoordinator ODStatus = "PENDING_COORDINATOR"
	StatusPendingHOD         ODStatus = "PENDING_HOD"
	StatusApproved           ODStatus = "APPROVED"
	StatusRejected           ODStatus = "REJECTED"
)

type ApprovalStep string

const (
	StepAdvisor     ApprovalStep = "ADVISOR"
	StepCoordinator ApprovalStep = "COORDINATOR"
	StepHOD         ApprovalStep = "HOD"
)

type ApprovalAction string

const (
	ActionApprove ApprovalAction = "APPROVE"
	ActionReject  ApprovalAction = "REJECT"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	RollNo       string    `json:"roll_no"`
	Department   string    `json:"department"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ODCategory struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	MaxDaysPerSem       int    `json:"max_days_per_sem"`
	RequiresCoordinator bool   `json:"requires_coordinator"`
}

type ODRequest struct {
	ID            string    `json:"id"`
	StudentID     string    `json:"student_id"`
	StudentName   string    `json:"student_name,omitempty"`
	StudentRollNo string    `json:"student_roll_no,omitempty"`
	CategoryID    int       `json:"category_id"`
	CategoryName  string    `json:"category_name,omitempty"`
	EventName     string    `json:"event_name"`
	EventVenue    string    `json:"event_venue"`
	StartDate     string    `json:"start_date"` // YYYY-MM-DD
	EndDate       string    `json:"end_date"`   // YYYY-MM-DD
	TotalDays     int       `json:"total_days"`
	AcademicYear  string    `json:"academic_year"`
	Semester      int       `json:"semester"`
	Reason        string    `json:"reason"`
	AttachmentURL string    `json:"attachment_url"`
	Status        ODStatus  `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ODApproval struct {
	ID           string         `json:"id"`
	RequestID    string         `json:"request_id"`
	ApproverID   string         `json:"approver_id"`
	ApproverName string         `json:"approver_name,omitempty"`
	ApproverRole string         `json:"approver_role,omitempty"`
	Step         ApprovalStep   `json:"step"`
	Action       ApprovalAction `json:"action"`
	Comments     string         `json:"comments"`
	DecidedAt    time.Time      `json:"decided_at"`
}

type AuditLog struct {
	ID          string    `json:"id"`
	Action      string    `json:"action"`
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	PerformedBy string    `json:"performed_by"`
	DetailsJSON string    `json:"details_json"`
	CreatedAt   time.Time `json:"created_at"`
}

type DepartmentStats struct {
	TotalRequests      int `json:"total_requests"`
	PendingRequests    int `json:"pending_requests"`
	ApprovedRequests   int `json:"approved_requests"`
	RejectedRequests   int `json:"rejected_requests"`
	TotalODDaysGranted int `json:"total_od_days_granted"`
}
