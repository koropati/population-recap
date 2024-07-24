package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
)

type forgotPasswordTokenUsecase struct {
	forgotPasswordTokenRepository domain.ForgotPasswordTokenRepository
	contextTimeout                time.Duration
}

func NewForgotPasswordTokenUsecase(forgotPasswordTokenRepository domain.ForgotPasswordTokenRepository, timeout time.Duration) domain.ForgotPasswordTokenUsecase {
	return &forgotPasswordTokenUsecase{
		forgotPasswordTokenRepository: forgotPasswordTokenRepository,
		contextTimeout:                timeout,
	}
}

func (a *forgotPasswordTokenUsecase) Create(c context.Context, forgotPasswordToken domain.ForgotPasswordToken) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.forgotPasswordTokenRepository.Create(ctx, forgotPasswordToken)
}

func (a *forgotPasswordTokenUsecase) Revoke(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.forgotPasswordTokenRepository.Revoke(ctx, token)
}

func (a *forgotPasswordTokenUsecase) RevokeByUserID(c context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.forgotPasswordTokenRepository.RevokeByUserID(ctx, userID)
}

func (a *forgotPasswordTokenUsecase) IsValid(c context.Context, token string) bool {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.forgotPasswordTokenRepository.IsValid(ctx, token)
}

func (a *forgotPasswordTokenUsecase) Delete(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.forgotPasswordTokenRepository.Delete(ctx, token)
}

func (a *forgotPasswordTokenUsecase) GetUserID(c context.Context, token string) (userID uuid.UUID, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.forgotPasswordTokenRepository.GetUserID(ctx, token)
}

func (a *forgotPasswordTokenUsecase) DeleteExpiredToken(c context.Context, millisDateTime int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.forgotPasswordTokenRepository.DeleteExpiredToken(ctx, millisDateTime)
}
