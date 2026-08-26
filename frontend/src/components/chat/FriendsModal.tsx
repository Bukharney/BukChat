import React, { useState } from "react";
import { UserPlus, Check, X, Search, ShieldCheck } from "lucide-react";
import { Friend } from "@/data";
import { getAvatarUrl } from "@/lib/avatar";

interface FriendRequest {
  id: number;
  username: string;
}

interface FriendsModalProps {
  isOpen: boolean;
  onClose: () => void;
  friendRequests: FriendRequest[];
  friendList: Friend[];
  onAcceptFriend: (username: string) => Promise<void>;
  onRejectFriend: (username: string) => Promise<void>;
  onAddFriend: (username: string) => Promise<void>;
}

export const FriendsModal: React.FC<FriendsModalProps> = ({
  isOpen,
  onClose,
  friendRequests,
  friendList,
  onAcceptFriend,
  onRejectFriend,
  onAddFriend,
}) => {
  const [newFriendUsername, setNewFriendUsername] = useState("");
  const [statusMessage, setStatusMessage] = useState<{ text: string; type: "success" | "error" } | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleAddSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newFriendUsername.trim()) return;

    setIsSubmitting(true);
    setStatusMessage(null);
    try {
      await onAddFriend(newFriendUsername.trim());
      setStatusMessage({ text: `Friend request sent to ${newFriendUsername}!`, type: "success" });
      setNewFriendUsername("");
    } catch (err: any) {
      setStatusMessage({ text: err.message || "Failed to send request", type: "error" });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/30 backdrop-blur-sm p-4 animate-in fade-in duration-200">
      <div className="bg-white rounded-3xl shadow-2xl border border-slate-200/80 w-full max-w-lg overflow-hidden flex flex-col max-h-[85vh]">
        {/* Modal Header */}
        <div className="flex items-center justify-between px-6 py-5 border-b border-slate-100 bg-slate-50/50">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-2xl bg-indigo-50 text-indigo-600 flex items-center justify-center">
              <UserPlus className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-slate-800">Manage Friends</h3>
              <p className="text-xs text-slate-500">Send requests & view pending invites</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Modal Body */}
        <div className="p-6 overflow-y-auto space-y-6">
          {/* Add Friend Form */}
          <div>
            <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">
              Add a Friend by Username
            </label>
            <form onSubmit={handleAddSubmit} className="flex gap-2">
              <div className="relative flex-1">
                <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
                <input
                  type="text"
                  placeholder="Enter username..."
                  value={newFriendUsername}
                  onChange={(e) => setNewFriendUsername(e.target.value)}
                  className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                />
              </div>
              <button
                type="submit"
                disabled={isSubmitting || !newFriendUsername.trim()}
                className="px-5 py-2.5 bg-indigo-600 text-white rounded-2xl text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-md shadow-indigo-100"
              >
                {isSubmitting ? "Sending..." : "Send Request"}
              </button>
            </form>
            {statusMessage && (
              <p
                className={`mt-2 text-xs font-medium ${
                  statusMessage.type === "success" ? "text-emerald-600" : "text-rose-600"
                }`}
              >
                {statusMessage.text}
              </p>
            )}
          </div>

          {/* Pending Friend Requests Section */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
                Pending Requests ({friendRequests.length})
              </span>
            </div>

            {friendRequests.length === 0 ? (
              <div className="p-4 rounded-2xl bg-slate-50 border border-dashed border-slate-200 text-center">
                <ShieldCheck className="w-6 h-6 text-slate-300 mx-auto mb-1" />
                <p className="text-xs text-slate-500">No incoming friend requests</p>
              </div>
            ) : (
              <div className="space-y-2">
                {friendRequests.map((req) => (
                  <div
                    key={req.id || req.username}
                    className="flex items-center justify-between p-3 bg-slate-50 border border-slate-200/70 rounded-2xl"
                  >
                    <div className="flex items-center gap-3">
                      <img
                        src={getAvatarUrl(req.username)}
                        alt={req.username}
                        className="w-9 h-9 rounded-xl object-cover border border-slate-200/60 shadow-sm"
                      />
                      <span className="text-sm font-medium text-slate-800">{req.username}</span>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <button
                        onClick={() => onAcceptFriend(req.username)}
                        className="p-2 rounded-xl bg-emerald-50 text-emerald-600 hover:bg-emerald-100 transition-colors"
                        title="Accept Request"
                      >
                        <Check className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => onRejectFriend(req.username)}
                        className="p-2 rounded-xl bg-rose-50 text-rose-600 hover:bg-rose-100 transition-colors"
                        title="Reject Request"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Current Friends List Section */}
          <div>
            <span className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">
              Your Friends ({friendList.length})
            </span>
            {friendList.length === 0 ? (
              <p className="text-xs text-slate-400 italic">No friends added yet.</p>
            ) : (
              <div className="space-y-2 max-h-48 overflow-y-auto">
                {friendList.map((friend) => (
                  <div
                    key={friend.id}
                    className="flex items-center justify-between p-3 bg-white border border-slate-100 rounded-2xl hover:bg-slate-50 transition-colors"
                  >
                    <div className="flex items-center gap-3">
                      <img
                        src={getAvatarUrl(friend.username)}
                        alt={friend.username}
                        className="w-9 h-9 rounded-xl object-cover border border-slate-200/60 shadow-sm"
                      />
                      <div>
                        <span className="text-sm font-medium text-slate-800 block">{friend.username}</span>
                        {friend.is_online !== false ? (
                          <span className="text-[11px] text-emerald-500 flex items-center gap-1 font-medium">
                            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 inline-block" />
                            Online
                          </span>
                        ) : (
                          <span className="text-[11px] text-slate-400 flex items-center gap-1">
                            <span className="w-1.5 h-1.5 rounded-full bg-slate-300 inline-block" />
                            Offline
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
