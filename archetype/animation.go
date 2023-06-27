package archetype

import (
	"amaru/component"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func NewAnimationComponent(w donburi.World, animationEntry *donburi.Entry, sprite *donburi.Entry, drawables []*ebiten.Image, rate float32) *donburi.Entry {
	animationData := component.AnimationData{
		Sprite:     sprite,
		Drawables:  drawables,
		Animations: make(map[string]*component.Animation),
		Rate:       rate,
	}
	component.AnimationComponent.SetValue(animationEntry, animationData)

	return animationEntry
}
