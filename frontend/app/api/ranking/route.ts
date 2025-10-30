// /api/ranking
// - クエリ ?source=amazon|gihyo に応じてランキングを返すモックAPI
// - 実運用では Amazon PA-API/RSS 等と差し替え、環境変数で認証を管理
import { NextResponse } from 'next/server';

type RankingItem = { rank: number; title: string; author?: string; cover?: string; url?: string };

async function fetchAmazonMock(): Promise<RankingItem[]> {
  return [
    { rank: 1, title: 'Clean Architecture', author: 'Robert C. Martin' },
    { rank: 2, title: 'Designing Data-Intensive Applications', author: 'Martin Kleppmann' },
    { rank: 3, title: 'Refactoring', author: 'Martin Fowler' },
    { rank: 4, title: "You Don't Know JS Yet", author: 'Kyle Simpson' },
  ];
}

async function fetchGihyoMock(): Promise<RankingItem[]> {
  return [
    { rank: 1, title: '達人に学ぶDB設計 徹底指南書', author: 'ミック' },
    { rank: 2, title: '詳解システムパフォーマンス', author: 'Brendan Gregg' },
    { rank: 3, title: 'Real-World HTTP 第3版', author: '渋川よしき' },
    { rank: 4, title: 'WEB+DB PRESS 総集編', author: '技術評論社編集部' },
  ];
}

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const source = (searchParams.get('source') || 'amazon').toLowerCase();

  // NOTE: 実API連携する場合はここで Amazon PA-API / RSS などを呼び出す
  // 認証情報は環境変数から取得し、失敗時はモックへフォールバック
  let items: RankingItem[] = [];
  try {
    if (source === 'gihyo') {
      items = await fetchGihyoMock();
    } else {
      items = await fetchAmazonMock();
    }
  } catch (e) {
    items = await fetchAmazonMock();
  }

  return NextResponse.json({ source, items });
}


