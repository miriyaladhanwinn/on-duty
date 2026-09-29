import React, { useState, useEffect } from 'react';
import { UserRole, ODRequest, ODCategory, DepartmentStats } from './types';
import { Navbar } from './components/Navbar';
import { StudentDashboard } from './pages/StudentDashboard';
import { FacultyDashboard } from './pages/FacultyDashboard';
import { HODDashboard } from './pages/HODDashboard';

const initialCategories: ODCategory[] = [
  { id: 1, name: 'Hackathon', description: 'National/International Hackathons', max_days_per_sem: 8, requires_coordinator: true },
  { id: 2, name: 'Conference / Paper Presentation', description: 'IEEE/ACM Research Conferences', max_days_per_sem: 10, requires_coordinator: true },
  { id: 3, name: 'Techfest / Anokha', description: 'Amrita Techfest Anokha Core Team', max_days_per_sem: 7, requires_coordinator: true },
  { id: 4, name: 'Sports & Athletics', description: 'Inter-collegiate and university representation', max_days_per_sem: 12, requires_coordinator: true },
  { id: 5, name: 'Club & Outreach Activity', description: 'The Init Club / ELITE Club community initiatives', max_days_per_sem: 5, requires_coordinator: true },
  { id: 6, name: 'Medical Leave', description: 'Certified Medical & Hospital Leaves', max_days_per_sem: 15, requires_coordinator: false },
];

const initialRequests: ODRequest[] = [
  {
    id: 'od-101',
    student_id: 'student-cb-u4cse25001',
    student_name: 'Miriyala Dhanwinn',
    student_roll_no: 'CB.SC.U4CSE25639',
    category_id: 1,
    category_name: 'Hackathon',
    event_name: 'Smart India Hackathon Grand Finale',
    event_venue: 'Amrita Coimbatore Campus',
    start_date: '2026-10-12',
    end_date: '2026-10-14',
    total_days: 3,
    academic_year: '2026-2027',
    semester: 3,
    reason: 'Selected as team lead representing Amrita CSE in national final round of AI hardware track.',
    status: 'PENDING_ADVISOR',
    created_at: '2026-09-28T10:00:00Z',
  },
  {
    id: 'od-102',
    student_id: 'student-cb-u4cse25001',
    student_name: 'Miriyala Dhanwinn',
    student_roll_no: 'CB.SC.U4CSE25639',
    category_id: 3,
    category_name: 'Techfest / Anokha',
    event_name: 'Anokha 2026 Techfest Platform Deployment',
    event_venue: 'Amrita Central Lab 3',
    start_date: '2026-09-15',
    end_date: '2026-09-16',
    total_days: 2,
    academic_year: '2026-2027',
    semester: 3,
    reason: 'Engineering high-concurrency web portals for 10,000+ national techfest participants.',
    status: 'APPROVED',
    created_at: '2026-09-10T14:30:00Z',
  },
];

export const App: React.FC = () => {
  const [role, setRole] = useState<UserRole>('student');
  const [requests, setRequests] = useState<ODRequest[]>(initialRequests);
  const [categories] = useState<ODCategory[]>(initialCategories);

  const stats: DepartmentStats = {
    total_requests: requests.length,
    pending_requests: requests.filter((r) => r.status.startsWith('PENDING')).length,
    approved_requests: requests.filter((r) => r.status === 'APPROVED').length,
    rejected_requests: requests.filter((r) => r.status === 'REJECTED').length,
    total_od_days_granted: requests
      .filter((r) => r.status === 'APPROVED')
      .reduce((sum, r) => sum + r.total_days, 0),
  };

  const handleCreateRequest = (data: any) => {
    const cat = categories.find((c) => c.id === data.category_id);
    const newReq: ODRequest = {
      ...data,
      id: `od-${Date.now().toString().slice(-4)}`,
      student_id: 'student-cb-u4cse25001',
      student_name: 'Miriyala Dhanwinn',
      student_roll_no: 'CB.SC.U4CSE25639',
      category_name: cat ? cat.name : 'General Event',
      status: 'PENDING_ADVISOR',
      created_at: new Date().toISOString(),
    };
    setRequests([newReq, ...requests]);
  };

  const handleAction = (id: string, action: 'APPROVE' | 'REJECT', comments: string) => {
    setRequests((prev) =>
      prev.map((r) => {
        if (r.id !== id) return r;
        if (action === 'REJECT') {
          return { ...r, status: 'REJECTED' };
        }
        // Hierarchy transitions
        if (r.status === 'PENDING_ADVISOR') {
          const cat = categories.find((c) => c.id === r.category_id);
          return { ...r, status: cat?.requires_coordinator ? 'PENDING_COORDINATOR' : 'PENDING_HOD' };
        }
        if (r.status === 'PENDING_COORDINATOR') {
          return { ...r, status: 'PENDING_HOD' };
        }
        if (r.status === 'PENDING_HOD') {
          return { ...r, status: 'APPROVED' };
        }
        return r;
      })
    );
  };

  const getFilteredRequestsForRole = () => {
    switch (role) {
      case 'faculty_advisor':
        return requests.filter((r) => r.status === 'PENDING_ADVISOR');
      case 'club_coordinator':
        return requests.filter((r) => r.status === 'PENDING_COORDINATOR');
      case 'hod':
        return requests.filter((r) => r.status === 'PENDING_HOD');
      default:
        return requests;
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col">
      <Navbar
        currentRole={role}
        onRoleChange={setRole}
        userName={
          role === 'student'
            ? 'Miriyala Dhanwinn (CB.SC.U4CSE25639)'
            : role === 'faculty_advisor'
            ? 'Dr. K. Ramanathan (Advisor)'
            : role === 'club_coordinator'
            ? 'Prof. S. Lakshmi (Coordinator)'
            : 'Dr. M. Sethumadhavan (HOD)'
        }
      />

      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {role === 'student' && (
          <StudentDashboard requests={requests} categories={categories} onSubmit={handleCreateRequest} />
        )}

        {role === 'faculty_advisor' && (
          <FacultyDashboard
            roleTitle="Faculty Advisor"
            requests={getFilteredRequestsForRole()}
            onAction={handleAction}
          />
        )}

        {role === 'club_coordinator' && (
          <FacultyDashboard
            roleTitle="Event / Club Coordinator"
            requests={getFilteredRequestsForRole()}
            onAction={handleAction}
          />
        )}

        {role === 'hod' && (
          <HODDashboard stats={stats} requests={getFilteredRequestsForRole()} onAction={handleAction} />
        )}
      </main>

      <footer className="bg-white border-t border-slate-200 py-4 text-center text-xs text-slate-500">
        © 2026 Amrita Vishwa Vidyapeetham • Department of Computer Science & Engineering • On-Duty Automation System
      </footer>
    </div>
  );
};
export default App;
