package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/miriyaladhanwinn/on-duty/internal/auth"
	"github.com/miriyaladhanwinn/on-duty/internal/db"
)

var (
	ErrInvalidDates       = errors.New("start date must be before or equal to end date")
	ErrExceedsMaxDays     = errors.New("requested duration exceeds category semester allowance")
	ErrUnauthorizedAction = errors.New("user is not authorized to perform approval at this workflow step")
	ErrInvalidTransition = errors.New("request is not in a valid state for this action")
)

type ODService struct {
	store db.Store
}

func NewODService(store db.Store) *ODService {
	return &ODService{store: store}
}

type CreateODInput struct {
	CategoryID    int    `json:"category_id" binding:"required"`
	EventName     string `json:"event_name" binding:"required"`
	EventVenue    string `json:"event_venue" binding:"required"`
	StartDate     string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate       string `json:"end_date" binding:"required"`   // YYYY-MM-DD
	TotalDays     int    `json:"total_days" binding:"required,min=1"`
	AcademicYear  string `json:"academic_year" binding:"required"`
	Semester      int    `json:"semester" binding:"required,min=1,max=8"`
	Reason        string `json:"reason" binding:"required"`
	AttachmentURL string `json:"attachment_url"`
}

func (s *ODService) SubmitRequest(ctx context.Context, claims *auth.Claims, input CreateODInput) (*db.ODRequest, error) {
	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	if end.Before(start) {
		return nil, ErrInvalidDates
	}

	category, err := s.store.GetCategoryByID(ctx, input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if input.TotalDays > category.MaxDaysPerSem {
		return nil, fmt.Errorf("%w: max allowed is %d days", ErrExceedsMaxDays, category.MaxDaysPerSem)
	}

	// Verify existing cumulative approved days for this semester
	approvedDays, err := s.store.GetStudentApprovedDays(ctx, claims.UserID, input.Semester, input.AcademicYear)
	if err == nil && (approvedDays+input.TotalDays) > 20 {
		return nil, fmt.Errorf("semester OD quota cap exceeded (approved: %d, requested: %d, max allowed: 20)", approvedDays, input.TotalDays)
	}

	req := &db.ODRequest{
		StudentID:     claims.UserID,
		StudentName:   claims.Name,
		StudentRollNo: claims.RollNo,
		CategoryID:    input.CategoryID,
		CategoryName:  category.Name,
		EventName:     input.EventName,
		EventVenue:    input.EventVenue,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		TotalDays:     input.TotalDays,
		AcademicYear:  input.AcademicYear,
		Semester:      input.Semester,
		Reason:        input.Reason,
		AttachmentURL: input.AttachmentURL,
		Status:        db.StatusPendingAdvisor, // Initial state in workflow
	}

	created, err := s.store.CreateODRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// Audit Log
	details, _ := json.Marshal(map[string]interface{}{
		"event":      req.EventName,
		"total_days": req.TotalDays,
		"category":   category.Name,
	})
	_ = s.store.CreateAuditLog(ctx, &db.AuditLog{
		Action:      "OD_REQUEST_SUBMITTED",
		EntityType:  "OD_REQUEST",
		EntityID:    created.ID,
		PerformedBy: claims.UserID,
		DetailsJSON: string(details),
	})

	return created, nil
}

type ActionInput struct {
	Action   db.ApprovalAction `json:"action" binding:"required"` // APPROVE / REJECT
	Comments string            `json:"comments"`
}

func (s *ODService) ProcessApproval(ctx context.Context, claims *auth.Claims, requestID string, input ActionInput) (*db.ODRequest, error) {
	req, err := s.store.GetODRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	category, err := s.store.GetCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}

	var step db.ApprovalStep
	var nextStatus db.ODStatus

	// Workflow State Machine Transition Checks
	switch req.Status {
	case db.StatusPendingAdvisor:
		if claims.Role != db.RoleFacultyAdvisor && claims.Role != db.RoleDeptAdmin {
			return nil, ErrUnauthorizedAction
		}
		step = db.StepAdvisor
		if input.Action == db.ActionReject {
			nextStatus = db.StatusRejected
		} else {
			if category.RequiresCoordinator {
				nextStatus = db.StatusPendingCoordinator
			} else {
				nextStatus = db.StatusPendingHOD
			}
		}

	case db.StatusPendingCoordinator:
		if claims.Role != db.RoleClubCoordinator && claims.Role != db.RoleDeptAdmin {
			return nil, ErrUnauthorizedAction
		}
		step = db.StepCoordinator
		if input.Action == db.ActionReject {
			nextStatus = db.StatusRejected
		} else {
			nextStatus = db.StatusPendingHOD
		}

	case db.StatusPendingHOD:
		if claims.Role != db.RoleHOD && claims.Role != db.RoleDeptAdmin {
			return nil, ErrUnauthorizedAction
		}
		step = db.StepHOD
		if input.Action == db.ActionReject {
			nextStatus = db.StatusRejected
		} else {
			nextStatus = db.StatusApproved
		}

	default:
		return nil, ErrInvalidTransition
	}

	// Record the approval decision
	approval := &db.ODApproval{
		RequestID:    req.ID,
		ApproverID:   claims.UserID,
		ApproverName: claims.Name,
		ApproverRole: string(claims.Role),
		Step:         step,
		Action:       input.Action,
		Comments:     input.Comments,
	}
	_, err = s.store.CreateApprovalAction(ctx, approval)
	if err != nil {
		return nil, fmt.Errorf("failed to record approval step: %w", err)
	}

	// Update request status
	err = s.store.UpdateODRequestStatus(ctx, req.ID, nextStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update request status: %w", err)
	}
	req.Status = nextStatus

	// Audit Log
	details, _ := json.Marshal(map[string]interface{}{
		"step":        step,
		"action":      input.Action,
		"next_status": nextStatus,
		"comments":    input.Comments,
	})
	_ = s.store.CreateAuditLog(ctx, &db.AuditLog{
		Action:      fmt.Sprintf("OD_%s_%s", step, input.Action),
		EntityType:  "OD_REQUEST",
		EntityID:    req.ID,
		PerformedBy: claims.UserID,
		DetailsJSON: string(details),
	})

	return req, nil
}

func (s *ODService) GetRequestWithTimeline(ctx context.Context, requestID string) (*db.ODRequest, []db.ODApproval, error) {
	req, err := s.store.GetODRequestByID(ctx, requestID)
	if err != nil {
		return nil, nil, err
	}

	approvals, err := s.store.ListApprovalsByRequestID(ctx, requestID)
	if err != nil {
		return nil, nil, err
	}

	return req, approvals, nil
}
