package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStorageGet(t *testing.T) {
	filename := "storage-read-test.txt"
	expectedString := "read-value-from-storage"

	err := os.WriteFile(filename, []byte(expectedString), 0644)
	defer os.Remove(filename)
	assert.NoErrorf(t, err, `Failed to write test file "%s" %s`, filename, err)

	storage := Storage{
		File: filename,
	}

	actualString, err := storage.Get()
	assert.NoErrorf(t, err, `storage.Get("") failed: file to read storage file: %s`, err)
	assert.Equalf(t, expectedString, actualString, "Expected %s, actual data in file: %s", expectedString, actualString)
}

func TestStorageGetNotExistsFile(t *testing.T) {
	filename := "storage-not-exists.txt"
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		t.Fatalf(`Failed to remove file "%s" %s`, filename, err)
	}

	storage := Storage{
		File: filename,
	}
	actualString, err := storage.Get()
	assert.NoErrorf(t, err, `storage.Get("") failed: file to read storage file: %s`, err)
	assert.Emptyf(t, actualString, `storage.Get("") = %q, want match for empty string`, actualString)
}

func TestStorageSet(t *testing.T) {
	filename := "storage-Set-test.txt"
	expectedString := "Set-value-to-storage"

	storage := Storage{
		File: filename,
	}

	defer os.Remove(filename)

	for i := 1; i < 3; i++ {
		err := storage.Set(expectedString)

		assert.NoErrorf(t, err, `storage.Set("") failed: %v`, err)
		assert.FileExists(t, filename, `Storage file not exists after execute storage.Set("")`)

		actualData, err := os.ReadFile(filename)
		actualString := string(actualData)
		assert.NoErrorf(t, err, `storage.Set("") failed: file to read storage file: %s`, err)
		assert.Equalf(t, expectedString, actualString, "Data in file is not match with excpected value: %s != %s", expectedString, actualString)
	}
}

func TestStorageGetSet(t *testing.T) {
	filename := "storage-Get-Set-test.txt"
	expectedString := "Set-value-to-Get-from-storage"

	storage := Storage{
		File: filename,
	}

	err := storage.Set(expectedString)
	defer os.Remove(filename)
	assert.NoErrorf(t, err, `storage.Set("") failed: %v`, err)

	// re init Storage for reset cache
	storage = Storage{
		File: filename,
	}

	actualString, err := storage.Get()
	assert.NoErrorf(t, err, `storage.Get("") failed: file to read storage file: %s`, err)
	assert.Equalf(
		t, expectedString, actualString,
		"Expected %s, actual data in file: %s",
		expectedString, actualString,
	)
}

func TestStorageSetWithWrongPath(t *testing.T) {
	filename := "not-exists-dir/not-exist/random\n&@random.txt"
	expectedString := "Set-value-to-storage"

	storage := Storage{
		File: filename,
	}

	err := storage.Set(expectedString)

	assert.Errorf(t, err, `storage.Set("") not failed`)
	var PathError *os.PathError
	assert.ErrorAs(t, err, &PathError, "Expect for fs.PathError, got %v", err)
}

func TestStorageGetWithWrongPath(t *testing.T) {
	storage := Storage{
		File: os.TempDir(),
	}

	value, err := storage.Get()

	assert.Errorf(t, err, `storage.Set("") not failed`)
	var PathError *os.PathError
	assert.ErrorAs(t, err, &PathError, "Expect for fs.PathError, got %v", err)
	assert.Emptyf(t, value, "Not empty value: %s", value)
}
