package main

import (
	"os"

	"github.com/dieklingel/core/internal/core"
	"gopkg.in/yaml.v3"
)

type StorageService struct {
	filename string
}

func NewStorageService(filename string) *StorageService {
	return &StorageService{
		filename: filename,
	}
}

func (storageService StorageService) Read() *core.Configuration {
	configuration := core.Configuration{}
	file, err := os.Open(storageService.filename)
	if err != nil {
		panic(err.Error())
	}

	if err := yaml.NewDecoder(file).Decode(&configuration); err != nil {
		panic(err)
	}

	return &configuration
}
