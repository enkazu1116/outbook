// Pagination：シンプルなページ送り（アクセシブルなナビ）
export function Pagination({ page, total, onChange }: { page: number; total: number; onChange: (p: number) => void }) {
  const pages = Array.from({ length: total }, (_, i) => i + 1);
  return (
    <nav aria-label="ページ送り" style={{ display: 'flex', gap: 8, marginTop: 16 }}>
      {pages.map((p) => (
        <button key={p} className="button" aria-current={p === page ? 'page' : undefined} onClick={() => onChange(p)} style={{ background: p === page ? '#e5e7eb' : '#fff', border: '1px solid #e5e7eb' }}>
          {p}
        </button>
      ))}
    </nav>
  );
}


