import ChatTopbar from "./chat-topbar";
import { ChatList } from "./chat-list";
import React from "react";
import { Message } from "@/data";
import { Friend } from "@/data";
interface ChatProps {
  messages?: Message[];
  selectedUser: Friend;
}

export function Chat({ selectedUser, messages }: ChatProps) {
  const [messagesState, setMessages] = React.useState<Message[]>(
    messages ?? []
  );

  React.useEffect(() => {
    setMessages(messages ?? []);
  }, [messages, selectedUser]);

  const sendMessage = (newMessage: Message) => {
    setMessages((prev) => [...prev, newMessage]);
  };

  return (
    <div className="flex flex-col justify-between w-full h-full">
      <ChatTopbar selectedUser={selectedUser} />

      <ChatList
        messages={messagesState}
        selectedUser={selectedUser}
        sendMessage={sendMessage}
      />
    </div>
  );
}
