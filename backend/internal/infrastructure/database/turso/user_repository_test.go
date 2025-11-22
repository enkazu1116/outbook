package database

import (
	"context"
	"testing"

	"example.com/backend/internal/domain/user/entity"
)

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:              "test@example.com",
		PasswordHash:       "hashed_password",
		Name:               "テストユーザー",
		SkillLevel:         "middle",
		YearsOfExperience:  3,
		Specialties:        []string{"backend", "frontend"},
		FavoriteLanguages:  []string{"Go", "TypeScript"},
		IsActive:           true,
		EmailNotification:  true,
	}

	ctx := context.Background()
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected ID to be set")
	}
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewUserRepository(db)

	// テストデータを作成
	user := &entity.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "テストユーザー",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// FindByIDで取得
	found, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, found.ID)
	}
	if found.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, found.Email)
	}
	if found.Name != user.Name {
		t.Errorf("Expected Name %s, got %s", user.Name, found.Name)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "テストユーザー",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	found, err := repo.FindByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if found.Email != "test@example.com" {
		t.Errorf("Expected Email test@example.com, got %s", found.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "テストユーザー",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 更新
	user.Name = "更新されたユーザー"
	user.Bio = "自己紹介文"

	err := repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// 検証
	updated, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if updated.Name != "更新されたユーザー" {
		t.Errorf("Expected Name 更新されたユーザー, got %s", updated.Name)
	}
	if updated.Bio != "自己紹介文" {
		t.Errorf("Expected Bio 自己紹介文, got %s", updated.Bio)
	}
	if updated.UpdatedAt.Before(user.CreatedAt) {
		t.Error("Expected UpdatedAt to be after CreatedAt")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "テストユーザー",
	}

	ctx := context.Background()
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 削除
	err := repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// 削除後は見つからないはず
	_, err = repo.FindByID(ctx, user.ID)
	if err == nil {
		t.Error("Expected error when finding deleted user")
	}
}

func TestUserRepository_JSONFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "テストユーザー",
		Specialties:  []string{"backend", "frontend"},
		FavoriteLanguages: []string{"Go", "TypeScript", "Python"},
	}

	ctx := context.Background()
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	found, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if len(found.Specialties) != 2 {
		t.Errorf("Expected 2 specialties, got %d", len(found.Specialties))
	}
	if len(found.FavoriteLanguages) != 3 {
		t.Errorf("Expected 3 favorite languages, got %d", len(found.FavoriteLanguages))
	}
	if found.FavoriteLanguages[0] != "Go" {
		t.Errorf("Expected first language Go, got %s", found.FavoriteLanguages[0])
	}
}

