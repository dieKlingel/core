package camerasrv

import (
	"log"

	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/io"
	"gorm.io/gorm"
)

type CameraService struct {
	database *gorm.DB
	camera   io.Camera
}

const (
	DefaultCameraPipeline = "videotestsrc ! video/x-raw, framerate=30/1, width=1280, height=720 ! appsink name=rawsink"
)

func NewService(db *gorm.DB) core.CameraService {
	db.AutoMigrate(&CameraServiceSettings{})

	return &CameraService{
		database: db,
	}
}

func (service *CameraService) CameraPipeline() string {
	settings := CameraServiceSettings{
		CameraPipeline: DefaultCameraPipeline,
	}
	service.database.FirstOrCreate(&settings)
	return settings.CameraPipeline
}

func (service *CameraService) SetCameraPipeline(pipeline string) {
	settings := CameraServiceSettings{
		CameraPipeline: pipeline,
	}
	service.database.FirstOrCreate(&settings)
	// TODO: shutdown camera and restart, but with migrating all streams
}

func (service *CameraService) NewCameraStream(codec io.CameraCodec) *io.Stream {
	if service.camera == nil {
		camera, err := io.NewCamera(service.CameraPipeline())
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		service.camera = camera
	}

	stream, err := service.camera.NewStream(codec)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return stream
}

func (service *CameraService) ReleaseCameraStream(stream *io.Stream) {
	service.camera.ReleaseStream(stream)
	if !service.camera.HasStream() {
		service.camera = nil
	}
}
