package service

import (
	"context"
	"errors"
	"example.com/backend/internal/domain/book/entity"
	"example.com/backend/internal/domain/book/repository"
)

var (
	ErrBookNotFound      = errors.New("書籍が見つかりません")
	ErrBookAlreadyExists = errors.New("書籍は既に存在します")
)

// Service 書籍ビジネスロジックを処理
type Service struct {
	repo repository.Repository
}

// NewService 新しい書籍サービスを作成
func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateBook 新しい書籍を作成
func (s *Service) CreateBook(ctx context.Context, book *entity.Book) error {
	// ビジネスロジック
	return s.repo.Create(ctx, book)
}

// GetBook IDで書籍を取得
func (s *Service) GetBook(ctx context.Context, id int64) (*entity.Book, error) {
	book, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrBookNotFound
	}
	return book, nil
}

// SearchBooks クエリとフィルターに基づいて書籍を検索
func (s *Service) SearchBooks(ctx context.Context, query string, filters map[string]interface{}) ([]*entity.Book, error) {
	return s.repo.Search(ctx, query, filters)
}

