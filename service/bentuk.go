package service

import (
	"context"
	"tim-b/entity"
	"tim-b/repository"
)

type BentukService struct {
	bRepo *repository.BentukRepository
}

func NewBentukService(bRepo *repository.BentukRepository) *BentukService {
	return &BentukService{
		bRepo: bRepo,
	}
}

func (b *BentukService) GetAllBentuk(ctx context.Context) ([]entity.Bentuk, error) {
	return b.bRepo.GetAllBentuk(ctx)
}

func (b *BentukService) AddBentuk(ctx context.Context, bentuk entity.Bentuk) (entity.Bentuk, error) {
	err := b.bRepo.AddBentuk(ctx, bentuk)
	if err != nil {
		return entity.Bentuk{}, err
	}
	return bentuk, nil
}

func (b *BentukService) GetBentukByID(ctx context.Context, id int) (entity.Bentuk, error) {
	return b.bRepo.GetBentukByID(ctx, id)
}

func (b *BentukService) UpdateBentuk(ctx context.Context, bentuk entity.Bentuk) (entity.Bentuk, error) {
	err := b.bRepo.UpdateBentuk(ctx, bentuk)
	if err != nil {
		return entity.Bentuk{}, err
	}

	return bentuk, nil
}

func (b *BentukService) DeleteBentuk(ctx context.Context, id int) error {
	return b.bRepo.DeleteBentuk(ctx, id)
}
