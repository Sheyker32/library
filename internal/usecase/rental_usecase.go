package usecase

import (
	"context"
	"library/internal/repository"
)

type Rentaler interface {
	RentBook(ctx context.Context, bookID, userID int) error
	ReturnBook(ctx context.Context, bookID int) error
}

type RentalUseCase struct {
	rentalRepo repository.Rentaler
}

func NewRentUseCase(rentRepo repository.Rentaler) Rentaler {
	return &RentalUseCase{
		rentalRepo: rentRepo,
	}
}

func (uc *RentalUseCase) RentBook(ctx context.Context, bookID, userID int) error {
	return uc.rentalRepo.RentBook(ctx, bookID, userID)
}

func (uc *RentalUseCase) ReturnBook(ctx context.Context, bookID int) error {
	return uc.rentalRepo.ReturnBook(ctx, bookID)
}
