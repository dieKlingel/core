package config

type Action struct {
	Trigger     string `yaml:"trigger"`
	Environment string `yaml:"environment"`
	Script      string `yaml:"script"`
}
