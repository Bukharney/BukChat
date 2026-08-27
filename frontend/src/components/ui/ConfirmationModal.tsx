import React, { useEffect } from "react";
import { AlertTriangle, Info, Loader2, X } from "lucide-react";

export interface ConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  description: string;
  confirmText?: string;
  cancelText?: string;
  variant?: "danger" | "warning" | "info";
  isLoading?: boolean;
  icon?: React.ReactNode;
}

export const ConfirmationModal: React.FC<ConfirmationModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  description,
  confirmText = "Confirm",
  cancelText = "Cancel",
  variant = "info",
  isLoading = false,
  icon,
}) => {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && isOpen && !isLoading) {
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [isOpen, isLoading, onClose]);

  if (!isOpen) return null;

  const getVariantStyles = () => {
    switch (variant) {
      case "danger":
        return {
          iconBg: "bg-rose-100 text-rose-600 border-rose-200",
          confirmBtn: "bg-rose-600 hover:bg-rose-700 text-white shadow-rose-200",
          accentColor: "text-rose-600",
        };
      case "warning":
        return {
          iconBg: "bg-amber-100 text-amber-600 border-amber-200",
          confirmBtn: "bg-amber-600 hover:bg-amber-700 text-white shadow-amber-200",
          accentColor: "text-amber-600",
        };
      case "info":
      default:
        return {
          iconBg: "bg-indigo-100 text-indigo-600 border-indigo-200",
          confirmBtn: "bg-slate-900 hover:bg-slate-800 text-white shadow-slate-200",
          accentColor: "text-indigo-600",
        };
    }
  };

  const styles = getVariantStyles();

  const defaultIcon =
    variant === "danger" ? (
      <AlertTriangle className="w-6 h-6 text-rose-600" />
    ) : (
      <Info className="w-6 h-6 text-indigo-600" />
    );

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center bg-slate-950/40 backdrop-blur-md p-4 animate-in fade-in duration-200 select-none">
      <div className="bg-white rounded-3xl shadow-2xl border border-slate-200/80 w-full max-w-md overflow-hidden transform animate-in zoom-in-95 duration-200 p-6 flex flex-col items-center text-center space-y-4">
        {/* Top Action Close */}
        <button
          onClick={onClose}
          disabled={isLoading}
          className="absolute top-4 right-4 p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors disabled:opacity-50"
        >
          <X className="w-4 h-4" />
        </button>

        {/* Icon Header */}
        <div className={`p-4 rounded-2xl border flex items-center justify-center ${styles.iconBg}`}>
          {icon || defaultIcon}
        </div>

        {/* Text Content */}
        <div className="space-y-1.5">
          <h3 className="text-base font-bold text-slate-800">{title}</h3>
          <p className="text-xs text-slate-500 max-w-xs leading-relaxed">{description}</p>
        </div>

        {/* Action Buttons */}
        <div className="flex items-center gap-3 w-full pt-2">
          <button
            onClick={onClose}
            disabled={isLoading}
            className="flex-1 py-2.5 px-4 bg-slate-100 border border-slate-200/80 text-slate-700 text-xs font-semibold rounded-2xl hover:bg-slate-200/70 transition-all disabled:opacity-50"
          >
            {cancelText}
          </button>
          <button
            onClick={onConfirm}
            disabled={isLoading}
            className={`flex-1 py-2.5 px-4 rounded-2xl text-xs font-semibold shadow-md transition-all flex items-center justify-center gap-2 disabled:opacity-50 ${styles.confirmBtn}`}
          >
            {isLoading && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
            <span>{isLoading ? "Processing..." : confirmText}</span>
          </button>
        </div>
      </div>
    </div>
  );
};
