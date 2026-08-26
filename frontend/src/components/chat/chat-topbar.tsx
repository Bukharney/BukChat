import { Friend } from "@/data";
import { Phone, Video, MoreVertical } from "lucide-react";
import { getAvatarUrl } from "@/lib/avatar";

interface ChatTopbarProps {
  selectedUser: Friend;
  isPartnerTyping?: boolean;
}

export default function ChatTopbar({ selectedUser, isPartnerTyping }: ChatTopbarProps) {
  return (
    <div className="w-full h-16 px-6 bg-white border-b border-slate-200/80 flex items-center justify-between shrink-0 shadow-sm">
      <div className="flex items-center gap-3">
        <div className="relative">
          <img
            src={getAvatarUrl(selectedUser.username)}
            alt={selectedUser.username}
            className="w-10 h-10 rounded-2xl object-cover border border-slate-200/60 shadow-sm"
          />
          <span
            className={`absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-white ${
              selectedUser.is_online !== false ? "bg-emerald-500" : "bg-slate-300"
            }`}
          />
        </div>
        <div>
          <h3 className="text-sm font-bold text-slate-800 tracking-tight">
            {selectedUser.username}
          </h3>
          {isPartnerTyping ? (
            <span className="text-[11px] text-indigo-600 font-semibold flex items-center gap-1.5 animate-pulse">
              <span className="w-1.5 h-1.5 rounded-full bg-indigo-500 inline-block animate-ping" />
              typing...
            </span>
          ) : selectedUser.is_online !== false ? (
            <span className="text-[11px] text-emerald-600 font-medium flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 inline-block" />
              Active Now
            </span>
          ) : (
            <span className="text-[11px] text-slate-400 font-medium flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-slate-300 inline-block" />
              Offline
            </span>
          )}
        </div>
      </div>

      <div className="flex items-center gap-1">
        <button
          className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-colors"
          title="Audio Call"
        >
          <Phone className="w-4 h-4" />
        </button>
        <button
          className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-colors"
          title="Video Call"
        >
          <Video className="w-4 h-4" />
        </button>
        <div className="w-px h-5 bg-slate-200 mx-1" />
        <button
          className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-colors"
          title="More Details"
        >
          <MoreVertical className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
