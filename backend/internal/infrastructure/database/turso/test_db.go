package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB テスト用のデータベース接続をセットアップ
func setupTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// ローカルテスト用のSQLiteデータベース
		dbURL = "file::memory:?cache=shared"
	}

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		t.Fatalf("データベースのオープンに失敗: %v", err)
	}

	// テーブルを作成
	if err := createTestTables(db); err != nil {
		t.Fatalf("テーブルの作成に失敗: %v", err)
	}

	return db
}

// createTestTables テスト用のテーブルを作成
func createTestTables(db *sql.DB) error {
	// SQLiteでは外部キー制約を有効にする必要がある
	_, _ = db.Exec("PRAGMA foreign_keys = ON;")

	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL,
			avatar_url TEXT,
			bio TEXT,
			skill_level TEXT,
			years_of_experience INTEGER,
			specialties TEXT,
			favorite_languages TEXT,
			is_active INTEGER NOT NULL DEFAULT 1,
			is_deleted INTEGER NOT NULL DEFAULT 0,
			email_notification INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME
		);

		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			publisher TEXT,
			published_date DATE,
			isbn TEXT UNIQUE,
			cover_image_url TEXT,
			description TEXT,
			primary_language TEXT,
			categories TEXT,
			level TEXT,
			page_count INTEGER,
			average_rating REAL DEFAULT 0.00,
			reviews_count INTEGER NOT NULL DEFAULT 0,
			outputs_count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS outputs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			book_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			visibility TEXT NOT NULL DEFAULT 'public',
			reading_finished_date DATE,
			tags TEXT,
			is_draft INTEGER NOT NULL DEFAULT 0,
			is_deleted INTEGER NOT NULL DEFAULT 0,
			views_count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			book_id INTEGER NOT NULL,
			rating INTEGER NOT NULL,
			content TEXT NOT NULL,
			recommended_levels TEXT,
			recommended_specialties TEXT,
			helpful_count INTEGER NOT NULL DEFAULT 0,
			is_hidden INTEGER NOT NULL DEFAULT 0,
			is_deleted INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME,
			UNIQUE(user_id, book_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
		);
	`

	_, err := db.Exec(schema)
	return err
}
