package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
	"gorm.io/gorm"
)

const (
	errMsgNoTokenUpdated = "no token was updated"
)

type refreshTokenRepository struct {
	database  *gorm.DB
	table     string
	pageInit  int64
	limitInit int64
}

func NewRefreshTokenRepository(db *gorm.DB, table string, pageInit int64, limitInit int64) domain.RefreshTokenRepository {
	return &refreshTokenRepository{
		database:  db,
		table:     table,
		pageInit:  pageInit,
		limitInit: limitInit,
	}
}

func (r *refreshTokenRepository) Create(c context.Context, refreshToken domain.RefreshToken) error {
	result := r.database.WithContext(c).Table(r.table).Create(&refreshToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *refreshTokenRepository) Revoke(c context.Context, token string) error {
	result := r.database.WithContext(c).Table(r.table).Where("token = ?", token).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New(errMsgNoTokenUpdated)
	}
	return nil
}

func (r *refreshTokenRepository) RevokeByPairToken(c context.Context, token string) error {
	result := r.database.WithContext(c).Table(r.table).Where("pair_token = ?", token).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New(errMsgNoTokenUpdated)
	}
	return nil
}

func (r *refreshTokenRepository) RevokeByUserID(c context.Context, userID uuid.UUID) error {
	result := r.database.WithContext(c).Table(r.table).Where("user_id = ?", userID).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New(errMsgNoTokenUpdated)
	}
	return nil
}

func (r *refreshTokenRepository) IsValid(c context.Context, token string) bool {
	var refreshToken domain.RefreshToken
	result := r.database.WithContext(c).Table(r.table).Where("token = ? AND revoked = ?", token, false).First(&refreshToken)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return time.Unix(refreshToken.ExpiresAt, 0).After(time.Now())
}

func (r *refreshTokenRepository) Delete(c context.Context, token string) error {
	result := r.database.WithContext(c).Table(r.table).Where("token = ?", token).Delete(&domain.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no token was deleted")
	}
	return nil
}

func (r *refreshTokenRepository) DeleteExpiredToken(c context.Context, millisDateTime int64) error {
	result := r.database.WithContext(c).Table(r.table).Where("expires_at < ?", millisDateTime).Or("revoked = ?", 1).Delete(&domain.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
