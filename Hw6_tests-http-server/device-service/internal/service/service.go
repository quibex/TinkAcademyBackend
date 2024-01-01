package service

import (
	"device-service/internal/device"
	"device-service/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Repository
type Repository interface {
	GetDevice(string) (device.Device, error)
	CreateDevice(device.Device) error
	DeleteDevice(string) error
	UpdateDevice(device.Device) error
}

type Service struct {
	Repository
}

func New() Service {
	service := Service{repository.New()}
	return service
}
