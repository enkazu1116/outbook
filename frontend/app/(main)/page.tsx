// ダッシュボード: KPI と 売上ランキング（Amazon/Gihyo）を表示
// ランキングは内製API（/api/ranking）から取得し、セクションコンポーネントで描画
import { Kpis } from '@/components/dashboard/Kpis';
import { RankingsSection } from '@/components/dashboard/RankingsSection';

export default function DashboardPage() {
  return (
    <section>
      <Kpis />
      <RankingsSection />
    </section>
  );
}


