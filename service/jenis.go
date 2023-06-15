package service

import (
	"context"
	"tim-b/entity"
	"tim-b/repository"
)

type JenisService struct {
	jRepo *repository.JenisRepository
}

func NewJenisService(jRepo *repository.JenisRepository) *JenisService {
	return &JenisService{
		jRepo: jRepo,
	}
}

func (j *JenisService) GetAllJenis(ctx context.Context) ([]entity.Jenis, error) {
	return j.jRepo.GetAllJenis(ctx)
}

func (j *JenisService) AddJenis(ctx context.Context, jenis entity.Jenis) (entity.Jenis, error) {
	err := j.jRepo.AddJenis(ctx, jenis)
	if err != nil {
		return entity.Jenis{}, err
	}
	return jenis, nil
}

func (j *JenisService) GetJenisByID(ctx context.Context, id int) (entity.Jenis, error) {
	return j.jRepo.GetJenisByID(ctx, id)
}

func (j *JenisService) UpdateJenis(ctx context.Context, jenis entity.Jenis) (entity.Jenis, error) {
	err := j.jRepo.UpdateJenis(ctx, jenis)
	if err != nil {
		return entity.Jenis{}, err
	}

	return jenis, nil
}

func (j *JenisService) DeleteJenis(ctx context.Context, id int) error {
	return j.jRepo.DeleteJenis(ctx, id)
}
