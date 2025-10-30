// 書籍一覧ページ：検索欄＋BookList（モックデータ）
import { BookList, type BookItem } from '@/components/books/BookList';

const mock: BookItem[] = [
  { id: '1', title: 'Clean Architecture', author: 'Robert C. Martin', tags: ['go', 'intermediate'] },
  { id: '2', title: 'Designing Data-Intensive Applications', author: 'Martin Kleppmann', tags: ['architecture', 'database'] },
  { id: '3', title: 'Refactoring', author: 'Martin Fowler', tags: ['refactoring'] },
  { id: '4', title: 'You Don\'t Know JS Yet', author: 'Kyle Simpson', tags: ['javascript'] },
];

export default function BooksPage() {
  return (
    <section>
      <h1>書籍一覧</h1>
      <div style={{ margin: '12px 0' }}>
        <input className="form__input" placeholder="キーワード検索" aria-label="書籍検索" />
      </div>
      <BookList items={mock} />
    </section>
  );
}


