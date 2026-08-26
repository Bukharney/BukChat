import React, { useState } from "react";
import { MessageSquare, Users, Settings, LogOut, MessageCircleCode, AlertCircle } from "lucide-react";
import { useNavigate } from "react-router-dom";

interface NavSidebarProps {
  activeTab: "chats" | "friends" | "settings";
  setActiveTab: (tab: "chats" | "friends" | "settings") => void;
  pendingCount: number;
}

export const NavSidebar: React.FC<NavSidebarProps> = ({
  activeTab,
  setActiveTab,
  pendingCount,
}) => {
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);
  const navigate = useNavigate();

  const handleConfirmLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <aside className="w-16 md:w-20 bg-white border-r border-slate-200/80 flex flex-col items-center justify-between py-5 z-20 shadow-sm shrink-0 relative select-none">
      {/* Top Logo */}
      <div className="flex flex-col items-center gap-6">
        <div className="w-10 h-10 rounded-2xl bg-indigo-600 flex items-center justify-center text-white shadow-md shadow-indigo-200 transition-transform hover:scale-105">
          <MessageCircleCode className="w-6 h-6" />
        </div>

        {/* Navigation Items */}
        <nav className="flex flex-col items-center gap-3">
          <button
            onClick={() => setActiveTab("chats")}
            title="Chats"
            className={`relative p-3 rounded-2xl transition-all duration-200 ${
              activeTab === "chats"
                ? "bg-indigo-50 text-indigo-600 font-medium shadow-sm"
                : "text-slate-400 hover:text-slate-600 hover:bg-slate-100/70"
            }`}
          >
            <MessageSquare className="w-5 h-5" />
          </button>

          <button
            onClick={() => setActiveTab("friends")}
            title="Friends & Requests"
            className={`relative p-3 rounded-2xl transition-all duration-200 ${
              activeTab === "friends"
                ? "bg-indigo-50 text-indigo-600 font-medium shadow-sm"
                : "text-slate-400 hover:text-slate-600 hover:bg-slate-100/70"
            }`}
          >
            <Users className="w-5 h-5" />
            {pendingCount > 0 && (
              <span className="absolute -top-1 -right-1 w-5 h-5 bg-indigo-600 text-white text-[10px] font-semibold flex items-center justify-center rounded-full border-2 border-white animate-pulse">
                {pendingCount}
              </span>
            )}
          </button>

          <button
            onClick={() => setActiveTab("settings")}
            title="Settings"
            className={`relative p-3 rounded-2xl transition-all duration-200 ${
              activeTab === "settings"
                ? "bg-indigo-50 text-indigo-600 font-medium shadow-sm"
                : "text-slate-400 hover:text-slate-600 hover:bg-slate-100/70"
            }`}
          >
            <Settings className="w-5 h-5" />
          </button>
        </nav>
      </div>

      {/* Bottom Logout Action */}
      <div className="flex flex-col items-center gap-3 relative">
        <button
          onClick={() => setShowLogoutConfirm((prev) => !prev)}
          title="Sign Out"
          className="p-3 rounded-2xl text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-all duration-200"
        >
          <LogOut className="w-5 h-5" />
        </button>

        {/* Logout Popover Confirmation */}
        {showLogoutConfirm && (
          <div className="absolute left-16 bottom-0 z-50 w-56 p-4 bg-white border border-slate-200 rounded-2xl shadow-xl animate-in slide-in-from-left-2 duration-150">
            <div className="flex items-center gap-2 mb-2 text-rose-600">
              <AlertCircle className="w-4 h-4" />
              <span className="text-xs font-bold text-slate-800">Sign Out?</span>
            </div>
            <p className="text-[11px] text-slate-500 mb-3 leading-tight">
              Are you sure you want to end your current session?
            </p>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setShowLogoutConfirm(false)}
                className="flex-1 py-1.5 bg-slate-100 text-slate-700 rounded-xl text-xs font-medium hover:bg-slate-200 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleConfirmLogout}
                className="flex-1 py-1.5 bg-rose-600 text-white rounded-xl text-xs font-semibold hover:bg-rose-700 transition-colors shadow-sm"
              >
                Sign Out
              </button>
            </div>
          </div>
        )}
      </div>
    </aside>
  );
};
