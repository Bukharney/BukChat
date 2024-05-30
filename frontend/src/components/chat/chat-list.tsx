import { Friend, Message } from "@/data";
import { cn } from "@/lib/utils";
import React, { useRef } from "react";
import ChatBottombar from "./chat-bottombar";

interface ChatListProps {
  messages?: Message[];
  selectedUser: Friend;
  sendMessage: (newMessage: Message) => void;
}

export function ChatList({
  messages,
  selectedUser,
  sendMessage,
}: ChatListProps) {
  const messagesContainerRef = useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (messagesContainerRef.current) {
      messagesContainerRef.current.scrollTop =
        messagesContainerRef.current.scrollHeight;
    }
  }, [messages, selectedUser]);

  return (
    <div className="overflow-y-auto overflow-x-hidden h-full flex flex-col">
      <div
        ref={messagesContainerRef}
        className="w-full overflow-y-auto overflow-x-hidden h-full flex flex-col"
      >
        {messages?.map((message, index) => (
          <div
            key={index}
            className={cn(
              "flex flex-col gap-2 p-4 whitespace-pre-wrap",
              message.user_id !== selectedUser.id ? "items-end" : "items-start"
            )}
          >
            <div className="flex gap-3 items-center">
              <span className=" bg-accent p-3 rounded-md max-w-xs">
                {message.message}
              </span>
            </div>
          </div>
        ))}
      </div>
      <ChatBottombar sendMessage={sendMessage} />
    </div>
  );
}
