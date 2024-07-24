package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
)

type userUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Create(c context.Context, user domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.Create(ctx, user)
}

func (u *userUsecase) Retrieve(c context.Context, filter domain.Filter) (users []domain.User, meta domain.MetaResponse, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.Retrieve(ctx, filter)
}

func (u *userUsecase) Update(c context.Context, id uuid.UUID, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.Update(ctx, id, user)
}

func (u *userUsecase) UpdatePassword(c context.Context, id uuid.UUID, newPasswordHash string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.UpdatePassword(ctx, id, newPasswordHash)
}

func (u *userUsecase) Delete(c context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.Delete(ctx, id)
}

func (u *userUsecase) GetByEmail(c context.Context, email string) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.GetByEmail(ctx, email)
}

func (u *userUsecase) GetById(c context.Context, id uuid.UUID) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.GetById(ctx, id)
}
