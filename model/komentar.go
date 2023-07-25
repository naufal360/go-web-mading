package model

import "time"

type Komentars struct {
	Id          string     `gorm:"primaryKey;type:varchar(255)" json:"id" form:"id"`
	Nama        string     `gorm:"not null;type:varchar(255)" json:"nama" form:"nama"`
	Email       string     `gorm:"not null;type:varchar(255)" json:"email" form:"email"`
	IsiKomentar string     `json:"isi_komentar" form:"isi_komentar"`
	ArtikelId   string     `gorm:"type:varchar(255)" json:"artikel_id"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Artikel     *Artikels  `json:"artikel"`
}
