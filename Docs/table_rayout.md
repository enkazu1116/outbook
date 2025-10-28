# OutBook テーブル設計書

## 概要
このドキュメントでは、OutBookプロジェクトのデータベース設計を定義します。
データベースはTurso（SQLite互換）を使用します。

## テーブル一覧

### 1. users（ユーザー）
ユーザーアカウント情報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ユーザーID（主キー） |
| email | VARCHAR(255) | NO | - | メールアドレス（ユニーク） |
| password_hash | VARCHAR(255) | NO | - | ハッシュ化されたパスワード |
| nickname | VARCHAR(100) | NO | - | ニックネーム |
| avatar_url | VARCHAR(500) | YES | NULL | プロフィール画像URL |
| bio | TEXT | YES | NULL | 自己紹介 |
| skill_level | VARCHAR(20) | YES | NULL | スキルレベル（junior/middle/senior） |
| years_of_experience | INTEGER | YES | NULL | 経験年数 |
| specialties | TEXT | YES | NULL | 専門分野（JSON配列形式） |
| favorite_languages | TEXT | YES | NULL | 得意な言語（JSON配列形式） |
| is_active | BOOLEAN | NO | true | アカウントの有効/無効 |
| is_deleted | BOOLEAN | NO | false | 論理削除フラグ |
| email_notification | BOOLEAN | NO | true | メール通知設定 |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | DATETIME | NO | CURRENT_TIMESTAMP | 更新日時 |
| deleted_at | DATETIME | YES | NULL | 削除日時 |

**インデックス:**
- PRIMARY KEY (id)
- UNIQUE INDEX idx_users_email ON users(email)
- INDEX idx_users_skill_level ON users(skill_level)
- INDEX idx_users_created_at ON users(created_at)

---

### 2. books（書籍）
技術書の情報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | 書籍ID（主キー） |
| title | VARCHAR(500) | NO | - | 書籍タイトル |
| author | VARCHAR(255) | NO | - | 著者名 |
| publisher | VARCHAR(255) | YES | NULL | 出版社 |
| published_date | DATE | YES | NULL | 出版日 |
| isbn | VARCHAR(20) | YES | NULL | ISBN（ユニーク） |
| cover_image_url | VARCHAR(500) | YES | NULL | カバー画像URL |
| description | TEXT | YES | NULL | 書籍概要 |
| primary_language | VARCHAR(50) | YES | NULL | 主要なプログラミング言語 |
| categories | TEXT | YES | NULL | カテゴリ（JSON配列形式） |
| level | VARCHAR(20) | YES | NULL | 難易度（beginner/intermediate/advanced） |
| page_count | INTEGER | YES | NULL | ページ数 |
| average_rating | DECIMAL(3,2) | YES | 0.00 | 平均評価（1.00〜5.00） |
| reviews_count | INTEGER | NO | 0 | 書評数 |
| outputs_count | INTEGER | NO | 0 | アウトプット数 |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | DATETIME | NO | CURRENT_TIMESTAMP | 更新日時 |

**インデックス:**
- PRIMARY KEY (id)
- UNIQUE INDEX idx_books_isbn ON books(isbn) WHERE isbn IS NOT NULL
- INDEX idx_books_title ON books(title)
- INDEX idx_books_author ON books(author)
- INDEX idx_books_language ON books(primary_language)
- INDEX idx_books_level ON books(level)
- INDEX idx_books_rating ON books(average_rating)

---

### 3. outputs（アウトプット）
ユーザーが投稿した技術書のアウトプットを管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | アウトプットID（主キー） |
| user_id | INTEGER | NO | - | ユーザーID（外部キー） |
| book_id | INTEGER | NO | - | 書籍ID（外部キー） |
| title | VARCHAR(100) | NO | - | アウトプットタイトル |
| content | TEXT | NO | - | アウトプット本文（Markdown） |
| visibility | VARCHAR(20) | NO | 'public' | 公開設定（public/title_only/private） |
| reading_finished_date | DATE | YES | NULL | 読了日 |
| tags | TEXT | YES | NULL | タグ（JSON配列形式、最大5個） |
| is_draft | BOOLEAN | NO | false | 下書きフラグ |
| is_deleted | BOOLEAN | NO | false | 論理削除フラグ |
| views_count | INTEGER | NO | 0 | 閲覧数 |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | DATETIME | NO | CURRENT_TIMESTAMP | 更新日時 |
| deleted_at | DATETIME | YES | NULL | 削除日時 |

**インデックス:**
- PRIMARY KEY (id)
- INDEX idx_outputs_user_id ON outputs(user_id)
- INDEX idx_outputs_book_id ON outputs(book_id)
- INDEX idx_outputs_visibility ON outputs(visibility)
- INDEX idx_outputs_created_at ON outputs(created_at)
- INDEX idx_outputs_is_draft ON outputs(is_draft)

**外部キー制約:**
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
- FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE

---

### 4. reviews（書評）
ユーザーが投稿した書評を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | 書評ID（主キー） |
| user_id | INTEGER | NO | - | ユーザーID（外部キー） |
| book_id | INTEGER | NO | - | 書籍ID（外部キー） |
| rating | INTEGER | NO | - | 評価スコア（1〜5） |
| content | TEXT | NO | - | 書評本文（最大1000文字） |
| recommended_levels | TEXT | YES | NULL | おすすめレベル（JSON配列形式） |
| recommended_specialties | TEXT | YES | NULL | おすすめ職域（JSON配列形式） |
| helpful_count | INTEGER | NO | 0 | 役立ち度カウント（いいね数） |
| is_hidden | BOOLEAN | NO | false | 非表示フラグ（モデレーション用） |
| is_deleted | BOOLEAN | NO | false | 論理削除フラグ |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | DATETIME | NO | CURRENT_TIMESTAMP | 更新日時 |
| deleted_at | DATETIME | YES | NULL | 削除日時 |

**インデックス:**
- PRIMARY KEY (id)
- UNIQUE INDEX idx_reviews_user_book ON reviews(user_id, book_id)
- INDEX idx_reviews_book_id ON reviews(book_id)
- INDEX idx_reviews_rating ON reviews(rating)
- INDEX idx_reviews_helpful_count ON reviews(helpful_count)
- INDEX idx_reviews_created_at ON reviews(created_at)

**外部キー制約:**
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
- FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE

---

### 5. reports（通報）
不適切な書評の通報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | 通報ID（主キー） |
| review_id | INTEGER | NO | - | 書評ID（外部キー） |
| reporter_user_id | INTEGER | NO | - | 通報者ユーザーID（外部キー） |
| reason | VARCHAR(50) | NO | - | 通報理由（spam/harassment/inappropriate/other） |
| detail | TEXT | YES | NULL | 通報詳細 |
| status | VARCHAR(20) | NO | 'pending' | ステータス（pending/reviewed/resolved/rejected） |
| admin_note | TEXT | YES | NULL | 管理者メモ |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 通報日時 |
| resolved_at | DATETIME | YES | NULL | 解決日時 |

**インデックス:**
- PRIMARY KEY (id)
- INDEX idx_reports_review_id ON reports(review_id)
- INDEX idx_reports_reporter_id ON reports(reporter_user_id)
- INDEX idx_reports_status ON reports(status)
- INDEX idx_reports_created_at ON reports(created_at)

**外部キー制約:**
- FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE
- FOREIGN KEY (reporter_user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 6. review_helpful（書評の役立ち度）
書評に対する「役立った」評価を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ID（主キー） |
| review_id | INTEGER | NO | - | 書評ID（外部キー） |
| user_id | INTEGER | NO | - | ユーザーID（外部キー） |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |

**インデックス:**
- PRIMARY KEY (id)
- UNIQUE INDEX idx_review_helpful_unique ON review_helpful(review_id, user_id)
- INDEX idx_review_helpful_review_id ON review_helpful(review_id)

**外部キー制約:**
- FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 7. user_api_keys（ユーザーAPIキー）
AIレビュー機能用のユーザーAPIキーを暗号化して保存するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ID（主キー） |
| user_id | INTEGER | NO | - | ユーザーID（外部キー） |
| provider | VARCHAR(50) | NO | - | APIプロバイダー（openai/anthropic） |
| encrypted_api_key | TEXT | NO | - | 暗号化されたAPIキー |
| is_active | BOOLEAN | NO | true | 有効/無効フラグ |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | DATETIME | NO | CURRENT_TIMESTAMP | 更新日時 |

**インデックス:**
- PRIMARY KEY (id)
- UNIQUE INDEX idx_user_api_keys_unique ON user_api_keys(user_id, provider)

**外部キー制約:**
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 8. ai_reviews（AIレビュー履歴）
AIによるアウトプットレビューの履歴を保存するテーブル（将来実装）

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ID（主キー） |
| output_id | INTEGER | NO | - | アウトプットID（外部キー） |
| user_id | INTEGER | NO | - | ユーザーID（外部キー） |
| provider | VARCHAR(50) | NO | - | AIプロバイダー（openai/anthropic） |
| model | VARCHAR(50) | NO | - | 使用モデル名 |
| feedback | TEXT | NO | - | フィードバック内容 |
| tokens_used | INTEGER | YES | NULL | 使用トークン数 |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |

**インデックス:**
- PRIMARY KEY (id)
- INDEX idx_ai_reviews_output_id ON ai_reviews(output_id)
- INDEX idx_ai_reviews_user_id ON ai_reviews(user_id)

**外部キー制約:**
- FOREIGN KEY (output_id) REFERENCES outputs(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 9. search_bookmarks（検索条件ブックマーク）
ユーザーが保存した検索条件を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ID（主キー） |
| user_id | INTEGER | NO | - | ユーザーID（外部キー） |
| name | VARCHAR(100) | NO | - | ブックマーク名 |
| search_params | TEXT | NO | - | 検索条件（JSON形式） |
| is_default | BOOLEAN | NO | false | デフォルト検索条件フラグ |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | DATETIME | NO | CURRENT_TIMESTAMP | 更新日時 |

**インデックス:**
- PRIMARY KEY (id)
- INDEX idx_search_bookmarks_user_id ON search_bookmarks(user_id)

**外部キー制約:**
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

## 将来実装予定のテーブル

### 10. follows（フォロー関係）
ユーザー間のフォロー関係を管理するテーブル（将来実装）

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ID（主キー） |
| follower_user_id | INTEGER | NO | - | フォローするユーザーID |
| following_user_id | INTEGER | NO | - | フォローされるユーザーID |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | フォロー日時 |

**インデックス:**
- PRIMARY KEY (id)
- UNIQUE INDEX idx_follows_unique ON follows(follower_user_id, following_user_id)
- INDEX idx_follows_follower ON follows(follower_user_id)
- INDEX idx_follows_following ON follows(following_user_id)

---

### 11. output_likes（アウトプットいいね）
アウトプットへのいいねを管理するテーブル（将来実装）

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ID（主キー） |
| output_id | INTEGER | NO | - | アウトプットID |
| user_id | INTEGER | NO | - | ユーザーID |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |

**インデックス:**
- PRIMARY KEY (id)
- UNIQUE INDEX idx_output_likes_unique ON output_likes(output_id, user_id)
- INDEX idx_output_likes_output ON output_likes(output_id)

---

### 12. comments（コメント）
アウトプットへのコメントを管理するテーブル（将来実装）

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | INTEGER | NO | AUTO_INCREMENT | ID（主キー） |
| output_id | INTEGER | NO | - | アウトプットID |
| user_id | INTEGER | NO | - | ユーザーID |
| content | TEXT | NO | - | コメント内容 |
| is_deleted | BOOLEAN | NO | false | 論理削除フラグ |
| created_at | DATETIME | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | DATETIME | NO | CURRENT_TIMESTAMP | 更新日時 |
| deleted_at | DATETIME | YES | NULL | 削除日時 |

**インデックス:**
- PRIMARY KEY (id)
- INDEX idx_comments_output_id ON comments(output_id)
- INDEX idx_comments_user_id ON comments(user_id)

---

## ER図（テキスト表現）

```
users
  ├── 1:N → outputs
  ├── 1:N → reviews
  ├── 1:N → reports
  ├── 1:N → review_helpful
  ├── 1:N → user_api_keys
  ├── 1:N → ai_reviews
  └── 1:N → search_bookmarks

books
  ├── 1:N → outputs
  └── 1:N → reviews

outputs
  ├── N:1 → users
  ├── N:1 → books
  └── 1:N → ai_reviews

reviews
  ├── N:1 → users
  ├── N:1 → books
  ├── 1:N → reports
  └── 1:N → review_helpful

reports
  ├── N:1 → reviews
  └── N:1 → users (reporter)

review_helpful
  ├── N:1 → reviews
  └── N:1 → users

user_api_keys
  └── N:1 → users

ai_reviews
  ├── N:1 → outputs
  └── N:1 → users

search_bookmarks
  └── N:1 → users
```

## 注意事項

### データ型について
- Turso（SQLite）では、一部のデータ型が実際には内部的にマッピングされます
  - `VARCHAR` → `TEXT`
  - `DATETIME` → `TEXT` (ISO8601形式)
  - `BOOLEAN` → `INTEGER` (0/1)
  - `DECIMAL` → `REAL`

### JSON配列の格納
以下のカラムはJSON配列として文字列を格納します：
- `users.specialties`: `["frontend", "backend"]`
- `users.favorite_languages`: `["Go", "TypeScript", "Python"]`
- `books.categories`: `["web", "infrastructure"]`
- `outputs.tags`: `["Go", "並行処理", "パフォーマンス"]`
- `reviews.recommended_levels`: `["beginner", "intermediate"]`
- `reviews.recommended_specialties`: `["backend", "infrastructure"]`

### インデックス戦略
- 検索頻度の高いカラムにインデックスを設定
- 複合ユニークキー制約で重複データを防止
- 外部キーカラムには必ずインデックスを設定

### セキュリティ考慮事項
- `password_hash`: bcryptでハッシュ化して保存
- `encrypted_api_key`: AES-256等で暗号化して保存
- 論理削除を採用し、データの完全性を保持

### パフォーマンス考慮事項
- `books.average_rating`, `books.reviews_count`, `books.outputs_count` は非正規化して保存（検索・ソート性能向上）
- `reviews.helpful_count` も同様に非正規化
- ランキング算出時はこれらのカラムを活用

---

## 変更履歴
| 日付 | バージョン | 変更内容 | 担当者 |
|------|-----------|---------|--------|
| 2025-10-28 | 1.0 | 初版作成 | - |

