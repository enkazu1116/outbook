package repository

import (
	"context"
	"example.com/backend/internal/domain/book/entity"
)

// Repository 書籍データアクセスのインターフェースを定義
type Repository interface {
	Create(ctx context.Context, book *entity.Book) error
	FindByID(ctx context.Context, id int64) (*entity.Book, error)
	FindByISBN(ctx context.Context, isbn string) (*entity.Book, error)
	Search(ctx context.Context, query string, filters map[string]interface{}) ([]*entity.Book, error)
	Update(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id int64) error
}

