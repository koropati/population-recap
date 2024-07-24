package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
)

type refreshTokenUsecase struct {
	refreshTokenRepository domain.RefreshTokenRepository
	contextTimeout         time.Duration
}

func NewRefreshTokenUsecase(refreshTokenRepository domain.RefreshTokenRepository, timeout time.Duration) domain.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		refreshTokenRepository: refreshTokenRepository,
		contextTimeout:         timeout,
	}
}

func (a *refreshTokenUsecase) Create(c context.Context, refreshToken domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.refreshTokenRepository.Create(ctx, refreshToken)
}

func (a *refreshTokenUsecase) Revoke(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.refreshTokenRepository.Revoke(ctx, token)
}

func (a *refreshTokenUsecase) RevokeByPairToken(c context.Context, pairToken string) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.refreshTokenRepository.RevokeByPairToken(ctx, pairToken)
}

func (a *refreshTokenUsecase) RevokeByUserID(c context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.refreshTokenRepository.RevokeByUserID(ctx, userID)
}

func (a *refreshTokenUsecase) IsValid(c context.Context, token string) bool {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.refreshTokenRepository.IsValid(ctx, token)
}

func (a *refreshTokenUsecase) Delete(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.refreshTokenRepository.Delete(ctx, token)
}

func (a *refreshTokenUsecase) DeleteExpiredToken(c context.Context, millisDateTime int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.refreshTokenRepository.DeleteExpiredToken(ctx, millisDateTime)
}
