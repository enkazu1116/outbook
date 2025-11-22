## Backend コンテナ実行手順

### 前提
- Docker / Docker Compose が導入済みであること

### 開発 (ホットリロード)
```bash
docker compose up --build backend-dev
```
- http://localhost:8080/healthz で疎通確認

### 本番用イメージのビルド
```bash
docker build -t backend:prod ./backend
```

### 本番用コンテナの実行
```bash
docker run --rm -p 8080:8080 -e PORT=8080 backend:prod
```

### 構成
- `backend/Dockerfile`: マルチステージ（builder → distroless）で軽量実行
- `backend/Dockerfile.dev`: 開発用。`air` によるホットリロード
- `docker-compose.yml`: 開発サービス定義（ボリュームでコード同期）
- `.dockerignore`: 不要ファイルの除外


