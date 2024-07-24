package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
	"gorm.io/gorm"
)

const (
	queryFindByID = "id = ?"
)

type userRepository struct {
	database  *gorm.DB
	table     string
	pageInit  int64
	limitInit int64
}

func NewUserRepository(db *gorm.DB, table string, pageInit int64, limitInit int64) domain.UserRepository {
	return &userRepository{
		database:  db,
		table:     table,
		pageInit:  pageInit,
		limitInit: limitInit,
	}
}

func (u *userRepository) Create(c context.Context, data domain.User) error {
	result := u.database.WithContext(c).Table(u.table).Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepository) Retrieve(c context.Context, filter domain.Filter) (users []domain.User, meta domain.MetaResponse, err error) {
	query := u.database.WithContext(c).Table(u.table)

	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}

	if filter.WithPagination {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(int(offset)).Limit(int(filter.Limit))
	}

	result := query.Find(&users)
	if result.Error != nil {
		return nil, domain.MetaResponse{}, result.Error
	}

	// Hitung total records
	var totalRecords int64
	u.database.Table(u.table).Count(&totalRecords)

	meta = domain.MetaResponse{
		TotalRecords:    totalRecords,
		FilteredRecords: int64(result.RowsAffected),
		Page:            filter.Page,
		PerPage:         filter.Limit,
		TotalPages:      (totalRecords + filter.Limit - 1) / filter.Limit,
	}

	return users, meta, nil
}

func (u *userRepository) Update(c context.Context, id uuid.UUID, data domain.User) (user domain.User, err error) {
	result := u.database.WithContext(c).Table(u.table).Where(queryFindByID, id).Updates(data)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return domain.User{}, errors.New("no user was updated")
	}
	err = u.database.WithContext(c).Table(u.table).Where(queryFindByID, id).First(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *userRepository) UpdatePassword(c context.Context, id uuid.UUID, newPasswordHash string) (err error) {
	err = u.database.WithContext(c).Table(u.table).Where(queryFindByID, id).Update("password", newPasswordHash).Error
	return err
}

func (u *userRepository) Delete(c context.Context, id uuid.UUID) error {
	result := u.database.WithContext(c).Table(u.table).Where(queryFindByID, id).Delete(&domain.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no user was deleted")
	}
	return nil
}

func (u *userRepository) GetByEmail(c context.Context, email string) (user domain.User, err error) {
	result := u.database.WithContext(c).Table(u.table).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (u *userRepository) GetById(c context.Context, id uuid.UUID) (user domain.User, err error) {
	result := u.database.WithContext(c).Table(u.table).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}
