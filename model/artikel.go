package model

import "time"

type Artikels struct {
	Id         string      `gorm:"primaryKey;type:varchar(255)" json:"id" form:"id"`
	Judul      string      `gorm:"not null;type:varchar(255)" json:"judul" form:"judul"`
	IsiArtikel string      `gorm:"not null" json:"isi_artikel" form:"isi_artikel"`
	CreatedAt  *time.Time  `json:"created_at,omitempty"`
	Komentar   []Komentars `gorm:"foreignKey:ArtikelId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"komentar"`
}

type ResponsArtikel struct {
	Id         string
	Judul      string
	Sneakpeak  string
	IsiArtikel string
	CreatedAt  *time.Time
	Komentar   []Komentars
}
