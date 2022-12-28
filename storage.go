package storage

import (
	"errors"
	"os"
)

type Interface interface {
	Get() (string, error)
	Set(value string) error
}

type Storage struct {
	File   string
	value  string
	loaded bool
}

func (storage Storage) Get() (string, error) {
	if storage.loaded == false {
		storage.loaded = true
		if _, err := os.Stat(storage.File); errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		data, err := os.ReadFile(storage.File)
		if err != nil {
			return "", err
		}

		storage.value = string(data)
	}

	return storage.value, nil
}

func (storage *Storage) Set(value string) error {
	if storage.loaded && storage.value == value {
		return nil
	}

	err := os.WriteFile(storage.File, []byte(value), 0644)
	if err != nil {
		return err
	}

	storage.value = value
	storage.loaded = true
	return nil
}
