package entity

import "time"

// Review ユーザーが書籍について投稿した書評を表す
type Review struct {
	ID                   int64      `json:"id"`
	UserID               int64      `json:"user_id"`
	BookID               int64      `json:"book_id"`
	Rating               int        `json:"rating"` // 評価スコア（1-5）
	Content              string     `json:"content"`
	RecommendedLevels    []string   `json:"recommended_levels"` // おすすめレベル（JSON配列）
	RecommendedSpecialties []string `json:"recommended_specialties"` // おすすめ職域（JSON配列）
	HelpfulCount         int        `json:"helpful_count"`
	IsHidden             bool       `json:"-"`        // 非表示フラグ（JSONに含めない）
	IsDeleted            bool       `json:"-"`        // 論理削除フラグ（JSONに含めない）
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"-"`        // 削除日時（JSONに含めない）
}

