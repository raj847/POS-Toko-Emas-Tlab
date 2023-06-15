package entity

import (
	"gorm.io/gorm"
)

type Jenis struct {
	gorm.Model
	Nama      string `json:"nama" gorm:"type:varchar(100);not null"`
	Kode      string `json:"kode" gorm:"type:varchar(10);not null"`
	CreatedBy string `json:"created_by" gorm:"type:varchar(50);not null"`
	UpdatedBy string `json:"updated_by" gorm:"type:varchar(50);not null"`
}

type JenisReq struct {
	Nama string `json:"nama" gorm:"type:varchar(100);not null"`
	Kode string `json:"kode" gorm:"type:varchar(10);not null"`
}
