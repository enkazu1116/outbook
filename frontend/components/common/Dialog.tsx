// Dialog：モーダルダイアログ
// - 背景クリック/ESCで閉じる。フォーカスは簡易制御
"use client";
import { useEffect, useRef } from 'react';

export function Dialog({ open, onClose, children }: { open: boolean; onClose: () => void; children: React.ReactNode }) {
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
    if (open) document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, [open, onClose]);
  if (!open) return null;
  return (
    <div role="dialog" aria-modal="true" style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.4)', display: 'grid', placeItems: 'center', zIndex: 50 }} onClick={onClose}>
      <div ref={ref} className="card" style={{ padding: 16, minWidth: 300, background: '#fff' }} onClick={(e) => e.stopPropagation()}>
        {children}
      </div>
    </div>
  );
}


