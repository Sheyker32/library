package repository

import (
	"context"
	"library/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Booker interface {
	Create(ctx context.Context, book *domain.Book) error
	GetByID(ctx context.Context, id int) (*domain.Book, error)
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id int) error
}

type BookRepository struct {
	db         *sqlx.DB
	authorRepo Authorer
}

func NewBookRepository(db *sqlx.DB, author Authorer) *BookRepository {
	return &BookRepository{db: db, authorRepo: author}
}
func (r *BookRepository) Create(ctx context.Context, book *domain.Book) error {
	query := `
		INSERT INTO books (title, author_id, available, created_at)
		VALUES (:title, :author_id, :available, :created_at)
		RETURNING id
	`

	_, err := r.db.NamedExecContext(ctx, query, book)

	if err != nil {
		return err
	}

	return nil
}

func (r *BookRepository) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	var book domain.Book
	query := `
		SELECT *
		FROM books b
		WHERE b.id = $1
	`
	err := r.db.GetContext(ctx, &book, query, id)
	if err != nil {
		return nil, err
	}
	author, err := r.authorRepo.GetByID(ctx, book.AuthorID)
	if err != nil {
		return nil, err
	}
	book.Author = author

	return &book, nil
}

func (r *BookRepository) Update(ctx context.Context, book *domain.Book) error {
	query := `
        UPDATE books 
        SET title = $1, 
            author_id = $2, 
            available = $3 
        WHERE id = $4
    `

	result, err := r.db.ExecContext(
		ctx,
		query,
		book.Title,
		book.AuthorID,
		book.Available,
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
		return &domain.ErrBookNotFound{BookID: book.ID}
	}

	return nil
}

func (r *BookRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
