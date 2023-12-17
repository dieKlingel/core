package config

type Camera struct {
	Src string `yaml:"src"`
}

type Audio struct {
	Src  string `yaml:"src"`
	Sink string `yaml:"sink"`
}

type Media struct {
	Camera Camera `yaml:"camera"`
	Audio  Audio  `yaml:"audio"`
}
