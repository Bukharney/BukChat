import React, { useState } from "react";
import { Search, UserPlus, Users, Sparkles } from "lucide-react";
import { Friend } from "@/data";
import { getAvatarUrl } from "@/lib/avatar";

interface ChannelSidebarProps {
  friendList: Friend[];
  selectedUser: Friend | null;
  setSelectedUser: (friend: Friend) => void;
  pendingCount: number;
  onOpenFriendsModal: () => void;
}

export const ChannelSidebar: React.FC<ChannelSidebarProps> = ({
  friendList,
  selectedUser,
  setSelectedUser,
  pendingCount,
  onOpenFriendsModal,
}) => {
  const [searchQuery, setSearchQuery] = useState("");

  const sortedFriends = [...friendList].sort((a, b) => {
    const aUnread = (a.unreadCount || 0) > 0 ? 1 : 0;
    const bUnread = (b.unreadCount || 0) > 0 ? 1 : 0;
    if (aUnread !== bUnread) {
      return bUnread - aUnread;
    }

    const aLastMsg = a.messages?.[a.messages.length - 1];
    const bLastMsg = b.messages?.[b.messages.length - 1];

    const aTime = aLastMsg?.timestamp ? new Date(aLastMsg.timestamp).getTime() : (aLastMsg?.id || 0);
    const bTime = bLastMsg?.timestamp ? new Date(bLastMsg.timestamp).getTime() : (bLastMsg?.id || 0);

    return bTime - aTime;
  });

  const filteredFriends = sortedFriends.filter((friend) =>
    friend.username.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <aside className="w-72 md:w-80 bg-slate-50/80 border-r border-slate-200/80 flex flex-col h-full shrink-0 select-none">
      {/* Top Header & Search */}
      <div className="p-4 border-b border-slate-200/60 bg-white/50 backdrop-blur-sm">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-lg font-bold text-slate-800 tracking-tight flex items-center gap-2">
            Messages
            <span className="text-xs px-2 py-0.5 rounded-full bg-indigo-50 text-indigo-600 font-semibold border border-indigo-100">
              {friendList.length}
            </span>
          </h2>

          <button
            onClick={onOpenFriendsModal}
            className="p-2 rounded-xl bg-indigo-50 text-indigo-600 hover:bg-indigo-100 transition-colors relative flex items-center gap-1 text-xs font-semibold"
            title="Manage Friends & Add Users"
          >
            <UserPlus className="w-4 h-4" />
            <span>Add</span>
            {pendingCount > 0 && (
              <span className="w-2 h-2 rounded-full bg-indigo-600 animate-ping absolute -top-0.5 -right-0.5" />
            )}
          </button>
        </div>

        {/* Search Bar */}
        <div className="relative">
          <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            type="text"
            placeholder="Search conversations..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-3 py-2 bg-white border border-slate-200 rounded-xl text-xs text-slate-700 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all shadow-sm"
          />
        </div>
      </div>

      {/* Friends & DMs List */}
      <div className="flex-1 overflow-y-auto p-3 space-y-1">
        {pendingCount > 0 && (
          <div
            onClick={onOpenFriendsModal}
            className="flex items-center justify-between p-3 mb-2 rounded-2xl bg-indigo-500 text-white shadow-md shadow-indigo-200 cursor-pointer hover:bg-indigo-600 transition-all"
          >
            <div className="flex items-center gap-2.5">
              <Users className="w-4 h-4" />
              <span className="text-xs font-semibold">Pending Requests</span>
            </div>
            <span className="px-2 py-0.5 rounded-full bg-white/20 text-white text-[11px] font-bold">
              {pendingCount} new
            </span>
          </div>
        )}

        <div className="px-2 py-1 text-[11px] font-bold uppercase tracking-wider text-slate-400">
          Direct Messages
        </div>

        {filteredFriends.length === 0 ? (
          <div className="py-12 text-center px-4">
            <Sparkles className="w-8 h-8 text-slate-300 mx-auto mb-2" />
            <p className="text-xs font-medium text-slate-500 mb-1">
              {searchQuery ? "No matching conversations" : "No friends added yet"}
            </p>
            <p className="text-[11px] text-slate-400 mb-3">
              Connect with friends to start messaging!
            </p>
            <button
              onClick={onOpenFriendsModal}
              className="px-4 py-2 bg-white border border-slate-200 text-slate-700 text-xs font-medium rounded-xl hover:bg-slate-50 transition-colors shadow-sm inline-flex items-center gap-1.5"
            >
              <UserPlus className="w-3.5 h-3.5 text-indigo-600" />
              Find Friends
            </button>
          </div>
        ) : (
          filteredFriends.map((friend) => {
            const isSelected = selectedUser?.id === friend.id;
            const lastMsg = friend.messages?.[friend.messages.length - 1]?.message || "No messages yet";
            const unread = friend.unreadCount || 0;
            const hasUnread = unread > 0;

            return (
              <button
                key={friend.id}
                onClick={() => setSelectedUser(friend)}
                className={`w-full flex items-center gap-3 p-3 rounded-2xl transition-all duration-150 text-left ${
                  isSelected
                    ? "bg-white shadow-sm border border-slate-200/80 text-slate-900"
                    : hasUnread
                    ? "bg-indigo-50/70 border border-indigo-100 text-slate-900 font-medium"
                    : "hover:bg-slate-100/80 text-slate-700"
                }`}
              >
                {/* Avatar */}
                <div className="relative shrink-0">
                  <img
                    src={getAvatarUrl(friend.username)}
                    alt={friend.username}
                    className="w-10 h-10 rounded-2xl object-cover border border-slate-200/60 shadow-sm"
                  />
                  <span className="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-emerald-500 border-2 border-white" />
                  {hasUnread && (
                    <span className="absolute -top-1 -right-1 w-3 h-3 rounded-full bg-indigo-600 animate-pulse border-2 border-white" />
                  )}
                </div>

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between mb-0.5">
                    <span className={`text-xs truncate ${hasUnread ? "font-bold text-indigo-950" : "font-semibold text-slate-800"}`}>
                      {friend.username}
                    </span>
                    {hasUnread ? (
                      <span className="px-1.5 py-0.5 rounded-full bg-indigo-600 text-white text-[10px] font-extrabold shadow-sm shrink-0">
                        {unread}
                      </span>
                    ) : (
                      <span className="text-[10px] text-slate-400">Online</span>
                    )}
                  </div>
                  <p className={`text-[11px] truncate ${hasUnread ? "font-bold text-indigo-900" : "text-slate-500"}`}>
                    {lastMsg}
                  </p>
                </div>
              </button>
            );
          })
        )}
      </div>
    </aside>
  );
};
