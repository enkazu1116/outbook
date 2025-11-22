package repository

import (
	"context"
	"example.com/backend/internal/domain/review/entity"
)

// Repository 書評データアクセスのインターフェースを定義
type Repository interface {
	Create(ctx context.Context, review *entity.Review) error
	FindByID(ctx context.Context, id int64) (*entity.Review, error)
	FindByBookID(ctx context.Context, bookID int64, filters map[string]interface{}) ([]*entity.Review, error)
	FindByUserIDAndBookID(ctx context.Context, userID, bookID int64) (*entity.Review, error)
	Update(ctx context.Context, review *entity.Review) error
	Delete(ctx context.Context, id int64) error
}

