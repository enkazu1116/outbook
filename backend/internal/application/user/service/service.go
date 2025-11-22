package service

import (
	"context"
	"errors"
	"example.com/backend/internal/domain/user/entity"
	"example.com/backend/internal/domain/user/repository"
)

var (
	ErrUserNotFound       = errors.New("ユーザーが見つかりません")
	ErrInvalidCredentials = errors.New("認証情報が無効です")
)

// Service ユーザービジネスロジックを処理
type Service struct {
	repo repository.Repository
}

// NewService 新しいユーザーサービスを作成
func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// RegisterUser 新しいユーザーアカウントを作成
func (s *Service) RegisterUser(ctx context.Context, email, passwordHash, name string) (*entity.User, error) {
	// ビジネスロジック
	user := &entity.User{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		IsActive:     true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser ユーザーの認証情報を検証
func (s *Service) AuthenticateUser(ctx context.Context, email, passwordHash string) (*entity.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.PasswordHash != passwordHash {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

