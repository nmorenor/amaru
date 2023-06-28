package archetype

import (
	"amaru/assets"
	"amaru/component"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

var (
	AnimalAnimationAction = &component.Animation{
		Name:   "default",
		Frames: []int{0, 1, 2, 3, 4, 5},
	}

	AnimalActions = []*component.Animation{AnimalAnimationAction}
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
				component.AnimationComponent,
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
			animationEntry := NewAnimationComponent(world, animal, animal, assets.TurtleSpriteSheet.Drawables(), 0.4)
			animationData := component.AnimationComponent.Get(animationEntry)
			animationData.AddAnimations(WasteActions)
			animationData.SelectAnimationByName("default")
			animationData.CurrentAnimation.Loop = true

			spriteData.Image = animationData.Cell()
		} else {
			animationEntry := NewAnimationComponent(world, animal, animal, assets.SealSpriteSheet.Drawables(), 0.4)
			animationData := component.AnimationComponent.Get(animationEntry)
			animationData.AddAnimations(WasteActions)
			animationData.SelectAnimationByName("default")
			animationData.CurrentAnimation.Loop = true

			spriteData.Image = animationData.Cell()
		}
		component.Sprite.SetValue(animal, spriteData)
		transform.Transform.Get(animal).LocalPosition = math.Vec2{X: vert.X - 32, Y: vert.Y + 32}
	}
}
