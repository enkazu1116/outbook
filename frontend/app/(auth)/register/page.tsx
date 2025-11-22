// 認証: 新規登録画面
"use client";
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { userApi, ApiError } from '@/lib/api';

export default function RegisterPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const result = await userApi.register(email, password, name);
      // 登録成功時の処理
      alert(result.message || '登録が完了しました');
      router.push('/login');
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('登録に失敗しました');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth">
      <div className="auth__card">
        <h1 className="auth__title">新規登録</h1>
        <form className="form form--stack" onSubmit={onSubmit}>
          {error && (
            <div className="form__error" role="alert" style={{ display: 'block', marginBottom: '1rem' }}>
              {error}
            </div>
          )}

          <label className="form__label" htmlFor="name">表示名</label>
          <input
            className="form__input"
            id="name"
            required
            disabled={loading}
            value={name}
            onChange={(e) => setName(e.target.value)}
          />

          <label className="form__label" htmlFor="email">メールアドレス</label>
          <input
            className="form__input"
            id="email"
            type="email"
            required
            disabled={loading}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />

          <label className="form__label" htmlFor="password">パスワード</label>
          <input
            className="form__input"
            id="password"
            type="password"
            minLength={8}
            required
            disabled={loading}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          <button
            className="button button--primary form__submit"
            type="submit"
            disabled={loading}
          >
            {loading ? '処理中...' : '登録'}
          </button>
          <div className="auth__links">
            <a href="/login">ログインへ戻る</a>
          </div>
        </form>
      </div>
    </div>
  );
}


