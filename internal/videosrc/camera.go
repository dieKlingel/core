package videosrc

import (
	"sync"

	"gocv.io/x/gocv"
)

var camera Camera = Camera{
	frame: gocv.NewMat(),
}

type Camera struct {
	camera *gocv.VideoCapture
	mutex  sync.Mutex
	frame  gocv.Mat
}

func (c *Camera) Open() {

	if c.camera == nil {
		print("create")
		c.camera, _ = gocv.VideoCaptureDevice(0)
		c.capture()
	}
}

func (c *Camera) Close() {
	c.mutex.Lock()
	c.camera.Close()
	c.camera = nil
	print("close")
	c.mutex.Unlock()
}

func (c *Camera) Read() gocv.Mat {

	c.mutex.Lock()
	mat := c.frame.Clone()
	c.mutex.Unlock()
	return mat
}

func (c *Camera) capture() {
	//ticker := time.NewTicker(16 * time.Millisecond)

	go func() {
		cancel := false
		for !cancel {
			//<-ticker.C
			c.mutex.Lock()
			cancel = camera.camera == nil
			if !cancel {
				c.camera.Read(&c.frame)
			}
			c.mutex.Unlock()
		}
		//ticker.Stop()
	}()

	c.camera.Read(&c.frame)
}
