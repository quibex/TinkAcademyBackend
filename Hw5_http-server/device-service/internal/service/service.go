package service

import (
	"device-service/internal/service/device"
	"device-service/internal/service/map_db"
)

type Service interface {
	GetDevice(string) (device.Device, error)
	CreateDevice(device.Device) error
	DeleteDevice(string) error
	UpdateDevice(device.Device) error
}

func New() Service {
	service := map_db.New()
	return service
}
