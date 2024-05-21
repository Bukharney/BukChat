import { SendHorizontal } from "lucide-react";
import React, { useRef, useState } from "react";
import { Textarea } from "@/components/ui/textarea";
import { Message } from "@/data";

interface ChatBottombarProps {
  sendMessage: (newMessage: Message) => void;
}

export default function ChatBottombar({ sendMessage }: ChatBottombarProps) {
  const [message, setMessage] = useState("");
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const handleInputChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(event.target.value);
  };

  const handleSend = () => {
    if (message.trim()) {
      const newMessage: Message = {
        id: message.length + 1,
        user_id: 0,
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
