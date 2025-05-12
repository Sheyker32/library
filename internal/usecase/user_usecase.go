package usecase

import (
	"context"
	"library/internal/domain"
	"library/internal/repository"
)

type Userer interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetByIDUser(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type UserUseCase struct {
	userRepo repository.Userer
}

func NewUserUseCase(userRepo repository.Userer) Userer {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u UserUseCase) CreateUser(ctx context.Context, user *domain.User) error {
	return u.userRepo.Create(ctx, user)
}

func (u UserUseCase) GetByIDUser(ctx context.Context, id int) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u UserUseCase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}
func (u UserUseCase) DeleteUser(ctx context.Context, id int) error {
	return u.userRepo.Delete(ctx, id)
}
