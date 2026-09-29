import React from 'react';
import { UserRole } from '../types';

interface NavbarProps {
  currentRole: UserRole;
  onRoleChange: (role: UserRole) => void;
  userName: string;
}

export const Navbar: React.FC<NavbarProps> = ({ currentRole, onRoleChange, userName }) => {
  return (
    <header className="bg-slate-900 border-b border-slate-800 text-white sticky top-0 z-50 shadow-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="h-10 w-10 rounded-lg bg-red-700 flex items-center justify-center font-bold text-lg text-white shadow-inner">
            OD
          </div>
          <div>
            <div className="font-bold text-lg tracking-tight flex items-center gap-2">
              <span>On Duty</span>
              <span className="text-xs bg-red-950/80 text-red-300 px-2 py-0.5 rounded-full border border-red-800/60">
                Amrita CSE
              </span>
            </div>
            <p className="text-xs text-slate-400">Department Leave & Attendance Workflow</p>
          </div>
        </div>

        <div className="flex items-center space-x-4">
          <div className="flex items-center bg-slate-800 rounded-lg p-1 border border-slate-700">
            <span className="text-xs text-slate-400 px-2 font-medium">Switch Persona:</span>
            <button
              onClick={() => onRoleChange('student')}
              className={`text-xs px-2.5 py-1 rounded transition-all font-medium ${
                currentRole === 'student' ? 'bg-red-700 text-white shadow' : 'text-slate-300 hover:text-white'
              }`}
            >
              Student
            </button>
            <button
              onClick={() => onRoleChange('faculty_advisor')}
              className={`text-xs px-2.5 py-1 rounded transition-all font-medium ${
                currentRole === 'faculty_advisor' ? 'bg-red-700 text-white shadow' : 'text-slate-300 hover:text-white'
              }`}
            >
              Advisor
            </button>
            <button
              onClick={() => onRoleChange('club_coordinator')}
              className={`text-xs px-2.5 py-1 rounded transition-all font-medium ${
                currentRole === 'club_coordinator' ? 'bg-red-700 text-white shadow' : 'text-slate-300 hover:text-white'
              }`}
            >
              Coordinator
            </button>
            <button
              onClick={() => onRoleChange('hod')}
              className={`text-xs px-2.5 py-1 rounded transition-all font-medium ${
                currentRole === 'hod' ? 'bg-red-700 text-white shadow' : 'text-slate-300 hover:text-white'
              }`}
            >
              HOD
            </button>
          </div>

          <div className="hidden md:flex items-center space-x-2 pl-2 border-l border-slate-700 text-xs">
            <span className="h-2 w-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span className="text-slate-300 font-medium">{userName}</span>
          </div>
        </div>
      </div>
    </header>
  );
};
