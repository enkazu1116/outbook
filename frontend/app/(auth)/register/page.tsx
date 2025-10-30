// 認証: 新規登録画面（モック送信）
"use client";
export default function RegisterPage() {
  return (
    <div className="auth">
      <div className="auth__card">
        <h1 className="auth__title">新規登録</h1>
        <form className="form form--stack" onSubmit={(e) => e.preventDefault()}>
          <label className="form__label" htmlFor="name">表示名</label>
          <input className="form__input" id="name" required />

          <label className="form__label" htmlFor="email">メールアドレス</label>
          <input className="form__input" id="email" type="email" required />

          <label className="form__label" htmlFor="password">パスワード</label>
          <input className="form__input" id="password" type="password" minLength={8} required />

          <button className="button button--primary form__submit" type="submit">登録</button>
          <div className="auth__links">
            <a href="/login">ログインへ戻る</a>
          </div>
        </form>
      </div>
    </div>
  );
}


