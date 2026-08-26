import React, { useState, useEffect } from "react";
import { X, User, ShieldAlert, LogOut, KeyRound, Check, AlertCircle, Trash2 } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { apiFetch } from "@/lib/api";
import { getAvatarUrl } from "@/lib/avatar";

interface UserProfile {
  id: number;
  username: string;
  email: string;
}

interface SettingsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const SettingsModal: React.FC<SettingsModalProps> = ({ isOpen, onClose }) => {
  const [activeTab, setActiveTab] = useState<"profile" | "security" | "account">("profile");
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [isLoadingProfile, setIsLoadingProfile] = useState(false);

  // Password change state
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [passMsg, setPassMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);
  const [isChangingPass, setIsChangingPass] = useState(false);

  // Confirmation prompts
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    if (isOpen) {
      fetchProfile();
    }
  }, [isOpen]);

  const fetchProfile = async () => {
    setIsLoadingProfile(true);
    try {
      const res = await apiFetch("/v1/users/");
      if (res.ok) {
        const data = await res.json();
        setProfile(data);
      }
    } catch (err) {
      console.error("Failed to fetch profile", err);
    } finally {
      setIsLoadingProfile(false);
    }
  };

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setPassMsg(null);

    if (newPassword !== confirmPassword) {
      setPassMsg({ text: "New passwords do not match", type: "error" });
      return;
    }
    if (newPassword.length < 6) {
      setPassMsg({ text: "Password must be at least 6 characters long", type: "error" });
      return;
    }

    setIsChangingPass(true);
    try {
      const res = await apiFetch("/v1/users/", {
        method: "PATCH",
        body: JSON.stringify({
          old_password: oldPassword,
          new_password: newPassword,
        }),
      });

      const data = await res.json();
      if (res.ok) {
        setPassMsg({ text: "Password changed successfully!", type: "success" });
        setOldPassword("");
        setNewPassword("");
        setConfirmPassword("");
      } else {
        setPassMsg({ text: data.message || data.error || "Failed to change password", type: "error" });
      }
    } catch {
      setPassMsg({ text: "Unable to connect to server", type: "error" });
    } finally {
      setIsChangingPass(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  const handleDeleteAccount = async () => {
    setIsDeleting(true);
    try {
      const res = await apiFetch("/v1/users/", {
        method: "DELETE",
      });
      if (res.ok) {
        localStorage.removeItem("token");
        navigate("/login");
      }
    } catch (err) {
      console.error("Failed to delete account", err);
    } finally {
      setIsDeleting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/30 backdrop-blur-sm p-4 animate-in fade-in duration-200">
      <div className="bg-white rounded-3xl shadow-2xl border border-slate-200/80 w-full max-w-xl overflow-hidden flex flex-col max-h-[85vh] select-none">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-5 border-b border-slate-100 bg-slate-50/50">
          <div>
            <h3 className="text-lg font-bold text-slate-800">Account Settings</h3>
            <p className="text-xs text-slate-500">Manage your profile, password, and security preferences</p>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Tab Navigation */}
        <div className="flex border-b border-slate-100 px-6 gap-2 bg-slate-50/30 text-xs font-semibold">
          <button
            onClick={() => setActiveTab("profile")}
            className={`py-3 px-4 flex items-center gap-2 border-b-2 transition-all ${
              activeTab === "profile"
                ? "border-indigo-600 text-indigo-600"
                : "border-transparent text-slate-400 hover:text-slate-600"
            }`}
          >
            <User className="w-4 h-4" />
            Profile Info
          </button>

          <button
            onClick={() => setActiveTab("security")}
            className={`py-3 px-4 flex items-center gap-2 border-b-2 transition-all ${
              activeTab === "security"
                ? "border-indigo-600 text-indigo-600"
                : "border-transparent text-slate-400 hover:text-slate-600"
            }`}
          >
            <KeyRound className="w-4 h-4" />
            Change Password
          </button>

          <button
            onClick={() => setActiveTab("account")}
            className={`py-3 px-4 flex items-center gap-2 border-b-2 transition-all ${
              activeTab === "account"
                ? "border-indigo-600 text-indigo-600"
                : "border-transparent text-slate-400 hover:text-slate-600"
            }`}
          >
            <ShieldAlert className="w-4 h-4" />
            Account Actions
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-6 overflow-y-auto flex-1 space-y-6">
          {/* PROFILE TAB */}
          {activeTab === "profile" && (
            <div className="space-y-4">
              {isLoadingProfile ? (
                <div className="py-8 text-center text-xs text-slate-400">Loading user details...</div>
              ) : profile ? (
                <div className="space-y-4">
                  <div className="flex items-center gap-4 p-4 bg-slate-50 border border-slate-200/70 rounded-2xl">
                    <img
                      src={getAvatarUrl(profile.username)}
                      alt={profile.username}
                      className="w-14 h-14 rounded-2xl object-cover border border-slate-200/80 shadow-md shadow-indigo-100/50"
                    />
                    <div>
                      <h4 className="text-base font-bold text-slate-800">{profile.username}</h4>
                      <p className="text-xs text-slate-500">{profile.email}</p>
                      <span className="inline-block mt-1 px-2.5 py-0.5 rounded-full bg-emerald-50 text-emerald-600 text-[10px] font-semibold border border-emerald-100">
                        Active Workspace Member
                      </span>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-3 text-xs">
                    <div className="p-3 bg-slate-50 border border-slate-100 rounded-2xl">
                      <span className="text-slate-400 block text-[11px] mb-0.5">User ID</span>
                      <span className="font-semibold text-slate-700">#{profile.id}</span>
                    </div>
                    <div className="p-3 bg-slate-50 border border-slate-100 rounded-2xl">
                      <span className="text-slate-400 block text-[11px] mb-0.5">Status</span>
                      <span className="font-semibold text-emerald-600 flex items-center gap-1">
                        <span className="w-2 h-2 rounded-full bg-emerald-500 inline-block" /> Online
                      </span>
                    </div>
                  </div>
                </div>
              ) : (
                <p className="text-xs text-rose-500">Failed to load profile information.</p>
              )}
            </div>
          )}

          {/* SECURITY TAB */}
          {activeTab === "security" && (
            <form onSubmit={handleChangePassword} className="space-y-4">
              {passMsg && (
                <div
                  className={`p-3 rounded-2xl text-xs font-medium flex items-center gap-2 ${
                    passMsg.type === "success"
                      ? "bg-emerald-50 border border-emerald-200 text-emerald-700"
                      : "bg-rose-50 border border-rose-200 text-rose-700"
                  }`}
                >
                  {passMsg.type === "success" ? <Check className="w-4 h-4" /> : <AlertCircle className="w-4 h-4" />}
                  <span>{passMsg.text}</span>
                </div>
              )}

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700 block">Current Password</label>
                <input
                  type="password"
                  required
                  placeholder="Enter current password"
                  value={oldPassword}
                  onChange={(e) => setOldPassword(e.target.value)}
                  className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700 block">New Password</label>
                <input
                  type="password"
                  required
                  placeholder="At least 6 characters"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700 block">Confirm New Password</label>
                <input
                  type="password"
                  required
                  placeholder="Re-enter new password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                />
              </div>

              <button
                type="submit"
                disabled={isChangingPass}
                className="w-full py-2.5 bg-indigo-600 text-white rounded-2xl text-xs font-semibold hover:bg-indigo-700 disabled:opacity-50 transition-all shadow-md shadow-indigo-100"
              >
                {isChangingPass ? "Updating Password..." : "Update Password"}
              </button>
            </form>
          )}

          {/* ACCOUNT & DANGER ZONE TAB */}
          {activeTab === "account" && (
            <div className="space-y-4">
              {/* Sign Out Card */}
              <div className="p-4 bg-slate-50 border border-slate-200/80 rounded-2xl flex items-center justify-between">
                <div>
                  <h4 className="text-xs font-bold text-slate-800">Sign Out of BukChat</h4>
                  <p className="text-[11px] text-slate-500">Safely log out of your active workspace session</p>
                </div>
                <button
                  onClick={() => setShowLogoutConfirm(true)}
                  className="px-4 py-2 bg-slate-800 text-white rounded-xl text-xs font-semibold hover:bg-slate-900 transition-colors shadow-sm flex items-center gap-1.5"
                >
                  <LogOut className="w-3.5 h-3.5" />
                  Sign Out
                </button>
              </div>

              {/* Danger Zone: Delete Account Card */}
              <div className="p-4 bg-rose-50/50 border border-rose-200/80 rounded-2xl flex items-center justify-between">
                <div>
                  <h4 className="text-xs font-bold text-rose-800">Delete Account</h4>
                  <p className="text-[11px] text-rose-600/80">Permanently erase your account and chat data</p>
                </div>
                <button
                  onClick={() => setShowDeleteConfirm(true)}
                  className="px-4 py-2 bg-rose-600 text-white rounded-xl text-xs font-semibold hover:bg-rose-700 transition-colors shadow-sm flex items-center gap-1.5"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  Delete
                </button>
              </div>

              {/* Sign Out Confirmation Modal */}
              {showLogoutConfirm && (
                <div className="p-4 bg-slate-100 border border-slate-200 rounded-2xl space-y-3 animate-in fade-in duration-150">
                  <p className="text-xs font-semibold text-slate-800">Are you sure you want to sign out?</p>
                  <div className="flex justify-end gap-2">
                    <button
                      onClick={() => setShowLogoutConfirm(false)}
                      className="px-3.5 py-1.5 bg-white border border-slate-200 text-slate-700 text-xs font-medium rounded-xl hover:bg-slate-50 transition-colors"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={handleLogout}
                      className="px-3.5 py-1.5 bg-slate-900 text-white text-xs font-semibold rounded-xl hover:bg-slate-800 transition-colors"
                    >
                      Yes, Sign Out
                    </button>
                  </div>
                </div>
              )}

              {/* Delete Account Confirmation Modal */}
              {showDeleteConfirm && (
                <div className="p-4 bg-rose-100/60 border border-rose-300 rounded-2xl space-y-3 animate-in fade-in duration-150">
                  <p className="text-xs font-bold text-rose-900">
                    Warning: This action is permanent and cannot be undone!
                  </p>
                  <div className="flex justify-end gap-2">
                    <button
                      onClick={() => setShowDeleteConfirm(false)}
                      className="px-3.5 py-1.5 bg-white border border-slate-200 text-slate-700 text-xs font-medium rounded-xl hover:bg-slate-50 transition-colors"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={handleDeleteAccount}
                      disabled={isDeleting}
                      className="px-3.5 py-1.5 bg-rose-600 text-white text-xs font-semibold rounded-xl hover:bg-rose-700 transition-colors"
                    >
                      {isDeleting ? "Deleting..." : "Permanently Delete"}
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
