package entity

import "time"

// User システム内のユーザーを表す
type User struct {
	ID                 int64     `json:"id"`
	Email              string    `json:"email"`
	PasswordHash       string    `json:"-"`                  // パスワードハッシュ（JSONに含めない）
	Name               string    `json:"name"`
	AvatarURL          string    `json:"avatar_url"`
	Bio                string    `json:"bio"`
	SkillLevel         string    `json:"skill_level"`         // スキルレベル（junior/middle/senior）
	YearsOfExperience  int       `json:"years_of_experience"`
	Specialties        []string  `json:"specialties"`         // 専門分野（JSON配列）
	FavoriteLanguages  []string  `json:"favorite_languages"` // 得意な言語（JSON配列）
	IsActive           bool      `json:"is_active"`
	IsDeleted          bool      `json:"-"`                  // 論理削除フラグ（JSONに含めない）
	EmailNotification  bool      `json:"email_notification"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DeletedAt          *time.Time `json:"-"`                 // 削除日時（JSONに含めない）
}

