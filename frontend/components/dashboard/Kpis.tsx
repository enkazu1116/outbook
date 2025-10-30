// KPIカード群：主要メトリクスをカードで表示（再利用可能）
import { memo } from 'react';
export const Kpis = memo(function Kpis() {
  return (
    <div className="kpis">
      <div className="kpi-card">
        <div className="kpi-card__label">総アウトプット</div>
        <div className="kpi-card__value">42</div>
      </div>
      <div className="kpi-card">
        <div className="kpi-card__label">今月の投稿</div>
        <div className="kpi-card__value">7</div>
      </div>
      <div className="kpi-card">
        <div className="kpi-card__label">ブックマーク</div>
        <div className="kpi-card__value">12</div>
      </div>
      <div className="kpi-card">
        <div className="kpi-card__label">レビュー</div>
        <div className="kpi-card__value">5</div>
      </div>
    </div>
  );
});


