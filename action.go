package main

type Action struct {
	Trigger string `json:"trigger" yaml:"trigger"`
	Lane    string `json:"lane" yaml:"lane"`
}
