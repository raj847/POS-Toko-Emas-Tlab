package repository

import (
	"context"
	"tim-b/entity"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db}
}

func (i *InventoryRepository) AddInventory(ctx context.Context, inv entity.Inventory) error {
	err := i.db.
		WithContext(ctx).
		Create(&inv).Error
	return err
}

func (i *InventoryRepository) ReadInventory() ([]entity.InventoryRead, error) {
	join := []entity.InventoryRead{}
	i.db.Table("inventories").Select("inventories.id,inventories.created_at,inventories.updated_at,inventories.deleted_at,inventories.id_jenis,inventories.id_bentuk,jenis.nama AS nama_jenis,bentuks.nama AS nama_bentuk,inventories.nama_barang,inventories.kode_barang,inventories.berat,inventories.kadar,inventories.harga_jual,inventories.notes,inventories.photo_url1,inventories.photo_url2,inventories.created_by,inventories.updated_by,inventories.deleted_by").
		Joins("inner join jenis on jenis.id = inventories.id_jenis").
		Joins("inner join bentuks on bentuks.id = inventories.id_bentuk").
		Where("inventories.deleted_at IS NULL").
		Scan(&join)
	return join, nil
}
func (i *InventoryRepository) ReadInventoryID(ctx context.Context, id int) (entity.InventoryRead, error) {
	join := entity.InventoryRead{}
	i.db.Table("inventories").Select("inventories.id,inventories.created_at,inventories.updated_at,inventories.deleted_at,inventories.id_jenis,inventories.id_bentuk,jenis.nama AS nama_jenis,bentuks.nama AS nama_bentuk,inventories.nama_barang,inventories.kode_barang,inventories.berat,inventories.kadar,inventories.harga_jual,inventories.notes,inventories.photo_url1,inventories.photo_url2,inventories.created_by,inventories.updated_by,inventories.deleted_by").
		Joins("inner join jenis on jenis.id = inventories.id_jenis").
		Joins("inner join bentuks on bentuks.id = inventories.id_bentuk").
		Where("inventories.id = ? AND inventories.deleted_at IS NULL", id).
		Scan(&join)
	return join, nil
}

func (i *InventoryRepository) DeleteInventory(ctx context.Context, id int) error {
	err := i.db.
		WithContext(ctx).
		Delete(&entity.Inventory{}, id).Error
	return err
}

func (i *InventoryRepository) UpdateInventory(ctx context.Context, inv entity.Inventory) error {
	err := i.db.
		WithContext(ctx).
		Table("inventories").
		Where("id = ?", inv.ID).
		Updates(&inv).Error
	return err
}
