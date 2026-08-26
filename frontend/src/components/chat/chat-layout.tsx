import { useEffect, useRef, useState } from "react";
import { NavSidebar } from "./NavSidebar";
import { ChannelSidebar } from "./ChannelSidebar";
import { FriendsModal } from "./FriendsModal";
import { SettingsModal } from "./SettingsModal";
import { Chat } from "./chat";
import { Friend, Message } from "@/data";
import { MessageSquareDashed, UserPlus } from "lucide-react";
import useWebSocket from "react-use-websocket";
import { apiFetch } from "@/lib/api";
import { jwtDecode } from "jwt-decode";

interface ChatLayoutProps {
  defaultLayout?: number[];
  defaultCollapsed?: boolean;
  navCollapsedSize?: number;
}

type FriendRequest = {
  id: number;
  username: string;
};

export function ChatLayout({ }: ChatLayoutProps) {
  const [friendList, setFriendList] = useState<Friend[]>([]);
  const [selectedUser, setSelectedUser] = useState<Friend | null>(null);
  const [friendRequestList, setFriendRequestList] = useState<FriendRequest[]>([]);
  const [activeTab, setActiveTab] = useState<"chats" | "friends" | "settings">("chats");
  const [isFriendsModalOpen, setIsFriendsModalOpen] = useState(false);
  const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);

  const selectedUserRef = useRef<Friend | null>(selectedUser);
  useEffect(() => {
    selectedUserRef.current = selectedUser;
  }, [selectedUser]);

  const handleSelectFriend = (friend: Friend) => {
    setSelectedUser({ ...friend, unreadCount: 0 });
    setFriendList((prevList) =>
      prevList.map((item) =>
        item.id === friend.id ? { ...item, unreadCount: 0 } : item
      )
    );
  };

  const token = localStorage.getItem("token");
  let currentUserId = 0;
  if (token) {
    try {
      const decoded: any = jwtDecode(token);
      currentUserId = decoded.user_id || decoded.id || 0;
    } catch (e) {
      console.error("Failed to decode token", e);
    }
  }

  const NOTIFICATION_WS_URL = token ? `ws://localhost:8080/ws/notifications?token=${token}` : null;

  // Real-time Notification WebSocket Listener
  useWebSocket(NOTIFICATION_WS_URL, {
    onMessage: (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "friend_request" || data.type === "friend_accepted" || data.type === "friend_rejected") {
          handleGetFriendRequests();
          handleGetFriends();
        } else if (data.type === "new_message") {
          const payload = data.payload || {};
          const roomId = Number(payload.room_id || data.id);
          const senderId = Number(data.sender || payload.sender);
          const content = payload.content || data.content;

          if (senderId && currentUserId && senderId === currentUserId) {
            return;
          }

          if (roomId && content) {
            const newMsg: Message = {
              id: Date.now(),
              user_id: senderId,
              message: content,
            };

            const isCurrentlyActive = selectedUserRef.current?.room_id === roomId;

            // If receiver is currently viewing this chat room, the room WebSocket in Chat component handles the message.
            if (isCurrentlyActive) {
              return;
            }

            setFriendList((prevList) =>
              prevList.map((item) =>
                item.room_id === roomId
                  ? {
                      ...item,
                      messages: [...(item.messages || []), newMsg],
                      unreadCount: isCurrentlyActive ? 0 : (item.unreadCount || 0) + 1,
                    }
                  : item
              )
            );

            setSelectedUser((prev) =>
              prev && prev.room_id === roomId
                ? { ...prev, messages: [...(prev.messages || []), newMsg] }
                : prev
            );
          }
        }
      } catch (err) {
        console.error("Failed to parse notification event", err);
      }
    },
    share: true,
    shouldReconnect: () => true,
  });

  const handleGetFriendRequests = async () => {
    try {
      const res = await apiFetch("/v1/users/friends-request");
      if (res.ok) {
        const data = await res.json();
        setFriendRequestList(Array.isArray(data) ? data : []);
      }
    } catch (err) {
      console.error("Failed to fetch friend requests", err);
    }
  };

  const handleGetFriends = async () => {
    try {
      const res = await apiFetch("/v1/users/friends");
      if (res.ok) {
        const data: Friend[] = await res.json();
        if (Array.isArray(data)) {
          const friendsWithMessages = await Promise.all(
            data.map(async (friend) => {
              try {
                const msgRes = await apiFetch(`/v1/chat/${friend.room_id}`);
                if (msgRes.ok) {
                  const msgs = await msgRes.json();
                  return { ...friend, messages: msgs };
                }
              } catch (e) {
                console.error("Failed to fetch messages for friend", friend.id, e);
              }
              return friend;
            })
          );
          setFriendList(friendsWithMessages);
          if (!selectedUser && friendsWithMessages.length > 0) {
            const sorted = [...friendsWithMessages].sort((a, b) => {
              const aLastMsg = a.messages?.[a.messages.length - 1];
              const bLastMsg = b.messages?.[b.messages.length - 1];
              const aTime = aLastMsg?.timestamp ? new Date(aLastMsg.timestamp).getTime() : (aLastMsg?.id || 0);
              const bTime = bLastMsg?.timestamp ? new Date(bLastMsg.timestamp).getTime() : (bLastMsg?.id || 0);
              return bTime - aTime;
            });
            setSelectedUser(sorted[0]);
          }
        }
      }
    } catch (err) {
      console.error("Failed to fetch friends", err);
    }
  };

  const handleAddFriend = async (username: string) => {
    const res = await apiFetch("/v1/users/add-friend", {
      method: "POST",
      body: JSON.stringify({ username }),
    });

    if (!res.ok) {
      const errData = await res.json();
      throw new Error(errData.message || "Failed to send request");
    }

    handleGetFriendRequests();
    handleGetFriends();
  };

  const handleAcceptFriend = async (username: string) => {
    const res = await apiFetch("/v1/users/add-friend", {
      method: "POST",
      body: JSON.stringify({ username }),
    });
    if (res.ok) {
      handleGetFriendRequests();
      handleGetFriends();
    }
  };

  const handleRejectFriend = async (username: string) => {
    const res = await apiFetch("/v1/users/reject-friend", {
      method: "POST",
      body: JSON.stringify({ username }),
    });
    if (res.ok) {
      handleGetFriendRequests();
    }
  };

  useEffect(() => {
    handleGetFriendRequests();
    handleGetFriends();
  }, []);

  useEffect(() => {
    if (activeTab === "friends") {
      setIsFriendsModalOpen(true);
      setIsSettingsModalOpen(false);
    } else if (activeTab === "settings") {
      setIsSettingsModalOpen(true);
      setIsFriendsModalOpen(false);
    }
  }, [activeTab]);

  useEffect(() => {
    const handleGetMessages = async (roomId: number, friendId: number) => {
      try {
        const res = await apiFetch(`/v1/chat/${roomId}`);
        if (res.ok) {
          const messagesData = await res.json();
          setFriendList((prevList) =>
            prevList.map((item) =>
              item.id === friendId ? { ...item, messages: messagesData } : item
            )
          );
          if (selectedUser && selectedUser.id === friendId) {
            setSelectedUser((prev) => (prev ? { ...prev, messages: messagesData } : prev));
          }
        }
      } catch (err) {
        console.error("Failed to fetch room messages", err);
      }
    };

    if (selectedUser) {
      handleGetMessages(selectedUser.room_id, selectedUser.id);
    }
  }, [selectedUser?.id]);

  const handleNewWsMessage = (roomId: number, msg: Message) => {
    setFriendList((prevList) =>
      prevList.map((item) =>
        item.room_id === roomId
          ? { ...item, messages: [...(item.messages || []), msg] }
          : item
      )
    );
    setSelectedUser((prev) =>
      prev && prev.room_id === roomId
        ? { ...prev, messages: [...(prev.messages || []), msg] }
        : prev
    );
  };

  return (
    <div className="flex h-screen w-screen bg-slate-100 overflow-hidden font-sans antialiased">
      {/* Primary Rail Navigation */}
      <NavSidebar
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        pendingCount={friendRequestList.length}
      />

      {/* Secondary Channel / Direct Message Sidebar */}
      <ChannelSidebar
        friendList={friendList}
        selectedUser={selectedUser}
        setSelectedUser={handleSelectFriend}
        pendingCount={friendRequestList.length}
        onOpenFriendsModal={() => setIsFriendsModalOpen(true)}
      />

      {/* Main Canvas / Chat Room View */}
      <main className="flex-1 h-full flex flex-col bg-white overflow-hidden">
        {selectedUser ? (
          <Chat
            selectedUser={selectedUser}
            messages={selectedUser.messages}
            onNewMessage={handleNewWsMessage}
          />
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center p-8 text-center bg-slate-50/50 select-none">
            <div className="w-16 h-16 rounded-3xl bg-indigo-50 text-indigo-600 flex items-center justify-center mb-4 shadow-sm">
              <MessageSquareDashed className="w-8 h-8" />
            </div>
            <h3 className="text-xl font-bold text-slate-800 mb-2">No Active Conversation Selected</h3>
            <p className="text-xs text-slate-500 max-w-sm mb-6 leading-relaxed">
              Choose a friend from your sidebar or add new connections to start chatting in real time!
            </p>
            <button
              onClick={() => setIsFriendsModalOpen(true)}
              className="px-5 py-2.5 bg-indigo-600 text-white rounded-2xl text-xs font-semibold hover:bg-indigo-700 transition-all shadow-md shadow-indigo-100 flex items-center gap-2"
            >
              <UserPlus className="w-4 h-4" />
              Manage Friends & Requests
            </button>
          </div>
        )}
      </main>

      {/* Friends Management Modal */}
      <FriendsModal
        isOpen={isFriendsModalOpen}
        onClose={() => {
          setIsFriendsModalOpen(false);
          setActiveTab("chats");
        }}
        friendRequests={friendRequestList}
        friendList={friendList}
        onAcceptFriend={handleAcceptFriend}
        onRejectFriend={handleRejectFriend}
        onAddFriend={handleAddFriend}
      />

      {/* Settings & Account Management Modal */}
      <SettingsModal
        isOpen={isSettingsModalOpen}
        onClose={() => {
          setIsSettingsModalOpen(false);
          setActiveTab("chats");
        }}
      />
    </div>
  );
}
