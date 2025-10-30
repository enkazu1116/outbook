// 書籍詳細ページ：メタ情報/タグ/関連アウトプット/レビューのセクション構成
import Link from 'next/link';
type Props = { params: { id: string } };

export default function BookDetailPage({ params }: Props) {
  const { id } = params;
  return (
    <section className="book-detail">
      <div className="book-detail__header">
        <img className="book-detail__cover" alt="書籍カバー" />
        <div className="book-detail__meta">
          <h1 className="book-detail__title">書籍タイトル {id}</h1>
          <p className="book-detail__author">著者名</p>
          <p className="book-detail__desc">説明テキスト...</p>
          <div className="book-detail__tags">
            <span className="tag">tag1</span>
            <span className="tag">tag2</span>
          </div>
          <Link className="button button--primary" href="/outputs/new?bookId=1">アウトプットを作成</Link>
        </div>
      </div>

      <div className="book-detail__section">
        <h2>関連アウトプット</h2>
        <div>出力一覧（後で実装）</div>
      </div>
      <div className="book-detail__section">
        <h2>レビュー</h2>
        <div>レビュー一覧（後で実装）</div>
      </div>
    </section>
  );
}


