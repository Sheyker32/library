package repository

import (
	"context"
	"library/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type Authorer interface {
	Create(ctx context.Context, author *domain.Author) error
	GetByID(ctx context.Context, id int) (*domain.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]*domain.Author, error)
	GetTopAuthors(ctx context.Context, limit int) ([]*domain.AuthorWithRentCount, error)
	GetByBooksAuthor(ctx context.Context, idAuthor int) ([]domain.Book, error)
}

type AuthorRepository struct {
	db *sqlx.DB
}

func NewAuthorRepository(db *sqlx.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) Create(ctx context.Context, author *domain.Author) error {
	query := `INSERT INTO authors (name, biography, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, author.Name, author.Biography, time.Now()).Scan(&author.ID)

	if len(author.Books) > 0 {
		bookQuery := `
			INSERT INTO books (title, author_id, available, created_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`

		for i := range author.Books {
			book := &author.Books[i]
			book.AuthorID = author.ID
			book.CreatedAt = time.Now()

			err = r.db.QueryRowContext(
				ctx,
				bookQuery,
				book.Title,
				book.AuthorID,
				book.Available,
				book.CreatedAt,
			).Scan(&book.ID)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (r *AuthorRepository) GetByID(ctx context.Context, id int) (*domain.Author, error) {
	query := `SELECT id, name, biography, created_at FROM authors WHERE id = $1`
	var author domain.Author

	err := r.db.GetContext(ctx, &author, query, id)
	if err != nil {
		return nil, err
	}
	books, err := r.GetByBooksAuthor(ctx, author.ID)
	if err != nil {
		return nil, err
	}

	if books != nil {
		author.Books = append(author.Books, books...)
	}

	return &author, err
}

func (r *AuthorRepository) GetAll(ctx context.Context) ([]*domain.Author, error) {
	query := `SELECT id, name, biography, created_at FROM authors`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*domain.Author
	for rows.Next() {
		author := &domain.Author{}
		if err := rows.Scan(&author.ID, &author.Name, &author.Biography, &author.CreatedAt); err != nil {
			return nil, err
		}
		books, err := r.GetByBooksAuthor(ctx, author.ID)
		if err != nil {
			return nil, err
		}

		if books != nil {
			author.Books = append(author.Books, books...)
		}

		authors = append(authors, author)
	}
	return authors, nil
}

func (r *AuthorRepository) GetTopAuthors(ctx context.Context, limit int) ([]*domain.AuthorWithRentCount, error) {
	query := `
		SELECT a.id, a.name, a.biography, a.created_at, COUNT(r.id) as rental_count
		FROM authors a
		LEFT JOIN books b ON b.author_id = a.id
		LEFT JOIN book_rental r ON r.book_id = b.id
		GROUP BY a.id
		ORDER BY rental_count DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.AuthorWithRentCount
	for rows.Next() {
		item := &domain.AuthorWithRentCount{}
		err := rows.Scan(
			&item.Author.ID,
			&item.Author.Name,
			&item.Author.Biography,
			&item.Author.CreatedAt,
			&item.RentCount,
		)
		if err != nil {
			return nil, err
		}
		books, err := r.GetByBooksAuthor(ctx, item.Author.ID)
		if err != nil {
			return nil, err
		}

		if books != nil {
			item.Author.Books = append(item.Author.Books, books...)
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *AuthorRepository) DeleteAuthor(ctx context.Context, id int) error {
	query := `DELETE FROM authors WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthorRepository) GetByBooksAuthor(ctx context.Context, idAuthor int) ([]domain.Book, error) {
	var books []domain.Book

	query := `
		SELECT b.id, b.title, b.author_id, b.available, b.created_at
		FROM books b
		WHERE b.author_id = $1
	`
	err := r.db.SelectContext(ctx, &books, query, idAuthor)
	if err != nil {
		return nil, err
	}

	return books, nil
}
