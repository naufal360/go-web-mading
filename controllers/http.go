package controllers

import "gorm.io/gorm"

type HttpServer struct {
	db *gorm.DB
}

func NewHttpServer(db *gorm.DB) HttpServer {
	return HttpServer{db: db}
}
