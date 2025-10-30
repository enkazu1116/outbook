export default function ProfilePage({ params }: { params: { username: string } }) {
  const { username } = params;
  return (
    <section>
      <header style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
        <div style={{ width: 64, height: 64, borderRadius: '50%', background: 'rgba(124,172,255,0.2)', border: '1px solid var(--color-border)' }} />
        <div>
          <h1 style={{ margin: 0 }}>{username}</h1>
          <p style={{ margin: 0, color: 'var(--color-muted)' }}>自己紹介がここに入ります</p>
        </div>
      </header>

      <nav aria-label="プロフィールタブ" style={{ display: 'flex', gap: 8, marginTop: 12 }}>
        <a className="button button--secondary" href="#outputs">アウトプット</a>
        <a className="button button--secondary" href="#reviews">レビュー</a>
      </nav>

      <div id="outputs" style={{ marginTop: 16 }}>
        <h2>アウトプット</h2>
        <p>コンテンツは未実装（モック）です。</p>
      </div>

      <div id="reviews" style={{ marginTop: 16 }}>
        <h2>レビュー</h2>
        <p>コンテンツは未実装（モック）です。</p>
      </div>
    </section>
  );
}


