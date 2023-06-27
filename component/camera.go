package component

import (
	"github.com/yohamta/donburi"
)

type CameraData struct {
	Zoom     float64
	Width    int
	Height   int
	Disabled bool
}

var Camera = donburi.NewComponentType[CameraData]()
