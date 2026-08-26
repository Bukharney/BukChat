import { SendHorizontal } from "lucide-react";
import React, { useRef, useState } from "react";
import { Textarea } from "@/components/ui/textarea";
import { Message } from "@/data";
import { jwtDecode } from "jwt-decode";

interface ChatBottombarProps {
  sendMessage: (newMessage: Message) => void;
}

export default function ChatBottombar({ sendMessage }: ChatBottombarProps) {
  const [message, setMessage] = useState("");
  const inputRef = useRef<HTMLTextAreaElement>(null);

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

  const handleInputChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(event.target.value);
  };

  const handleSend = () => {
    if (message.trim()) {
      const newMessage: Message = {
        id: message.length + 1,
        user_id: currentUserId,
        message: message.trim(),
      };
      sendMessage(newMessage);
      setMessage("");

      if (inputRef.current) {
        inputRef.current.focus();
      }
    }
  };

  const handleKeyPress = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      handleSend();
    }

    if (event.key === "Enter" && event.shiftKey) {
      event.preventDefault();
      setMessage((prev) => prev + "\n");
    }
  };

  return (
    <div className="p-2 flex justify-between w-full items-center gap-2">
      <Textarea
        autoComplete="off"
        value={message}
        ref={inputRef}
        onKeyDown={handleKeyPress}
        onChange={handleInputChange}
        name="message"
        placeholder="Aa"
        className=" w-full border rounded-full flex items-center h-2 resize-none overflow-hidden bg-background"
      ></Textarea>
      <a onClick={handleSend}>
        <SendHorizontal size={20} className="text-muted-foreground" />
      </a>
    </div>
  );
}
