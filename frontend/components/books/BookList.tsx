"use client";
// 書籍リスト：BookCard をグリッドで並べ、クリックで詳細遷移
import { useRouter } from 'next/navigation';
import { BookCard } from './BookCard';

export type BookItem = { id: string; title: string; author: string; cover?: string; tags?: string[] };

export function BookList({ items, isLoading = false }: { items: BookItem[]; isLoading?: boolean }) {
  const router = useRouter();
  if (isLoading) return <p>読み込み中...</p>;
  return (
    <div className="books">
      {items.map((b) => (
        <BookCard key={b.id} title={b.title} author={b.author} cover={b.cover} tags={b.tags} onClick={() => router.push(`/books/${b.id}`)} />
      ))}
    </div>
  );
}


