import ChatTopbar from "./chat-topbar";
import { ChatList } from "./chat-list";
import React, { useEffect } from "react";
import { Message } from "@/data";
import { Friend } from "@/data";
import useWebSocket, { ReadyState } from "react-use-websocket";

interface ChatProps {
  messages?: Message[];
  selectedUser: Friend;
}

interface Chat {
  id: string;
  content: string;
}

export function Chat({ selectedUser, messages }: ChatProps) {
  const [messagesState, setMessages] = React.useState<Message[]>([]);
  const [newMessage, setNewMessage] = React.useState<Chat | null>(null);

  const token = localStorage.getItem("token");
  const id = selectedUser.room_id;
  const host = "gigachat.bukharney.tech";
  const WS_URL = `wss://${host}/ws/${id}?token=${token}`;
  console.log(messages);

  const { sendJsonMessage, readyState } = useWebSocket(WS_URL, {
    onOpen: () => {
      console.log("WebSocket connection established." + id);
    },

    onMessage: (event) => {
      const data = JSON.parse(event.data);
      const message: Message = {
        id: data.id,
        user_id: data.sender,
        message: data.content,
      };
      console.log(message);
      if (data.sender != 0) {
        setMessages((prev) => [...prev, message]);
      }
    },

    onClose: () => {
      setNewMessage(null);
      console.log("WebSocket connection closed.");
    },
    share: true,
    filter: () => false,
    retryOnError: true,
    shouldReconnect: () => true,
  });

  useEffect(() => {
    if (newMessage && readyState === ReadyState.OPEN) {
      sendJsonMessage(newMessage);
      setNewMessage(null);
    }
  }, [sendJsonMessage, readyState, newMessage]);

  const sendMessage = (newMessage: Message) => {
    setNewMessage({
      id: id.toString(),
      content: newMessage.message,
    });
  };

  React.useEffect(() => {
    setMessages(messages ?? []);
  }, [messages, selectedUser.messages]);

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
