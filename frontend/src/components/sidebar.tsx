"use client";

import { Friend, Message } from "../data";

interface SidebarProps {
  links: {
    id: number;
    username: string;
    messages: Message[];
    room_id: number;
  }[];

  setSelectedUser: React.Dispatch<React.SetStateAction<Friend | null>>;
}

export function Sidebar({ links, setSelectedUser }: SidebarProps) {
  return (
    <div className="relative group flex flex-col h-full gap-4 p-2">
      <div className="flex justify-between p-2 items-center">
        <div className="flex gap-2 items-center text-2xl">
          <p className="font-medium">Chats</p>
          <span className="text-zinc-300">({links.length})</span>
        </div>
      </div>

      <nav className="grid gap-1 px-2">
        {links.map((link, index) => (
          <button
            onClick={(e) => {
              e.preventDefault();
              setSelectedUser(link);
            }}
            className="flex items-center gap-2 p-2 rounded-lg hover:bg-zinc-100 transition-colors"
            key={index}
          >
            <div className="flex flex-col max-w-28">
              <span>{link.username}</span>
            </div>
          </button>
        ))}
      </nav>
    </div>
  );
}
