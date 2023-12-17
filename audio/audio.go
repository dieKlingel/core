package audio

import (
	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

func init() {
	gst.Init(nil)
	glib.NewMainLoop(glib.MainContextDefault(), true)
}
