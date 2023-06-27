package system

import (
	"amaru/archetype"
	"amaru/component"
	"amaru/net"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Player struct {
	game  *component.GameData
	query *query.Query
	space *cp.Space
}

func NewPlayer(space *cp.Space) *Player {
	return &Player{
		query: query.NewQuery(filter.Contains(
			component.Player,
		)),
	}
}

func (p *Player) Update(w donburi.World) {
	if p.game == nil {
		p.game = component.MustFindGame(w)
		if p.game == nil {
			return
		}
	}
	if p.space == nil {
		physics, _ := archetype.MustFindPhysics(w)
		if physics == nil {
			return
		}
		p.space = physics.Space
	}

	p.query.Each(w, func(entry *donburi.Entry) {
		player := component.Player.Get(entry)
		if player == nil || player.OutOfBounds {
			vector, changed := archetype.GetSpeed(w, p.game, player)
			dir := vector.Dot(*player.LastDirection)
			if changed && dir <= 0 {
				player.OutOfBounds = false
			}
			// if changing direction let the player move
			if player.OutOfBounds {
				player.Body.SetVelocityVector(cp.Vector{X: 0, Y: 0})
				return
			}
		}
		if player.Local {
			p.updateLocalPlayer(w, entry, player)
		} else {
			p.updateRemotePlayer(w, entry, player)
		}
	})
	p.space.Step(1.0 / 60.0)
}

func (p *Player) updateRemotePlayer(w donburi.World, entry *donburi.Entry, player *component.PlayerData) {
	if p.game.Session.PlayerMessage[player.ID] == nil {
		if player.PlayerCollision {
			player.Body.SetVelocityVector(cp.Vector{X: 0, Y: 0})
			player.PlayerCollision = false
		}
		return
	}
	player.PlayerCollision = false
	message := *p.game.Session.PlayerMessage[player.ID]

	vector := cp.Vector{
		X: message.Vector.X,
		Y: message.Vector.Y,
	}

	if message.Position != nil {
		animation := archetype.ActionsByKey[message.Animation]
		component.AnimationComponent.Get(entry).SelectAnimationByAction(animation)
		newPosition := cp.Vector{X: float64(message.Position.X), Y: float64(message.Position.Y)}

		player.Body.SetPosition(newPosition)
		p.game.Session.PlayerMessage[player.ID].Position = nil
	}
	if message.Vector.X == 0 && message.Vector.Y == 0 {
		player.Body.SetVelocityVector(cp.Vector{X: 0, Y: 0})
		pos := player.Body.Position()
		transform.Transform.Get(entry).LocalPosition = math.Vec2{X: pos.X - 16, Y: pos.Y + 16}
		transform.Transform.Get(player.Label).LocalPosition = math.Vec2{X: pos.X - 16, Y: pos.Y + 16}
		p.game.Session.PlayerMessage[player.ID] = nil
		return
	}

	sprite := component.Sprite.Get(entry)
	newVelocity := cp.Vector{X: float64(vector.X) * float64(sprite.Image.Bounds().Dx()), Y: float64(vector.Y) * float64(sprite.Image.Bounds().Dx())}
	player.Body.SetVelocityVector(newVelocity)

	pos := player.Body.Position()
	transform.Transform.Get(entry).LocalPosition = math.Vec2{X: pos.X - 16, Y: pos.Y + 16}
	transform.Transform.Get(player.Label).LocalPosition = math.Vec2{X: pos.X - 16, Y: pos.Y + 16}
}

func (p *Player) updateLocalPlayer(w donburi.World, entry *donburi.Entry, player *component.PlayerData) {
	player.PlayerCollision = false
	if player.Collision {
		player.Body.SetVelocityVector(cp.Vector{X: 0, Y: 0})

		vector, changed := archetype.GetSpeed(w, p.game, player)
		if player.LastDirection == nil {
			player.Collision = false
			pos := transform.Transform.Get(entry).LocalPosition
			anim := component.AnimationComponent.Get(entry).CurrentAnimation
			if anim != nil && changed {
				go p.game.Session.RemoteClient.SendMessage(net.Point{X: 0, Y: 0}, net.Point{X: pos.X, Y: pos.Y}, anim.Name)
			}
			var animname *string
			if anim != nil {
				animname = &anim.Name
			}
			go p.game.Session.RemoteClient.SetLocalPosition(&net.Point{X: pos.X, Y: pos.Y}, animname)
			return
		}
		dir := vector.Dot(*player.LastDirection)
		if changed && dir <= 0 {
			player.Collision = false
		}
		if player.Collision {
			archetype.SetAnimation(w, p.game, entry, player)
			pos := transform.Transform.Get(entry).LocalPosition
			anim := component.AnimationComponent.Get(entry).CurrentAnimation
			if anim != nil && changed {
				go p.game.Session.RemoteClient.SendMessage(net.Point{X: 0, Y: 0}, net.Point{X: pos.X, Y: pos.Y}, anim.Name)
			}
			var animname *string
			if anim != nil {
				animname = &anim.Name
			}
			go p.game.Session.RemoteClient.SetLocalPosition(&net.Point{X: pos.X, Y: pos.Y}, animname)
			return
		}
	}

	archetype.SetAnimation(w, p.game, entry, player)
	sprite := component.Sprite.Get(entry)
	vector, changed := archetype.GetSpeed(w, p.game, player)
	if vector.X != 0 || vector.Y != 0 {
		player.LastDirection = &cp.Vector{X: vector.X, Y: vector.Y}
	}
	vector = vector.Mult(p.game.Speed + 5)
	player.Body.SetVelocityVector(cp.Vector{X: float64(vector.X) * float64(sprite.Image.Bounds().Dx()), Y: float64(vector.Y) * float64(sprite.Image.Bounds().Dx())})

	pos := player.Body.Position()
	transform.Transform.Get(entry).LocalPosition = math.Vec2{X: pos.X - 16, Y: pos.Y + 16}
	transform.Transform.Get(player.Label).LocalPosition = math.Vec2{X: pos.X - 16, Y: pos.Y + 16}
	anim := component.AnimationComponent.Get(entry).CurrentAnimation
	if changed && anim != nil {
		go p.game.Session.RemoteClient.SendMessage(net.Point{X: vector.X, Y: vector.Y}, net.Point{X: pos.X, Y: pos.Y}, anim.Name)
	}
	var animname *string
	if anim != nil {
		animname = &anim.Name
	}
	go p.game.Session.RemoteClient.SetLocalPosition(&net.Point{X: pos.X - 16, Y: pos.Y + 16}, animname)
}
