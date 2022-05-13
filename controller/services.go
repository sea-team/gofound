package controller

import "gofound/service"

var srv *Services

type Services struct {
	Base     *service.Base
	Index    *service.Index
	Database *service.Database
	Word     *service.Word
}

func NewServices() {
	srv = &Services{
		Base:     service.NewBase(),
		Index:    service.NewIndex(),
		Database: service.NewDatabase(),
		Word:     service.NewWord(),
	}
}
