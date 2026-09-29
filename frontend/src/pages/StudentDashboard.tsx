import React, { useState } from 'react';
import { ODRequest, ODCategory } from '../types';
import { ODRequestCard } from '../components/ODRequestCard';

interface StudentDashboardProps {
  requests: ODRequest[];
  categories: ODCategory[];
  onSubmit: (data: any) => void;
}

export const StudentDashboard: React.FC<StudentDashboardProps> = ({ requests, categories, onSubmit }) => {
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    category_id: 1,
    event_name: '',
    event_venue: '',
    start_date: '',
    end_date: '',
    total_days: 1,
    academic_year: '2026-2027',
    semester: 3,
    reason: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
    setShowModal(false);
    setFormData({
      category_id: 1,
      event_name: '',
      event_venue: '',
      start_date: '',
      end_date: '',
      total_days: 1,
      academic_year: '2026-2027',
      semester: 3,
      reason: '',
    });
  };

  const approvedDays = requests
    .filter((r) => r.status === 'APPROVED')
    .reduce((acc, curr) => acc + curr.total_days, 0);

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-gradient-to-r from-red-900 to-slate-900 text-white p-6 rounded-2xl shadow-sm">
        <div>
          <h1 className="text-2xl font-bold">Student On-Duty Portal</h1>
          <p className="text-sm text-red-200 mt-1">Submit OD applications for Hackathons, Techfests, and Academic Representations</p>
        </div>
        <button
          onClick={() => setShowModal(true)}
          className="bg-amber-400 hover:bg-amber-300 text-slate-900 font-bold px-4 py-2.5 rounded-lg shadow transition flex items-center gap-2"
        >
          <span>+ Apply for On-Duty</span>
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="text-xs text-slate-500 font-medium">Semester OD Days Approved</div>
          <div className="text-2xl font-extrabold text-slate-800 mt-1">{approvedDays} / 20 <span className="text-xs font-normal text-slate-400">Days</span></div>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="text-xs text-slate-500 font-medium">Pending In Hierarchy</div>
          <div className="text-2xl font-extrabold text-amber-600 mt-1">
            {requests.filter((r) => r.status.startsWith('PENDING')).length}
          </div>
        </div>
        <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <div className="text-xs text-slate-500 font-medium">Total Applications</div>
          <div className="text-2xl font-extrabold text-slate-800 mt-1">{requests.length}</div>
        </div>
      </div>

      <div className="space-y-3">
        <h2 className="text-lg font-bold text-slate-800">My Leave Applications</h2>
        {requests.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-xl border border-dashed border-slate-300 text-slate-400">
            No on-duty applications submitted yet. Click "+ Apply for On-Duty" above to start!
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-4">
            {requests.map((r) => (
              <ODRequestCard key={r.id} request={r} />
            ))}
          </div>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl border border-slate-200">
            <h2 className="text-xl font-bold text-slate-900 mb-4 border-b pb-2">Apply for On-Duty (OD)</h2>
            <form onSubmit={handleSubmit} className="space-y-4 text-xs">
              <div>
                <label className="block font-semibold text-slate-700 mb-1">OD Category</label>
                <select
                  value={formData.category_id}
                  onChange={(e) => setFormData({ ...formData, category_id: Number(e.target.value) })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                >
                  {categories.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name} (Max: {c.max_days_per_sem} days)
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block font-semibold text-slate-700 mb-1">Event / Activity Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Smart India Hackathon, Anokha Core Team"
                  value={formData.event_name}
                  onChange={(e) => setFormData({ ...formData, event_name: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                />
              </div>

              <div>
                <label className="block font-semibold text-slate-700 mb-1">Venue / Location</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Amrita Coimbatore Campus, Bangalore"
                  value={formData.event_venue}
                  onChange={(e) => setFormData({ ...formData, event_venue: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block font-semibold text-slate-700 mb-1">Start Date</label>
                  <input
                    type="date"
                    required
                    value={formData.start_date}
                    onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block font-semibold text-slate-700 mb-1">End Date</label>
                  <input
                    type="date"
                    required
                    value={formData.end_date}
                    onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block font-semibold text-slate-700 mb-1">Total Days</label>
                  <input
                    type="number"
                    min="1"
                    max="15"
                    required
                    value={formData.total_days}
                    onChange={(e) => setFormData({ ...formData, total_days: Number(e.target.value) })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block font-semibold text-slate-700 mb-1">Semester</label>
                  <input
                    type="number"
                    min="1"
                    max="8"
                    required
                    value={formData.semester}
                    onChange={(e) => setFormData({ ...formData, semester: Number(e.target.value) })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label className="block font-semibold text-slate-700 mb-1">Detailed Reason / Purpose</label>
                <textarea
                  required
                  rows={3}
                  placeholder="Explain event significance, team representation, and schedule details..."
                  value={formData.reason}
                  onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-red-600 focus:outline-none"
                />
              </div>

              <div className="flex justify-end gap-2 pt-2 border-t">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 border border-slate-300 rounded-lg text-slate-600 hover:bg-slate-50 font-semibold"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-red-700 hover:bg-red-800 text-white rounded-lg font-bold shadow-md transition"
                >
                  Submit Application
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
