import React from 'react';
import { ODRequest, ODStatus } from '../types';

interface ODRequestCardProps {
  request: ODRequest;
  onAction?: (action: 'APPROVE' | 'REJECT', comments: string) => void;
  showActions?: boolean;
}

const getStatusBadge = (status: ODStatus) => {
  switch (status) {
    case 'APPROVED':
      return <span className="bg-emerald-100 text-emerald-800 text-xs font-semibold px-2.5 py-1 rounded-full border border-emerald-300">Approved</span>;
    case 'REJECTED':
      return <span className="bg-rose-100 text-rose-800 text-xs font-semibold px-2.5 py-1 rounded-full border border-rose-300">Rejected</span>;
    case 'PENDING_ADVISOR':
      return <span className="bg-amber-100 text-amber-800 text-xs font-semibold px-2.5 py-1 rounded-full border border-amber-300">Pending Advisor</span>;
    case 'PENDING_COORDINATOR':
      return <span className="bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-1 rounded-full border border-blue-300">Pending Coordinator</span>;
    case 'PENDING_HOD':
      return <span className="bg-purple-100 text-purple-800 text-xs font-semibold px-2.5 py-1 rounded-full border border-purple-300">Pending HOD</span>;
    default:
      return <span className="bg-slate-100 text-slate-800 text-xs font-semibold px-2.5 py-1 rounded-full">Draft</span>;
  }
};

export const ODRequestCard: React.FC<ODRequestCardProps> = ({ request, onAction, showActions = false }) => {
  const [comments, setComments] = React.useState('');

  return (
    <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-5 hover:shadow-md transition-shadow">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-100 pb-3">
        <div>
          <div className="flex items-center gap-2">
            <h3 className="font-bold text-slate-800 text-base">{request.event_name}</h3>
            {getStatusBadge(request.status)}
          </div>
          <p className="text-xs text-slate-500 mt-0.5">
            Category: <span className="font-medium text-slate-700">{request.category_name || 'General Event'}</span> • Venue: <span className="font-medium text-slate-700">{request.event_venue}</span>
          </p>
        </div>
        <div className="text-right">
          <span className="inline-block bg-slate-100 text-slate-800 text-xs px-2 py-0.5 rounded font-mono font-bold">
            {request.total_days} {request.total_days === 1 ? 'Day' : 'Days'} OD
          </span>
          <div className="text-xs text-slate-400 mt-0.5">
            {request.start_date} to {request.end_date}
          </div>
        </div>
      </div>

      <div className="mt-3 text-sm text-slate-600">
        <p className="line-clamp-2"><span className="font-semibold text-slate-700">Reason:</span> {request.reason}</p>
        {request.student_name && (
          <p className="text-xs text-slate-400 mt-2">
            Submitted by: <span className="font-medium text-slate-600">{request.student_name}</span> ({request.student_roll_no})
          </p>
        )}
      </div>

      {showActions && onAction && (
        <div className="mt-4 pt-3 border-t border-slate-100">
          <input
            type="text"
            placeholder="Review comments / approval remarks..."
            value={comments}
            onChange={(e) => setComments(e.target.value)}
            className="w-full text-xs px-3 py-1.5 border border-slate-300 rounded mb-2 focus:ring-1 focus:ring-red-600 focus:outline-none"
          />
          <div className="flex justify-end gap-2">
            <button
              onClick={() => onAction('REJECT', comments)}
              className="px-3 py-1 bg-rose-50 text-rose-700 hover:bg-rose-100 border border-rose-200 rounded text-xs font-semibold transition"
            >
              Reject
            </button>
            <button
              onClick={() => onAction('APPROVE', comments)}
              className="px-3 py-1 bg-emerald-600 text-white hover:bg-emerald-700 rounded text-xs font-semibold shadow-sm transition"
            >
              Approve
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
