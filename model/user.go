package model

import "time"

type Users struct {
	Id        string     `gorm:"primaryKey;type:varchar(255)" json:"id" form:"id"`
	Username  string     `gorm:"not null;type:varchar(100);uniqueIndex" json:"username" form:"username"`
	Password  string     `gorm:"not null;type:varchar(100)" json:"password" form:"password"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
