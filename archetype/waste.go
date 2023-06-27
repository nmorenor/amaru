package archetype

import (
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"fmt"
	"math/rand"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func PlaceWasteComponents(world donburi.World, space *cp.Space, numWaste int, debug *component.DebugData, shapeList []cp.BB, mapWidth float64, mapHeight float64) []*component.WasteData {
	wasteList := []*component.WasteData{}
	for i := 0; i < numWaste; i++ {
		for {
			x := rand.Float64() * mapWidth
			y := rand.Float64() * mapHeight

			points := []math.Vec2{}
			points = append(points, math.Vec2{X: x, Y: y - 32})
			points = append(points, math.Vec2{X: x + 32, Y: y - 32})
			points = append(points, math.Vec2{X: x + 32, Y: y})
			points = append(points, math.Vec2{X: x, Y: y})

			wastePath := assets.Path{
				Points: points,
				Loops:  true,
			}
			wasteShape := CreateBoxFromPath(space, wastePath, component.WasteCollisionType)

			overlapping := false
			for _, waste := range wasteList {
				if engine.BBIntersects(waste.Shape.BB(), wasteShape.BB()) {
					overlapping = true
					break
				}
			}

			for _, box := range shapeList {
				if engine.BBIntersects(box, wasteShape.BB()) {
					overlapping = true
					break
				}
			}

			if !overlapping {
				newWaste := placeWasteFromPath(world, space, debug, fmt.Sprint(i), wastePath, false, wasteShape)
				wasteList = append(wasteList, newWaste)
				break
			}
		}
	}
	return wasteList
}

func PlaceRemoteWasteFromPath(world donburi.World, space *cp.Space, debug *component.DebugData, id string, wastePath assets.Path, collected bool) *component.WasteData {
	return placeWasteFromPath(world, space, debug, id, wastePath, collected, nil)
}

func placeWasteFromPath(world donburi.World, space *cp.Space, debug *component.DebugData, id string, wastePath assets.Path, collected bool, wasteShape *cp.Shape) *component.WasteData {
	if wasteShape == nil {
		wasteShape = CreateBoxFromPath(space, wastePath, component.WasteCollisionType)
	}
	cameraEntry := MustFindCamera(world)
	camera := component.Camera.Get(cameraEntry)
	position := wasteShape.Body().Position()
	// position is center of shape, shape is 32x32
	x := position.X - 16
	y := position.Y + 16
	newWaste := &component.WasteData{
		Id:        id,
		Path:      wastePath,
		Shape:     wasteShape,
		Collected: false,
	}
	if !camera.Disabled {
		debug.Shapes = append(debug.Shapes, newWaste.Shape)
		waste := world.Entry(
			world.Create(
				component.Sprite,
				transform.Transform,
				component.Waste,
			),
		)
		newWaste.Entry = waste
		newWaste.Shape.UserData = waste
		component.Waste.SetValue(waste, *newWaste)

		spriteData := component.SpriteData{}
		spriteData.Layer = component.SpriteLayerDefault
		spriteData.Pivot = component.SpritePivotTopLeft
		spriteData.Image = assets.WasteImage
		component.Sprite.SetValue(waste, spriteData)
		transform.Transform.Get(waste).LocalPosition = math.Vec2{X: x, Y: y}
	}
	return newWaste
}
