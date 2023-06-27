package archetype

import (
	"amaru/assets"
	"amaru/component"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func PlaceAnimalComponents(world donburi.World, space *cp.Space, debug *component.DebugData, animalPaths []assets.Path, mapWidth float64, mapHeight float64) {
	for index, animalPath := range animalPaths {
		animalShape := CreateBoxFromPath(space, animalPath, component.AnimalCollisionType)
		box := animalShape.Class.(*cp.PolyShape)
		vert := animalShape.Body().LocalToWorld(box.Vert(0))

		newAnimal := &component.AnimalData{
			X:         vert.X,
			Y:         vert.Y,
			Shape:     animalShape,
			Collected: false,
		}
		debug.Shapes = append(debug.Shapes, newAnimal.Shape)
		animal := world.Entry(
			world.Create(
				component.Sprite,
				transform.Transform,
				component.Animal,
			),
		)
		newAnimal.Entry = animal
		newAnimal.Shape.UserData = animal
		component.Animal.SetValue(animal, *newAnimal)

		spriteData := component.SpriteData{}
		spriteData.Layer = component.SpriteLayerDefault
		spriteData.Pivot = component.SpritePivotTopLeft
		if index == 0 || index%2 == 0 {
			spriteData.Image = assets.TurtleImage
		} else {
			spriteData.Image = assets.SealImage
		}
		component.Sprite.SetValue(animal, spriteData)
		transform.Transform.Get(animal).LocalPosition = math.Vec2{X: vert.X - 32, Y: vert.Y + 32}
	}
}
