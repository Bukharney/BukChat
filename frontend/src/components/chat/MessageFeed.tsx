import React, { useEffect, useRef } from "react";
import { Friend, Message } from "@/data";
import { MessageSquare } from "lucide-react";
import { getAvatarUrl } from "@/lib/avatar";
import { jwtDecode } from "jwt-decode";

interface MessageFeedProps {
  messages?: Message[];
  selectedUser: Friend;
  isPartnerTyping?: boolean;
}

export const MessageFeed: React.FC<MessageFeedProps> = ({ messages = [], selectedUser, isPartnerTyping }) => {
  const containerRef = useRef<HTMLDivElement>(null);

  const token = localStorage.getItem("token");
  let currentUserId = 0;
  if (token) {
    try {
      const decoded = jwtDecode<{ user_id?: number; id?: number }>(token);
      currentUserId = decoded.user_id || decoded.id || 0;
    } catch (e) {
      console.error("Failed to decode token", e);
    }
  }

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, [messages, selectedUser, isPartnerTyping]);

  return (
    <div ref={containerRef} className="flex-1 overflow-y-auto p-4 md:p-6 space-y-4 bg-slate-50/40">
      {messages.length === 0 ? (
        <div className="h-full flex flex-col items-center justify-center text-center p-6 select-none">
          <div className="w-12 h-12 rounded-3xl bg-indigo-50 text-indigo-600 flex items-center justify-center mb-3">
            <MessageSquare className="w-6 h-6" />
          </div>
          <h4 className="text-sm font-semibold text-slate-800 mb-1">
            Start a conversation with {selectedUser.username}
          </h4>
          <p className="text-xs text-slate-400 max-w-xs">
            Send a friendly message to kick off your conversation!
          </p>
        </div>
      ) : (
        messages.map((msg, index) => {
          const isCurrentUser = currentUserId !== 0 ? msg.user_id === currentUserId : msg.user_id !== selectedUser.id;

          return (
            <div
              key={msg.id || index}
              className={`flex items-end gap-2.5 ${isCurrentUser ? "justify-end" : "justify-start"}`}
            >
              {!isCurrentUser && (
                <img
                  src={getAvatarUrl(selectedUser.username)}
                  alt={selectedUser.username}
                  className="w-8 h-8 rounded-xl object-cover border border-slate-200/60 shrink-0 mb-1 shadow-sm"
                />
              )}

              <div
                className={`max-w-[75%] md:max-w-[65%] px-4 py-2.5 rounded-2xl text-xs leading-relaxed shadow-sm ${
                  isCurrentUser
                    ? "bg-indigo-600 text-white rounded-br-none"
                    : "bg-white border border-slate-200/80 text-slate-800 rounded-bl-none"
                }`}
              >
                <p className="whitespace-pre-wrap break-words">{msg.message}</p>
              </div>

              {isCurrentUser && (
                <img
                  src={getAvatarUrl("Me")}
                  alt="Me"
                  className="w-8 h-8 rounded-xl object-cover border border-indigo-200 shrink-0 mb-1 shadow-sm"
                />
              )}
            </div>
          );
        })
      )}

      {/* Partner Typing Indicator Bubble */}
      {isPartnerTyping && (
        <div className="flex items-end gap-2.5 justify-start animate-in fade-in duration-200">
          <img
            src={getAvatarUrl(selectedUser.username)}
            alt={selectedUser.username}
            className="w-8 h-8 rounded-xl object-cover border border-slate-200/60 shrink-0 mb-1 shadow-sm"
          />
          <div className="px-4 py-3 bg-white border border-slate-200/80 rounded-2xl rounded-bl-none shadow-sm flex items-center gap-1.5">
            <span className="w-1.5 h-1.5 rounded-full bg-slate-400 animate-bounce [animation-delay:-0.3s]" />
            <span className="w-1.5 h-1.5 rounded-full bg-slate-400 animate-bounce [animation-delay:-0.15s]" />
            <span className="w-1.5 h-1.5 rounded-full bg-slate-400 animate-bounce" />
          </div>
        </div>
      )}
    </div>
  );
};
