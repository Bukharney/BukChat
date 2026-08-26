import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { User, Mail, Lock, AlertCircle } from "lucide-react";

export const Form = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);
    try {
      const res = await fetch("/v1/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username,
          password,
          email,
        }),
      });
      const data = await res.json();
      if (res.ok) {
        localStorage.setItem("token", data.access_token || data.token);
        navigate("/");
      } else {
        setError(data.message || data.error || "Failed to create account");
      }
    } catch {
      setError("Unable to connect to server");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={onSubmit} className="space-y-4 w-full">
      {error && (
        <div className="p-3 bg-rose-50 border border-rose-200/80 rounded-2xl flex items-center gap-2 text-xs text-rose-600 font-medium">
          <AlertCircle className="w-4 h-4 shrink-0" />
          <span>{error}</span>
        </div>
      )}

      <div className="space-y-1">
        <label className="text-xs font-semibold text-slate-700 block">Username</label>
        <div className="relative">
          <User className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            type="text"
            required
            placeholder="Choose username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
          />
        </div>
      </div>

      <div className="space-y-1">
        <label className="text-xs font-semibold text-slate-700 block">Email Address</label>
        <div className="relative">
          <Mail className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            type="email"
            required
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
          />
        </div>
      </div>

      <div className="space-y-1">
        <label className="text-xs font-semibold text-slate-700 block">Password</label>
        <div className="relative">
          <Lock className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            type="password"
            required
            placeholder="Create password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
          />
        </div>
      </div>

      <button
        type="submit"
        disabled={isLoading}
        className="w-full py-3 bg-indigo-600 text-white rounded-2xl text-xs font-semibold hover:bg-indigo-700 disabled:opacity-50 transition-all shadow-md shadow-indigo-100 mt-2"
      >
        {isLoading ? "Creating Account..." : "Create Account"}
      </button>
    </form>
  );
};
