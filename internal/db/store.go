package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
)

type Store interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, u *User) (*User, error)

	ListCategories(ctx context.Context) ([]ODCategory, error)
	GetCategoryByID(ctx context.Context, id int) (*ODCategory, error)

	CreateODRequest(ctx context.Context, req *ODRequest) (*ODRequest, error)
	GetODRequestByID(ctx context.Context, id string) (*ODRequest, error)
	ListODRequestsByStudent(ctx context.Context, studentID string) ([]ODRequest, error)
	ListPendingODRequestsByStatus(ctx context.Context, status ODStatus) ([]ODRequest, error)
	UpdateODRequestStatus(ctx context.Context, id string, status ODStatus) error

	CreateApprovalAction(ctx context.Context, approval *ODApproval) (*ODApproval, error)
	ListApprovalsByRequestID(ctx context.Context, requestID string) ([]ODApproval, error)

	CreateAuditLog(ctx context.Context, log *AuditLog) error
	GetStudentApprovedDays(ctx context.Context, studentID string, semester int, academicYear string) (int, error)
	GetDepartmentStats(ctx context.Context) (*DepartmentStats, error)
}

type MemoryStore struct {
	users      map[string]*User
	categories map[int]*ODCategory
	requests   map[string]*ODRequest
	approvals  map[string][]*ODApproval
	auditLogs  []*AuditLog
	mu         sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	ms := &MemoryStore{
		users:      make(map[string]*User),
		categories: make(map[int]*ODCategory),
		requests:   make(map[string]*ODRequest),
		approvals:  make(map[string][]*ODApproval),
		auditLogs:  make([]*AuditLog, 0),
	}
	ms.seedDefaults()
	return ms
}

func (ms *MemoryStore) seedDefaults() {
	// Categories
	cats := []ODCategory{
		{ID: 1, Name: "Hackathon", Description: "National/International Hackathons", MaxDaysPerSem: 8, RequiresCoordinator: true},
		{ID: 2, Name: "Conference / Paper Presentation", Description: "IEEE/ACM Research Conferences", MaxDaysPerSem: 10, RequiresCoordinator: true},
		{ID: 3, Name: "Techfest / Anokha", Description: "Amrita National Techfest Anokha Organizing", MaxDaysPerSem: 7, RequiresCoordinator: true},
		{ID: 4, Name: "Sports & Athletics", Description: "Inter-collegiate and university representation", MaxDaysPerSem: 12, RequiresCoordinator: true},
		{ID: 5, Name: "Club & Outreach Activity", Description: "The Init Club / ELITE Club authorized community events", MaxDaysPerSem: 5, RequiresCoordinator: true},
		{ID: 6, Name: "Medical Leave", Description: "Certified Medical & Hospital Leaves", MaxDaysPerSem: 15, RequiresCoordinator: false},
	}
	for _, c := range cats {
		cp := c
		ms.categories[c.ID] = &cp
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
	passHash := string(hash)

	// Seed Sample Users for Amrita CSE
	uStudent := &User{
		ID:           "student-cb-u4cse25001",
		Email:        "dhanwinn@cb.students.amrita.edu",
		PasswordHash: passHash,
		Name:         "Miriyala Dhanwinn",
		RollNo:       "CB.SC.U4CSE25639",
		Department:   "CSE",
		Role:         RoleStudent,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	uAdvisor := &User{
		ID:           "faculty-advisor-01",
		Email:        "advisor.cse@cb.amrita.edu",
		PasswordHash: passHash,
		Name:         "Dr. K. Ramanathan (Advisor)",
		RollNo:       "FAC-CSE-104",
		Department:   "CSE",
		Role:         RoleFacultyAdvisor,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	uCoordinator := &User{
		ID:           "club-coord-01",
		Email:        "anokha.coord@cb.amrita.edu",
		PasswordHash: passHash,
		Name:         "Prof. S. Lakshmi (Anokha / Techfest Coordinator)",
		RollNo:       "FAC-CSE-212",
		Department:   "CSE",
		Role:         RoleClubCoordinator,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	uHOD := &User{
		ID:           "hod-cse-01",
		Email:        "hod.cse@cb.amrita.edu",
		PasswordHash: passHash,
		Name:         "Dr. M. Sethumadhavan (HOD CSE)",
		RollNo:       "HOD-CSE-001",
		Department:   "CSE",
		Role:         RoleHOD,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	ms.users[uStudent.Email] = uStudent
	ms.users[uAdvisor.Email] = uAdvisor
	ms.users[uCoordinator.Email] = uCoordinator
	ms.users[uHOD.Email] = uHOD
}

func (ms *MemoryStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	user, ok := ms.users[email]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (ms *MemoryStore) GetUserByID(ctx context.Context, id string) (*User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for _, u := range ms.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (ms *MemoryStore) CreateUser(ctx context.Context, u *User) (*User, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.users[u.Email]; exists {
		return nil, ErrAlreadyExists
	}
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	ms.users[u.Email] = u
	return u, nil
}

func (ms *MemoryStore) ListCategories(ctx context.Context) ([]ODCategory, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []ODCategory
	for _, c := range ms.categories {
		result = append(result, *c)
	}
	return result, nil
}

func (ms *MemoryStore) GetCategoryByID(ctx context.Context, id int) (*ODCategory, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	cat, ok := ms.categories[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cat, nil
}

func (ms *MemoryStore) CreateODRequest(ctx context.Context, req *ODRequest) (*ODRequest, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if req.ID == "" {
		req.ID = fmt.Sprintf("od-%s", uuid.New().String()[:8])
	}
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	ms.requests[req.ID] = req
	return req, nil
}

func (ms *MemoryStore) GetODRequestByID(ctx context.Context, id string) (*ODRequest, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	req, ok := ms.requests[id]
	if !ok {
		return nil, ErrNotFound
	}
	return req, nil
}

func (ms *MemoryStore) ListODRequestsByStudent(ctx context.Context, studentID string) ([]ODRequest, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var res []ODRequest
	for _, r := range ms.requests {
		if r.StudentID == studentID {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (ms *MemoryStore) ListPendingODRequestsByStatus(ctx context.Context, status ODStatus) ([]ODRequest, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var res []ODRequest
	for _, r := range ms.requests {
		if r.Status == status {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (ms *MemoryStore) UpdateODRequestStatus(ctx context.Context, id string, status ODStatus) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	req, ok := ms.requests[id]
	if !ok {
		return ErrNotFound
	}
	req.Status = status
	req.UpdatedAt = time.Now()
	return nil
}

func (ms *MemoryStore) CreateApprovalAction(ctx context.Context, approval *ODApproval) (*ODApproval, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if approval.ID == "" {
		approval.ID = uuid.New().String()
	}
	approval.DecidedAt = time.Now()
	ms.approvals[approval.RequestID] = append(ms.approvals[approval.RequestID], approval)
	return approval, nil
}

func (ms *MemoryStore) ListApprovalsByRequestID(ctx context.Context, requestID string) ([]ODApproval, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	list := ms.approvals[requestID]
	var res []ODApproval
	for _, a := range list {
		res = append(res, *a)
	}
	return res, nil
}

func (ms *MemoryStore) CreateAuditLog(ctx context.Context, log *AuditLog) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if log.ID == "" {
		log.ID = uuid.New().String()
	}
	log.CreatedAt = time.Now()
	ms.auditLogs = append(ms.auditLogs, log)
	return nil
}

func (ms *MemoryStore) GetStudentApprovedDays(ctx context.Context, studentID string, semester int, academicYear string) (int, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	total := 0
	for _, r := range ms.requests {
		if r.StudentID == studentID && r.Semester == semester && r.AcademicYear == academicYear && r.Status == StatusApproved {
			total += r.TotalDays
		}
	}
	return total, nil
}

func (ms *MemoryStore) GetDepartmentStats(ctx context.Context) (*DepartmentStats, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	stats := &DepartmentStats{}
	stats.TotalRequests = len(ms.requests)
	for _, r := range ms.requests {
		switch r.Status {
		case StatusApproved:
			stats.ApprovedRequests++
			stats.TotalODDaysGranted += r.TotalDays
		case StatusRejected:
			stats.RejectedRequests++
		default:
			stats.PendingRequests++
		}
	}
	return stats, nil
}
