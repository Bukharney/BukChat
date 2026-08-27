import React, { useState, useRef } from "react";
import { SendHorizontal, Smile, Paperclip, X } from "lucide-react";
import { Message } from "@/data";
import { jwtDecode } from "jwt-decode";

interface MessageComposerProps {
  sendMessage: (newMessage: Message) => void;
  onTyping?: (isTyping: boolean) => void;
}

const EMOJIS = ["👍", "❤️", "😊", "🔥", "🎉", "🙌", "😂", "✨", "🚀", "💡", "😎", "👋"];

export const MessageComposer: React.FC<MessageComposerProps> = ({ sendMessage, onTyping }) => {
  const [text, setText] = useState("");
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const inputRef = useRef<HTMLTextAreaElement>(null);
  const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

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

  const notifyTyping = () => {
    if (!onTyping) return;
    onTyping(true);
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }
    typingTimeoutRef.current = setTimeout(() => {
      onTyping(false);
    }, 2000);
  };

  const handleSend = () => {
    if (!text.trim() && !imagePreview) return;

    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }
    if (onTyping) {
      onTyping(false);
    }

    let content = text.trim();
    if (imagePreview) {
      content = content ? `${content}\n[Image attached]` : "[Image attached]";
    }

    const newMessage: Message = {
      id: Date.now(),
      user_id: currentUserId,
      message: content,
      timestamp: new Date().toISOString(),
    };

    sendMessage(newMessage);
    setText("");
    setImagePreview(null);
    setShowEmojiPicker(false);

    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const addEmoji = (emoji: string) => {
    setText((prev) => prev + emoji);
  };

  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        setImagePreview(event.target?.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div className="p-4 bg-white border-t border-slate-200/80 relative">
      {/* Image Preview Thumbnail */}
      {imagePreview && (
        <div className="mb-3 relative inline-block">
          <img
            src={imagePreview}
            alt="Attachment preview"
            className="w-20 h-20 object-cover rounded-xl border border-slate-200 shadow-sm"
          />
          <button
            onClick={() => setImagePreview(null)}
            className="absolute -top-2 -right-2 p-1 bg-slate-800 text-white rounded-full hover:bg-slate-900 transition-colors shadow"
          >
            <X className="w-3 h-3" />
          </button>
        </div>
      )}

      {/* Emoji Picker Popover */}
      {showEmojiPicker && (
        <div className="absolute bottom-16 left-4 z-30 p-2.5 bg-white border border-slate-200 rounded-2xl shadow-xl grid grid-cols-6 gap-1.5 animate-in slide-in-from-bottom-2 duration-150">
          {EMOJIS.map((emoji) => (
            <button
              key={emoji}
              onClick={() => addEmoji(emoji)}
              className="w-8 h-8 text-lg flex items-center justify-center rounded-xl hover:bg-slate-100 transition-colors"
            >
              {emoji}
            </button>
          ))}
        </div>
      )}

      {/* Input Row */}
      <div className="flex items-center gap-2 bg-slate-50 border border-slate-200 rounded-2xl px-3 py-2 focus-within:ring-2 focus-within:ring-indigo-500/20 focus-within:border-indigo-500 transition-all">
        {/* Attachment Button */}
        <button
          type="button"
          onClick={() => fileInputRef.current?.click()}
          className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-white rounded-xl transition-colors shrink-0"
          title="Attach Image"
        >
          <Paperclip className="w-5 h-5" />
        </button>
        <input
          type="file"
          accept="image/*"
          ref={fileInputRef}
          onChange={handleImageSelect}
          className="hidden"
        />

        {/* Emoji Button */}
        <button
          type="button"
          onClick={() => setShowEmojiPicker((prev) => !prev)}
          className={`p-2 rounded-xl transition-colors shrink-0 ${
            showEmojiPicker
              ? "text-indigo-600 bg-indigo-50"
              : "text-slate-400 hover:text-indigo-600 hover:bg-white"
          }`}
          title="Pick Emoji"
        >
          <Smile className="w-5 h-5" />
        </button>

        {/* Text Area */}
        <textarea
          ref={inputRef}
          value={text}
          onChange={(e) => {
            setText(e.target.value);
            notifyTyping();
          }}
          onKeyDown={handleKeyDown}
          placeholder="Type your message..."
          rows={1}
          className="w-full bg-transparent border-0 focus:outline-none text-xs text-slate-800 placeholder:text-slate-400 resize-none py-1 min-h-[28px] max-h-32"
        />

        {/* Send Button */}
        <button
          onClick={handleSend}
          disabled={!text.trim() && !imagePreview}
          className="p-2.5 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all shadow-sm shrink-0"
          title="Send Message"
        >
          <SendHorizontal className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
};
