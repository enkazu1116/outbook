package entity

import "time"

// Book システム内の技術書を表す
type Book struct {
	ID             int64      `json:"id"`
	Title          string     `json:"title"`
	Author         string     `json:"author"`
	Publisher      string     `json:"publisher"`
	PublishedDate  *time.Time `json:"published_date"`
	ISBN           string     `json:"isbn"`
	CoverImageURL  string     `json:"cover_image_url"`
	Description    string     `json:"description"`
	PrimaryLanguage string    `json:"primary_language"`
	Categories     []string   `json:"categories"` // カテゴリ（JSON配列）
	Level          string     `json:"level"`      // 難易度（beginner/intermediate/advanced）
	PageCount      int        `json:"page_count"`
	AverageRating  float64    `json:"average_rating"`
	ReviewsCount   int        `json:"reviews_count"`
	OutputsCount   int        `json:"outputs_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

