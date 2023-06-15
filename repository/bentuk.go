package repository

import (
	"context"
	"tim-b/entity"

	"gorm.io/gorm"
)

type BentukRepository struct {
	db *gorm.DB
}

func NewBentukRepository(db *gorm.DB) *BentukRepository {
	return &BentukRepository{db}
}

func (b *BentukRepository) GetAllBentuk(ctx context.Context) ([]entity.Bentuk, error) {
	var bentukResult []entity.Bentuk

	bentuk, err := b.db.
		WithContext(ctx).
		Table("bentuks").
		Select("*").
		Where("deleted_at IS NULL").
		Rows()
	if err != nil {
		return []entity.Bentuk{}, err
	}
	defer bentuk.Close()

	for bentuk.Next() {
		b.db.ScanRows(bentuk, &bentukResult)
	}

	return bentukResult, nil
}

func (b *BentukRepository) AddBentuk(ctx context.Context, bentuk entity.Bentuk) error {
	err := b.db.
		WithContext(ctx).
		Create(&bentuk).Error
	return err
}

func (b *BentukRepository) GetBentukByID(ctx context.Context, id int) (entity.Bentuk, error) {
	var bentukResult entity.Bentuk

	err := b.db.
		WithContext(ctx).
		Table("bentuks").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&bentukResult).Error
	if err != nil {
		return entity.Bentuk{}, err
	}

	return bentukResult, nil
}

func (b *BentukRepository) DeleteBentuk(ctx context.Context, id int) error {
	err := b.db.
		WithContext(ctx).
		Delete(&entity.Bentuk{}, id).Error
	return err
}

func (b *BentukRepository) UpdateBentuk(ctx context.Context, bentuk entity.Bentuk) error {
	err := b.db.
		WithContext(ctx).
		Table("bentuks").
		Where("id = ?", bentuk.ID).
		Updates(&bentuk).Error
	return err
}
