import { ChatLayout } from "@/components/chat/chat-layout";
import { jwtDecode } from "jwt-decode";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Home() {
  const nevigate = useNavigate();
  const token = localStorage.getItem("token");
  if (!token) {
    window.location.href = "/login";
  }

  useEffect(() => {
    const checkToken = (token: string) => {
      const payload: { role: string; username: string; exp: number } =
        jwtDecode(token);
      if (payload.exp * 1000 < Date.now()) {
        localStorage.removeItem("token");
        nevigate("/login");
      }
    };

    if (token) checkToken(token);
  }, [nevigate, token]);

  return (
    <div className="flex h-[calc(100dvh)] flex-col items-center justify-center p-4 md:px-24 py-32 gap-4 ">
      <div className="z-10 border rounded-lg max-w-5xl w-full h-full text-sm">
        <ChatLayout defaultLayout={[]} navCollapsedSize={1} />
      </div>
    </div>
  );
}
