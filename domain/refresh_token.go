package domain

import (
	"context"

	"github.com/google/uuid"
)

const (
	RefreshTokenTable = "refresh_tokens"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"type:longtext" json:"token"`
	PairToken string    `gorm:"index" json:"pair_token"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index;foreignKey:ID" json:"user_id"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt int64     `gorm:"index" json:"expires_at"`
}

type RefreshTokenRepository interface {
	Create(c context.Context, refreshToken RefreshToken) error
	Revoke(c context.Context, token string) error
	RevokeByPairToken(c context.Context, pairToken string) error
	RevokeByUserID(c context.Context, userID uuid.UUID) error
	IsValid(c context.Context, token string) bool
	Delete(c context.Context, token string) error
	DeleteExpiredToken(c context.Context, millisDateTime int64) error
}

type RefreshTokenUsecase interface {
	Create(c context.Context, refreshToken RefreshToken) error
	Revoke(c context.Context, token string) error
	RevokeByPairToken(c context.Context, pairToken string) error
	RevokeByUserID(c context.Context, userID uuid.UUID) error
	IsValid(c context.Context, token string) bool
	Delete(c context.Context, token string) error
	DeleteExpiredToken(c context.Context, millisDateTime int64) error
}
