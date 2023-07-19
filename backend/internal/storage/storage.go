package storage

import (
	"bar/internal/config"
	"os"
	"path"
)

func SaveFile(name string, data []byte) error {
	config := config.GetConfig()
	name = path.Join(config.StorageConfig.StoragePath, name)

	dir := path.Dir(name)
	os.MkdirAll(dir, 0755)

	return os.WriteFile(name, data, 0644)
}

func GetFile(name string) ([]byte, error) {
	config := config.GetConfig()
	name = path.Join(config.StorageConfig.StoragePath, name)
	return os.ReadFile(name)
}

func DeleteFile(name string) error {
	config := config.GetConfig()
	name = path.Join(config.StorageConfig.StoragePath, name)
	return os.Remove(name)
}
