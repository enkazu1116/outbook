// アウトプット一覧：カードグリッド表示＋新規作成ボタン
"use client";
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { OutputCard } from '@/components/outputs/OutputCard';

const mock = [
  { id: 'o1', title: 'DDDの集約境界まとめ', user: 'alice', likeCount: 12, publishedAt: '2024-04-01' },
  { id: 'o2', title: 'CQRS入門ノート', user: 'bob', likeCount: 5, publishedAt: '2024-04-02' },
];

export default function OutputsPage() {
  const router = useRouter();
  return (
    <section>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
        <h1>アウトプット一覧</h1>
        <Link className="button button--primary" href="/outputs/new">新規作成</Link>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, minmax(0,1fr))', gap: 16 }}>
        {mock.map((o) => (
          <OutputCard key={o.id} title={o.title} user={o.user} likeCount={o.likeCount} publishedAt={o.publishedAt} onClick={() => router.push(`/outputs/${o.id}`)} />
        ))}
      </div>
    </section>
  );
}


