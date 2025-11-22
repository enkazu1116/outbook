package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql" // Tursoドライバー
)

// Client データベース接続を管理
type Client struct {
	db *sql.DB
}

// NewClient 新しいデータベースクライアントを作成
func NewClient() (*Client, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL環境変数が設定されていません")
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("データベースのオープンに失敗しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベースへのpingに失敗しました: %w", err)
	}

	return &Client{db: db}, nil
}

// DB 基盤となるデータベース接続を返す
func (c *Client) DB() *sql.DB {
	return c.db
}

// Close データベース接続を閉じる
func (c *Client) Close() error {
	return c.db.Close()
}

