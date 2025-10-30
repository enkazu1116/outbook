// Toast：一時的な通知（info/success/warning/error）
// - duration 経過で自動クローズ
"use client";
import { useEffect, useState } from 'react';

export function Toast({ type = 'info', message, duration = 3000 }: { type?: 'info' | 'success' | 'warning' | 'error'; message: string; duration?: number }) {
  const [open, setOpen] = useState(true);
  useEffect(() => {
    const t = setTimeout(() => setOpen(false), duration);
    return () => clearTimeout(t);
  }, [duration]);
  if (!open) return null;
  const bg = type === 'success' ? '#d1fae5' : type === 'error' ? '#fee2e2' : type === 'warning' ? '#fef3c7' : '#e5e7eb';
  return (
    <div style={{ position: 'fixed', top: 12, right: 12, background: bg, padding: '8px 12px', borderRadius: 8, boxShadow: '0 2px 6px rgba(0,0,0,.1)', zIndex: 60 }}>
      {message}
    </div>
  );
}


