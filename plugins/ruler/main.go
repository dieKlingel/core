package main

import "github.com/dieklingel/core/internal/core"

type Ruler struct{}

func NewPlugin() *Ruler {
	return &Ruler{}
}

func Run(ruler *Ruler, storageService core.StorageService) {
	// TODO: implement run
}
