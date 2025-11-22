package database

import (
	"context"
	"testing"
	"time"

	bookentity "example.com/backend/internal/domain/book/entity"
	outputentity "example.com/backend/internal/domain/output/entity"
	userentity "example.com/backend/internal/domain/user/entity"
)

func TestOutputRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// 依存するデータを作成
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	
	user := &struct {
		Email        string
		PasswordHash string
		Name         string
	}{
		Email:        "test@example.com",
		PasswordHash: "hash",
		Name:         "Test User",
	}
	
	book := &struct {
		Title  string
		Author string
	}{
		Title:  "Test Book",
		Author: "Author",
	}
	
	ctx := context.Background()
	
	// ユーザーと書籍を作成してIDを取得
	userEntity := &userentity.User{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	}
	if err := userRepo.Create(ctx, userEntity); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	bookEntity := &bookentity.Book{
		Title:  book.Title,
		Author: book.Author,
	}
	if err := bookRepo.Create(ctx, bookEntity); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	// アウトプットのテスト
	outputRepo := NewOutputRepository(db)
	finishedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	output := &outputentity.Output{
		UserID:              userEntity.ID,
		BookID:              bookEntity.ID,
		Title:               "テストアウトプット",
		Content:             "# 内容\n\nこの本のアウトプットです。",
		Visibility:          "public",
		ReadingFinishedDate: &finishedDate,
		Tags:                []string{"学習", "Go", "並行処理"},
		IsDraft:             false,
		ViewsCount:          0,
	}
	
	err := outputRepo.Create(ctx, output)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if output.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func TestOutputRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	outputRepo := NewOutputRepository(db)
	
	// データセットアップ
	user := &userentity.User{
		Email:        "test@example.com",
		PasswordHash: "hash",
		Name:         "Test User",
	}
	book := &bookentity.Book{
		Title:  "Test Book",
		Author: "Author",
	}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	// アウトプットを作成
	outputs := []*outputentity.Output{
		{UserID: user.ID, BookID: book.ID, Title: "Output 1", Content: "Content 1", Visibility: "public"},
		{UserID: user.ID, BookID: book.ID, Title: "Output 2", Content: "Content 2", Visibility: "public"},
	}
	
	for _, output := range outputs {
		if err := outputRepo.Create(ctx, output); err != nil {
			t.Fatalf("Failed to create output: %v", err)
		}
	}
	
	// 検索
	found, err := outputRepo.FindByUserID(ctx, user.ID, nil)
	if err != nil {
		t.Fatalf("FindByUserID failed: %v", err)
	}
	
	if len(found) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(found))
	}
}

func TestOutputRepository_FindByBookID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	outputRepo := NewOutputRepository(db)
	
	user1 := &userentity.User{
		Email:        "user1@example.com",
		PasswordHash: "hash1",
		Name:         "User 1",
	}
	user2 := &userentity.User{
		Email:        "user2@example.com",
		PasswordHash: "hash2",
		Name:         "User 2",
	}
	book := &bookentity.Book{
		Title:  "Test Book",
		Author: "Author",
	}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user1); err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}
	if err := userRepo.Create(ctx, user2); err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	// 各ユーザーが同じ書籍のアウトプットを作成
	if err := outputRepo.Create(ctx, &outputentity.Output{UserID: user1.ID, BookID: book.ID, Title: "Output 1", Content: "Content 1", Visibility: "public"}); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}
	if err := outputRepo.Create(ctx, &outputentity.Output{UserID: user2.ID, BookID: book.ID, Title: "Output 2", Content: "Content 2", Visibility: "public"}); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}
	
	// 非公開のアウトプットは取得されない
	if err := outputRepo.Create(ctx, &outputentity.Output{UserID: user1.ID, BookID: book.ID, Title: "Private Output", Content: "Private Content", Visibility: "private"}); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}
	
	found, err := outputRepo.FindByBookID(ctx, book.ID, nil)
	if err != nil {
		t.Fatalf("FindByBookID failed: %v", err)
	}
	
	if len(found) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(found))
	}
}

func TestOutputRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	outputRepo := NewOutputRepository(db)
	
	user := &userentity.User{Email: "test@example.com", PasswordHash: "hash", Name: "Test User"}
	book := &bookentity.Book{Title: "Test Book", Author: "Author"}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	output := &outputentity.Output{
		UserID:     user.ID,
		BookID:     book.ID,
		Title:      "Original Title",
		Content:    "Original Content",
		Visibility: "public",
	}
	if err := outputRepo.Create(ctx, output); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}
	
	output.Title = "Updated Title"
	err := outputRepo.Update(ctx, output)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	updated, err := outputRepo.FindByID(ctx, output.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	
	if updated.Title != "Updated Title" {
		t.Errorf("Expected Title Updated Title, got %s", updated.Title)
	}
}

func TestOutputRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	outputRepo := NewOutputRepository(db)
	
	user := &userentity.User{Email: "test@example.com", PasswordHash: "hash", Name: "Test User"}
	book := &bookentity.Book{Title: "Test Book", Author: "Author"}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	output := &outputentity.Output{
		UserID:     user.ID,
		BookID:     book.ID,
		Title:      "Test Output",
		Content:    "Test Content",
		Visibility: "public",
	}
	if err := outputRepo.Create(ctx, output); err != nil {
		t.Fatalf("Failed to create output: %v", err)
	}
	
	err := outputRepo.Delete(ctx, output.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	_, err = outputRepo.FindByID(ctx, output.ID)
	if err == nil {
		t.Error("Expected error when finding deleted output")
	}
}

