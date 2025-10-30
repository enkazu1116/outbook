// MarkdownRenderer（簡易版）
// - 行ごとの段落化のみ。本番ではサニタイズ/パーサ導入を推奨
// NOTE: 簡易表示（セキュリティ対策として本番はサニタイズ必須）
export function MarkdownRenderer({ markdown }: { markdown: string }) {
  // 極簡易：コードブロックやリンクはプレーンに
  return (
    <div className="md">
      {markdown.split('\n').map((line, i) => (
        <p key={i}>{line}</p>
      ))}
    </div>
  );
}


