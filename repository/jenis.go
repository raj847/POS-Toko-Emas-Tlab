package repository

import (
	"context"
	"tim-b/entity"

	"gorm.io/gorm"
)

type JenisRepository struct {
	db *gorm.DB
}

func NewJenisRepository(db *gorm.DB) *JenisRepository {
	return &JenisRepository{db}
}

func (j *JenisRepository) GetAllJenis(ctx context.Context) ([]entity.Jenis, error) {
	var jenisResult []entity.Jenis

	jenis, err := j.db.
		WithContext(ctx).
		Table("jenis").
		Select("*").
		Where("deleted_at IS NULL").
		Rows()
	if err != nil {
		return []entity.Jenis{}, err
	}
	defer jenis.Close()

	for jenis.Next() {
		j.db.ScanRows(jenis, &jenisResult)
	}

	return jenisResult, nil
}

func (j *JenisRepository) AddJenis(ctx context.Context, jenis entity.Jenis) error {
	err := j.db.
		WithContext(ctx).
		Create(&jenis).Error
	return err
}

func (j *JenisRepository) GetJenisByID(ctx context.Context, id int) (entity.Jenis, error) {
	var jenisResult entity.Jenis

	err := j.db.
		WithContext(ctx).
		Table("jenis").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&jenisResult).Error
	if err != nil {
		return entity.Jenis{}, err
	}

	return jenisResult, nil
}

func (j *JenisRepository) DeleteJenis(ctx context.Context, id int) error {
	err := j.db.
		WithContext(ctx).
		Delete(&entity.Jenis{}, id).Error
	return err
}

func (j *JenisRepository) UpdateJenis(ctx context.Context, jenis entity.Jenis) error {
	err := j.db.
		WithContext(ctx).
		Table("jenis").
		Where("id = ?", jenis.ID).
		Updates(&jenis).Error
	return err
}
