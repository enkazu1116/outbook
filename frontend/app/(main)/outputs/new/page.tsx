// アウトプット作成：エディタ（左入力/右プレビュー）とメタ情報
"use client";
import { useRouter, useSearchParams } from 'next/navigation';
import { OutputEditor } from '@/components/outputs/OutputEditor';

export default function OutputNewPage() {
  const router = useRouter();
  const sp = useSearchParams();
  const bookId = sp.get('bookId');

  return (
    <section>
      <h1>アウトプット作成 {bookId ? `(bookId: ${bookId})` : ''}</h1>
      <OutputEditor onSubmit={() => router.push('/(main)/outputs')} />
    </section>
  );
}


