package storage

import (
	"bytes"
	"errors"
	"os"
)

type Interface interface {
	Get() ([]byte, error)
	Set(value []byte) error
}

type Storage struct {
	File  string
	value []byte
}

func (storage *Storage) Get() ([]byte, error) {
	if storage.value == nil {
		_, err := os.Stat(storage.File)
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		storage.value, err = os.ReadFile(storage.File)
		if err != nil {
			return nil, err
		}
	}

	return storage.value, nil
}

func (storage *Storage) Set(value []byte) error {
	if bytes.Equal(storage.value, value) {
		return nil
	}

	err := os.WriteFile(storage.File, value, 0644)
	if err != nil {
		return err
	}

	storage.value = value
	return nil
}
