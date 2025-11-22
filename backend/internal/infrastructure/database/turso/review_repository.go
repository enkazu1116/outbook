package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"example.com/backend/internal/domain/review/entity"
	"example.com/backend/internal/domain/review/repository"
)

// reviewRepository Turso用の書評リポジトリ実装
type reviewRepository struct {
	db *sql.DB
}

// NewReviewRepository 新しい書評リポジトリを作成
func NewReviewRepository(db *sql.DB) repository.Repository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *entity.Review) error {
	query := `
		INSERT INTO reviews (user_id, book_id, rating, content, recommended_levels, 
		                     recommended_specialties, helpful_count, is_hidden, 
		                     created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		review.UserID,
		review.BookID,
		review.Rating,
		review.Content,
		serializeStringArray(review.RecommendedLevels),
		serializeStringArray(review.RecommendedSpecialties),
		review.HelpfulCount,
		review.IsHidden,
		review.CreatedAt,
		review.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	review.ID = id
	return nil
}

func (r *reviewRepository) FindByID(ctx context.Context, id int64) (*entity.Review, error) {
	query := `
		SELECT id, user_id, book_id, rating, content, recommended_levels, 
		       recommended_specialties, helpful_count, is_hidden, is_deleted, 
		       created_at, updated_at, deleted_at
		FROM reviews
		WHERE id = ? AND is_deleted = 0
	`

	var review entity.Review
	var recommendedLevels, recommendedSpecialties string
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.UserID,
		&review.BookID,
		&review.Rating,
		&review.Content,
		&recommendedLevels,
		&recommendedSpecialties,
		&review.HelpfulCount,
		&review.IsHidden,
		&review.IsDeleted,
		&review.CreatedAt,
		&review.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("書評が見つかりません")
		}
		return nil, err
	}

	review.RecommendedLevels = parseStringArray(recommendedLevels)
	review.RecommendedSpecialties = parseStringArray(recommendedSpecialties)
	if deletedAt.Valid {
		review.DeletedAt = &deletedAt.Time
	}

	return &review, nil
}

func (r *reviewRepository) FindByBookID(ctx context.Context, bookID int64, filters map[string]interface{}) ([]*entity.Review, error) {
	query := `
		SELECT id, user_id, book_id, rating, content, recommended_levels, 
		       recommended_specialties, helpful_count, is_hidden, is_deleted, 
		       created_at, updated_at, deleted_at
		FROM reviews
		WHERE book_id = ? AND is_deleted = 0 AND is_hidden = 0
		ORDER BY helpful_count DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *reviewRepository) FindByUserIDAndBookID(ctx context.Context, userID, bookID int64) (*entity.Review, error) {
	query := `
		SELECT id, user_id, book_id, rating, content, recommended_levels, 
		       recommended_specialties, helpful_count, is_hidden, is_deleted, 
		       created_at, updated_at, deleted_at
		FROM reviews
		WHERE user_id = ? AND book_id = ? AND is_deleted = 0
	`

	var review entity.Review
	var recommendedLevels, recommendedSpecialties string
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, bookID).Scan(
		&review.ID,
		&review.UserID,
		&review.BookID,
		&review.Rating,
		&review.Content,
		&recommendedLevels,
		&recommendedSpecialties,
		&review.HelpfulCount,
		&review.IsHidden,
		&review.IsDeleted,
		&review.CreatedAt,
		&review.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("書評が見つかりません")
		}
		return nil, err
	}

	review.RecommendedLevels = parseStringArray(recommendedLevels)
	review.RecommendedSpecialties = parseStringArray(recommendedSpecialties)
	if deletedAt.Valid {
		review.DeletedAt = &deletedAt.Time
	}

	return &review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *entity.Review) error {
	query := `
		UPDATE reviews
		SET user_id = ?, book_id = ?, rating = ?, content = ?, recommended_levels = ?, 
		    recommended_specialties = ?, helpful_count = ?, is_hidden = ?, updated_at = ?
		WHERE id = ? AND is_deleted = 0
	`

	review.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		review.UserID,
		review.BookID,
		review.Rating,
		review.Content,
		serializeStringArray(review.RecommendedLevels),
		serializeStringArray(review.RecommendedSpecialties),
		review.HelpfulCount,
		review.IsHidden,
		review.UpdatedAt,
		review.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("書評が見つかりません")
	}

	return nil
}

func (r *reviewRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE reviews
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
		return errors.New("書評が見つかりません")
	}

	return nil
}

// scanRows 行をスキャンしてエンティティのスライスを返す
func (r *reviewRepository) scanRows(rows *sql.Rows) ([]*entity.Review, error) {
	var reviews []*entity.Review

	for rows.Next() {
		var review entity.Review
		var recommendedLevels, recommendedSpecialties string
		var deletedAt sql.NullTime

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.BookID,
			&review.Rating,
			&review.Content,
			&recommendedLevels,
			&recommendedSpecialties,
			&review.HelpfulCount,
			&review.IsHidden,
			&review.IsDeleted,
			&review.CreatedAt,
			&review.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		review.RecommendedLevels = parseStringArray(recommendedLevels)
		review.RecommendedSpecialties = parseStringArray(recommendedSpecialties)
		if deletedAt.Valid {
			review.DeletedAt = &deletedAt.Time
		}

		reviews = append(reviews, &review)
	}

	return reviews, rows.Err()
}
