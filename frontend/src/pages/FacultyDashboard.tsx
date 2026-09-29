import React from 'react';
import { ODRequest } from '../types';
import { ODRequestCard } from '../components/ODRequestCard';

interface FacultyDashboardProps {
  roleTitle: string;
  requests: ODRequest[];
  onAction: (id: string, action: 'APPROVE' | 'REJECT', comments: string) => void;
}

export const FacultyDashboard: React.FC<FacultyDashboardProps> = ({ roleTitle, requests, onAction }) => {
  return (
    <div className="space-y-6">
      <div className="bg-slate-900 text-white p-6 rounded-2xl shadow-sm border border-slate-800">
        <h1 className="text-2xl font-bold">{roleTitle} Action Queue</h1>
        <p className="text-sm text-slate-300 mt-1">Review pending duty-leave applications awaiting verification and hierarchy sign-off</p>
      </div>

      <div className="space-y-3">
        <div className="flex justify-between items-center">
          <h2 className="text-lg font-bold text-slate-800">Pending Review ({requests.length})</h2>
          <span className="text-xs text-slate-500">Sorted by submission date</span>
        </div>

        {requests.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-xl border border-dashed border-slate-300 text-slate-400">
            No pending applications in your review queue. All clear!
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
