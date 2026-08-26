import { ChatLayout } from "@/components/chat/chat-layout";
import { jwtDecode } from "jwt-decode";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Home() {
  const navigate = useNavigate();
  const token = localStorage.getItem("token");

  useEffect(() => {
    if (!token) {
      navigate("/login");
      return;
    }

    try {
      const payload: { exp: number } = jwtDecode(token);
      if (payload.exp * 1000 < Date.now()) {
        localStorage.removeItem("token");
        navigate("/login");
      }
    } catch {
      localStorage.removeItem("token");
      navigate("/login");
    }
  }, [navigate, token]);

  if (!token) return null;

  return (
    <div className="w-screen h-screen overflow-hidden">
      <ChatLayout />
    </div>
  );
}
