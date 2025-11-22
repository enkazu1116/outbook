package database

import (
	"context"
	"database/sql"
	"errors"
	"example.com/backend/internal/domain/book/entity"
	"example.com/backend/internal/domain/book/repository"
	"time"
)

// bookRepository Turso用の書籍リポジトリ実装
type bookRepository struct {
	db *sql.DB
}

// NewBookRepository 新しい書籍リポジトリを作成
func NewBookRepository(db *sql.DB) repository.Repository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	query := `
		INSERT INTO books (title, author, publisher, published_date, isbn, cover_image_url, 
		                   description, primary_language, categories, level, page_count, 
		                   average_rating, reviews_count, outputs_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now
	
	result, err := r.db.ExecContext(ctx, query,
		book.Title,
		book.Author,
		book.Publisher,
		book.PublishedDate,
		book.ISBN,
		book.CoverImageURL,
		book.Description,
		book.PrimaryLanguage,
		serializeStringArray(book.Categories),
		book.Level,
		book.PageCount,
		book.AverageRating,
		book.ReviewsCount,
		book.OutputsCount,
		book.CreatedAt,
		book.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	book.ID = id
	return nil
}

func (r *bookRepository) FindByID(ctx context.Context, id int64) (*entity.Book, error) {
	query := `
		SELECT id, title, author, publisher, published_date, isbn, cover_image_url, 
		       description, primary_language, categories, level, page_count, 
		       average_rating, reviews_count, outputs_count, created_at, updated_at
		FROM books
		WHERE id = ?
	`
	
	var book entity.Book
	var categories string
	var publishedDate sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Publisher,
		&publishedDate,
		&book.ISBN,
		&book.CoverImageURL,
		&book.Description,
		&book.PrimaryLanguage,
		&categories,
		&book.Level,
		&book.PageCount,
		&book.AverageRating,
		&book.ReviewsCount,
		&book.OutputsCount,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("書籍が見つかりません")
		}
		return nil, err
	}
	
	book.Categories = parseStringArray(categories)
	if publishedDate.Valid {
		book.PublishedDate = &publishedDate.Time
	}
	
	return &book, nil
}

func (r *bookRepository) FindByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	query := `
		SELECT id, title, author, publisher, published_date, isbn, cover_image_url, 
		       description, primary_language, categories, level, page_count, 
		       average_rating, reviews_count, outputs_count, created_at, updated_at
		FROM books
		WHERE isbn = ?
	`
	
	var book entity.Book
	var categories string
	var publishedDate sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, isbn).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Publisher,
		&publishedDate,
		&book.ISBN,
		&book.CoverImageURL,
		&book.Description,
		&book.PrimaryLanguage,
		&categories,
		&book.Level,
		&book.PageCount,
		&book.AverageRating,
		&book.ReviewsCount,
		&book.OutputsCount,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("書籍が見つかりません")
		}
		return nil, err
	}
	
	book.Categories = parseStringArray(categories)
	if publishedDate.Valid {
		book.PublishedDate = &publishedDate.Time
	}
	
	return &book, nil
}

func (r *bookRepository) Search(ctx context.Context, query string, filters map[string]interface{}) ([]*entity.Book, error) {
	// TODO: 高度な検索クエリを実装
	sqlQuery := `
		SELECT id, title, author, publisher, published_date, isbn, cover_image_url, 
		       description, primary_language, categories, level, page_count, 
		       average_rating, reviews_count, outputs_count, created_at, updated_at
		FROM books
		WHERE (? = '' OR title LIKE '%' || ? || '%' OR author LIKE '%' || ? || '%')
	`
	
	rows, err := r.db.QueryContext(ctx, sqlQuery, query, query, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var books []*entity.Book
	for rows.Next() {
		var book entity.Book
		var categories string
		var publishedDate sql.NullTime
		
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Publisher,
			&publishedDate,
			&book.ISBN,
			&book.CoverImageURL,
			&book.Description,
			&book.PrimaryLanguage,
			&categories,
			&book.Level,
			&book.PageCount,
			&book.AverageRating,
			&book.ReviewsCount,
			&book.OutputsCount,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		book.Categories = parseStringArray(categories)
		if publishedDate.Valid {
			book.PublishedDate = &publishedDate.Time
		}
		
		books = append(books, &book)
	}
	
	return books, rows.Err()
}

func (r *bookRepository) Update(ctx context.Context, book *entity.Book) error {
	query := `
		UPDATE books
		SET title = ?, author = ?, publisher = ?, published_date = ?, isbn = ?, 
		    cover_image_url = ?, description = ?, primary_language = ?, 
		    categories = ?, level = ?, page_count = ?, average_rating = ?, 
		    reviews_count = ?, outputs_count = ?, updated_at = ?
		WHERE id = ?
	`
	
	book.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		book.Title,
		book.Author,
		book.Publisher,
		book.PublishedDate,
		book.ISBN,
		book.CoverImageURL,
		book.Description,
		book.PrimaryLanguage,
		serializeStringArray(book.Categories),
		book.Level,
		book.PageCount,
		book.AverageRating,
		book.ReviewsCount,
		book.OutputsCount,
		book.UpdatedAt,
		book.ID,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("書籍が見つかりません")
	}
	
	return nil
}

func (r *bookRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM books WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("書籍が見つかりません")
	}
	
	return nil
}

