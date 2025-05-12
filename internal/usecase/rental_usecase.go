package usecase

import (
	"context"
	"fmt"
	"library/internal/repository"
)

type Rentaler interface {
	RentBook(ctx context.Context, bookID, userID int) error
	ReturnBook(ctx context.Context, bookID int) error
}

type RentalUseCase struct {
	user       repository.Userer
	bookRepo   repository.Booker
	rentalRepo repository.Rentaler
}

func NewRentUseCase(user repository.Userer, bookRepo repository.Booker, rentRepo repository.Rentaler) Rentaler {
	return &RentalUseCase{
		user:       user,
		bookRepo:   bookRepo,
		rentalRepo: rentRepo,
	}
}

func (uc *RentalUseCase) RentBook(ctx context.Context, bookID, userID int) error {
	book, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}
	if !book.Available {
		return fmt.Errorf("book is not available")
	}

	_, err = uc.user.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	err = uc.rentalRepo.RentBook(ctx, bookID, userID)
	if err != nil {
		return err
	}

	book.Available = false
	err = uc.bookRepo.Update(ctx, book)
	if err != nil {
		return err
	}

	return nil
}

func (uc *RentalUseCase) ReturnBook(ctx context.Context, bookID int) error {
	book, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}
	if book.Available {
		return fmt.Errorf("book was not issued")
	}

	err = uc.rentalRepo.ReturnBook(ctx, bookID)
	if err != nil {
		return err
	}

	book.Available = true
	err = uc.bookRepo.Update(ctx, book)
	if err != nil {
		return err
	}

	return nil
}
