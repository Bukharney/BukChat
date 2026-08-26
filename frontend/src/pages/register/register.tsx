import { Form as RegisterForm } from "./form";
import { MessageCircleCode } from "lucide-react";

export default function RegisterPage() {
  return (
    <div className="h-screen w-screen flex justify-center items-center bg-slate-50 p-4 font-sans select-none">
      <div className="w-full max-w-md bg-white border border-slate-200/80 shadow-xl rounded-3xl p-8 md:p-10 space-y-6">
        <div className="flex flex-col items-center text-center space-y-2">
          <div className="w-12 h-12 rounded-2xl bg-indigo-600 flex items-center justify-center text-white shadow-lg shadow-indigo-200 mb-2">
            <MessageCircleCode className="w-7 h-7" />
          </div>
          <h1 className="font-bold text-2xl text-slate-900 tracking-tight">Create an Account</h1>
          <p className="text-xs text-slate-500">Join BukChat workspace to start messaging</p>
        </div>

        <RegisterForm />

        <div className="pt-2 text-center text-xs text-slate-500 border-t border-slate-100">
          Already have an account?{" "}
          <a href="/login" className="font-semibold text-indigo-600 hover:text-indigo-700 transition-colors">
            Sign in
          </a>
        </div>
      </div>
    </div>
  );
}
