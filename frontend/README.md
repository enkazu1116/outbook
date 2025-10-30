# Next.js アプリ作成手順（詳細ガイド）

本書は、Next.js(App Router) を使ったフロントエンドアプリをゼロから作成するための手順書です。Node.js の導入から、TypeScript 化、ディレクトリ構成、CSS、状態管理、実行・ビルド・デプロイ、トラブルシュートまでカバーします。

---

## 1. 前提・準備
- OS: Windows/Mac/Linux（クロスプラットフォーム想定）
- Node.js 18 以上（LTS 推奨）
  - Windows: https://nodejs.org/ja からインストール
  - バージョン確認: `node -v` / `npm -v`
- エディタ: VS Code 推奨（ESLint / Prettier 拡張が便利）

---

## 2. プロジェクト作成
1) 作業ディレクトリへ移動
```bash
cd <任意の作業ディレクトリ>
```
2) Next.js アプリを作成（TypeScript/ESLint/App Router）
```bash
npx create-next-app@latest my-app \
  --typescript \
  --eslint \
  --app \
  --src-dir false \
  --import-alias "@/*" \
  --use-npm
cd my-app
npm install
```

---

## 3. ディレクトリ構成（基本）
```
my-app/
├── app/
│   ├── layout.tsx        # ルートレイアウト
│   └── page.tsx          # トップページ
├── public/               # 静的ファイル
├── styles/               # CSS（グローバル/モジュール）
├── next.config.mjs
├── package.json
└── tsconfig.json
```
- ページ: `app/<segment>/page.tsx`
- レイアウト: `app/<segment>/layout.tsx`
- URLに出ないグループ: `(group)` ディレクトリ

---

## 4. 起動・ビルド・Lint
- 開発起動
```bash
npm run dev
# 3000 が埋まっている場合
npm run dev -- -p 3001
```
- 本番ビルド
```bash
npm run build
```
- 本番起動
```bash
npm run start
```
- Lint
```bash
npm run lint
```

---

## 5. TypeScript とパスエイリアス
`tsconfig.json`（抜粋）
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "jsx": "preserve",
    "strict": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./*"],
      "@/components/*": ["./components/*"],
      "@/styles/*": ["./styles/*"],
      "@/store": ["./store/index.ts"]
    },
    "types": ["node"],
    "plugins": [{ "name": "next" }]
  }
}
```

---

## 6. App Router の基本
- クライアントコンポーネント: 先頭に `"use client";`
- サーバコンポーネント: デフォルト（`"use client";` が無い）
- 例: ログインページ
```tsx
// app/(auth)/login/page.tsx
"use client";
export default function LoginPage() {
  return <div>ログイン</div>;
}
```

---

## 7. CSS 設計
- グローバル: `app/layout.tsx` で `styles/*.css` を読み込む
- コンポーネント: CSS Modules（`*.module.css`）や BEM で管理
- テーマ: CSS 変数で `:root` / `.theme-dark` を切替
```tsx
// app/layout.tsx
import "@/styles/base.css";
import "@/styles/layout.css";
import "@/styles/components.css";
import "@/styles/pages.css";
import "@/styles/themes.css";
```

---

## 8. 状態管理（Redux Toolkit 例）
1) 依存追加
```bash
npm install @reduxjs/toolkit react-redux
```
2) ストア作成
```ts
// store/index.ts
import { configureStore, createSlice } from "@reduxjs/toolkit";
const uiSlice = createSlice({
  name: "ui",
  initialState: { theme: "light" as "light" | "dark" },
  reducers: { toggleTheme(s){ s.theme = s.theme === "dark" ? "light" : "dark"; } }
});
export const { toggleTheme } = uiSlice.actions;
export const store = configureStore({ reducer: { ui: uiSlice.reducer } });
export type RootState = ReturnType<typeof store.getState>;
```
3) Provider 設置
```tsx
// app/layout.tsx
import StoreProvider from "@/components/providers/StoreProvider";
```
4) 利用例（テーマ切替）
```tsx
"use client";
import { useDispatch, useSelector } from "react-redux";
import { toggleTheme } from "@/store";
import type { RootState } from "@/store";
export function ThemeToggle(){
  const dispatch = useDispatch();
  const mode = useSelector((s: RootState) => s.ui.theme);
  return <button onClick={() => dispatch(toggleTheme())}>{mode === "dark" ? "ライト" : "ダーク"}</button>;
}
```

---

## 9. API ルート（Route Handler）
- `app/api/<name>/route.ts` に `GET/POST` を実装
```ts
// app/api/ping/route.ts
import { NextResponse } from "next/server";
export async function GET(){
  return NextResponse.json({ ok: true, ts: Date.now() });
}
```

---

## 10. デプロイ
- Vercel: リポジトリを接続 → Framework: Next.js → 自動ビルド
- Node サーバ: `npm run build` → `npm run start`
- 静的エクスポート（制約あり）: `next.config.mjs` に `output: 'export'`

---

## 11. トラブルシュート
- ポート衝突（EADDRINUSE: 3000）
  - `npm run dev -- -p 3001`
- 依存解決エラー（Module not found）
  - `npm install`、`tsconfig.json` の `paths`、`next-env.d.ts` の存在、再起動
- Client Component エラー
  - `"use client";` を付与（イベントハンドラやHooks使用時）

---

## 12. 推奨ライブラリ
- バリデーション: `zod` / `react-hook-form`
- データ取得: `swr` / `@tanstack/react-query`
- Markdown: `remark/rehype` + サニタイズ（`dompurify` 等）
- UI アイコン: `@radix-ui/react-icons` / `lucide-react`

---

## 13. 参考
- Next.js: https://nextjs.org/
- App Router: https://nextjs.org/docs/app
- Route Handlers: https://nextjs.org/docs/app/building-your-application/routing/route-handlers
- TypeScript: https://nextjs.org/docs/app/building-your-application/configuring/typescript
- Styling: https://nextjs.org/docs/app/building-your-application/styling
