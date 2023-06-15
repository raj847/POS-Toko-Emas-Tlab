package entity

import (
	"time"

	"gorm.io/gorm"
)

type Statusx string

const (
	Aktif    Statusx = "aktif"
	NonAktif Statusx = "non-aktif"
)

type User struct {
	gorm.Model
	NamaLengkap  string    `json:"nama_lengkap" gorm:"type:varchar(255);not null"`
	Username     string    `json:"username" gorm:"type:varchar(50);not null"`
	Password     string    `json:"-" gorm:"type:varchar(255);not null"`
	NoHp         string    `json:"no_hp" gorm:"type:varchar(20);not null"`
	Email        string    `json:"email" gorm:"type:text;not null"`
	TanggalMasuk time.Time `json:"tanggal_masuk" gorm:"type:date;not null"`
	Status       Statusx   `json:"status" gorm:"type:varchar(255);not null"`
	CreatedBy    string    `json:"created_by" gorm:"type:varchar(50);not null"`
	UpdatedBy    string    `json:"updated_by" gorm:"type:varchar(50);not null"`
	Role         string    `json:"role" gorm:"type:varchar(50);not null"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	NamaLengkap     string `json:"nama_lengkap" gorm:"type:varchar(255);not null"`
	Username        string `json:"username" gorm:"type:varchar(50);not null"`
	Password        string `json:"password" gorm:"type:varchar(255);not null"`
	NoHp            string `json:"no_hp" gorm:"type:varchar(20);not null"`
	Email           string `json:"email" gorm:"type:text;not null"`
	TanggalMasukStr string `json:"tanggal_masuk" gorm:"type:date;not null"`
	Status          string `json:"status" gorm:"type:varchar(255);not null"`
	CreatedBy       string `json:"-" gorm:"type:varchar(50);not null"`
	UpdatedBy       string `json:"-" gorm:"type:varchar(50);not null"`
	Role            string `json:"role" gorm:"type:varchar(50);not null"`
}
