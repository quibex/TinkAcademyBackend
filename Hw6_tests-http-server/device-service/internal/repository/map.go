package repository

import (
	"device-service/internal/device"
	"fmt"
	"sync"
)

type Repository struct {
	db sync.Map
}

func New() *Repository {
	return &Repository{sync.Map{}}
}

func (s *Repository) GetDevice(serialNum string) (device.Device, error) {
	d, ok := s.db.Load(serialNum)
	if !ok {
		return device.Device{}, fmt.Errorf("device not found: %s", serialNum)
	}
	return d.(device.Device), nil
}

func (s *Repository) CreateDevice(d device.Device) error {
	if _, ok := s.db.Load(d.SerialNum); ok {
		return fmt.Errorf("device already exists: %s", d.SerialNum)
	}
	s.db.Store(d.SerialNum, d)
	return nil
}

func (s *Repository) DeleteDevice(serialNum string) error {
	if _, ok := s.db.Load(serialNum); !ok {
		return fmt.Errorf("device not found: %s", serialNum)
	}
	s.db.Delete(serialNum)
	return nil
}

func (s *Repository) UpdateDevice(d device.Device) error {
	if _, ok := s.db.Load(d.SerialNum); !ok {
		return fmt.Errorf("device not found: %s", d.SerialNum)
	}
	s.db.Store(d.SerialNum, d)
	return nil
}
