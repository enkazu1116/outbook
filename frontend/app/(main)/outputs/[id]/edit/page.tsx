// アウトプット編集：既存内容を初期値としてエディタ表示
"use client";
import { useRouter } from 'next/navigation';
import { OutputEditor } from '@/components/outputs/OutputEditor';

type Props = { params: { id: string } };

export default function OutputEditPage({ params }: Props) {
  const router = useRouter();
  return (
    <section>
      <h1>アウトプット編集 {params.id}</h1>
      <OutputEditor initial={{ title: '既存タイトル', content: '既存本文' }} onSubmit={() => router.push(`/(main)/outputs/${params.id}`)} />
    </section>
  );
}


