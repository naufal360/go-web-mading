package config

import (
	"fmt"
	"jewepe/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm() *gorm.DB {

	host := "localhost"
	username := "root"
	password := ""
	dbname := "jewepe_db"
	port := "3306"

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbname,
	)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		logrus.Error("Tidak dapat terkoneksi database!")
		panic(err)
	}

	fmt.Println("Berhasil terkoneksi dengan database!")

	if err := MigrateDB(db); err != nil {
		logrus.Error("Gagal migrasi database!")
	}
	fmt.Println("Migrasi database berhasil!")

	return db
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		model.Artikels{},
		model.Users{},
		model.Komentars{},
	)
}
