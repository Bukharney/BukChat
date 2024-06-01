import { Label } from "@radix-ui/react-label";

import { Button } from "./ui/button";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "./ui/dialog";
import { Input } from "./ui/input";
import { useState } from "react";

export default function SidebarFooter() {
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
    }
  };
  return (
    <div className="h-6/12 flex justify-end">
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
  );
}
