package domain

import (
	"context"

	"github.com/google/uuid"
)

const (
	ForgotPasswordTokenTable = "forgot_password_tokens"
)

type ForgotPasswordToken struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"type:longtext" json:"token"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index;foreignKey:ID" json:"user_id"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt int64     `gorm:"index" json:"expires_at"`
}

type ForgotPasswordTokenRepository interface {
	Create(c context.Context, forgotPasswordToken ForgotPasswordToken) error
	Revoke(c context.Context, token string) error
	RevokeByUserID(c context.Context, userID uuid.UUID) error
	IsValid(c context.Context, token string) bool
	GetUserID(c context.Context, token string) (userID uuid.UUID, err error)
	Delete(c context.Context, token string) error
	DeleteExpiredToken(c context.Context, millisDateTime int64) error
}

type ForgotPasswordTokenUsecase interface {
	Create(c context.Context, forgotPasswordToken ForgotPasswordToken) error
	Revoke(c context.Context, token string) error
	RevokeByUserID(c context.Context, userID uuid.UUID) error
	IsValid(c context.Context, token string) bool
	GetUserID(c context.Context, token string) (userID uuid.UUID, err error)
	Delete(c context.Context, token string) error
	DeleteExpiredToken(c context.Context, millisDateTime int64) error
}
