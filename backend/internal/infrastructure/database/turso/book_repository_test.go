package database

import (
	"context"
	"example.com/backend/internal/domain/book/entity"
	"testing"
	"time"
)

func TestBookRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewBookRepository(db)

	publishedDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	book := &entity.Book{
		Title:          "Go言語で学ぶ並行処理",
		Author:         "John Doe",
		Publisher:      "技術出版社",
		PublishedDate:  &publishedDate,
		ISBN:           "978-4-1234-5678-9",
		CoverImageURL:  "https://example.com/cover.jpg",
		Description:    "Go言語の並行処理について詳しく解説",
		PrimaryLanguage: "Go",
		Categories:     []string{"programming", "concurrency"},
		Level:          "intermediate",
		PageCount:      300,
		AverageRating:  4.5,
		ReviewsCount:   10,
		OutputsCount:   5,
	}

	ctx := context.Background()
	err := repo.Create(ctx, book)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if book.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func TestBookRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewBookRepository(db)

	book := &entity.Book{
		Title:  "テスト書籍",
		Author: "著者名",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	found, err := repo.FindByID(ctx, book.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.Title != book.Title {
		t.Errorf("Expected Title %s, got %s", book.Title, found.Title)
	}
	if found.Author != book.Author {
		t.Errorf("Expected Author %s, got %s", book.Author, found.Author)
	}
}

func TestBookRepository_FindByISBN(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewBookRepository(db)

	book := &entity.Book{
		Title: "テスト書籍",
		Author: "著者名",
		ISBN:  "978-4-1234-5678-9",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	found, err := repo.FindByISBN(ctx, "978-4-1234-5678-9")
	if err != nil {
		t.Fatalf("FindByISBN failed: %v", err)
	}

	if found.ISBN != "978-4-1234-5678-9" {
		t.Errorf("Expected ISBN 978-4-1234-5678-9, got %s", found.ISBN)
	}
}

func TestBookRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewBookRepository(db)

	books := []*entity.Book{
		{Title: "Go言語入門", Author: "Go Masters", ISBN: "978-4-1000-0001-0"},
		{Title: "TypeScript完全ガイド", Author: "JS Experts", ISBN: "978-4-1000-0002-0"},
		{Title: "Pythonデータサイエンス", Author: "Data Scientists", ISBN: "978-4-1000-0003-0"},
	}

	ctx := context.Background()
	for _, book := range books {
		if err := repo.Create(ctx, book); err != nil {
			t.Fatalf("Failed to create book: %v", err)
		}
	}

	// "Go"で検索
	results, err := repo.Search(ctx, "Go", nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected at least one result for 'Go'")
	}

	found := false
	for _, book := range results {
		if book.Title == "Go言語入門" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'Go言語入門'")
	}
}

func TestBookRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewBookRepository(db)

	book := &entity.Book{
		Title:  "テスト書籍",
		Author: "著者名",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	book.Title = "更新されたタイトル"
	book.AverageRating = 4.8

	err := repo.Update(ctx, book)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.FindByID(ctx, book.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if updated.Title != "更新されたタイトル" {
		t.Errorf("Expected Title 更新されたタイトル, got %s", updated.Title)
	}
	if updated.AverageRating != 4.8 {
		t.Errorf("Expected AverageRating 4.8, got %f", updated.AverageRating)
	}
}

func TestBookRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewBookRepository(db)

	book := &entity.Book{
		Title:  "テスト書籍",
		Author: "著者名",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	err := repo.Delete(ctx, book.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByID(ctx, book.ID)
	if err == nil {
		t.Error("Expected error when finding deleted book")
	}
}

