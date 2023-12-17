package config

type IceServer struct {
	Urls       string `yaml:"urls"`
	Username   string `yaml:"username"`
	Credential string `yaml:"credential"`
}

type Rtc struct {
	IceServers []IceServer `yaml:"ice-servers"`
}
