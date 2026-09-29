export type UserRole = 'student' | 'faculty_advisor' | 'club_coordinator' | 'hod' | 'dept_admin';

export type ODStatus =
  | 'DRAFT'
  | 'PENDING_ADVISOR'
  | 'PENDING_COORDINATOR'
  | 'PENDING_HOD'
  | 'APPROVED'
  | 'REJECTED';

export interface User {
  id: string;
  name: string;
  email: string;
  roll_no: string;
  department: string;
  role: UserRole;
}

export interface ODCategory {
  id: number;
  name: string;
  description: string;
  max_days_per_sem: number;
  requires_coordinator: boolean;
}

export interface ODApproval {
  id: string;
  request_id: string;
  approver_id: string;
  approver_name: string;
  approver_role: string;
  step: 'ADVISOR' | 'COORDINATOR' | 'HOD';
  action: 'APPROVE' | 'REJECT';
  comments: string;
  decided_at: string;
}

export interface ODRequest {
  id: string;
  student_id: string;
  student_name?: string;
  student_roll_no?: string;
  category_id: number;
  category_name?: string;
  event_name: string;
  event_venue: string;
  start_date: string;
  end_date: string;
  total_days: number;
  academic_year: string;
  semester: number;
  reason: string;
  attachment_url?: string;
  status: ODStatus;
  created_at: string;
}

export interface DepartmentStats {
  total_requests: number;
  pending_requests: number;
  approved_requests: number;
  rejected_requests: number;
  total_od_days_granted: number;
}
