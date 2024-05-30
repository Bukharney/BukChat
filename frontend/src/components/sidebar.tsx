import { useEffect, useState } from "react";
import { Friend, Message } from "../data";
import Footer from "./sidebarFooter";
import React from "react";
import { Button } from "./ui/button";

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

type FriendRequest = {
  id: number;
  username: string;
};
export function Sidebar({
  links,
  setSelectedUser,
  selectedRoom,
}: SidebarProps) {
  const [friendRequestList, setFriendRequestList] = useState<FriendRequest[]>(
    []
  );
  const handleGetFriendRequests = async () => {
    await fetch("/v1/users/friends-request", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    })
      .then((res) => res.json())
      .then((data) => {
        if (!data) {
          setFriendRequestList([]);
        }
        setFriendRequestList(data);
        console.log(data);
      });
  };

  const handleAcceptFriend = async (username: string) => {
    await fetch("/v1/users/add-friend", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify({ username: username }),
    })
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
      });
  };

  useEffect(() => {
    handleGetFriendRequests();
  }, []);

  return (
    <div className="flex flex-col h-full justify-between">
      <div className="flex flex-col h-6/12 gap-4 p-2">
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
      <div className="flex flex-col h-6/12 p-2 overflow-auto">
        <div className="flex justify-between p-2 items-centerbg-slate-300">
          <div className="flex gap-2 items-center text-2xl">
            <p className="font-medium">Pending requests</p>
            <span className="text-zinc-300">
              ({(friendRequestList && friendRequestList.length) || 0})
            </span>
          </div>
        </div>

        <div className="grid gap-1 px-2 overflow-auto h-48 items-start ">
          {!friendRequestList ? (
            <div className="flex justify-center items-center h-full">
              <p className="text-center font-medium text-zinc-500">
                No pending requests
              </p>
            </div>
          ) : (
            <>
              {friendRequestList &&
                friendRequestList.map((link, index) => (
                  <div
                    className="flex justify-between p-2 items-center "
                    key={index}
                  >
                    <div className="flex flex-col max-w-28">
                      <span>{link.username}</span>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        onClick={(e) => {
                          e.preventDefault();
                          console.log("Accepting friend request");
                        }}
                        variant={"outline"}
                        size={"sm"}
                      >
                        Reject
                      </Button>
                      <Button
                        size={"sm"}
                        onClick={(e) => {
                          e.preventDefault();
                          handleAcceptFriend(link.username);
                          handleGetFriendRequests();
                          console.log("Accepting friend request");
                        }}
                      >
                        Accept
                      </Button>
                    </div>
                  </div>
                ))}
            </>
          )}
        </div>
      </div>
      <Footer />
    </div>
  );
}
