package config_test

import (
	"device-service/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func FuzzMustLoad(f *testing.F) {
	f.Fuzz(func(t *testing.T, port uint8) {
		err := os.Setenv("HOST", "localhost")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = os.Setenv("PORT", strconv.Itoa(int(port)))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		cfg := config.MustLoad()

		assert.Equal(t, "localhost:"+strconv.Itoa(int(port)), cfg.Address)
	})
}

func TestMustLoad(t *testing.T) {
	err := os.Setenv("HOST", "127.0.0.1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = os.Setenv("PORT", "123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	cfg := config.MustLoad()

	assert.Equal(t, "127.0.0.1:123", cfg.Address)
}

func TestMustLoadEmpty(t *testing.T) {
	err := os.Setenv("HOST", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = os.Setenv("PORT", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Panics(t,
		func() {
			_ = config.MustLoad()
		},
	)
}

func TestMustLoadValidate(t *testing.T) {
	err := os.Setenv("HOST", "fg&s.")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = os.Setenv("PORT", "99999")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Panics(t,
		func() {
			_ = config.MustLoad()
		},
	)
}
