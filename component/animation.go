package component

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

const (
	DefaultPlayerAnimation = "downstop"
)

type Animation struct {
	Name   string
	Frames []int
	Loop   bool
}

type AnimationData struct {
	Sprite           *donburi.Entry        // Sprite
	Drawables        []*ebiten.Image       // Renderables
	Animations       map[string]*Animation // All possible animations
	CurrentAnimation *Animation            // The current animation
	CurrentFrame     int                   // The current animation frame number
	Rate             float32               // How often frames should increment, in seconds.
	index            int                   // What frame in the is being used
	Change           float32               // The time since the last incrementation
	Def              *Animation            // The default animation to play when nothing else is playing
}

// SelectAnimationByName sets the current animation. The name must be
// registered.
func (ac *AnimationData) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
	ac.index = 0
}

// SelectAnimationByAction sets the current animation.
// An nil action value selects the default animation.
func (ac *AnimationData) SelectAnimationByAction(action *Animation) {
	if ac.CurrentAnimation != nil && ac.CurrentAnimation.Name == action.Name {
		return
	}
	ac.CurrentAnimation = action
	ac.index = 0
}

// AddDefaultAnimation adds an animation which is used when no other animation is playing.
func (ac *AnimationData) AddDefaultAnimation(action *Animation) {
	ac.AddAnimation(action)
	ac.Def = action
}

// AddAnimation registers an animation under its name, making it available
// through SelectAnimationByName.
func (ac *AnimationData) AddAnimation(action *Animation) {
	ac.Animations[action.Name] = action
}

// AddAnimations registers all given animations.
func (ac *AnimationData) AddAnimations(actions []*Animation) {
	for _, action := range actions {
		ac.AddAnimation(action)
	}
}

// Cell returns the drawable for the current frame.
func (ac *AnimationData) Cell() *ebiten.Image {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No frame data for this animation. Selecting zeroth drawable. If this is incorrect, add an action to the animation.")
		return ac.Drawables[0]
	}
	idx := ac.CurrentAnimation.Frames[ac.index]
	ac.CurrentFrame = idx
	return ac.Drawables[idx]
}

// NextFrame advances the current animation by one frame.
func (ac *AnimationData) NextFrame() {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No frame data for this animation")
		return
	}

	ac.index++
	ac.Change = 0
	if ac.index >= len(ac.CurrentAnimation.Frames) {
		ac.index = 0

		if !ac.CurrentAnimation.Loop {
			ac.CurrentAnimation = nil
			return
		}
	}
}

var AnimationComponent = donburi.NewComponentType[AnimationData]()
