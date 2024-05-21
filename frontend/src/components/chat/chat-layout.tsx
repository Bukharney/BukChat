import { useEffect, useState } from "react";
import { Sidebar } from "@/components/sidebar";
import { Chat } from "./chat";
import { Friend } from "@/data";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";

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

  useEffect(() => {
    handleGetFriend();
  }, []);

  useEffect(() => {
    if (users && users.length > 0) {
      setSelectedUser(users[0]);
    }
  }, [users]);

  useEffect(() => {
    const handleGetMessages = async (id: number, userId: number) => {
      await fetch(`/v1/chat/${id}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      })
        .then((res) => res.json())
        .then((data) => {
          if (!data || !users) {
            return;
          }
          const user = users.find((user) => user.id === userId);
          if (user) {
            user.messages = data;
            setSelectedUser(user);
          }
        });
    };

    if (!selectedUser || !users) {
      return;
    }

    if (selectedUser.id && selectedUser.room_id) {
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
        {users && users.length > 0 && (
          <Sidebar
            links={users.map((user) => ({
              id: user.id,
              username: user.username,
              messages: user.messages ?? [],
              room_id: user.room_id,
            }))}
            setSelectedUser={setSelectedUser}
          />
        )}
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
