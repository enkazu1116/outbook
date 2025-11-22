package repository

import (
	"context"
	"example.com/backend/internal/domain/user/entity"
)

// Repository ユーザーデータアクセスのインターフェースを定義
type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
}

