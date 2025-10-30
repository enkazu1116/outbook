// サイドバー：主要ナビ（Dashboard/Books/Outputs/Profile）
"use client";
import Link from 'next/link';
import { usePathname } from 'next/navigation';

export function Sidebar() {
  const pathname = usePathname();
  const isActive = (href: string) => href === '/' ? pathname === '/' : pathname?.startsWith(href);

  return (
    <aside className="sidebar">
      <nav className="nav" aria-label="メインナビゲーション">
        <Link href="/" aria-current={isActive('/') ? 'page' : undefined}>Dashboard</Link>
        <Link href="/books" aria-current={isActive('/books') ? 'page' : undefined}>Books</Link>
        <Link href="/outputs" aria-current={isActive('/outputs') ? 'page' : undefined}>Outputs</Link>
        <Link href="/profile/you" aria-current={isActive('/profile') ? 'page' : undefined}>Profile</Link>
      </nav>
    </aside>
  );
}


