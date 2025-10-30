// テーマ切替ボタン
// - Reduxの ui.theme をトグル
// - body クラス（theme-light/theme-dark）と localStorage に同期
"use client";
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import type { RootState } from '@/store';
import { toggleTheme } from '@/store';

export function ThemeToggle() {
  const dispatch = useDispatch();
  const mode = useSelector((s: RootState) => s.ui.theme);

  useEffect(() => {
    const body = document.body;
    body.classList.remove('theme-light', 'theme-dark');
    body.classList.add(mode === 'dark' ? 'theme-dark' : 'theme-light');
    try { localStorage.setItem('theme', mode); } catch {}
  }, [mode]);

  return (
    <button className="button button--secondary" onClick={() => dispatch(toggleTheme())} aria-pressed={mode === 'dark'}>
      {mode === 'dark' ? 'ライト' : 'ダーク'}
    </button>
  );
}


