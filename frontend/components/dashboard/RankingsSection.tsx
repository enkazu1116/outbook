// 売上ランキングセクション
// - Redux(ranking)の source/items/loading/error を使用
// - タブ（Amazon/Gihyo）操作で source を切替 → fetchRanking を発火
"use client";
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import type { RootState, AppDispatch } from '@/store';
import { fetchRanking, setSource } from '@/store';

export function RankingsSection() {
  const dispatch = useDispatch<AppDispatch>();
  const source = useSelector((s: RootState) => s.ranking.source);
  const items = useSelector((s: RootState) => s.ranking.items);
  const loading = useSelector((s: RootState) => s.ranking.loading);
  const error = useSelector((s: RootState) => s.ranking.error);

  useEffect(() => {
    dispatch(fetchRanking(source));
  }, [dispatch, source]);

  return (
    <section style={{ marginTop: 16 }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>売上ランキング</h2>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="button button--secondary" aria-pressed={source==='amazon'} onClick={() => dispatch(setSource('amazon'))}>Amazon</button>
          <button className="button button--secondary" aria-pressed={source==='gihyo'} onClick={() => dispatch(setSource('gihyo'))}>Gihyo</button>
        </div>
      </header>

      {loading && <p>読み込み中...</p>}
      {error && <p className="form__error">{error}</p>}

      <ol style={{ listStyle: 'none', padding: 0, margin: '12px 0 0', display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: 12 }}>
        {items.map((b) => (
          <li key={`${source}-${b.rank}`} className="card" style={{ padding: 12, display: 'flex', gap: 12, alignItems: 'center' }}>
            <div style={{ width: 28, textAlign: 'right', fontWeight: 700 }}>{b.rank}</div>
            <img alt="cover" src={b.cover || '/covers/placeholder.jpg'} style={{ width: 56, height: 80, objectFit: 'cover', borderRadius: 6, background: '#0d1128' }} />
            <div style={{ flex: 1 }}>
              <div style={{ fontWeight: 600 }}>{b.title}</div>
              {b.author && <div style={{ color: 'var(--color-muted)', fontSize: 14 }}>{b.author}</div>}
            </div>
            {b.url && <a className="button button--primary" href={b.url} target="_blank" rel="noreferrer">見る</a>}
          </li>
        ))}
      </ol>
    </section>
  );
}


