import { Friend } from "@/data";

interface ChatTopbarProps {
  selectedUser: Friend;
}

export default function ChatTopbar({ selectedUser }: ChatTopbarProps) {
  return (
    <div className="w-full h-20 flex p-4 justify-between items-center border-b">
      <div className="flex items-center gap-2">
        <div className="flex flex-col">
          <span className="font-medium">{selectedUser.username}</span>
        </div>
      </div>
    </div>
  );
}
