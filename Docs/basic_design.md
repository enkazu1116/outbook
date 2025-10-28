# OutBook 基本機能設計書

## 1. システムアーキテクチャ

### 1.1 全体構成
```
┌─────────────────┐
│   Frontend      │
│   (Next.js)     │
│   Vercel        │
└────────┬────────┘
         │ HTTPS
         │
┌────────▼────────┐
│   Backend API   │
│   (Go + Echo)   │
│   AWS           │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼──┐  ┌──▼────────┐
│ Turso│  │External   │
│  DB  │  │APIs       │
└──────┘  │・OpenBD    │
          │・OpenAI   │
          │・Claude   │
          └───────────┘
```

### 1.2 レイヤー構成（バックエンド）

```
┌────────────────────────────────┐
│     Presentation Layer         │
│  (HTTPハンドラー、ミドルウェア)  │
└────────────────────────────────┘
                │
┌────────────────────────────────┐
│      Application Layer         │
│     (ビジネスロジック)           │
└────────────────────────────────┘
                │
┌────────────────────────────────┐
│      Domain Layer              │
│  (エンティティ、リポジトリIF)    │
└────────────────────────────────┘
                │
┌────────────────────────────────┐
│    Infrastructure Layer        │
│  (DB接続、外部API連携)          │
└────────────────────────────────┘
```

### 1.3 ディレクトリ構成

#### バックエンド（Go）
```
backend/
├── cmd/
│   └── server/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── domain/                  # ドメイン層
│   │   ├── entity/             # エンティティ
│   │   │   ├── user.go
│   │   │   ├── book.go
│   │   │   ├── output.go
│   │   │   └── review.go
│   │   └── repository/         # リポジトリインターフェース
│   │       ├── user.go
│   │       ├── book.go
│   │       ├── output.go
│   │       └── review.go
│   ├── application/            # アプリケーション層
│   │   ├── usecase/           # ユースケース
│   │   │   ├── user/
│   │   │   ├── book/
│   │   │   ├── output/
│   │   │   └── review/
│   │   └── service/           # ドメインサービス
│   ├── infrastructure/        # インフラ層
│   │   ├── persistence/      # データベース実装
│   │   │   ├── turso/
│   │   │   └── migration/
│   │   └── external/         # 外部API連携
│   │       ├── openbd/
│   │       └── ai/
│   └── presentation/         # プレゼンテーション層
│       ├── handler/          # HTTPハンドラー
│       │   ├── user.go
│       │   ├── book.go
│       │   ├── output.go
│       │   └── review.go
│       ├── middleware/       # ミドルウェア
│       │   ├── auth.go
│       │   ├── cors.go
│       │   └── ratelimit.go
│       └── router/           # ルーティング
│           └── router.go
├── pkg/                      # 共通パッケージ
│   ├── validator/
│   ├── password/
│   └── jwt/
├── config/                   # 設定ファイル
│   └── config.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

#### フロントエンド（Next.js）
```
frontend/
├── src/
│   ├── app/                    # App Router
│   │   ├── (auth)/            # 認証関連
│   │   │   ├── login/
│   │   │   ├── register/
│   │   │   └── layout.tsx
│   │   ├── (main)/            # メイン機能
│   │   │   ├── books/         # 書籍一覧・詳細
│   │   │   ├── outputs/       # アウトプット一覧・詳細
│   │   │   ├── my/            # マイページ
│   │   │   ├── search/        # 検索
│   │   │   ├── ranking/       # ランキング
│   │   │   └── layout.tsx
│   │   ├── layout.tsx         # ルートレイアウト
│   │   └── page.tsx           # トップページ
│   ├── components/            # コンポーネント
│   │   ├── ui/               # 共通UIコンポーネント
│   │   ├── features/         # 機能別コンポーネント
│   │   │   ├── auth/
│   │   │   ├── book/
│   │   │   ├── output/
│   │   │   └── review/
│   │   └── layout/           # レイアウトコンポーネント
│   ├── lib/                  # ユーティリティ
│   │   ├── api/             # APIクライアント
│   │   ├── hooks/           # カスタムフック
│   │   └── utils/           # 汎用関数
│   └── types/               # 型定義
├── public/
├── tailwind.config.ts
├── next.config.js
├── package.json
└── tsconfig.json
```

## 2. データベース詳細設計

### 2.1 テーブル定義

#### users（ユーザー）
```sql
CREATE TABLE users (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    email               TEXT NOT NULL UNIQUE,
    password_hash       TEXT NOT NULL,
    nickname            TEXT NOT NULL,
    avatar_url          TEXT,
    bio                 TEXT,
    skill_level         TEXT CHECK(skill_level IN ('junior', 'middle', 'senior')),
    years_of_experience INTEGER,
    specialties         TEXT, -- JSON配列形式 ["frontend", "backend"]
    favorite_languages  TEXT, -- JSON配列形式
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at          DATETIME
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);
```

#### books（書籍）
```sql
CREATE TABLE books (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    title             TEXT NOT NULL,
    author            TEXT NOT NULL,
    publisher         TEXT,
    published_date    DATE,
    isbn              TEXT UNIQUE,
    cover_image_url   TEXT,
    description       TEXT,
    created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_books_isbn ON books(isbn);
CREATE INDEX idx_books_title ON books(title);
CREATE INDEX idx_books_author ON books(author);
```

#### book_tags（書籍タグ）
```sql
CREATE TABLE book_tags (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id         INTEGER NOT NULL,
    tag_type        TEXT NOT NULL, -- 'language', 'level', 'specialty'
    tag_value       TEXT NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE INDEX idx_book_tags_book_id ON book_tags(book_id);
CREATE INDEX idx_book_tags_type_value ON book_tags(tag_type, tag_value);
```

#### outputs（アウトプット）
```sql
CREATE TABLE outputs (
    id                      INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id                 INTEGER NOT NULL,
    book_id                 INTEGER NOT NULL,
    title                   TEXT NOT NULL,
    content                 TEXT NOT NULL,
    visibility              TEXT NOT NULL CHECK(visibility IN ('public', 'title_only', 'private')),
    reading_finished_date   DATE,
    is_draft                BOOLEAN DEFAULT FALSE,
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at              DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE INDEX idx_outputs_user_id ON outputs(user_id);
CREATE INDEX idx_outputs_book_id ON outputs(book_id);
CREATE INDEX idx_outputs_visibility ON outputs(visibility);
CREATE INDEX idx_outputs_created_at ON outputs(created_at DESC);
```

#### output_tags（アウトプットタグ）
```sql
CREATE TABLE output_tags (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    output_id   INTEGER NOT NULL,
    tag_name    TEXT NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (output_id) REFERENCES outputs(id) ON DELETE CASCADE
);

CREATE INDEX idx_output_tags_output_id ON output_tags(output_id);
CREATE INDEX idx_output_tags_tag_name ON output_tags(tag_name);
```

#### reviews（書評）
```sql
CREATE TABLE reviews (
    id                          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id                     INTEGER NOT NULL,
    book_id                     INTEGER NOT NULL,
    rating                      INTEGER NOT NULL CHECK(rating >= 1 AND rating <= 5),
    content                     TEXT NOT NULL,
    recommended_levels          TEXT, -- JSON配列 ["beginner", "intermediate"]
    recommended_specialties     TEXT, -- JSON配列 ["frontend", "backend"]
    helpful_count               INTEGER DEFAULT 0,
    created_at                  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at                  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    UNIQUE(user_id, book_id) -- 1ユーザー1書籍1レビュー
);

CREATE INDEX idx_reviews_book_id ON reviews(book_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);
```

#### review_helpfulness（書評の役立ち度）
```sql
CREATE TABLE review_helpfulness (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    review_id   INTEGER NOT NULL,
    user_id     INTEGER NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(review_id, user_id)
);

CREATE INDEX idx_review_helpfulness_review_id ON review_helpfulness(review_id);
```

#### reports（通報）
```sql
CREATE TABLE reports (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    review_id           INTEGER NOT NULL,
    reporter_user_id    INTEGER NOT NULL,
    reason              TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'reviewed', 'dismissed', 'actioned')),
    admin_note          TEXT,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    resolved_at         DATETIME,
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (reporter_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_review_id ON reports(review_id);
```

#### user_settings（ユーザー設定）
```sql
CREATE TABLE user_settings (
    id                      INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id                 INTEGER NOT NULL UNIQUE,
    email_notification      BOOLEAN DEFAULT TRUE,
    encrypted_openai_key    TEXT,
    encrypted_claude_key    TEXT,
    default_search_filters  TEXT, -- JSON形式
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);
```

### 2.2 マスターデータ

#### プログラミング言語
```json
[
  "Go", "Python", "JavaScript", "TypeScript", "Java", "C#", "C++", "C",
  "Ruby", "PHP", "Swift", "Kotlin", "Rust", "Scala", "R", "SQL", "Shell"
]
```

#### レベル
```json
{
  "beginner": "初級",
  "intermediate": "中級",
  "advanced": "上級"
}
```

#### 職域
```json
[
  "frontend", "backend", "infrastructure", "network", "security",
  "data_science", "mobile", "devops", "qa", "full_stack"
]
```

## 3. API設計（REST API）

### 3.1 エンドポイント一覧

#### 認証系
| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| POST | /api/v1/auth/register | ユーザー登録 | 不要 |
| POST | /api/v1/auth/login | ログイン | 不要 |
| POST | /api/v1/auth/logout | ログアウト | 必要 |
| POST | /api/v1/auth/refresh | トークンリフレッシュ | 必要 |

#### ユーザー系
| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | /api/v1/users/me | 自分の情報取得 | 必要 |
| PUT | /api/v1/users/me | プロフィール更新 | 必要 |
| PUT | /api/v1/users/me/password | パスワード変更 | 必要 |
| DELETE | /api/v1/users/me | アカウント削除 | 必要 |
| GET | /api/v1/users/:id | ユーザー情報取得 | 不要 |

#### 書籍系
| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | /api/v1/books | 書籍一覧取得 | 不要 |
| GET | /api/v1/books/:id | 書籍詳細取得 | 不要 |
| POST | /api/v1/books | 書籍登録 | 必要 |
| GET | /api/v1/books/search | 書籍検索 | 不要 |
| GET | /api/v1/books/isbn/:isbn | ISBN検索 | 不要 |

#### アウトプット系
| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | /api/v1/outputs | アウトプット一覧 | 不要 |
| GET | /api/v1/outputs/:id | アウトプット詳細 | 不要 |
| POST | /api/v1/outputs | アウトプット作成 | 必要 |
| PUT | /api/v1/outputs/:id | アウトプット更新 | 必要 |
| DELETE | /api/v1/outputs/:id | アウトプット削除 | 必要 |
| GET | /api/v1/users/:id/outputs | ユーザーのアウトプット一覧 | 不要 |
| GET | /api/v1/books/:id/outputs | 書籍のアウトプット一覧 | 不要 |

#### 書評系
| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | /api/v1/reviews | 書評一覧 | 不要 |
| GET | /api/v1/books/:id/reviews | 書籍の書評一覧 | 不要 |
| POST | /api/v1/reviews | 書評作成 | 必要 |
| PUT | /api/v1/reviews/:id | 書評更新 | 必要 |
| DELETE | /api/v1/reviews/:id | 書評削除 | 必要 |
| POST | /api/v1/reviews/:id/helpful | 役立ち度評価 | 必要 |

#### ランキング系
| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | /api/v1/rankings/overall | 総合ランキング | 不要 |
| GET | /api/v1/rankings/category | カテゴリ別ランキング | 不要 |

#### 通報系
| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| POST | /api/v1/reports | 通報作成 | 必要 |

### 3.2 リクエスト・レスポンス例

#### POST /api/v1/auth/register
**Request**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "nickname": "tech_reader",
  "skill_level": "junior",
  "years_of_experience": 2,
  "specialties": ["frontend", "backend"],
  "favorite_languages": ["JavaScript", "TypeScript"]
}
```

**Response (201 Created)**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "nickname": "tech_reader",
    "avatar_url": null,
    "skill_level": "junior",
    "years_of_experience": 2,
    "specialties": ["frontend", "backend"],
    "favorite_languages": ["JavaScript", "TypeScript"],
    "created_at": "2025-10-28T10:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### POST /api/v1/auth/login
**Request**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK)**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "nickname": "tech_reader"
  },
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2025-10-29T10:00:00Z"
}
```

#### POST /api/v1/outputs
**Request**
```json
{
  "book_id": 123,
  "title": "Clean Architectureを読んで学んだこと",
  "content": "## 概要\n\nClean Architectureは...",
  "visibility": "public",
  "reading_finished_date": "2025-10-27",
  "tags": ["アーキテクチャ", "設計", "Go"],
  "is_draft": false
}
```

**Response (201 Created)**
```json
{
  "id": 456,
  "user_id": 1,
  "book_id": 123,
  "title": "Clean Architectureを読んで学んだこと",
  "content": "## 概要\n\nClean Architectureは...",
  "visibility": "public",
  "reading_finished_date": "2025-10-27",
  "tags": ["アーキテクチャ", "設計", "Go"],
  "is_draft": false,
  "created_at": "2025-10-28T10:00:00Z",
  "updated_at": "2025-10-28T10:00:00Z"
}
```

#### GET /api/v1/books/search
**Request**
```
GET /api/v1/books/search?q=Clean&language=Go&level=intermediate&page=1&limit=20
```

**Response (200 OK)**
```json
{
  "books": [
    {
      "id": 123,
      "title": "Clean Architecture 達人に学ぶソフトウェアの構造と設計",
      "author": "Robert C. Martin",
      "publisher": "KADOKAWA",
      "published_date": "2018-07-27",
      "isbn": "9784048930659",
      "cover_image_url": "https://...",
      "description": "ソフトウェア開発の名著...",
      "tags": {
        "languages": ["Go", "Java"],
        "levels": ["intermediate", "advanced"],
        "specialties": ["backend", "architecture"]
      },
      "stats": {
        "output_count": 145,
        "review_count": 89,
        "average_rating": 4.5
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_count": 95,
    "limit": 20
  }
}
```

#### POST /api/v1/reviews
**Request**
```json
{
  "book_id": 123,
  "rating": 5,
  "content": "設計の原則について非常に分かりやすく解説されている。実務経験2年以上あると理解しやすい。",
  "recommended_levels": ["intermediate", "advanced"],
  "recommended_specialties": ["backend", "full_stack"]
}
```

**Response (201 Created)**
```json
{
  "id": 789,
  "user_id": 1,
  "book_id": 123,
  "rating": 5,
  "content": "設計の原則について非常に分かりやすく解説されている...",
  "recommended_levels": ["intermediate", "advanced"],
  "recommended_specialties": ["backend", "full_stack"],
  "helpful_count": 0,
  "created_at": "2025-10-28T10:00:00Z",
  "updated_at": "2025-10-28T10:00:00Z"
}
```

### 3.3 エラーレスポンス

#### 標準エラーフォーマット
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力値が不正です",
    "details": [
      {
        "field": "email",
        "message": "メールアドレスの形式が正しくありません"
      }
    ]
  }
}
```

#### エラーコード一覧
| コード | HTTPステータス | 説明 |
|-------|---------------|------|
| VALIDATION_ERROR | 400 | バリデーションエラー |
| UNAUTHORIZED | 401 | 認証エラー |
| FORBIDDEN | 403 | 権限エラー |
| NOT_FOUND | 404 | リソースが見つからない |
| CONFLICT | 409 | リソースの競合 |
| RATE_LIMIT_EXCEEDED | 429 | レート制限超過 |
| INTERNAL_ERROR | 500 | サーバー内部エラー |

## 4. 画面設計

### 4.1 画面一覧

#### 認証系
1. **ログイン画面** (`/login`)
2. **ユーザー登録画面** (`/register`)

#### メイン機能
3. **トップページ** (`/`)
   - 最新のアウトプット一覧
   - 人気の書籍
   - ランキング（簡易版）

4. **書籍一覧画面** (`/books`)
   - 書籍の一覧表示
   - フィルター（言語、レベル、職域）
   - ソート（人気順、新着順、評価順）

5. **書籍詳細画面** (`/books/:id`)
   - 書籍情報
   - 書評一覧
   - アウトプット一覧
   - 書評投稿ボタン

6. **書籍検索画面** (`/search`)
   - 検索フォーム
   - 高度なフィルター
   - 検索結果一覧

7. **アウトプット一覧画面** (`/outputs`)
   - 公開されたアウトプット一覧
   - フィルター（タグ、書籍、ユーザー）

8. **アウトプット詳細画面** (`/outputs/:id`)
   - アウトプット本文（Markdownレンダリング）
   - 書籍情報
   - 著者情報

9. **アウトプット作成/編集画面** (`/outputs/new`, `/outputs/:id/edit`)
   - Markdownエディター
   - プレビュー機能
   - 公開設定
   - 下書き保存

10. **ランキング画面** (`/ranking`)
    - 期間別タブ（週間/月間/年間）
    - カテゴリ別タブ

11. **マイページ** (`/my`)
    - 自分のアウトプット一覧
    - 自分の書評一覧
    - 読書記録（タイトルのみ公開含む）

12. **プロフィール編集画面** (`/my/profile`)
    - ユーザー情報編集
    - パスワード変更
    - APIキー設定

13. **ユーザーページ** (`/users/:id`)
    - ユーザー情報
    - 公開アウトプット一覧
    - 書評一覧

### 4.2 画面遷移図

```
[トップページ] ────┬──→ [書籍一覧] ──→ [書籍詳細] ──┬──→ [アウトプット詳細]
      │           │                               │
      │           └──→ [検索]                     └──→ [書評投稿]
      │
      ├──→ [アウトプット一覧] ──→ [アウトプット詳細]
      │
      ├──→ [ランキング] ──→ [書籍詳細]
      │
      └──→ [マイページ] ──┬──→ [プロフィール編集]
                         │
                         ├──→ [アウトプット作成/編集]
                         │
                         └──→ [自分のアウトプット詳細]

[ログイン] ←─→ [ユーザー登録]
```

### 4.3 主要画面ワイヤーフレーム

#### 書籍詳細画面
```
┌──────────────────────────────────────┐
│ Header (ロゴ、検索バー、ユーザーメニュー)│
├──────────────────────────────────────┤
│ ┌────────────┬───────────────────┐   │
│ │            │ タイトル            │   │
│ │  書籍画像   │ 著者名              │   │
│ │            │ 出版社・出版日       │   │
│ │            │ ★★★★☆ (4.5)      │   │
│ └────────────┤ [書評を書く]        │   │
│              │ タグ: Go 中級 設計  │   │
│              └───────────────────┘   │
├──────────────────────────────────────┤
│ [書籍情報] [書評(89)] [アウトプット(145)]│
├──────────────────────────────────────┤
│ 書評一覧                              │
│ ┌──────────────────────────────────┐ │
│ │ ★★★★★ by tech_reader          │ │
│ │ 非常に分かりやすい...             │ │
│ │ 役立った: 25                      │ │
│ └──────────────────────────────────┘ │
│ ┌──────────────────────────────────┐ │
│ │ ★★★★☆ by developer_1          │ │
│ │ 実践的な内容で...                │ │
│ │ 役立った: 18                      │ │
│ └──────────────────────────────────┘ │
└──────────────────────────────────────┘
```

## 5. 主要機能フロー

### 5.1 ユーザー登録フロー

```
[ユーザー] → [登録フォーム入力]
              ↓
          [バリデーション]
              ↓
          [重複チェック]
              ↓
          [パスワードハッシュ化]
              ↓
          [DB登録]
              ↓
          [JWTトークン発行]
              ↓
          [ログイン状態でリダイレクト]
```

### 5.2 アウトプット投稿フロー

```
[ユーザー] → [書籍選択/登録]
              ↓
          [アウトプット作成画面]
              ↓
          [下書き保存] ←→ [プレビュー確認]
              ↓
          [公開設定選択]
              ↓
          [投稿]
              ↓
          [バリデーション]
              ↓
          [DB保存]
              ↓
          [アウトプット詳細ページへ]
```

### 5.3 書籍検索フロー

```
[ユーザー] → [検索条件入力]
              ↓
          [フィルター選択]
          (言語/レベル/職域)
              ↓
          [API呼び出し]
              ↓
          [DB検索実行]
          (複合条件クエリ)
              ↓
          [結果集計・ソート]
              ↓
          [ページネーション処理]
              ↓
          [検索結果表示]
```

### 5.4 書評投稿フロー

```
[ユーザー] → [書籍詳細ページ]
              ↓
          [書評を書くボタン]
              ↓
          [既存レビューチェック]
              ↓ (なし)
          [書評フォーム表示]
              ↓
          [評価・本文入力]
              ↓
          [おすすめ設定]
          (レベル/職域選択)
              ↓
          [投稿]
              ↓
          [バリデーション]
              ↓
          [DB保存]
              ↓
          [書籍の平均評価更新]
              ↓
          [書籍詳細ページへ戻る]
```

## 6. セキュリティ設計

### 6.1 認証・認可

#### JWT認証
- **トークン形式**: Bearer Token
- **有効期限**: 24時間
- **リフレッシュトークン**: 30日間
- **格納場所**: HTTPOnly Cookie (CSRF対策済み)

#### パスワードポリシー
- 最低8文字
- 英大文字・小文字・数字・記号のうち3種類以上
- よくあるパスワードのブラックリストチェック

#### ハッシュ化
- **アルゴリズム**: bcrypt
- **コストファクター**: 12

### 6.2 API保護

#### レート制限
- 認証API: 5回/分/IP
- 一般API: 100回/分/ユーザー
- 検索API: 30回/分/ユーザー

#### CORS設定
```go
AllowOrigins: []string{"https://outbook.example.com"}
AllowMethods: []string{"GET", "POST", "PUT", "DELETE"}
AllowHeaders: []string{"Authorization", "Content-Type"}
AllowCredentials: true
```

### 6.3 データ保護

#### 暗号化対象
- ユーザーのAPIキー（AES-256-GCM）
- パスワード（bcrypt）

#### SQL インジェクション対策
- プリペアドステートメントの使用
- ORMの適切な使用

#### XSS対策
- 入力値のサニタイズ
- Content-Security-Policy ヘッダー設定
- HTMLエスケープ処理

## 7. パフォーマンス最適化

### 7.1 データベース最適化

#### インデックス戦略
- 検索条件に使われるカラムにインデックス作成
- 複合インデックスの活用
- カーディナリティの高いカラムを優先

#### クエリ最適化
- N+1問題の回避（Eager Loading）
- 必要なカラムのみSELECT
- ページネーション実装

### 7.2 キャッシュ戦略

#### Redis活用（将来実装）
- セッション情報
- ランキングデータ（1時間キャッシュ）
- 書籍情報（24時間キャッシュ）
- 検索結果（15分キャッシュ）

### 7.3 フロントエンド最適化

#### Next.js最適化
- ISR (Incremental Static Regeneration) 活用
  - 書籍詳細ページ: 1時間ごとに再生成
  - ランキングページ: 6時間ごとに再生成
- 画像最適化（next/image）
- コード分割（dynamic import）
- フォント最適化

## 8. 監視・ログ設計

### 8.1 ログレベル

| レベル | 用途 |
|-------|------|
| ERROR | エラー発生時 |
| WARN | 警告（APIレート制限到達など） |
| INFO | 重要なイベント（ユーザー登録、投稿作成など） |
| DEBUG | デバッグ情報（開発環境のみ） |

### 8.2 ログ出力形式

```json
{
  "timestamp": "2025-10-28T10:00:00Z",
  "level": "INFO",
  "service": "outbook-api",
  "trace_id": "abc123",
  "user_id": 1,
  "method": "POST",
  "path": "/api/v1/outputs",
  "status": 201,
  "duration_ms": 45,
  "message": "Output created successfully"
}
```

### 8.3 監視項目

#### アプリケーション監視
- APIレスポンスタイム
- エラーレート
- リクエスト数（エンドポイント別）
- アクティブユーザー数

#### インフラ監視
- CPU使用率
- メモリ使用率
- ディスク使用率
- ネットワークトラフィック

## 9. テスト戦略

### 9.1 テストレベル

#### ユニットテスト
- カバレッジ目標: 70%以上
- 対象: ビジネスロジック、ユーティリティ関数
- ツール: Go testing, testify

#### 統合テスト
- API エンドポイントのテスト
- DB連携テスト
- ツール: httptest, testcontainers

#### E2Eテスト（将来実装）
- 主要ユーザーフローのテスト
- ツール: Playwright

### 9.2 テストデータ

#### シードデータ
- 開発環境用のサンプルデータ
- テスト用ユーザー（各スキルレベル）
- サンプル書籍データ（20冊程度）
- サンプルアウトプット・書評

## 10. デプロイメント

### 10.1 環境構成

| 環境 | 用途 | インフラ |
|------|------|---------|
| Local | ローカル開発 | Docker Compose |
| Development | 開発環境 | AWS（または同等） |
| Staging | 検証環境 | AWS（または同等） |
| Production | 本番環境 | Vercel (FE) + AWS (BE) |

### 10.2 CI/CDパイプライン

```
[Git Push] → [GitHub Actions]
              ↓
          [Lint & Format Check]
              ↓
          [Unit Tests]
              ↓
          [Build]
              ↓
          [Integration Tests]
              ↓
        [develop branch]
              ↓
          [Deploy to Staging]
              ↓
        [main branch]
              ↓
          [Deploy to Production]
```

### 10.3 環境変数

#### バックエンド
```env
# Database
DATABASE_URL=
DATABASE_AUTH_TOKEN=

# JWT
JWT_SECRET=
JWT_EXPIRATION=24h

# External APIs
OPENBD_API_URL=https://api.openbd.jp/v1
GOOGLE_BOOKS_API_KEY=

# Server
PORT=8080
ENV=production
```

#### フロントエンド
```env
# API
NEXT_PUBLIC_API_URL=https://api.outbook.example.com
API_URL=https://api.outbook.example.com

# Auth
NEXTAUTH_SECRET=
NEXTAUTH_URL=https://outbook.example.com
```

## 11. MVP開発優先順位

### Phase 1: 基盤構築（Week 1-2）
1. プロジェクト環境セットアップ
2. データベース設計・マイグレーション
3. 認証機能実装
4. 基本API構築

### Phase 2: コア機能（Week 3-4）
1. 書籍管理機能
2. アウトプット投稿機能（基本）
3. 書評機能（基本）
4. ユーザープロフィール

### Phase 3: 検索・表示（Week 5-6）
1. 書籍検索機能
2. フィルタリング機能
3. 一覧表示の最適化
4. ランキング機能（簡易版）

### Phase 4: UI/UX改善（Week 7-8）
1. レスポンシブデザイン実装
2. ユーザビリティ改善
3. パフォーマンス最適化
4. エラーハンドリング改善

---

## 付録

### A. 用語集

| 用語 | 説明 |
|------|------|
| アウトプット | ユーザーが投稿する読書記録・学習記録 |
| 書評 | 書籍に対する評価・レビュー |
| 可視性 | アウトプットの公開範囲設定 |
| 職域 | エンジニアの専門分野（フロント、バックエンドなど） |

### B. 参考リンク

- [Echo Framework](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- [Turso Documentation](https://docs.turso.tech/)
- [Next.js Documentation](https://nextjs.org/docs)
- [OpenBD API](https://openbd.jp/)

