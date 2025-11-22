package repository

import (
	"context"
	"example.com/backend/internal/domain/output/entity"
)

// Repository アウトプットデータアクセスのインターフェースを定義
type Repository interface {
	Create(ctx context.Context, output *entity.Output) error
	FindByID(ctx context.Context, id int64) (*entity.Output, error)
	FindByUserID(ctx context.Context, userID int64, filters map[string]interface{}) ([]*entity.Output, error)
	FindByBookID(ctx context.Context, bookID int64, filters map[string]interface{}) ([]*entity.Output, error)
	Update(ctx context.Context, output *entity.Output) error
	Delete(ctx context.Context, id int64) error
}

