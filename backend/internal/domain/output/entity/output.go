package entity

import "time"

// Visibility 公開設定の定数
const (
	VisibilityPublic   = "public"    // 公開
	VisibilityPrivate  = "private"   // 非公開
	VisibilityTitleOnly = "title_only" // タイトルのみ公開
)

// Output ユーザーが書籍について投稿したアウトプット/記事を表す
// Markdown形式で技術書の内容をまとめる
type Output struct {
	ID                 int64      `json:"id"`
	UserID             int64      `json:"user_id"`
	BookID             int64      `json:"book_id"`
	Title              string     `json:"title"`
	Content            string     `json:"content"`             // Markdown形式のコンテンツ
	Visibility         string     `json:"visibility"`          // 公開設定（public/private/title_only）
	ReadingFinishedDate *time.Time `json:"reading_finished_date"`
	Tags               []string   `json:"tags"`                // タグ（JSON配列、最大5個）
	IsDraft            bool       `json:"is_draft"`
	IsDeleted          bool       `json:"-"`                   // 論理削除フラグ（JSONに含めない）
	ViewsCount         int        `json:"views_count"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"-"`                   // 削除日時（JSONに含めない）
}

// IsPublic 公開されているかどうかを判定
func (o *Output) IsPublic() bool {
	return o.Visibility == VisibilityPublic && !o.IsDraft
}

// IsPrivate 非公開かどうかを判定
func (o *Output) IsPrivate() bool {
	return o.Visibility == VisibilityPrivate
}

