package service

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	bookentity "example.com/backend/internal/domain/book/entity"
	outputentity "example.com/backend/internal/domain/output/entity"
	userentity "example.com/backend/internal/domain/user/entity"
	database "example.com/backend/internal/infrastructure/database/turso"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestService(t *testing.T) (*Service, *sql.DB, func()) {
	db := setupTestDB(t)
	repo := database.NewOutputRepository(db)
	service := NewService(repo)
	return service, db, func() { db.Close() }
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

func TestService_CreateOutput(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	finishedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	output := &outputentity.Output{
		UserID:              userID,
		BookID:              bookID,
		Title:               "テストアウトプット",
		Content:             "# 内容\n\nこの本のアウトプットです。",
		Visibility:          outputentity.VisibilityPublic,
		ReadingFinishedDate: &finishedDate,
		Tags:                []string{"学習", "Go"},
		IsDraft:             false,
		ViewsCount:          0,
	}

	err := svc.CreateOutput(ctx, output)
	if err != nil {
		t.Fatalf("CreateOutput failed: %v", err)
	}

	if output.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func TestService_CreateOutput_InvalidVisibility(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	output := &outputentity.Output{
		UserID:     userID,
		BookID:     bookID,
		Title:      "テストアウトプット",
		Content:    "内容",
		Visibility: "invalid",
	}

	err := svc.CreateOutput(ctx, output)
	if err != ErrInvalidVisibility {
		t.Errorf("Expected ErrInvalidVisibility, got %v", err)
	}
}

func TestService_CreateOutput_TooManyTags(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	output := &outputentity.Output{
		UserID:     userID,
		BookID:     bookID,
		Title:      "テストアウトプット",
		Content:    "内容",
		Visibility: outputentity.VisibilityPublic,
		Tags:       []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"}, // 6個
	}

	err := svc.CreateOutput(ctx, output)
	if err != ErrTooManyTags {
		t.Errorf("Expected ErrTooManyTags, got %v", err)
	}
}

func TestService_GetOutput_Public(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	repo := database.NewOutputRepository(db)
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

	// 認証なしでも取得可能
	found, err := svc.GetOutput(ctx, output.ID, nil)
	if err != nil {
		t.Fatalf("GetOutput failed: %v", err)
	}

	if found.ID != output.ID {
		t.Errorf("Expected ID %d, got %d", output.ID, found.ID)
	}
}

func TestService_GetOutput_Private_WithAuth(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	repo := database.NewOutputRepository(db)
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

	// 作成者は取得可能
	found, err := svc.GetOutput(ctx, output.ID, &userID)
	if err != nil {
		t.Fatalf("GetOutput failed: %v", err)
	}

	if found.ID != output.ID {
		t.Errorf("Expected ID %d, got %d", output.ID, found.ID)
	}
}

func TestService_GetOutput_Private_WithoutAuth(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	repo := database.NewOutputRepository(db)
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

	// 認証なしでは取得不可
	_, err := svc.GetOutput(ctx, output.ID, nil)
	if err != ErrOutputUnauthorized {
		t.Errorf("Expected ErrOutputUnauthorized, got %v", err)
	}
}

func TestService_GetOutput_Private_WrongUser(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	// 別ユーザーを作成
	userRepo := database.NewUserRepository(db)
	otherUser := &userentity.User{
		Email:        "other@example.com",
		PasswordHash: "hash",
		Name:         "Other User",
	}
	ctx := context.Background()
	if err := userRepo.Create(ctx, otherUser); err != nil {
		t.Fatalf("Failed to create other user: %v", err)
	}

	repo := database.NewOutputRepository(db)
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

	// 別ユーザーでは取得不可
	otherUserID := otherUser.ID
	_, err := svc.GetOutput(ctx, output.ID, &otherUserID)
	if err != ErrOutputUnauthorized {
		t.Errorf("Expected ErrOutputUnauthorized, got %v", err)
	}
}

func TestService_UpdateOutput_Success(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	repo := database.NewOutputRepository(db)
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

	output.Title = "更新後のタイトル"
	output.Content = "更新後の内容"

	err := svc.UpdateOutput(ctx, output, userID)
	if err != nil {
		t.Fatalf("UpdateOutput failed: %v", err)
	}

	updated, err := svc.GetOutput(ctx, output.ID, &userID)
	if err != nil {
		t.Fatalf("GetOutput failed: %v", err)
	}

	if updated.Title != "更新後のタイトル" {
		t.Errorf("Expected Title 更新後のタイトル, got %s", updated.Title)
	}
}

func TestService_UpdateOutput_Unauthorized(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	// 別ユーザーを作成
	userRepo := database.NewUserRepository(db)
	otherUser := &userentity.User{
		Email:        "other@example.com",
		PasswordHash: "hash",
		Name:         "Other User",
	}
	ctx := context.Background()
	if err := userRepo.Create(ctx, otherUser); err != nil {
		t.Fatalf("Failed to create other user: %v", err)
	}

	repo := database.NewOutputRepository(db)
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

	output.Title = "更新後のタイトル"

	// 別ユーザーでは更新不可
	otherUserID := otherUser.ID
	err := svc.UpdateOutput(ctx, output, otherUserID)
	if err != ErrOutputUnauthorized {
		t.Errorf("Expected ErrOutputUnauthorized, got %v", err)
	}
}

func TestService_DeleteOutput_Success(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	ctx := context.Background()
	repo := database.NewOutputRepository(db)
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

	err := svc.DeleteOutput(ctx, output.ID, userID)
	if err != nil {
		t.Fatalf("DeleteOutput failed: %v", err)
	}

	// 削除後は取得不可
	_, err = svc.GetOutput(ctx, output.ID, &userID)
	if err != ErrOutputNotFound {
		t.Errorf("Expected ErrOutputNotFound, got %v", err)
	}
}

func TestService_DeleteOutput_Unauthorized(t *testing.T) {
	svc, db, cleanup := setupTestService(t)
	defer cleanup()
	userID, bookID := createTestUserAndBook(t, db)

	// 別ユーザーを作成
	userRepo := database.NewUserRepository(db)
	otherUser := &userentity.User{
		Email:        "other@example.com",
		PasswordHash: "hash",
		Name:         "Other User",
	}
	ctx := context.Background()
	if err := userRepo.Create(ctx, otherUser); err != nil {
		t.Fatalf("Failed to create other user: %v", err)
	}

	repo := database.NewOutputRepository(db)
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

	// 別ユーザーでは削除不可
	otherUserID := otherUser.ID
	err := svc.DeleteOutput(ctx, output.ID, otherUserID)
	if err != ErrOutputUnauthorized {
		t.Errorf("Expected ErrOutputUnauthorized, got %v", err)
	}
}
