import React from 'react';
import { ODRequest, DepartmentStats } from '../types';
import { ODRequestCard } from '../components/ODRequestCard';

interface HODDashboardProps {
  stats: DepartmentStats;
  requests: ODRequest[];
  onAction: (id: string, action: 'APPROVE' | 'REJECT', comments: string) => void;
}

export const HODDashboard: React.FC<HODDashboardProps> = ({ stats, requests, onAction }) => {
  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-slate-950 via-slate-900 to-red-950 text-white p-6 rounded-2xl shadow-sm border border-slate-800">
        <h1 className="text-2xl font-bold">Amrita CSE Department Executive Portal</h1>
        <p className="text-sm text-slate-300 mt-1">Head of Department (HOD) Final Approval Authority & Attendance Impact Metrics</p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="text-xs text-slate-500 font-medium">Total OD Applications</div>
          <div className="text-2xl font-extrabold text-slate-800 mt-1">{stats.total_requests}</div>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="text-xs text-slate-500 font-medium">Pending Final Sign-Off</div>
          <div className="text-2xl font-extrabold text-purple-600 mt-1">{requests.length}</div>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="text-xs text-slate-500 font-medium">Approved ODs</div>
          <div className="text-2xl font-extrabold text-emerald-600 mt-1">{stats.approved_requests}</div>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="text-xs text-slate-500 font-medium">Total OD Days Granted</div>
          <div className="text-2xl font-extrabold text-amber-600 mt-1">{stats.total_od_days_granted}</div>
        </div>
      </div>

      <div className="space-y-3">
        <h2 className="text-lg font-bold text-slate-800">Final Executive Sign-Off Queue ({requests.length})</h2>
        {requests.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-xl border border-dashed border-slate-300 text-slate-400">
            No applications pending HOD final review. All department workflows up to date!
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-4">
            {requests.map((r) => (
              <ODRequestCard
                key={r.id}
                request={r}
                showActions={true}
                onAction={(action, comments) => onAction(r.id, action, comments)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
