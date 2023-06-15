package entity

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	IDJenis    int     `json:"id_jenis_barang"`
	IDBentuk   int     `json:"id_bentuk_barang"`
	NamaBarang string  `json:"nama_barang" gorm:"type:varchar(255);not null"`
	KodeBarang string  `json:"kode" gorm:"type:varchar(20);not null"`
	Berat      float64 `json:"berat"`
	Kadar      float64 `json:"kadar"`
	HargaJual  float64 `json:"harga_jual"`
	Notes      string  `json:catatan`
	PhotoURL1  string  `json:"foto_1" gorm:"type:varchar(255)"`
	PhotoURL2  string  `json:"foto_2" gorm:"type:varchar(255)"`
	CreatedBy  string  `json:"created_by" gorm:"type:varchar(50);not null"`
	UpdatedBy  string  `json:"updated_by" gorm:"type:varchar(50);not null"`
	DeletedBy  string  `json:"deleted_by" gorm:"type:varchar(50);not null"`
}

type InventoryReq struct {
	IDJenis    int     `json:"id_jenis_barang"`
	IDBentuk   int     `json:"id_bentuk_barang"`
	NamaBarang string  `json:"nama_barang" gorm:"type:varchar(255);not null"`
	Berat      float64 `json:"berat"`
	Kadar      float64 `json:"kadar"`
	HargaJual  float64 `json:"harga_jual"`
	Notes      string  `json:catatan`
	PhotoURL1  string  `json:"foto_1" gorm:"type:varchar(255)"`
	PhotoURL2  string  `json:"foto_2" gorm:"type:varchar(255)"`
}

type InventoryRead struct {
	gorm.Model
	IDJenis    int     `json:"id_jenis_barang"`
	IDBentuk   int     `json:"id_bentuk_barang"`
	NamaJenis  string  `json:"nama_jenis"`
	NamaBentuk string  `json:"nama_bentuk"`
	NamaBarang string  `json:"nama_barang" gorm:"type:varchar(255);not null"`
	KodeBarang string  `json:"kode" gorm:"type:varchar(20);not null"`
	Berat      float64 `json:"berat"`
	Kadar      float64 `json:"kadar"`
	HargaJual  float64 `json:"harga_jual"`
	Notes      string  `json:catatan`
	PhotoURL1  string  `json:"foto_1" gorm:"type:varchar(255)"`
	PhotoURL2  string  `json:"foto_2" gorm:"type:varchar(255)"`
	CreatedBy  string  `json:"created_by" gorm:"type:varchar(50);not null"`
	UpdatedBy  string  `json:"updated_by" gorm:"type:varchar(50);not null"`
	DeletedBy  string  `json:"deleted_by" gorm:"type:varchar(50);not null"`
}
