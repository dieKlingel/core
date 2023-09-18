package service

type SystemService struct {
}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (system *SystemService) Version() string {
	return "Version: unknown"
}
