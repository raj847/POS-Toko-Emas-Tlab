package service

import (
	"context"
	"tim-b/entity"
	"tim-b/repository"
)

type InventoryService struct {
	iRepo *repository.InventoryRepository
}

func NewInventoryService(iRepo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		iRepo: iRepo,
	}
}

func (i *InventoryService) AddInventory(ctx context.Context, inv entity.Inventory) (entity.Inventory, error) {
	err := i.iRepo.AddInventory(ctx, inv)
	if err != nil {
		return entity.Inventory{}, err
	}
	return inv, nil
}
func (i *InventoryService) ReadInventory() ([]entity.InventoryRead, error) {
	inv, err := i.iRepo.ReadInventory()
	if err != nil {
		return []entity.InventoryRead{}, err
	}
	return inv, nil
}
func (i *InventoryService) ReadInventoryID(ctx context.Context, id int) (entity.InventoryRead, error) {
	inv, err := i.iRepo.ReadInventoryID(ctx, id)
	if err != nil {
		return entity.InventoryRead{}, err
	}
	return inv, nil
}

func (i *InventoryService) UpdateInventory(ctx context.Context, inv entity.Inventory) (entity.Inventory, error) {
	err := i.iRepo.UpdateInventory(ctx, inv)
	if err != nil {
		return entity.Inventory{}, err
	}

	return inv, nil
}

func (i *InventoryService) DeleteInventory(ctx context.Context, id int) error {
	return i.iRepo.DeleteInventory(ctx, id)
}
