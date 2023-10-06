package camerasrv

type CameraServiceSettings struct {
	Id             uint64 `gorm:"primaryKey"`
	CameraPipeline string
}
