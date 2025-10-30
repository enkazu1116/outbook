// 認証: ログイン画面（onSubmit でトップへ遷移）
"use client";
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function LoginPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    // ここで本来は認証APIを呼ぶ。ひとまず成功として遷移
    router.push('/');
  };
  return (
    <div className="auth">
      <div className="auth__card">
        <h1 className="auth__title">ログイン</h1>
        <form className="form form--stack" onSubmit={onSubmit}>
          <label className="form__label" htmlFor="email">メールアドレス</label>
          <input className="form__input" id="email" type="email" required disabled={loading} />
          <p className="form__error" aria-live="polite" style={{ display: 'none' }}>メール形式が不正です</p>

          <label className="form__label" htmlFor="password">パスワード</label>
          <input className="form__input" id="password" type="password" minLength={8} required disabled={loading} />
          <p className="form__error" aria-live="polite" style={{ display: 'none' }}>8文字以上で入力してください</p>

          <button className="button button--primary form__submit" type="submit" disabled={loading}>{loading ? '処理中...' : 'ログイン'}</button>
          <div className="auth__links">
            <a href="/register">新規登録</a>
            <a href="/forgot">パスワードをお忘れの方</a>
          </div>
        </form>
      </div>
    </div>
  );
}


