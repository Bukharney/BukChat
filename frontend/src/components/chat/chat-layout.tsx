import { useEffect, useState } from "react";
import { Sidebar } from "@/components/sidebar";
import { Chat } from "./chat";
import { Friend } from "@/data";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import LogoutPage from "../logout";

interface ChatLayoutProps {
  defaultLayout: number[] | undefined;
  defaultCollapsed?: boolean;
  navCollapsedSize: number;
}

export function ChatLayout({
  defaultLayout = [],
  navCollapsedSize,
}: ChatLayoutProps) {
  const [users, setUsers] = useState<Friend[] | null>(null);
  const [selectedUser, setSelectedUser] = useState<Friend | null>(null);

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
        .then((data) => setUsers(data));
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

      if (res && users) {
        const user = users.find((user) => user.id === userId);
        if (user) {
          user.messages = res;
          setSelectedUser(user);
        }
      }
    };

    if (selectedUser) {
      handleGetMessages(selectedUser.room_id, selectedUser.id);
    }
  }, [selectedUser, users]);

  return (
    <ResizablePanelGroup
      direction="horizontal"
      onLayout={(sizes: number[]) => {
        document.cookie = `react-resizable-panels:layout=${JSON.stringify(
          sizes
        )}`;
      }}
      className="h-full items-stretch"
    >
      <ResizablePanel
        defaultSize={defaultLayout[0]}
        collapsedSize={navCollapsedSize}
        minSize={30}
      >
        {users &&
          (users.length > 0 ? (
            <Sidebar
              links={users.map((user) => ({
                id: user.id,
                username: user.username,
                messages: user.messages ?? [],
                room_id: user.room_id,
              }))}
              setSelectedUser={setSelectedUser}
              selectedRoom={selectedUser}
            />
          ) : (
            <>
              <div className="flex flex-col items-center justify-center h-4/6">
                <p className="text-center text-lg font-medium text-zinc-500">
                  You have no friends yet
                </p>
              </div>
              <LogoutPage />
            </>
          ))}
      </ResizablePanel>
      <ResizableHandle withHandle />
      <ResizablePanel defaultSize={defaultLayout[1]} minSize={70}>
        {selectedUser && (
          <Chat messages={selectedUser.messages} selectedUser={selectedUser} />
        )}
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}
