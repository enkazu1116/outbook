package dto

import (
	"time"

	"example.com/backend/internal/domain/output/entity"
)

// StandardResponse は標準的なAPIレスポンス構造
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationRequest はページネーションパラメータ
type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// PaginationResponse はページネーション付きレスポンス
type PaginationResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// Timestamps は共通のタイムスタンプフィールド
type Timestamps struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CreateOutputRequest アウトプット作成リクエスト（UserIDは認証から取得）
type CreateOutputRequest struct {
	BookID              int64    `json:"book_id"`
	Title               string   `json:"title"`
	Content             string   `json:"content"`
	Visibility          string   `json:"visibility"`
	ReadingFinishedDate *string  `json:"reading_finished_date,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	IsDraft             bool     `json:"is_draft"`
}

// ToEntity CreateOutputRequestからOutputエンティティを作成
func (r *CreateOutputRequest) ToEntity(userID int64) *entity.Output {
	output := &entity.Output{
		UserID:     userID,
		BookID:     r.BookID,
		Title:      r.Title,
		Content:    r.Content,
		Visibility: r.Visibility,
		Tags:       r.Tags,
		IsDraft:    r.IsDraft,
		ViewsCount: 0,
	}

	// デフォルトの公開設定を設定
	if output.Visibility == "" {
		output.Visibility = entity.VisibilityPublic
	}

	return output
}

// UpdateOutputRequest アウトプット更新リクエスト
type UpdateOutputRequest struct {
	BookID              *int64    `json:"book_id,omitempty"`
	Title               *string   `json:"title,omitempty"`
	Content             *string   `json:"content,omitempty"`
	Visibility          *string   `json:"visibility,omitempty"`
	ReadingFinishedDate *string   `json:"reading_finished_date,omitempty"`
	Tags                *[]string `json:"tags,omitempty"`
	IsDraft             *bool     `json:"is_draft,omitempty"`
}

// ApplyToEntity UpdateOutputRequestのフィールドを既存のOutputエンティティに適用
func (r *UpdateOutputRequest) ApplyToEntity(output *entity.Output) {
	if r.BookID != nil {
		output.BookID = *r.BookID
	}
	if r.Title != nil && *r.Title != "" {
		output.Title = *r.Title
	}
	if r.Content != nil {
		output.Content = *r.Content
	}
	if r.Visibility != nil && *r.Visibility != "" {
		output.Visibility = *r.Visibility
	}
	if r.Tags != nil {
		output.Tags = *r.Tags
	}
	if r.IsDraft != nil {
		output.IsDraft = *r.IsDraft
	}
}
