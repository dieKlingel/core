package config

type Mqtt struct {
	Uri      string `yaml:"uri"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
