package facade

import (
	"context"
	"errors"
	"fmt"
	"library/internal/domain"
	"library/internal/repository"
	"library/internal/usecase"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/jmoiron/sqlx"
)

type Facader interface {
	CreateAuthor(ctx context.Context, author *domain.Author) error
	GetAuthor(ctx context.Context, id int) (*domain.Author, error)
	ListAuthors(ctx context.Context) ([]*domain.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	GetTopAuthors(ctx context.Context, limit int) ([]*domain.AuthorWithRentCount, error)
	AddBook(ctx context.Context, book *domain.Book) error
	GetBook(ctx context.Context, id int) (*domain.Book, error)
	GetByBooksAuthor(ctx context.Context, idAuthor int) ([]domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) error
	DeleteBook(ctx context.Context, id int) error
	RentBook(ctx context.Context, bookID, userID int) error
	ReturnBook(ctx context.Context, bookID int) error
	CreateUser(ctx context.Context, user *domain.User) error
	GetByIDUser(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	InitializeDataIfEmpty(ctx context.Context) error
	DeleteUser(ctx context.Context, id int) error
}

type LibraryFacade struct {
	db     *sqlx.DB
	author usecase.Authorer
	book   usecase.Booker
	rental usecase.Rentaler
	user   usecase.Userer
}

func NewLibraryFacade(
	db *sqlx.DB,
	author usecase.Authorer,
	book usecase.Booker,
	rental usecase.Rentaler,
	user usecase.Userer,
) *LibraryFacade {
	return &LibraryFacade{
		db:     db,
		author: author,
		book:   book,
		rental: rental,
		user:   user,
	}
}

func (l LibraryFacade) CreateAuthor(ctx context.Context, author *domain.Author) error {
	return l.author.CreateAuthor(ctx, author)
}
func (l LibraryFacade) GetAuthor(ctx context.Context, id int) (*domain.Author, error) {
	return l.author.GetAuthor(ctx, id)
}
func (l LibraryFacade) ListAuthors(ctx context.Context) ([]*domain.Author, error) {
	return l.author.ListAuthors(ctx)
}
func (l LibraryFacade) GetTopAuthors(ctx context.Context, limit int) ([]*domain.AuthorWithRentCount, error) {
	return l.author.GetTopAuthors(ctx, limit)
}
func (l LibraryFacade) DeleteAuthor(ctx context.Context, id int) error {
	return l.author.DeleteAuthor(ctx, id)
}

func (l LibraryFacade) AddBook(ctx context.Context, book *domain.Book) error {
	return l.book.AddBook(ctx, book)
}
func (l LibraryFacade) GetBook(ctx context.Context, id int) (*domain.Book, error) {
	return l.book.GetBook(ctx, id)
}
func (l LibraryFacade) GetByBooksAuthor(ctx context.Context, idAuthor int) ([]domain.Book, error) {
	return l.author.GetByBooksAuthor(ctx, idAuthor)
}
func (l LibraryFacade) UpdateBook(ctx context.Context, book *domain.Book) error {
	return l.book.UpdateBook(ctx, book)
}
func (l LibraryFacade) DeleteBook(ctx context.Context, id int) error {
	return l.book.DeleteBook(ctx, id)
}
func (l LibraryFacade) RentBook(ctx context.Context, bookID, userID int) error {
	return l.rental.RentBook(ctx, bookID, userID)
}
func (l LibraryFacade) ReturnBook(ctx context.Context, bookID int) error {
	return l.rental.ReturnBook(ctx, bookID)
}

func (l LibraryFacade) CreateUser(ctx context.Context, user *domain.User) error {
	return l.user.CreateUser(ctx, user)
}
func (l LibraryFacade) GetByIDUser(ctx context.Context, id int) (*domain.User, error) {
	return l.user.GetByIDUser(ctx, id)
}
func (l LibraryFacade) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return l.user.GetAllUsers(ctx)
}
func (l LibraryFacade) DeleteUser(ctx context.Context, id int) error {
	return l.user.DeleteUser(ctx, id)
}

func (lf LibraryFacade) InitializeDataIfEmpty(ctx context.Context) error {
	ok, err := repository.CheckIfTableHasRecords(lf.db, "authors")
	if err != nil {
		return err
	}
	if !ok {
		authors := make([]domain.Author, 10)
		for i := 0; i < 10; i++ {
			biography := fmt.Sprintf(
				"%s (род. %s) — %s из %s, известный своими работами в жанре %s.",
				gofakeit.Name(),
				gofakeit.Date().Format("2006-01-02"),
				gofakeit.JobTitle(),
				gofakeit.Country(),
				gofakeit.Book().Genre,
			)
			authors[i] = domain.Author{
				Name:      gofakeit.Name(),
				Biography: biography,
				CreatedAt: time.Now(),
			}
			err := lf.author.CreateAuthor(ctx, &authors[i])
			if err != nil {
				return err
			}
		}

	}

	ok, err = repository.CheckIfTableHasRecords(lf.db, "users")
	if err != nil {
		return err
	}
	if !ok {
		users := make([]domain.User, 55) // Больше 50
		for i := 0; i < 55; i++ {
			users[i] = domain.User{
				Name:      gofakeit.Name(),
				Email:     gofakeit.Email(),
				CreatedAt: time.Now(),
			}
			err := lf.user.CreateUser(ctx, &users[i])
			if err != nil {
				return err
			}
		}

	}

	ok, err = repository.CheckIfTableHasRecords(lf.db, "books")
	if err != nil {
		return err
	}

	if !ok {
		authors, err := lf.author.ListAuthors(ctx)
		if err != nil {
			return err
		}

		if len(authors) == 0 {
			return errors.New("no authors found to assign books")
		}

		books := make([]domain.Book, 100)
		for i := 0; i < 100; i++ {
			author := authors[gofakeit.Number(0, len(authors)-1)]
			books[i] = domain.Book{
				Title:     gofakeit.BookTitle(),
				AuthorID:  author.ID,
				CreatedAt: time.Now(),
				Available: true,
			}
			err := lf.book.AddBook(ctx, &books[i])
			if err != nil {
				return err
			}
		}

	}

	return nil
}
