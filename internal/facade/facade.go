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
	RentBook(ctx context.Context, bookID, userID int) error
	ReturnBook(ctx context.Context, bookID int) error
	InitializeDataIfEmpty(ctx context.Context) error
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

func (l LibraryFacade) RentBook(ctx context.Context, bookID, userID int) error {
	book, err := l.book.GetBook(ctx, bookID)
	if err != nil {
		return err
	}
	if !book.Available {
		return fmt.Errorf("book is not available")
	}

	_, err = l.user.GetByIDUser(ctx, userID)
	if err != nil {
		return err
	}

	err = l.rental.RentBook(ctx, bookID, userID)
	if err != nil {
		return err
	}

	book.Available = false
	err = l.book.UpdateBook(ctx, book)
	if err != nil {
		return err
	}

	return nil
}
func (l LibraryFacade) ReturnBook(ctx context.Context, bookID int) error {
	book, err := l.book.GetBook(ctx, bookID)
	if err != nil {
		return err
	}
	if book.Available {
		return fmt.Errorf("book was not issued")
	}

	err = l.rental.ReturnBook(ctx, bookID)
	if err != nil {
		return err
	}

	book.Available = true
	err = l.book.UpdateBook(ctx, book)
	if err != nil {
		return err
	}

	return nil
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
