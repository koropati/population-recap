package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
	"gorm.io/gorm"
)

type accessTokenRepository struct {
	database  *gorm.DB
	table     string
	pageInit  int64
	limitInit int64
}

func NewAccessTokenRepository(db *gorm.DB, table string, pageInit int64, limitInit int64) domain.AccessTokenRepository {
	return &accessTokenRepository{
		database:  db,
		table:     table,
		pageInit:  pageInit,
		limitInit: limitInit,
	}
}

func (r *accessTokenRepository) Create(c context.Context, accessToken domain.AccessToken) error {
	result := r.database.WithContext(c).Table(r.table).Create(&accessToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *accessTokenRepository) Revoke(c context.Context, token string) error {
	result := r.database.WithContext(c).Table(r.table).Where("token = ?", token).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no token was updated")
	}
	return nil
}

func (r *accessTokenRepository) RevokeByUserID(c context.Context, userID uuid.UUID) error {
	result := r.database.WithContext(c).Table(r.table).Where("user_id = ?", userID).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no token was updated")
	}
	return nil
}

func (r *accessTokenRepository) IsValid(c context.Context, token string) bool {
	var accessToken domain.AccessToken
	result := r.database.WithContext(c).Table(r.table).Where("token = ? AND revoked = ?", token, false).First(&accessToken)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return time.Unix(accessToken.ExpiresAt, 0).After(time.Now())
}

func (r *accessTokenRepository) Delete(c context.Context, token string) error {
	result := r.database.WithContext(c).Table(r.table).Where("token = ?", token).Delete(&domain.AccessToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no token was deleted")
	}
	return nil
}

func (r *accessTokenRepository) DeleteExpiredToken(c context.Context, millisDateTime int64) error {
	result := r.database.WithContext(c).Table(r.table).Where("expires_at < ?", millisDateTime).Or("revoked = ?", 1).Delete(&domain.AccessToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
