package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"example.com/backend/internal/domain/user/entity"
	"example.com/backend/internal/domain/user/repository"
)

// userRepository Turso用のユーザーリポジトリ実装
type userRepository struct {
	db *sql.DB
}

// NewUserRepository 新しいユーザーリポジトリを作成
func NewUserRepository(db *sql.DB) repository.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, avatar_url, bio, skill_level, 
		                   years_of_experience, specialties, favorite_languages, 
		                   is_active, email_notification, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.AvatarURL,
		user.Bio,
		user.SkillLevel,
		user.YearsOfExperience,
		serializeStringArray(user.Specialties),
		serializeStringArray(user.FavoriteLanguages),
		user.IsActive,
		user.EmailNotification,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, avatar_url, bio, skill_level, 
		       years_of_experience, specialties, favorite_languages, is_active, 
		       email_notification, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND is_deleted = 0
	`

	var user entity.User
	var specialties, favoriteLangs string
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.AvatarURL,
		&user.Bio,
		&user.SkillLevel,
		&user.YearsOfExperience,
		&specialties,
		&favoriteLangs,
		&user.IsActive,
		&user.EmailNotification,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ユーザーが見つかりません")
		}
		return nil, err
	}

	user.Specialties = parseStringArray(specialties)
	user.FavoriteLanguages = parseStringArray(favoriteLangs)
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, avatar_url, bio, skill_level, 
		       years_of_experience, specialties, favorite_languages, is_active, 
		       email_notification, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND is_deleted = 0
	`

	var user entity.User
	var specialties, favoriteLangs string
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.AvatarURL,
		&user.Bio,
		&user.SkillLevel,
		&user.YearsOfExperience,
		&specialties,
		&favoriteLangs,
		&user.IsActive,
		&user.EmailNotification,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("ユーザーが見つかりません")
		}
		return nil, err
	}

	user.Specialties = parseStringArray(specialties)
	user.FavoriteLanguages = parseStringArray(favoriteLangs)
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email = ?, password_hash = ?, name = ?, avatar_url = ?, bio = ?, 
		    skill_level = ?, years_of_experience = ?, specialties = ?, 
		    favorite_languages = ?, is_active = ?, email_notification = ?, 
		    updated_at = ?
		WHERE id = ? AND is_deleted = 0
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.AvatarURL,
		user.Bio,
		user.SkillLevel,
		user.YearsOfExperience,
		serializeStringArray(user.Specialties),
		serializeStringArray(user.FavoriteLanguages),
		user.IsActive,
		user.EmailNotification,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("ユーザーが見つかりません")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE users
		SET is_deleted = 1, deleted_at = ?
		WHERE id = ? AND is_deleted = 0
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("ユーザーが見つかりません")
	}

	return nil
}
