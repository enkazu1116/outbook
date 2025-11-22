# OutBook Backend

モジュラーモノリス構成のGoバックエンドアプリケーション

## ディレクトリ構成

```
backend/
├── cmd/
│   └── server/           # アプリケーションエントリーポイント
│       └── main.go
├── internal/
│   ├── domain/           # ドメイン層（ビジネスロジックの核心）
│   │   ├── user/         # ユーザードメイン
│   │   │   ├── entity/
│   │   │   └── repository/
│   │   ├── book/         # 書籍ドメイン
│   │   │   ├── entity/
│   │   │   └── repository/
│   │   ├── output/       # アウトプットドメイン
│   │   │   ├── entity/
│   │   │   └── repository/
│   │   └── review/       # 書評ドメイン
│   │       ├── entity/
│   │       └── repository/
│   ├── application/      # アプリケーション層（ユースケース）
│   │   ├── user/
│   │   │   └── service/
│   │   ├── book/
│   │   │   └── service/
│   │   ├── output/
│   │   │   └── service/
│   │   └── review/
│   │       └── service/
│   ├── presentation/     # プレゼンテーション層（インターフェース）
│   │   └── http/
│   │       ├── handler/
│   │       └── router.go
│   ├── infrastructure/   # インフラストラクチャ層（外部依存）
│   │   └── database/
│   │       └── turso/
│   └── shared/           # 共通コード
│       ├── dto/
│       └── middleware/
├── pkg/                  # 再利用可能なパッケージ
├── Dockerfile
├── Dockerfile.dev
├── go.mod
└── go.sum
```

## アーキテクチャ

### クリーンアーキテクチャ

```
┌─────────────────────────────────────────┐
│  Presentation Layer (HTTP)              │
│  - Handlers                             │
│  - Router                               │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│  Application Layer                      │
│  - Services (Use Cases)                 │
│  - Business Logic                       │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│  Domain Layer                           │
│  - Entities                             │
│  - Repository Interfaces                │
│  - Domain Logic                         │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│  Infrastructure Layer                   │
│  - Database Implementations             │
│  - External Services                    │
└─────────────────────────────────────────┘
```

### 依存関係の方向

- **Presentation → Application → Domain ← Infrastructure**
- **Domain層は他の層に依存しない**
- **Infrastructure層はDomain層のインターフェースを実装**

## 実行方法

### 開発環境

```bash
docker compose up --build backend-dev
```

### 本番ビルド

```bash
docker build -t backend:prod ./backend
docker run --rm -p 8080:8080 -e PORT=8080 backend:prod
```

## 環境変数

- `PORT`: サーバーポート（デフォルト: 8080）
- `DATABASE_URL`: TursoデータベースURL

## API エンドポイント

### Health Check

```
GET /healthz
```

### ユーザー管理

```
POST   /api/users              # ユーザー登録
GET    /api/users/:id          # ユーザー情報取得
PUT    /api/users/:id          # ユーザー情報更新
DELETE /api/users/:id          # ユーザー削除
POST   /api/auth/login         # ログイン
POST   /api/auth/logout        # ログアウト
```

### 書籍管理

```
GET    /api/books              # 書籍一覧取得
GET    /api/books/:id          # 書籍詳細取得
POST   /api/books              # 書籍登録
GET    /api/books/search       # 書籍検索
```

### アウトプット管理

```
GET    /api/outputs            # アウトプット一覧取得
GET    /api/outputs/:id        # アウトプット詳細取得
POST   /api/outputs            # アウトプット投稿
PUT    /api/outputs/:id        # アウトプット更新
DELETE /api/outputs/:id        # アウトプット削除
```

### 書評管理

```
GET    /api/reviews            # 書評一覧取得
GET    /api/reviews/:id        # 書評詳細取得
POST   /api/reviews            # 書評投稿
PUT    /api/reviews/:id        # 書評更新
DELETE /api/reviews/:id        # 書評削除
```

## テスト

```bash
go test ./...
```

## 開発ガイドライン

1. **ドメイン駆動設計**: 各ドメイン（user, book, output, review）は独立して管理
2. **インターフェース**: リポジトリはインターフェースで定義し、実装はInfrastructure層
3. **依存性注入**: 各レイヤーは依存をインターフェースで受け取る
4. **エラーハンドリング**: エラーは適切なHTTPステータスコードで返す

