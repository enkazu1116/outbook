# OutBook 詳細UI設計（Next.js App Router + CSS）

## 1. 目的と前提
- 対象：Next.js（App Router）フロントエンドの画面・コンポーネント・UIの詳細設計。
- デザイン制約：すべてのスタイリングは純粋なCSSで行う（CSS Modules / グローバルCSS / BEM推奨）。
- 情報源：OutBook システム設計書 v2.0 を踏まえた要件（Books / Outputs / Reviews / Profile）。
- ペルソナ：学習アウトプットを投稿・閲覧するユーザー、管理者。

---

## 2. ナビゲーション設計（サイトマップ・画面遷移）

### 2.1 サイトマップ
```
/(login)
/(main)
  ├─ /books
  │   └─ /books/[id]
  │       └─ /books/[id]/outputs
  ├─ /outputs
  │   ├─ /outputs/new
  │   └─ /outputs/[id]
  │       └─ /outputs/[id]/edit
  └─ /profile/[username]
```

### 2.2 主要遷移
- ログイン成功 → /(main)
- サイドバーから書籍一覧へ → /(main)/books
- 書籍カード押下 → /books/[id]
- 書籍詳細の「アウトプットを作成」 → /outputs/new（bookIdパラメータ付与）
- 一覧のアイテム押下 → /outputs/[id]
- 投稿者本人が「編集」 → /outputs/[id]/edit

---

## 3. レイアウト・レスポンシブ設計

### 3.1 共通レイアウト
- 固定ヘッダー（高さ 64px）＋左サイドバー（幅 240px）＋メインコンテンツ（可変）。
- 最小幅 320px を想定、ブレークポイントは 768px / 1024px / 1280px。

### 3.2 ブレークポイント
| 幅 | レイアウト |
|----|-----------|
| 1280px 以上 | サイドバー常時表示、グリッド3〜4列 |
| 1024〜1279px | サイドバー縮小（アイコンのみ）、グリッド2〜3列 |
| 768〜1023px | サイドバー折りたたみ、グリッド2列 |
| 767px 以下 | ヘッダー固定、ボトムナビ検討、リストは1列 |

---

## 4. 共通UI仕様

### 4.1 タイポグラフィ・カラー
- フォント：Noto Sans JP, sans-serif
- ベース文字色：#1f2937
- 背景：#f9fafb
- 強調色：#2563eb
- 警告色：#dc2626
- 成功色：#16a34a
- ボーダー：#e5e7eb
- 影：0 2px 4px rgba(0,0,0,0.1)
- 角丸：6px（カード・入力・ボタン共通）

### 4.2 要素別コンポーネント
- Header：ロゴ、グローバル検索、プロフィールメニュー。
- Sidebar：主要ナビゲーション（Dashboard / Books / Outputs / Profile）。
- Footer：コピーライト、ポリシーリンク。
- Card：影・角丸付きの情報コンテナ。
- Table：一覧表示。ヘッダー固定、行ホバー。
- Dialog：中央固定モーダル。背景スクラム防止（body scroll lock）。
- Toast：右上に出る通知（情報、成功、警告、エラー）。
- Form：ラベル上、入力下、エラー文は赤文字・小サイズ。

### 4.3 入力・ボタン共通
- フォーカスリング：outline: 2px solid rgba(37,99,235,0.5)
- 入力：高さ 40px、左右 12px、border 1px #d1d5db、角丸 6px。
- プライマリボタン：背景 #2563eb、hover #1d4ed8、disabled 0.5透明。
- セカンダリボタン：背景 #e5e7eb、hover #d1d5db。

---

## 5. 画面別UI仕様

### 5.1 ログイン（/login）
- 構成：中央カード配置。ロゴ、タイトル、入力（email/password）、送信、補助リンク（登録・パスワード再発行）。
- バリデーション：email形式、password最小8文字。失敗時はカード内上部にエラーバナー。
- 状態：送信中はボタンにスピナー、入力はdisabled。

レイアウト例（BEM）:
```html
<div class="auth">
  <div class="auth__card">
    <h1 class="auth__title">ログイン</h1>
    <form class="form form--stack">
      <label class="form__label" for="email">メールアドレス</label>
      <input class="form__input" id="email" type="email" />
      <p class="form__error">メール形式が不正です</p>

      <label class="form__label" for="password">パスワード</label>
      <input class="form__input" id="password" type="password" />
      <p class="form__error">8文字以上で入力してください</p>

      <button class="button button--primary form__submit">ログイン</button>
      <div class="auth__links">
        <a href="/register">新規登録</a>
        <a href="/forgot">パスワードをお忘れの方</a>
      </div>
    </form>
  </div>
</div>
```

主要CSS（抜粋）:
```css
.auth { min-height: 100vh; display: grid; place-items: center; background:#f9fafb; }
.auth__card { width: 100%; max-width: 420px; background: #fff; border:1px solid #e5e7eb; border-radius:8px; padding:24px; box-shadow:0 2px 4px rgba(0,0,0,0.08); }
.auth__title { margin:0 0 16px; font-size:20px; }
.form--stack .form__label { display:block; font-size:14px; margin:12px 0 6px; }
.form__input { width:100%; height:40px; padding:0 12px; border:1px solid #d1d5db; border-radius:6px; }
.form__input:focus { outline:2px solid rgba(37,99,235,0.3); }
.form__error { color:#dc2626; font-size:12px; margin:6px 0 0; }
.form__submit { width:100%; margin-top:16px; }
.auth__links { display:flex; justify-content:space-between; margin-top:12px; font-size:14px; }
```

---

### 5.2 ダッシュボード（/(main)）
- ヘッダー：検索、ユーザーアイコン、通知。
- サイドバー：現在地強調。キーボード操作可（Tab順）。
- メイン：KPIカード、最近のアウトプット、ブックマーク。

KPIカード（Grid 4列、縮小で2列/1列）:
```html
<section class="kpis">
  <div class="kpi-card">
    <div class="kpi-card__label">総アウトプット</div>
    <div class="kpi-card__value">42</div>
  </div>
  <div class="kpi-card">
    <div class="kpi-card__label">今月の投稿</div>
    <div class="kpi-card__value">7</div>
  </div>
  <!-- ... -->
</section>
```

---

### 5.3 書籍一覧（/books）
- 検索バー＋フィルター（言語、レベル、タグ）。
- グリッドカード表示、カード押下で詳細に遷移。
- ページネーション：下部に配置、現在ページを強調。

BookCard:
```html
<article class="book-card">
  <img class="book-card__cover" src="/covers/xxx.jpg" alt="書籍カバー">
  <div class="book-card__body">
    <h3 class="book-card__title">Clean Architecture</h3>
    <p class="book-card__author">Robert C. Martin</p>
    <div class="book-card__meta">
      <span class="tag">go</span>
      <span class="tag">intermediate</span>
    </div>
  </div>
</article>
```

主要CSS（抜粋）:
```css
.books { display:grid; grid-template-columns: repeat(4, minmax(0,1fr)); gap:16px; }
@media (max-width: 1279px){ .books { grid-template-columns: repeat(3,1fr); } }
@media (max-width: 1023px){ .books { grid-template-columns: repeat(2,1fr); } }
@media (max-width: 767px){ .books { grid-template-columns: 1fr; } }

.book-card { background:#fff; border:1px solid #e5e7eb; border-radius:8px; overflow:hidden; box-shadow:0 2px 4px rgba(0,0,0,.06); transition:transform .15s; }
.book-card:hover { transform: translateY(-2px); }
.book-card__cover { width:100%; height:180px; object-fit:cover; background:#f3f4f6; }
.book-card__body { padding:12px; }
.book-card__title { margin:0 0 4px; font-size:16px; line-height:1.3; }
.book-card__author { margin:0 0 8px; color:#6b7280; font-size:14px; }
.tag { display:inline-block; padding:2px 8px; background:#eef2ff; color:#3730a3; border-radius:9999px; font-size:12px; margin-right:6px; }
```

---

### 5.4 書籍詳細（/books/[id]）
- 上部：カバー、タイトル、著者、説明、タグ。
- 右側（または下部）：統計（アウトプット数、レビュー数、平均評価）。
- ボタン：「アウトプットを作成」→ /outputs/new?bookId=[id]

構造:
```html
<section class="book-detail">
  <div class="book-detail__header">
    <img class="book-detail__cover" />
    <div class="book-detail__meta">
      <h1 class="book-detail__title">Clean Architecture</h1>
      <p class="book-detail__author">Robert C. Martin</p>
      <p class="book-detail__desc">説明テキスト...</p>
      <div class="book-detail__tags">...</div>
      <button class="button button--primary">アウトプットを作成</button>
    </div>
  </div>

  <div class="book-detail__section">
    <h2>関連アウトプット</h2>
    <!-- OutputList -->
  </div>
  <div class="book-detail__section">
    <h2>レビュー</h2>
    <!-- ReviewList -->
  </div>
</section>
```

---

### 5.5 アウトプット一覧（/outputs）
- タブ：すべて / 自分 / 下書き。
- 並び替え：新着 / 人気 / 評価。
- カードリスト、右上に新規作成ボタン。

---

### 5.6 アウトプット詳細（/outputs/[id]）
- タイトル、メタ（著者、投稿日、いいね数）。
- Markdown表示。見出し・コード・表・画像のスタイル最適化。
- 右側にAIレビュー結果パネル（存在する場合）。
- フッターにアクション（いいね、共有、通報）。

Markdown表示CSS（抜粋）:
```css
.md { line-height: 1.8; font-size: 16px; }
.md h1,.md h2,.md h3 { margin: 24px 0 12px; line-height:1.3; }
.md pre { background:#111827; color:#e5e7eb; padding:12px; border-radius:6px; overflow:auto; }
.md code { background:#f3f4f6; padding:2px 6px; border-radius:4px; }
.md a { color:#2563eb; text-decoration: underline; }
.md table { width:100%; border-collapse: collapse; margin:12px 0; }
.md th,.md td { border:1px solid #e5e7eb; padding:8px; text-align:left; }
.md img { max-width:100%; border-radius:6px; }
```

---

### 5.7 アウトプット作成・編集（/outputs/new, /outputs/[id]/edit）
- 二画面エディタ：左にMarkdown入力、右にプレビュー。
- ツールバー：見出し、太字、斜体、コード、リスト、リンク、画像、プレビュー同期。
- 入力項目：タイトル、本文、タグ（最大10個）、種類（summary/note/review/tutorial）、公開範囲（public/private/followers）。
- バリデーション：タイトル必須、本文10文字以上、タグ最大10。

HTML構造:
```html
<form class="editor">
  <div class="editor__toolbar">
    <button type="button" data-md="h2">H2</button>
    <button type="button" data-md="bold">B</button>
    <button type="button" data-md="italic">I</button>
    <button type="button" data-md="code">Code</button>
    <button type="button" data-md="ul">List</button>
    <button type="button" data-md="link">Link</button>
    <button type="button" data-md="image">Image</button>
  </div>

  <label class="form__label">タイトル</label>
  <input class="form__input" name="title" />

  <div class="editor__panes">
    <textarea class="editor__input" name="content"></textarea>
    <div class="editor__preview md"></div>
  </div>

  <div class="editor__meta">
    <div class="field">
      <label class="form__label">種類</label>
      <select class="form__select" name="type">
        <option value="summary">要約</option>
        <option value="note">ノート</option>
        <option value="review">レビュー</option>
        <option value="tutorial">チュートリアル</option>
      </select>
    </div>
    <div class="field">
      <label class="form__label">公開範囲</label>
      <select class="form__select" name="visibility">
        <option value="public">公開</option>
        <option value="followers">フォロワー</option>
        <option value="private">非公開</option>
      </select>
    </div>
    <div class="field">
      <label class="form__label">タグ</label>
      <input class="form__input" name="tags" placeholder="カンマ区切り 例: architecture, pattern" />
      <p class="form__help">最大10個まで</p>
    </div>
  </div>

  <div class="editor__actions">
    <button class="button button--secondary" type="button">下書き保存</button>
    <button class="button button--primary" type="submit">投稿</button>
  </div>
</form>
```

CSS（抜粋）:
```css
.editor__toolbar { display:flex; gap:8px; margin-bottom:8px; }
.editor__toolbar button { background:#e5e7eb; border:none; padding:6px 10px; border-radius:4px; cursor:pointer; }
.editor__toolbar button:hover { background:#d1d5db; }

.editor__panes { display:grid; grid-template-columns: 1fr 1fr; gap:16px; }
@media (max-width: 1023px){ .editor__panes { grid-template-columns: 1fr; } }

.editor__input, .editor__preview { height: 520px; border:1px solid #d1d5db; border-radius:6px; padding:12px; background:#fff; }
.editor__preview { background:#f8fafc; overflow:auto; }

.editor__meta { display:grid; grid-template-columns: repeat(3,1fr); gap:16px; margin-top:16px; }
@media (max-width: 1023px){ .editor__meta { grid-template-columns: 1fr; } }

.editor__actions { display:flex; gap:12px; justify-content:flex-end; margin-top:16px; }
```

---

### 5.8 レビューリスト（/books/[id] 内セクション）
- 並び替え：新着順、役立ち度。
- レーティング表示（星）。
- 各カードはタイトル、本文要約、役立ち度、投稿者、投稿日。

---

### 5.9 プロフィール（/profile/[username]）
- ヘッダー：アイコン、表示名、自己紹介、統計。
- タブ：アウトプット、レビュー。
- 自己プロフィール編集はモーダルで行う。

---

## 6. アクセシビリティ・操作性

- コントラスト比 WCAG AA 以上を確保（本文 4.5:1 以上）。
- フォーカスインジケータを全操作部品に付与。
- ページ内見出し階層を適正化（h1→h2→h3）。
- フォーム要素は label for と id を対応。
- モーダルはフォーカストラップ・Escクローズ・背景クリックで閉じる。
- 画像 alt 属性必須。装飾は空 alt。

---

## 7. バリデーション・エラーハンドリング

- クライアント側：入力直後に即時（onBlur）、送信時に総合チェック。
- サーバー側：エラーJSONを RFC7807 風に統一。メッセージはフォーム下に表示。
- エラーバナー：フォーム上部に要約、各フィールド下に詳細。

---

## 8. フィードバック・通知

- 成功：緑のトースト（3秒）。
- 警告・確認：ダイアログで明示的に確認（削除など破壊的操作）。
- 長時間処理：ボタンにローディングスピナー、再押下抑止。

---

## 9. CSS 設計指針

- BEM命名：`.block__element--modifier`
- レイヤ分割：`base/`（リセット・タイポ）、`layout/`（ヘッダー・サイドバー）、`components/`（カード・ボタン）、`pages/`（ページ固有）。
- 変数：`:root` に色・間隔スケールを定義しテーマ切り替え容易化。
- ダークテーマ：`.theme-dark` を body に付与して色を上書き。

CSS変数例:
```css
:root {
  --color-text: #1f2937;
  --color-bg: #f9fafb;
  --color-primary: #2563eb;
  --color-danger: #dc2626;
  --radius: 6px;
  --space-1: 8px;
  --space-2: 16px;
  --space-3: 24px;
}
.theme-dark {
  --color-text: #e5e7eb;
  --color-bg: #111827;
}
```

---

## 10. コンポーネント一覧（仕様サマリ）

| 名称 | 役割 | 主要プロパティ | 補足 |
|-----|------|---------------|------|
| Header | 全ページ共通 | none | 検索・ユーザーメニュー |
| Sidebar | ナビ | activePath | アクセシビリティ考慮 |
| BookCard | 書籍カード | title, author, cover, tags | カードクリックで詳細 |
| BookList | 一覧 | items[], isLoading | グリッドレスポンシブ |
| OutputCard | 投稿カード | title, user, likeCount, publishedAt | タグ表示 |
| OutputEditor | エディタ | value, onChange, onSubmit | ツールバー操作 |
| MarkdownRenderer | 表示 | markdown | セキュリティ（サニタイズ） |
| ReviewCard | レビュー項目 | rating, title, content | helpful ボタン |
| Pagination | ページ送り | page, total, onChange | キーボード操作 |
| Dialog | モーダル | open, onClose | フォーカストラップ |
| Toast | 通知 | type, message | 自動クローズ |

---

## 11. 画面一覧の完了条件（DoD）

- 全画面でキーボード操作・スクリーンリーダーの基本動作を満たす。
- 全主要操作のフォーカステストを通過。
- スマートフォン幅での崩れがない。
- 入力エラーの説明がフィールド下に表示される。
- ローディング・成功・失敗の状態が視覚化される。

---

## 12. 参考ファイル構成（フロント側）

```
frontend/
├── app/
│   ├── (auth)/login/page.tsx
│   ├── (main)/layout.tsx
│   ├── (main)/page.tsx
│   ├── (main)/books/page.tsx
│   ├── (main)/books/[id]/page.tsx
│   ├── (main)/outputs/page.tsx
│   ├── (main)/outputs/new/page.tsx
│   ├── (main)/outputs/[id]/page.tsx
│   └── (main)/outputs/[id]/edit/page.tsx
├── components/
│   ├── layout/{Header,Sidebar,Footer}.tsx
│   ├── books/{BookCard,BookList,BookSearchBar}.tsx
│   ├── outputs/{OutputCard,OutputEditor,MarkdownRenderer}.tsx
│   ├── common/{Dialog,Toast,Pagination}.tsx
│   └── forms/{Field,Help,Error}.tsx
├── styles/
│   ├── base.css
│   ├── layout.css
│   ├── components.css
│   ├── pages.css
│   └── themes.css
└── public/
```

---

## 13. 今後の拡張
- 表示パフォーマンス向上：画像遅延読み込み、コード分割。
- テーマ対応：ダークテーマ、アクセントカラー変更。
- エディタ機能拡張：ドラッグ＆ドロップ画像、表挿入、絵文字禁則のまま記法支援。
- 国際化：日本語・英語の切替。

