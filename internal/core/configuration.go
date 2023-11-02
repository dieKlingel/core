package core

import "gopkg.in/yaml.v3"

type Configuration struct {
	Actions []Action           `yaml:"actions"`
	Sip     SipConfiguration   `yaml:"sip"`
	Mqtt    MqttConfiguration  `yaml:"mqtt"`
	Plugins yaml.Node          `yaml:"plugins"`
	Media   MediaConfiguration `yaml:"media"`
	Redis   RedisConfiguration `yaml:"redis"`
}

type SipConfiguration struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MqttConfiguration struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MediaConfiguration struct {
	Camera struct {
		Src string `yaml:"src"`
	} `yaml:"camera"`
}

type RedisConfiguration struct {
	Host string `yaml:"host"`
}
