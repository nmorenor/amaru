package archetype

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"amaru/component"
)

func NewCamera(w donburi.World, width, height int, startPosition math.Vec2) *donburi.Entry {
	camera := w.Entry(
		w.Create(
			transform.Transform,
			component.Camera,
		),
	)

	cameraCamera := component.Camera.Get(camera)
	cameraCamera.Zoom = 1.0
	cameraCamera.Width = width
	cameraCamera.Height = height
	transform.Transform.Get(camera).LocalPosition = startPosition

	return camera
}

func MustFindCamera(w donburi.World) *donburi.Entry {
	camera, ok := query.NewQuery(filter.Contains(component.Camera)).First(w)
	if !ok {
		panic("no camera found")
	}

	return camera
}
