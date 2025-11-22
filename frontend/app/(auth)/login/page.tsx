// 認証: ログイン画面（onSubmit でトップへ遷移）
"use client";
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { ApiError, healthCheck } from '@/lib/api';

export default function LoginPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      // TODO: backendのログインAPIエンドポイントが実装されたらここで呼び出す
      // 現在はヘルスチェックのみ実行して接続確認
      const isHealthy = await healthCheck();
      if (!isHealthy) {
        throw new Error('サーバーに接続できませんでした');
      }

      // 認証成功時の処理（トークン保存など）
      // 例: const token = response.token;
      // localStorage.setItem('token', token);
      
      router.push('/');
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('ログインに失敗しました');
      }
    } finally {
      setLoading(false);
    }
  };
  return (
    <div className="auth">
      <div className="auth__card">
        <h1 className="auth__title">ログイン</h1>
        <form className="form form--stack" onSubmit={onSubmit}>
          {error && (
            <div className="form__error" role="alert" style={{ display: 'block', marginBottom: '1rem' }}>
              {error}
            </div>
          )}
          
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
          <p className="form__error" aria-live="polite" style={{ display: 'none' }}>メール形式が不正です</p>

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


