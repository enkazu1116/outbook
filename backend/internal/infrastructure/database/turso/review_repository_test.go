package database

import (
	"context"
	"testing"

	bookentity "example.com/backend/internal/domain/book/entity"
	reviewentity "example.com/backend/internal/domain/review/entity"
	userentity "example.com/backend/internal/domain/user/entity"
)

func TestReviewRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	reviewRepo := NewReviewRepository(db)
	
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
	
	review := &reviewentity.Review{
		UserID:                 user.ID,
		BookID:                 book.ID,
		Rating:                 5,
		Content:                "素晴らしい本でした",
		RecommendedLevels:      []string{"beginner", "intermediate"},
		RecommendedSpecialties: []string{"backend", "frontend"},
		HelpfulCount:           0,
		IsHidden:               false,
	}
	
	err := reviewRepo.Create(ctx, review)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	if review.ID == 0 {
		t.Error("Expected ID to be set")
	}
}

func TestReviewRepository_FindByBookID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	reviewRepo := NewReviewRepository(db)
	
	user1 := &userentity.User{Email: "user1@example.com", PasswordHash: "hash1", Name: "User 1"}
	user2 := &userentity.User{Email: "user2@example.com", PasswordHash: "hash2", Name: "User 2"}
	book := &bookentity.Book{Title: "Test Book", Author: "Author"}
	
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
	
	// 複数の書評を作成
	reviews := []*reviewentity.Review{
		{UserID: user1.ID, BookID: book.ID, Rating: 5, Content: "Excellent", HelpfulCount: 10},
		{UserID: user2.ID, BookID: book.ID, Rating: 4, Content: "Good", HelpfulCount: 5},
	}
	
	for _, review := range reviews {
		if err := reviewRepo.Create(ctx, review); err != nil {
			t.Fatalf("Failed to create review: %v", err)
		}
	}
	
	found, err := reviewRepo.FindByBookID(ctx, book.ID, nil)
	if err != nil {
		t.Fatalf("FindByBookID failed: %v", err)
	}
	
	if len(found) != 2 {
		t.Errorf("Expected 2 reviews, got %d", len(found))
	}
	
	// helpful_count順にソートされていることを確認
	if found[0].HelpfulCount < found[1].HelpfulCount {
		t.Error("Expected reviews to be sorted by helpful_count DESC")
	}
}

func TestReviewRepository_FindByUserIDAndBookID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	reviewRepo := NewReviewRepository(db)
	
	user := &userentity.User{Email: "test@example.com", PasswordHash: "hash", Name: "Test User"}
	book := &bookentity.Book{Title: "Test Book", Author: "Author"}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	review := &reviewentity.Review{
		UserID:  user.ID,
		BookID:  book.ID,
		Rating:  4,
		Content: "Good book",
	}
	if err := reviewRepo.Create(ctx, review); err != nil {
		t.Fatalf("Failed to create review: %v", err)
	}
	
	found, err := reviewRepo.FindByUserIDAndBookID(ctx, user.ID, book.ID)
	if err != nil {
		t.Fatalf("FindByUserIDAndBookID failed: %v", err)
	}
	
	if found.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, found.UserID)
	}
	if found.BookID != book.ID {
		t.Errorf("Expected BookID %d, got %d", book.ID, found.BookID)
	}
}

func TestReviewRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	reviewRepo := NewReviewRepository(db)
	
	user := &userentity.User{Email: "test@example.com", PasswordHash: "hash", Name: "Test User"}
	book := &bookentity.Book{Title: "Test Book", Author: "Author"}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	review := &reviewentity.Review{
		UserID:  user.ID,
		BookID:  book.ID,
		Rating:  3,
		Content: "Average",
	}
	if err := reviewRepo.Create(ctx, review); err != nil {
		t.Fatalf("Failed to create review: %v", err)
	}
	
	review.Rating = 5
	review.Content = "Excellent after re-reading"
	
	err := reviewRepo.Update(ctx, review)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	
	updated, err := reviewRepo.FindByID(ctx, review.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	
	if updated.Rating != 5 {
		t.Errorf("Expected Rating 5, got %d", updated.Rating)
	}
	if updated.Content != "Excellent after re-reading" {
		t.Errorf("Expected Content 'Excellent after re-reading', got '%s'", updated.Content)
	}
}

func TestReviewRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	reviewRepo := NewReviewRepository(db)
	
	user := &userentity.User{Email: "test@example.com", PasswordHash: "hash", Name: "Test User"}
	book := &bookentity.Book{Title: "Test Book", Author: "Author"}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	review := &reviewentity.Review{
		UserID:  user.ID,
		BookID:  book.ID,
		Rating:  4,
		Content: "Good book",
	}
	if err := reviewRepo.Create(ctx, review); err != nil {
		t.Fatalf("Failed to create review: %v", err)
	}
	
	err := reviewRepo.Delete(ctx, review.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	_, err = reviewRepo.FindByID(ctx, review.ID)
	if err == nil {
		t.Error("Expected error when finding deleted review")
	}
}

func TestReviewRepository_JSONFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	userRepo := NewUserRepository(db)
	bookRepo := NewBookRepository(db)
	reviewRepo := NewReviewRepository(db)
	
	user := &userentity.User{Email: "test@example.com", PasswordHash: "hash", Name: "Test User"}
	book := &bookentity.Book{Title: "Test Book", Author: "Author"}
	
	ctx := context.Background()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := bookRepo.Create(ctx, book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	
	review := &reviewentity.Review{
		UserID:                 user.ID,
		BookID:                 book.ID,
		Rating:                 5,
		Content:                "Excellent",
		RecommendedLevels:      []string{"beginner", "intermediate", "advanced"},
		RecommendedSpecialties: []string{"backend"},
	}
	
	if err := reviewRepo.Create(ctx, review); err != nil {
		t.Fatalf("Failed to create review: %v", err)
	}
	
	found, err := reviewRepo.FindByID(ctx, review.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	
	if len(found.RecommendedLevels) != 3 {
		t.Errorf("Expected 3 recommended levels, got %d", len(found.RecommendedLevels))
	}
	if len(found.RecommendedSpecialties) != 1 {
		t.Errorf("Expected 1 recommended specialty, got %d", len(found.RecommendedSpecialties))
	}
}

