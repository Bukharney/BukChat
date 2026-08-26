import { useEffect, useRef, useState } from "react";
import ChatTopbar from "./chat-topbar";
import { MessageFeed } from "./MessageFeed";
import { MessageComposer } from "./MessageComposer";
import { Friend, Message } from "@/data";
import useWebSocket from "react-use-websocket";
import { jwtDecode } from "jwt-decode";

interface ChatProps {
  messages?: Message[];
  selectedUser: Friend;
  onNewMessage?: (roomId: number, msg: Message) => void;
}

export function Chat({ selectedUser, messages, onNewMessage }: ChatProps) {
  const [messagesState, setMessages] = useState<Message[]>([]);
  const [isPartnerTyping, setIsPartnerTyping] = useState(false);
  const typingTimerRef = useRef<any>(null);

  const token = localStorage.getItem("token");
  let currentUserId = 0;
  if (token) {
    try {
      const decoded: any = jwtDecode(token);
      currentUserId = decoded.user_id || decoded.id || 0;
    } catch (e) {
      console.error("Failed to decode token", e);
    }
  }

  const id = selectedUser.room_id;
  const WS_URL = `ws://localhost:8080/ws/${id}?token=${token}`;

  const { sendJsonMessage } = useWebSocket(WS_URL, {
    onOpen: () => {
      console.log(`[Chat WS] Connected to room: ${id} (URL: ${WS_URL})`);
    },
    onMessage: (event) => {
      try {
        const data = JSON.parse(event.data);

        // Handle typing events
        if (data.type === "typing" && data.sender !== currentUserId) {
          const isTyping = Boolean(data.is_typing);
          setIsPartnerTyping(isTyping);

          if (typingTimerRef.current) {
            clearTimeout(typingTimerRef.current);
          }

          if (isTyping) {
            typingTimerRef.current = setTimeout(() => {
              setIsPartnerTyping(false);
            }, 3000);
          }
          return;
        }

        // Only process incoming messages sent by other users (non-zero and not self)
        if (data.sender && data.sender !== 0 && data.sender !== currentUserId) {
          setIsPartnerTyping(false);
          const incomingMsg: Message = {
            id: data.id || Date.now(),
            user_id: data.sender,
            message: data.content,
            timestamp: data.timestamp || new Date().toISOString(),
          };
          setMessages((prev) => [...prev, incomingMsg]);
          if (onNewMessage) {
            onNewMessage(id, incomingMsg);
          }
        }
      } catch (err) {
        console.error("[Chat WS] Failed to parse websocket message", err);
      }
    },
    onClose: (event) => {
      console.log(`[Chat WS] Closed connection for room ${id}`, event);
    },
    share: true,
    retryOnError: true,
    shouldReconnect: () => true,
  });

  const handleTyping = (isTyping: boolean) => {
    sendJsonMessage({
      type: "typing",
      is_typing: isTyping,
      id: id.toString(),
    });
  };

  const handleSendMessage = (msg: Message) => {
    setIsPartnerTyping(false);
    // Append locally for immediate optimistic UI update
    setMessages((prev) => [...prev, msg]);
    if (onNewMessage) {
      onNewMessage(id, msg);
    }

    const payload = {
      id: id.toString(),
      content: msg.message,
    };

    sendJsonMessage(payload);
  };

  useEffect(() => {
    setMessages(messages ?? []);
    setIsPartnerTyping(false);
  }, [messages, selectedUser.id]);

  return (
    <div className="flex flex-col w-full h-full bg-white select-none overflow-hidden">
      <ChatTopbar selectedUser={selectedUser} isPartnerTyping={isPartnerTyping} />
      <MessageFeed messages={messagesState} selectedUser={selectedUser} isPartnerTyping={isPartnerTyping} />
      <MessageComposer sendMessage={handleSendMessage} onTyping={handleTyping} />
    </div>
  );
}
