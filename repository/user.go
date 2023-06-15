package repository

import (
	"context"
	"fmt"
	"tim-b/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) AddAnggota(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var res entity.User
	err := r.db.WithContext(ctx).Table("users").Where("username = ?", username).Find(&res).Error
	if err != nil {
		return entity.User{}, err
	}

	return res, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var listUsers []entity.User

	user, err := r.db.
		WithContext(ctx).
		Table("users").
		Select("*,users.tanggal_masuk").
		Where("deleted_at IS NULL").
		Rows()
	if err != nil {
		return []entity.User{}, err
	}
	defer user.Close()

	for user.Next() {
		r.db.ScanRows(user, &listUsers)
	}

	return listUsers, nil
}

func (r *UserRepository) GetUsersbyAdder(ctx context.Context, adder string) ([]entity.User, error) {
	var listUsers []entity.User

	user, err := r.db.
		WithContext(ctx).
		Table("users").
		Select("*").
		Where("created_by = ? AND deleted_at IS NULL", adder).
		Rows()
	if err != nil {
		return []entity.User{}, err
	}
	defer user.Close()

	for user.Next() {
		r.db.ScanRows(user, &listUsers)
	}

	return listUsers, nil
}

func (r *UserRepository) GetUsersbyID(ctx context.Context, id uint) (entity.User, error) {
	res := entity.User{}
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).Find(&res).Error
	if err != nil {
		return entity.User{}, err
	}

	return res, nil
}

func (r *UserRepository) SearchAnggota(ctx context.Context, namalengkap string) ([]entity.User, error) {
	var userResult []entity.User

	err := r.db.
		WithContext(ctx).
		Table("products").
		Where("name LIKE ? AND deleted_at IS NULL", fmt.Sprintf("%s%s%s", "%", namalengkap, "%")).
		Find(&userResult).Error
	if err != nil {
		return []entity.User{}, err
	}
	return userResult, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
