import { ThemeToggle } from '@/components/common/ThemeToggle';

// ヘッダー：ロゴ/検索/通知/プロフィール/テーマ切替
export function Header() {
  return (
    <header className="header">
      <div className="header__logo">OutBook</div>
      <div className="header__right">
        <input className="form__input" placeholder="検索" aria-label="検索" />
        <button className="button button--secondary">通知</button>
        <button className="button button--secondary" aria-label="プロフィール">👤</button>
        {/* theme toggle */}
        <ThemeToggle />
      </div>
    </header>
  );
}


