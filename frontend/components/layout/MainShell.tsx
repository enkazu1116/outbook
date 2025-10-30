// アプリ共通のレイアウト枠（Header/Sidebar/Main/Footer）
// Ctrl + B（または ⌘ + B）でサイドバーの表示/非表示をトグル
"use client";
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { toggleSidebar } from '@/store';
import type { RootState } from '@/store';
import { Header } from '@/components/layout/Header';
import { Sidebar } from '@/components/layout/Sidebar';
import { Footer } from '@/components/layout/Footer';

export function MainShell({ children }: { children: React.ReactNode }) {
  const collapsed = useSelector((s: RootState) => s.ui.sidebarCollapsed);
  const dispatch = useDispatch();
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && (e.key.toLowerCase() === 'b')) {
        e.preventDefault();
        dispatch(toggleSidebar());
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [dispatch]);

  return (
    <div className={`layout ${collapsed ? 'layout--collapsed' : ''}`}>
      <Header />
      <Sidebar />
      <main className="main">
        {children}
        <Footer />
      </main>
    </div>
  );
}


