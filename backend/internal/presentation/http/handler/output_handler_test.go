package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"example.com/backend/internal/application/output/service"
	bookentity "example.com/backend/internal/domain/book/entity"
	outputentity "example.com/backend/internal/domain/output/entity"
	userentity "example.com/backend/internal/domain/user/entity"
	database "example.com/backend/internal/infrastructure/database/turso"
	"example.com/backend/internal/shared/dto"
	"example.com/backend/internal/shared/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestHandler(t *testing.T) (*OutputHandler, *sql.DB, int64, int64, func()) {
	db := setupTestDB(t)
	repo := database.NewOutputRepository(db)
	svc := service.NewService(repo)
	handler := NewOutputHandler(svc)

	// テスト用のユーザーと書籍を作成
	userID, bookID := createTestUserAndBook(t, db)

	return handler, db, userID, bookID, func() { db.Close() }
}

func setupTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "file::memory:?cache=shared"
	}

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		t.Fatalf("データベースのオープンに失敗: %v", err)
	}

	if err := createTestTables(db); err != nil {
		t.Fatalf("テーブルの作成に失敗: %v", err)
	}

	return db
}

func createTestTables(db *sql.DB) error {
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
	`
	_, err := db.Exec(schema)
	return err
}

func createTestUserAndBook(t *testing.T, db *sql.DB) (int64, int64) {
	userRepo := database.NewUserRepository(db)
	bookRepo := database.NewBookRepository(db)

	user := &userentity.User{
		Email:        "test@example.com",
		PasswordHash: "hash",
		Name:         "Test User",
	}
	book := &bookentity.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}

	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	return user.ID, book.ID
}

func createContextWithUserID(userID int64) context.Context {
	return context.WithValue(context.Background(), middleware.UserIDKey, userID)
}

func TestOutputHandler_CreateOutput(t *testing.T) {
	handler, _, userID, bookID, cleanup := setupTestHandler(t)
	defer cleanup()

	reqBody := dto.CreateOutputRequest{
		BookID:     bookID,
		Title:      "テストアウトプット",
		Content:    "# 内容\n\nMarkdown形式のコンテンツです。",
		Visibility: outputentity.VisibilityPublic,
		Tags:       []string{"学習", "Go"},
		IsDraft:    false,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/outputs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := createContextWithUserID(userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateOutput(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response dto.StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}

	outputData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected output data, got %T", response.Data)
	}

	if outputData["title"] != "テストアウトプット" {
		t.Errorf("Expected title テストアウトプット, got %v", outputData["title"])
	}
}

func TestOutputHandler_CreateOutput_Unauthorized(t *testing.T) {
	handler, _, _, bookID, cleanup := setupTestHandler(t)
	defer cleanup()

	reqBody := dto.CreateOutputRequest{
		BookID:     bookID,
		Title:      "テストアウトプット",
		Content:    "内容",
		Visibility: outputentity.VisibilityPublic,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/outputs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// 認証なしのリクエスト

	w := httptest.NewRecorder()
	handler.CreateOutput(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOutputHandler_GetOutput_Public(t *testing.T) {
	handler, db, userID, bookID, cleanup := setupTestHandler(t)
	defer cleanup()

	// アウトプットを作成
	repo := database.NewOutputRepository(db)
	ctx := context.Background()
	output := &outputentity.Output{
		UserID:     userID,
		BookID:     bookID,
		Title:      "公開アウトプット",
		Content:    "内容",
		Visibility: outputentity.VisibilityPublic,
	}
	if err := repo.Create(ctx, output); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}

	// 作成されたアウトプットのIDを使用
	idStr := strconv.FormatInt(output.ID, 10)

	req := httptest.NewRequest(http.MethodGet, "/api/outputs/"+idStr, nil)
	req.URL.Path = "/outputs/" + idStr

	w := httptest.NewRecorder()
	handler.GetOutput(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}
}

func TestOutputHandler_GetOutput_Private_WithAuth(t *testing.T) {
	handler, db, userID, bookID, cleanup := setupTestHandler(t)
	defer cleanup()

	// 非公開アウトプットを作成
	repo := database.NewOutputRepository(db)
	ctx := context.Background()
	output := &outputentity.Output{
		UserID:     userID,
		BookID:     bookID,
		Title:      "非公開アウトプット",
		Content:    "内容",
		Visibility: outputentity.VisibilityPrivate,
	}
	if err := repo.Create(ctx, output); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}

	idStr := strconv.FormatInt(output.ID, 10)
	req := httptest.NewRequest(http.MethodGet, "/api/outputs/"+idStr, nil)
	req.URL.Path = "/outputs/" + idStr
	ctxWithAuth := createContextWithUserID(userID)
	req = req.WithContext(ctxWithAuth)

	w := httptest.NewRecorder()
	handler.GetOutput(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOutputHandler_GetOutput_Private_WithoutAuth(t *testing.T) {
	handler, db, userID, bookID, cleanup := setupTestHandler(t)
	defer cleanup()

	// 非公開アウトプットを作成
	repo := database.NewOutputRepository(db)
	ctx := context.Background()
	output := &outputentity.Output{
		UserID:     userID,
		BookID:     bookID,
		Title:      "非公開アウトプット",
		Content:    "内容",
		Visibility: outputentity.VisibilityPrivate,
	}
	if err := repo.Create(ctx, output); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}

	idStr := strconv.FormatInt(output.ID, 10)
	req := httptest.NewRequest(http.MethodGet, "/api/outputs/"+idStr, nil)
	req.URL.Path = "/outputs/" + idStr
	// 認証なしのリクエスト

	w := httptest.NewRecorder()
	handler.GetOutput(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestOutputHandler_UpdateOutput(t *testing.T) {
	handler, db, userID, bookID, cleanup := setupTestHandler(t)
	defer cleanup()

	// アウトプットを作成
	repo := database.NewOutputRepository(db)
	ctx := context.Background()
	output := &outputentity.Output{
		UserID:     userID,
		BookID:     bookID,
		Title:      "元のタイトル",
		Content:    "元の内容",
		Visibility: outputentity.VisibilityPublic,
	}
	if err := repo.Create(ctx, output); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}

	newTitle := "更新後のタイトル"
	updateReq := dto.UpdateOutputRequest{
		Title: &newTitle,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/outputs/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.URL.Path = "/outputs/1"
	ctxWithAuth := createContextWithUserID(userID)
	req = req.WithContext(ctxWithAuth)

	w := httptest.NewRecorder()
	handler.UpdateOutput(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.StandardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}
}

func TestOutputHandler_DeleteOutput(t *testing.T) {
	handler, db, userID, bookID, cleanup := setupTestHandler(t)
	defer cleanup()

	// アウトプットを作成
	repo := database.NewOutputRepository(db)
	ctx := context.Background()
	output := &outputentity.Output{
		UserID:     userID,
		BookID:     bookID,
		Title:      "削除対象",
		Content:    "内容",
		Visibility: outputentity.VisibilityPublic,
	}
	if err := repo.Create(ctx, output); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/outputs/1", nil)
	req.URL.Path = "/outputs/1"
	ctxWithAuth := createContextWithUserID(userID)
	req = req.WithContext(ctxWithAuth)

	w := httptest.NewRecorder()
	handler.DeleteOutput(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
