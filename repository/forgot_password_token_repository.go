package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
	"gorm.io/gorm"
)

type forgotPasswordTokenRepository struct {
	database  *gorm.DB
	table     string
	pageInit  int64
	limitInit int64
}

func NewForgotPasswordTokenRepository(db *gorm.DB, table string, pageInit int64, limitInit int64) domain.ForgotPasswordTokenRepository {
	return &forgotPasswordTokenRepository{
		database:  db,
		table:     table,
		pageInit:  pageInit,
		limitInit: limitInit,
	}
}

func (r *forgotPasswordTokenRepository) Create(c context.Context, forgotPasswordToken domain.ForgotPasswordToken) error {
	result := r.database.WithContext(c).Table(r.table).Create(&forgotPasswordToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *forgotPasswordTokenRepository) Revoke(c context.Context, token string) error {
	result := r.database.WithContext(c).Table(r.table).Where("token = ?", token).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no token was updated")
	}
	return nil
}

func (r *forgotPasswordTokenRepository) RevokeByUserID(c context.Context, userID uuid.UUID) error {
	result := r.database.WithContext(c).Table(r.table).Where("user_id = ?", userID).Update("revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no token was updated")
	}
	return nil
}

func (r *forgotPasswordTokenRepository) IsValid(c context.Context, token string) bool {
	var forgotPasswordToken domain.ForgotPasswordToken
	result := r.database.WithContext(c).Table(r.table).Where("token = ? AND revoked = ?", token, false).First(&forgotPasswordToken)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return time.Unix(forgotPasswordToken.ExpiresAt, 0).After(time.Now())
}

func (r *forgotPasswordTokenRepository) Delete(c context.Context, token string) error {
	result := r.database.WithContext(c).Table(r.table).Where("token = ?", token).Delete(&domain.ForgotPasswordToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no token was deleted")
	}
	return nil
}

func (r *forgotPasswordTokenRepository) GetUserID(c context.Context, token string) (userID uuid.UUID, err error) {
	var forgotPasswordToken domain.ForgotPasswordToken
	result := r.database.WithContext(c).Table(r.table).Where("token = ? AND revoked = ?", token, false).First(&forgotPasswordToken)
	if result.Error != nil || result.RowsAffected == 0 {
		return uuid.Nil, errors.New("failed get user id")
	}
	return forgotPasswordToken.UserID, nil
}

func (r *forgotPasswordTokenRepository) DeleteExpiredToken(c context.Context, millisDateTime int64) error {
	result := r.database.WithContext(c).Table(r.table).Where("expires_at < ?", millisDateTime).Or("revoked = ?", 1).Delete(&domain.ForgotPasswordToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
