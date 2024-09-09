import { useEffect, useState } from "react";
import { Sidebar } from "@/components/sidebar";
import { Chat } from "./chat";
import { Friend } from "@/data";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import LogoutPage from "../sidebarFooter";
import { Button } from "../ui/button";
import Footer from "../sidebarFooter";

interface ChatLayoutProps {
  defaultLayout: number[] | undefined;
  defaultCollapsed?: boolean;
  navCollapsedSize: number;
}

type FriendRequest = {
  id: number;
  username: string;
};

export function ChatLayout({
  defaultLayout = [],
  navCollapsedSize,
}: ChatLayoutProps) {
  const [friendList, setFriendList] = useState<Friend[]>([]);
  const [selectedUser, setSelectedUser] = useState<Friend | null>(null);
  const [friendRequestList, setFriendRequestList] = useState<FriendRequest[]>(
    []
  );

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

  useEffect(() => {
    handleGetFriendRequests();
  }, []);

  useEffect(() => {
    const handleGetFriend = async () => {
      await fetch("/v1/users/friends", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      })
        .then((res) => res.json())
        .then((data) => setFriendList(data));
    };

    handleGetFriend();
  }, []);

  useEffect(() => {
    const handleGetMessages = async (id: number, userId: number) => {
      const res = await fetch(`/v1/chat/${id}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      })
        .then((res) => res.json())
        .then((data) => {
          return data;
        });

      if (res && friendList) {
        const user = friendList.find((user) => user.id === userId);
        if (user) {
          user.messages = res;
          setSelectedUser(user);
        }
      }
    };

    if (selectedUser) {
      handleGetMessages(selectedUser.room_id, selectedUser.id);
    }
  }, [selectedUser, friendList]);

  return (
    <ResizablePanelGroup
      direction="horizontal"
      onLayout={(sizes: number[]) => {
        document.cookie = `react-resizable-panels:layout=${JSON.stringify(
          sizes
        )}`;
      }}
      className="h-full items-stretch w-[600px]"
    >
      <ResizablePanel
        defaultSize={defaultLayout[0]}
        collapsedSize={navCollapsedSize}
        minSize={30}
      >
        {friendList.length > 0 ? (
          <div className="flex flex-col h-full justify-between">
            <Sidebar
              links={friendList.map((user) => ({
                id: user.id,
                username: user.username,
                messages: user.messages ?? [],
                room_id: user.room_id,
              }))}
              setSelectedUser={setSelectedUser}
              selectedRoom={selectedUser}
            />
            <PendingReq
              friendRequestList={friendRequestList}
              handleAcceptFriend={handleAcceptFriend}
              handleGetFriendRequests={handleGetFriendRequests}
            />
            <Footer />
          </div>
        ) : (
          <div className="flex flex-col h-full justify-between">
            <PendingReq
              friendRequestList={friendRequestList}
              handleAcceptFriend={handleAcceptFriend}
              handleGetFriendRequests={handleGetFriendRequests}
            />
            <div className="flex flex-col h-full justify-between">
              <div className="flex flex-col items-center justify-center h-full">
                <p className="text-center text-lg font-medium text-zinc-500">
                  You have no friends yet
                </p>
              </div>
              <LogoutPage />
            </div>
          </div>
        )}
      </ResizablePanel>
      <ResizableHandle withHandle />
      <ResizablePanel defaultSize={defaultLayout[1]} minSize={70}>
        {selectedUser ? (
          <Chat messages={selectedUser.messages} selectedUser={selectedUser} />
        ) : (
          <div className="flex flex-col h-full justify-between">
            <div className="flex flex-col items-center justify-center h-full">
              <p className="text-center text-lg font-medium text-zinc-500">
                Select a friend to chat with
              </p>
            </div>
          </div>
        )}
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

function PendingReq({
  friendRequestList,
  handleAcceptFriend,
  handleGetFriendRequests,
}: {
  friendRequestList: FriendRequest[];
  handleAcceptFriend: (username: string) => void;
  handleGetFriendRequests: () => void;
}) {
  return (
    <div className="flex flex-col h-6/12 p-2 overflow-auto">
      <div className="flex justify-between p-2 items-centerbg-slate-300">
        <div className="flex gap-2 items-center text-2xl">
          <p className="font-medium">Pending requests</p>
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
  );
}
