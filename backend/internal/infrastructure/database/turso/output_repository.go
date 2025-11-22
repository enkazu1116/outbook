package database

import (
	"context"
	"database/sql"
	"errors"
	"example.com/backend/internal/domain/output/entity"
	"example.com/backend/internal/domain/output/repository"
	"time"
)

// outputRepository Turso用のアウトプットリポジトリ実装
type outputRepository struct {
	db *sql.DB
}

// NewOutputRepository 新しいアウトプットリポジトリを作成
func NewOutputRepository(db *sql.DB) repository.Repository {
	return &outputRepository{db: db}
}

func (r *outputRepository) Create(ctx context.Context, output *entity.Output) error {
	query := `
		INSERT INTO outputs (user_id, book_id, title, content, visibility, 
		                     reading_finished_date, tags, is_draft, views_count, 
		                     created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	output.CreatedAt = now
	output.UpdatedAt = now
	
	result, err := r.db.ExecContext(ctx, query,
		output.UserID,
		output.BookID,
		output.Title,
		output.Content,
		output.Visibility,
		output.ReadingFinishedDate,
		serializeStringArray(output.Tags),
		output.IsDraft,
		output.ViewsCount,
		output.CreatedAt,
		output.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	output.ID = id
	return nil
}

func (r *outputRepository) FindByID(ctx context.Context, id int64) (*entity.Output, error) {
	query := `
		SELECT id, user_id, book_id, title, content, visibility, reading_finished_date, 
		       tags, is_draft, is_deleted, views_count, created_at, updated_at, deleted_at
		FROM outputs
		WHERE id = ? AND is_deleted = 0
	`
	
	var output entity.Output
	var tags string
	var readingFinishedDate, deletedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&output.ID,
		&output.UserID,
		&output.BookID,
		&output.Title,
		&output.Content,
		&output.Visibility,
		&readingFinishedDate,
		&tags,
		&output.IsDraft,
		&output.IsDeleted,
		&output.ViewsCount,
		&output.CreatedAt,
		&output.UpdatedAt,
		&deletedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("アウトプットが見つかりません")
		}
		return nil, err
	}
	
	output.Tags = parseStringArray(tags)
	if readingFinishedDate.Valid {
		output.ReadingFinishedDate = &readingFinishedDate.Time
	}
	if deletedAt.Valid {
		output.DeletedAt = &deletedAt.Time
	}
	
	return &output, nil
}

func (r *outputRepository) FindByUserID(ctx context.Context, userID int64, filters map[string]interface{}) ([]*entity.Output, error) {
	query := `
		SELECT id, user_id, book_id, title, content, visibility, reading_finished_date, 
		       tags, is_draft, is_deleted, views_count, created_at, updated_at, deleted_at
		FROM outputs
		WHERE user_id = ? AND is_deleted = 0
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanRows(rows)
}

func (r *outputRepository) FindByBookID(ctx context.Context, bookID int64, filters map[string]interface{}) ([]*entity.Output, error) {
	query := `
		SELECT id, user_id, book_id, title, content, visibility, reading_finished_date, 
		       tags, is_draft, is_deleted, views_count, created_at, updated_at, deleted_at
		FROM outputs
		WHERE book_id = ? AND is_deleted = 0 AND visibility != 'private'
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanRows(rows)
}

func (r *outputRepository) Update(ctx context.Context, output *entity.Output) error {
	query := `
		UPDATE outputs
		SET user_id = ?, book_id = ?, title = ?, content = ?, visibility = ?, 
		    reading_finished_date = ?, tags = ?, is_draft = ?, views_count = ?, 
		    updated_at = ?
		WHERE id = ? AND is_deleted = 0
	`
	
	output.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		output.UserID,
		output.BookID,
		output.Title,
		output.Content,
		output.Visibility,
		output.ReadingFinishedDate,
		serializeStringArray(output.Tags),
		output.IsDraft,
		output.ViewsCount,
		output.UpdatedAt,
		output.ID,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("アウトプットが見つかりません")
	}
	
	return nil
}

func (r *outputRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE outputs
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
		return errors.New("アウトプットが見つかりません")
	}
	
	return nil
}

// scanRows 行をスキャンしてエンティティのスライスを返す
func (r *outputRepository) scanRows(rows *sql.Rows) ([]*entity.Output, error) {
	var outputs []*entity.Output
	
	for rows.Next() {
		var output entity.Output
		var tags string
		var readingFinishedDate, deletedAt sql.NullTime
		
		err := rows.Scan(
			&output.ID,
			&output.UserID,
			&output.BookID,
			&output.Title,
			&output.Content,
			&output.Visibility,
			&readingFinishedDate,
			&tags,
			&output.IsDraft,
			&output.IsDeleted,
			&output.ViewsCount,
			&output.CreatedAt,
			&output.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		
		output.Tags = parseStringArray(tags)
		if readingFinishedDate.Valid {
			output.ReadingFinishedDate = &readingFinishedDate.Time
		}
		if deletedAt.Valid {
			output.DeletedAt = &deletedAt.Time
		}
		
		outputs = append(outputs, &output)
	}
	
	return outputs, rows.Err()
}

