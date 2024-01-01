package service

import (
	. "device-service/internal/service/device"
	"testing"
)

func TestCreateDevice(t *testing.T) {
	service := New()
	wantDevice := Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gotDevice, err := service.GetDevice(wantDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if wantDevice != gotDevice {
		t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
	}
}

func TestCreateMultipleDevices(t *testing.T) {
	service := New()
	devices := []Device{
		{
			SerialNum: "123",
			Model:     "model1",
			IP:        "1.1.1.1",
		},
		{
			SerialNum: "124",
			Model:     "model2",
			IP:        "1.1.1.2",
		},
		{
			SerialNum: "125",
			Model:     "model3",
			IP:        "1.1.1.3",
		},
	}

	for _, d := range devices {
		err := service.CreateDevice(d)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	for _, wantDevice := range devices {
		gotDevice, err := service.GetDevice(wantDevice.SerialNum)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wantDevice != gotDevice {
			t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
		}
	}
}

func TestCreateDuplicate(t *testing.T) {
	service := New()
	wantDevice := Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.CreateDevice(wantDevice)
	if err == nil {
		t.Errorf("want error, but got nil")
	}

}

func TestGetDeviceUnexisting(t *testing.T) {
	service := New()
	wantDevice := Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = service.GetDevice("1")
	if err == nil {
		t.Error("want error, but got nil")
	}
}

func TestDeleteDevice(t *testing.T) {
	service := New()
	newDevice := Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(newDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.DeleteDevice(newDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = service.GetDevice(newDevice.SerialNum)
	if err == nil {
		t.Error("want error, but got nil")
	}
}

func TestDeleteDeviceUnexisting(t *testing.T) {
	service := New()

	err := service.DeleteDevice("123")
	if err == nil {
		t.Errorf("want error, but got nil")
	}
}

func TestUpdateDevice(t *testing.T) {
	service := New()
	device := Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(device)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	newDevice := Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.2",
	}
	err = service.UpdateDevice(newDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gotDevice, err := service.GetDevice(newDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if gotDevice != newDevice {
		t.Errorf("new device %+#v not equal got device %+#v", newDevice, gotDevice)
	}
}

func TestUpdateDeviceUnexsting(t *testing.T) {
	service := New()
	device := Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(device)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	newDevice := Device{
		SerialNum: "124",
		Model:     "model1",
		IP:        "1.1.1.2",
	}
	err = service.UpdateDevice(newDevice)
	if err == nil {
		t.Errorf("want err, but got nil")
	}
}
