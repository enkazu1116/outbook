// アウトプット詳細：Markdown表示とアクション（編集/いいね/共有/通報）
import Link from 'next/link';
import { MarkdownRenderer } from '@/components/outputs/MarkdownRenderer';

type Props = { params: { id: string } };

export default function OutputDetailPage({ params }: Props) {
  const content = `# タイトル\n\n本文の例\n\n- 箇条書き\n- 箇条書き2`;
  return (
    <article>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h1>アウトプット {params.id}</h1>
          <p style={{ color: '#6b7280' }}>by alice ・ 2024/04/02 ・ 👍 12</p>
        </div>
        <Link href={`/outputs/${params.id}/edit`} className="button button--secondary">編集</Link>
      </header>

      <section className="md" style={{ marginTop: 16 }}>
        <MarkdownRenderer markdown={content} />
      </section>

      <footer style={{ marginTop: 16, display: 'flex', gap: 8 }}>
        <button className="button button--secondary">いいね</button>
        <button className="button button--secondary">共有</button>
        <button className="button button--secondary">通報</button>
      </footer>
    </article>
  );
}


