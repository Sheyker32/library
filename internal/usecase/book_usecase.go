package usecase

import (
	"context"
	"library/internal/domain"
	"library/internal/repository"
)

type Booker interface {
	AddBook(ctx context.Context, book *domain.Book) error
	GetBook(ctx context.Context, id int) (*domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) error
	DeleteBook(ctx context.Context, id int) error
}
type BookUseCase struct {
	bookRepo repository.Booker
}

func NewBookUseCase(bookRepo repository.Booker) Booker {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (uc *BookUseCase) AddBook(ctx context.Context, book *domain.Book) error {
	return uc.bookRepo.Create(ctx, book)
}

func (uc *BookUseCase) GetBook(ctx context.Context, id int) (*domain.Book, error) {
	return uc.bookRepo.GetByID(ctx, id)
}

func (uc *BookUseCase) UpdateBook(ctx context.Context, book *domain.Book) error {
	return uc.bookRepo.Update(ctx, book)
}

func (uc *BookUseCase) DeleteBook(ctx context.Context, id int) error {
	return uc.bookRepo.Delete(ctx, id)
}
