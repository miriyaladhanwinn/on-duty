package service

import (
	"context"
	"testing"

	"github.com/miriyaladhanwinn/on-duty/internal/auth"
	"github.com/miriyaladhanwinn/on-duty/internal/db"
)

func setupTestContext() (*ODService, *auth.Claims, *auth.Claims, *auth.Claims, *auth.Claims) {
	store := db.NewMemoryStore()
	svc := NewODService(store)

	studentClaims := &auth.Claims{
		UserID:     "student-cb-u4cse25001",
		Email:      "dhanwinn@cb.students.amrita.edu",
		Name:       "Miriyala Dhanwinn",
		RollNo:     "CB.SC.U4CSE25639",
		Department: "CSE",
		Role:       db.RoleStudent,
	}

	advisorClaims := &auth.Claims{
		UserID:     "faculty-advisor-01",
		Email:      "advisor.cse@cb.amrita.edu",
		Name:       "Dr. K. Ramanathan",
		Department: "CSE",
		Role:       db.RoleFacultyAdvisor,
	}

	coordinatorClaims := &auth.Claims{
		UserID:     "club-coord-01",
		Email:      "anokha.coord@cb.amrita.edu",
		Name:       "Prof. S. Lakshmi",
		Department: "CSE",
		Role:       db.RoleClubCoordinator,
	}

	hodClaims := &auth.Claims{
		UserID:     "hod-cse-01",
		Email:      "hod.cse@cb.amrita.edu",
		Name:       "Dr. M. Sethumadhavan",
		Department: "CSE",
		Role:       db.RoleHOD,
	}

	return svc, studentClaims, advisorClaims, coordinatorClaims, hodClaims
}

func TestCompleteApprovalWorkflow(t *testing.T) {
	svc, student, advisor, coordinator, hod := setupTestContext()
	ctx := context.Background()

	// 1. Student submits Hackathon OD request (Category 1, requires coordinator)
	input := CreateODInput{
		CategoryID:   1, // Hackathon
		EventName:    "Smart India Hackathon Grand Finale",
		EventVenue:   "Amrita Coimbatore Campus",
		StartDate:    "2026-10-10",
		EndDate:      "2026-10-12",
		TotalDays:    3,
		AcademicYear: "2026-2027",
		Semester:     3,
		Reason:       "Selected for SIH hardware & software round representation",
	}

	req, err := svc.SubmitRequest(ctx, student, input)
	if err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}

	if req.Status != db.StatusPendingAdvisor {
		t.Fatalf("Expected StatusPendingAdvisor, got: %s", req.Status)
	}

	// 2. Student tries to approve their own request (should fail)
	_, err = svc.ProcessApproval(ctx, student, req.ID, ActionInput{Action: db.ActionApprove})
	if err == nil {
		t.Fatal("Expected unauthorized error for student approving, got nil")
	}

	// 3. Faculty advisor approves
	req, err = svc.ProcessApproval(ctx, advisor, req.ID, ActionInput{Action: db.ActionApprove, Comments: "Recommended for SIH"})
	if err != nil {
		t.Fatalf("Advisor approval failed: %v", err)
	}
	if req.Status != db.StatusPendingCoordinator {
		t.Fatalf("Expected StatusPendingCoordinator, got: %s", req.Status)
	}

	// 4. Coordinator approves
	req, err = svc.ProcessApproval(ctx, coordinator, req.ID, ActionInput{Action: db.ActionApprove, Comments: "Event verified with organizer"})
	if err != nil {
		t.Fatalf("Coordinator approval failed: %v", err)
	}
	if req.Status != db.StatusPendingHOD {
		t.Fatalf("Expected StatusPendingHOD, got: %s", req.Status)
	}

	// 5. HOD approves
	req, err = svc.ProcessApproval(ctx, hod, req.ID, ActionInput{Action: db.ActionApprove, Comments: "Approved by HOD. Attendance exemption granted."})
	if err != nil {
		t.Fatalf("HOD approval failed: %v", err)
	}
	if req.Status != db.StatusApproved {
		t.Fatalf("Expected StatusApproved, got: %s", req.Status)
	}

	// 6. Verify full timeline
	fetchedReq, timeline, err := svc.GetRequestWithTimeline(ctx, req.ID)
	if err != nil {
		t.Fatalf("GetRequestWithTimeline failed: %v", err)
	}

	if fetchedReq.Status != db.StatusApproved {
		t.Errorf("Expected status approved, got: %s", fetchedReq.Status)
	}

	if len(timeline) != 3 {
		t.Errorf("Expected 3 approval milestones in timeline, got: %d", len(timeline))
	}
}

func TestCategoryDaysQuotaEnforcement(t *testing.T) {
	svc, student, _, _, _ := setupTestContext()
	ctx := context.Background()

	// Hackathon max is 8 days. Request 10 days -> should fail
	input := CreateODInput{
		CategoryID:   1, // Hackathon
		EventName:    "Overlimit Event",
		EventVenue:   "Bangalore",
		StartDate:    "2026-11-01",
		EndDate:      "2026-11-10",
		TotalDays:    10,
		AcademicYear: "2026-2027",
		Semester:     3,
		Reason:       "Excessive leave test",
	}

	_, err := svc.SubmitRequest(ctx, student, input)
	if err == nil {
		t.Fatal("Expected error exceeding max days allowance, got nil")
	}
}
