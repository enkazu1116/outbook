package service

import (
	"context"
	"errors"
	"example.com/backend/internal/domain/output/entity"
	"example.com/backend/internal/domain/output/repository"
	"time"
)

var (
	ErrOutputNotFound      = errors.New("アウトプットが見つかりません")
	ErrOutputUnauthorized  = errors.New("このアウトプットへのアクセス権限がありません")
	ErrInvalidVisibility   = errors.New("無効な公開設定です")
	ErrTooManyTags         = errors.New("タグは最大5個までです")
)

// Service アウトプットビジネスロジックを処理
type Service struct {
	repo repository.Repository
}

// NewService 新しいアウトプットサービスを作成
func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateOutput 新しいアウトプットを作成
func (s *Service) CreateOutput(ctx context.Context, output *entity.Output) error {
	// バリデーション
	if err := s.validateOutput(output); err != nil {
		return err
	}

	return s.repo.Create(ctx, output)
}

// GetOutput IDでアウトプットを取得
func (s *Service) GetOutput(ctx context.Context, id int64, userID *int64) (*entity.Output, error) {
	output, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrOutputNotFound
	}

	// 非公開の場合は、作成者のみアクセス可能
	if output.IsPrivate() && (userID == nil || *userID != output.UserID) {
		return nil, ErrOutputUnauthorized
	}

	return output, nil
}

// GetOutputsByUserID ユーザーIDでアウトプットを取得
func (s *Service) GetOutputsByUserID(ctx context.Context, userID int64, filters map[string]interface{}) ([]*entity.Output, error) {
	return s.repo.FindByUserID(ctx, userID, filters)
}

// GetOutputsByBookID 書籍IDでアウトプットを取得（公開されているもののみ）
func (s *Service) GetOutputsByBookID(ctx context.Context, bookID int64, filters map[string]interface{}) ([]*entity.Output, error) {
	return s.repo.FindByBookID(ctx, bookID, filters)
}

// UpdateOutput アウトプットを更新
func (s *Service) UpdateOutput(ctx context.Context, output *entity.Output, userID int64) error {
	// 既存のアウトプットを取得
	existing, err := s.repo.FindByID(ctx, output.ID)
	if err != nil {
		return ErrOutputNotFound
	}

	// 作成者のみ更新可能
	if existing.UserID != userID {
		return ErrOutputUnauthorized
	}

	// バリデーション
	if err := s.validateOutput(output); err != nil {
		return err
	}

	// ユーザーIDを上書きしない
	output.UserID = userID
	output.UpdatedAt = time.Now()

	return s.repo.Update(ctx, output)
}

// DeleteOutput アウトプットを削除（論理削除）
func (s *Service) DeleteOutput(ctx context.Context, id int64, userID int64) error {
	// 既存のアウトプットを取得
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrOutputNotFound
	}

	// 作成者のみ削除可能
	if existing.UserID != userID {
		return ErrOutputUnauthorized
	}

	return s.repo.Delete(ctx, id)
}

// validateOutput アウトプットのバリデーション
func (s *Service) validateOutput(output *entity.Output) error {
	// 公開設定のバリデーション
	if output.Visibility != entity.VisibilityPublic &&
		output.Visibility != entity.VisibilityPrivate &&
		output.Visibility != entity.VisibilityTitleOnly {
		return ErrInvalidVisibility
	}

	// タグの数のバリデーション
	if len(output.Tags) > 5 {
		return ErrTooManyTags
	}

	return nil
}

