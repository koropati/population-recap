package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
)

type accessTokenUsecase struct {
	accessTokenRepository domain.AccessTokenRepository
	contextTimeout        time.Duration
}

func NewAccessTokenUsecase(accessTokenRepository domain.AccessTokenRepository, timeout time.Duration) domain.AccessTokenUsecase {
	return &accessTokenUsecase{
		accessTokenRepository: accessTokenRepository,
		contextTimeout:        timeout,
	}
}

func (a *accessTokenUsecase) Create(c context.Context, accessToken domain.AccessToken) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.accessTokenRepository.Create(ctx, accessToken)
}

func (a *accessTokenUsecase) Revoke(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.accessTokenRepository.Revoke(ctx, token)
}

func (a *accessTokenUsecase) RevokeByUserID(c context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.accessTokenRepository.RevokeByUserID(ctx, userID)
}

func (a *accessTokenUsecase) IsValid(c context.Context, token string) bool {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.accessTokenRepository.IsValid(ctx, token)
}

func (a *accessTokenUsecase) Delete(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.accessTokenRepository.Delete(ctx, token)
}

func (a *accessTokenUsecase) DeleteExpiredToken(c context.Context, millisDateTime int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.accessTokenRepository.DeleteExpiredToken(ctx, millisDateTime)
}
