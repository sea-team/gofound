package controller

import (
	service2 "gofound/web/service"
)

var srv *Services

type Services struct {
	Base     *service2.Base
	Index    *service2.Index
	Database *service2.Database
	Word     *service2.Word
}

func NewServices() {
	srv = &Services{
		Base:     service2.NewBase(),
		Index:    service2.NewIndex(),
		Database: service2.NewDatabase(),
		Word:     service2.NewWord(),
	}
}
