// 認証: パスワード再発行（モック送信）
"use client";
export default function ForgotPasswordPage() {
  return (
    <div className="auth">
      <div className="auth__card">
        <h1 className="auth__title">パスワード再発行</h1>
        <form className="form form--stack" onSubmit={(e) => e.preventDefault()}>
          <label className="form__label" htmlFor="email">メールアドレス</label>
          <input className="form__input" id="email" type="email" required />
          <p className="form__help">登録済みのメールアドレスに再設定リンクを送信します。</p>
          <button className="button button--primary form__submit" type="submit">送信</button>
          <div className="auth__links">
            <a href="/login">ログインへ戻る</a>
          </div>
        </form>
      </div>
    </div>
  );
}


