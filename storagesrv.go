package main

import (
	"os"

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

func (storageService StorageService) Read() *Configuration {
	configuration := Configuration{}
	file, err := os.Open(storageService.filename)
	if err != nil {
		panic(err.Error())
	}

	if err := yaml.NewDecoder(file).Decode(storageService); err != nil {
		panic(configuration)
	}

	return &configuration
}
