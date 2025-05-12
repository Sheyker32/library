package usecase

import (
	"context"
	"library/internal/domain"
	"library/internal/repository"
)

type Authorer interface {
	CreateAuthor(ctx context.Context, author *domain.Author) error
	GetAuthor(ctx context.Context, id int) (*domain.Author, error)
	ListAuthors(ctx context.Context) ([]*domain.Author, error)
	GetTopAuthors(ctx context.Context, limit int) ([]*domain.AuthorWithRentCount, error)
	DeleteAuthor(ctx context.Context, id int) error
	GetByBooksAuthor(ctx context.Context, idAuthor int) ([]domain.Book, error)
}

type AuthorUseCase struct {
	authorRepo repository.Authorer
	bookRepo   repository.Booker
}

func NewAuthorUseCase(authorRepo repository.Authorer, bookRepo repository.Booker) Authorer {
	return &AuthorUseCase{
		authorRepo: authorRepo,
		bookRepo:   bookRepo,
	}
}

func (uc *AuthorUseCase) CreateAuthor(ctx context.Context, author *domain.Author) error {
	return uc.authorRepo.Create(ctx, author)
}

func (uc *AuthorUseCase) GetAuthor(ctx context.Context, id int) (*domain.Author, error) {
	author, err := uc.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	books, err := uc.GetByBooksAuthor(ctx, id)

	if err != nil {
		return nil, err
	}

	author.Books = append(author.Books, books...)
	return author, nil
}

func (uc *AuthorUseCase) ListAuthors(ctx context.Context) ([]*domain.Author, error) {
	return uc.authorRepo.GetAll(ctx)
}

func (uc *AuthorUseCase) GetTopAuthors(ctx context.Context, limit int) ([]*domain.AuthorWithRentCount, error) {
	return uc.authorRepo.GetTopAuthors(ctx, limit)
}

func (uc *AuthorUseCase) DeleteAuthor(ctx context.Context, id int) error {
	return uc.authorRepo.DeleteAuthor(ctx, id)
}

func (uc *AuthorUseCase) GetByBooksAuthor(ctx context.Context, idAuthor int) ([]domain.Book, error) {
	return uc.authorRepo.GetByBooksAuthor(ctx, idAuthor)
}
