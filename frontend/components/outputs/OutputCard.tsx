// OutputCard：タイトル/投稿者/いいね/投稿日を表示する一覧用カード
export type OutputCardProps = {
  title: string;
  user: string;
  likeCount: number;
  publishedAt: string;
  onClick?: () => void;
};

export function OutputCard({ title, user, likeCount, publishedAt, onClick }: OutputCardProps) {
  return (
    <article className="card" role="button" tabIndex={0} style={{ padding: 12 }} onClick={onClick} onKeyDown={(e) => { if (e.key === 'Enter') onClick?.(); }}>
      <h3 style={{ margin: '0 0 8px' }}>{title}</h3>
      <p style={{ color: '#6b7280', margin: 0 }}>{user} ・ {new Date(publishedAt).toLocaleDateString('ja-JP')}</p>
      <p style={{ margin: '8px 0 0' }}>👍 {likeCount}</p>
    </article>
  );
}


