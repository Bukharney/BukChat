import { Friend, Message } from "../data";
import { Button } from "./ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { useState } from "react";

interface SidebarProps {
  links: {
    id: number;
    username: string;
    messages: Message[];
    room_id: number;
  }[];
  selectedRoom: Friend | null;
  setSelectedUser: React.Dispatch<React.SetStateAction<Friend | null>>;
}

export function Sidebar({
  links,
  setSelectedUser,
  selectedRoom,
}: SidebarProps) {
  const [username, setUsername] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const handleAddFriend = async (username: string) => {
    try {
      const res = await fetch("/v1/users/add-friend", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ username }),
      });
      if (res.ok) {
        await res.json().then((data) => {
          setError(null);
          console.log(data);
          window.location.reload();
        });
      } else {
        await res.json().then((data) => {
          setError(data.error);
        });
      }
    } catch (error) {
      setError("An error occurred");
      console.log("An error occurred");
    }
  };

  return (
    <>
      <div className="relative group flex flex-col h-4/6 gap-4 p-2 ">
        <div className="flex justify-between p-2 items-center ">
          <div className="flex gap-2 items-center text-2xl">
            <p className="font-medium">Chats</p>
            <span className="text-zinc-300">({links.length})</span>
          </div>
        </div>

        {links.length > 0 && (
          <nav className="grid gap-1 px-2">
            {links.map((link, index) => (
              <button
                onClick={(e) => {
                  e.preventDefault();
                  if (link.room_id == selectedRoom?.room_id) {
                    return;
                  }
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
        )}
      </div>
      <div className="h-2/6 flex justify-end">
        <div className="flex w-full items-end p-4 gap-2 end-0 justify-between">
          <Button
            onClick={() => {
              localStorage.removeItem("token");
              window.location.reload();
            }}
            variant="destructive"
          >
            <div className="flex flex-col max-w-28">
              <span>Logout</span>
            </div>
          </Button>

          <Dialog>
            <DialogTrigger asChild>
              <Button variant="outline">Add friend</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>Add frined</DialogTitle>
                <DialogDescription>
                  Add a friend by entering their username below and clicking the
                  Add button.
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="username" className="text-right">
                    Username
                  </Label>
                  <Input
                    id="username"
                    className="col-span-3"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                  />
                </div>
              </div>
              {error && (
                <div className="p-2 bg-red-100 text-red-500 rounded-lg">
                  {error}
                </div>
              )}
              <DialogFooter>
                <Button
                  onClick={(e: { preventDefault: () => void }) => {
                    e.preventDefault();
                    console.log(username);
                    handleAddFriend(username);
                  }}
                >
                  Add Friend
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </>
  );
}