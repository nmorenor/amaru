package archetype

import (
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"time"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

func MustFindPhysics(w donburi.World) (*component.PhysicsData, *donburi.Entry) {
	physics, ok := query.NewQuery(filter.Contains(component.Physics)).First(w)
	if !ok {
		panic("no physics found")
	}

	return component.Physics.Get(physics), physics
}

func SetupColliders(world donburi.World) {
	physics, _ := MustFindPhysics(world)
	game := component.MustFindGame(world)
	boxCollisionHandler := physics.Space.NewCollisionHandler(component.PlayerCollisionType, component.BoxCollisionType)
	boxCollisionHandler.PreSolveFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		playerShape, boxShape := arb.Shapes()
		playerPoly, ok := playerShape.Class.(*cp.PolyShape)
		if !ok {
			// Handle the case where the player shape is not a polygon shape
			return false
		}
		boxPoly, ok := boxShape.Class.(*cp.PolyShape)
		if !ok {
			// Handle the case where the box shape is not a polygon shape
			return false
		}

		if playerShape.UserData != nil {
			playerEntry := playerShape.UserData.(*donburi.Entry)
			if playerEntry == nil || !world.Valid(playerEntry.Entity()) {
				return false
			}
			player := component.Player.Get(playerEntry)
			if player.Collision || !player.Local {
				return true
			}
			vector, _ := GetSpeed(world, game, player)

			if vector.X == 0 && vector.Y == 0 {
				player.Collision = true
				return true
			}

			overlapDirection := engine.OverlapSide(boxPoly.BB(), playerPoly.BB())
			if vector.X == overlapDirection.X && overlapDirection.X != 0 {
				player.Collision = true
				return player.Collision
			} else if vector.Y == overlapDirection.Y && overlapDirection.Y != 0 {
				player.Collision = true
				return player.Collision
			}
		}

		return false
	}
	wasteCollisionHandler := physics.Space.NewCollisionHandler(component.PlayerCollisionType, component.WasteCollisionType)
	wasteCollisionHandler.PreSolveFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		playerShape, wasteShape := arb.Shapes()
		if playerShape.UserData == nil || wasteShape.UserData == nil {
			return false
		}
		playerEntry := playerShape.UserData.(*donburi.Entry)
		wasteEntry := wasteShape.UserData.(*donburi.Entry)

		if playerEntry == nil || wasteEntry == nil || !world.Valid(playerEntry.Entity()) || !world.Valid(wasteEntry.Entity()) {
			return false
		}

		player := component.Player.Get(playerEntry)
		waste := component.Waste.Get(wasteEntry)

		if waste.Collected {
			component.Sprite.Get(wasteEntry).Hidden = true
			return false
		}
		wLocation := game.Session.RemoteClient.GameData.WasteLocations[waste.Id]
		if wLocation.Collected {
			component.Sprite.Get(wasteEntry).Hidden = true
			return false
		}
		waste.Collected = true
		wLocation.Collected = true
		component.Sprite.Get(wasteEntry).Hidden = true

		game.Session.RemoteClient.GameData.SessionParticipants[player.ID].Score += component.WastePoints
		game.CollectedWaste += 1

		if player.Local {
			if !game.Muted {
				PlayCollectedAudio()
			}
		}

		return false
	}

	animalCollisionHandler := physics.Space.NewCollisionHandler(component.PlayerCollisionType, component.AnimalCollisionType)
	animalCollisionHandler.PreSolveFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		playerShape, animalShape := arb.Shapes()
		if playerShape.UserData == nil || animalShape.UserData == nil {
			return false
		}
		playerEntry := playerShape.UserData.(*donburi.Entry)
		animalEntry := animalShape.UserData.(*donburi.Entry)

		if playerEntry == nil || animalEntry == nil || !world.Valid(playerEntry.Entity()) || !world.Valid(animalEntry.Entity()) {
			return false
		}

		player := component.Player.Get(playerEntry)
		animal := component.Animal.Get(animalEntry)

		if animal.Collected {
			return false
		}

		animal.Collected = true
		component.Sprite.Get(animalEntry).Hidden = true
		game.Session.RemoteClient.GameData.SessionParticipants[player.ID].Score += component.AnimalPoints

		if player.Local {
			if !game.Muted {
				PlayHeronAudio()
			}
		}

		return false
	}
	playersCollisionHandler := physics.Space.NewCollisionHandler(component.PlayerCollisionType, component.PlayerCollisionType)
	playersCollisionHandler.PreSolveFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		onePlayerShape, otherPlayerShape := arb.Shapes()
		if onePlayerShape.UserData == nil || otherPlayerShape.UserData == nil {
			return false
		}
		onePlayerEntry := onePlayerShape.UserData.(*donburi.Entry)
		otherPlayerEntry := otherPlayerShape.UserData.(*donburi.Entry)

		if onePlayerEntry == nil || otherPlayerEntry == nil || !world.Valid(onePlayerEntry.Entity()) || !world.Valid(otherPlayerEntry.Entity()) {
			return false
		}

		onePlayer := component.Player.Get(onePlayerEntry)
		otherPlayer := component.Player.Get(otherPlayerEntry)

		if onePlayer.LastPlayerCollision == nil {
			game.Session.RemoteClient.GameData.SessionParticipants[onePlayer.ID].Score -= component.PlayerCollisionPoints
			onePlayer.LastPlayerCollision = engine.Ptr(time.Now())
			onePlayer.PlayerCollision = true
			otherPlayer.PlayerCollision = true
		} else if onePlayer.LastPlayerCollision != nil && time.Since(*onePlayer.LastPlayerCollision).Seconds() > 2 {
			game.Session.RemoteClient.GameData.SessionParticipants[onePlayer.ID].Score -= component.PlayerCollisionPoints
			onePlayer.LastPlayerCollision = engine.Ptr(time.Now())
			onePlayer.PlayerCollision = true
			otherPlayer.PlayerCollision = true
		}
		if otherPlayer.LastPlayerCollision == nil {
			game.Session.RemoteClient.GameData.SessionParticipants[otherPlayer.ID].Score -= component.PlayerCollisionPoints
			otherPlayer.LastPlayerCollision = engine.Ptr(time.Now())
			onePlayer.PlayerCollision = true
			otherPlayer.PlayerCollision = true
		} else if otherPlayer.LastPlayerCollision != nil && time.Since(*otherPlayer.LastPlayerCollision).Seconds() > 2 {
			game.Session.RemoteClient.GameData.SessionParticipants[otherPlayer.ID].Score -= component.PlayerCollisionPoints
			otherPlayer.LastPlayerCollision = engine.Ptr(time.Now())
			onePlayer.PlayerCollision = true
			otherPlayer.PlayerCollision = true
		}
		if onePlayer.Local {
			if !game.Muted {
				PlayShipAudio()
			}
		}

		return true
	}
}

func CreateBoxFromPath(space *cp.Space, path assets.Path, collisionType cp.CollisionType) *cp.Shape {
	// Find the minimum and maximum X and Y coordinates in the path
	minX, minY := path.Points[0].X, path.Points[0].Y
	maxX, maxY := path.Points[0].X, path.Points[0].Y

	for _, point := range path.Points[1:] {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	// Calculate the center, width, and height of the box
	center := cp.Vector{X: (minX + maxX) / 2, Y: (minY + maxY) / 2}
	width := maxX - minX
	height := maxY - minY

	// Create a static body (non-moving)
	body := cp.NewStaticBody()
	body.SetPosition(center)

	// Create a box shape and add it to the static body
	shape := cp.NewBox(body, width, height, 0)
	shape.SetElasticity(1)
	shape.SetFriction(1)
	shape.SetCollisionType(collisionType)

	// Add the shape to the Chipmunk physics space
	space.AddShape(shape)

	return shape
}

func SetupSpaceForLevel(level assets.Level) (*cp.Space, []*cp.Shape) {
	space := cp.NewSpace()
	space.SetGravity(cp.Vector{X: 0, Y: 0})
	shapes := []*cp.Shape{}
	// Add paths to the space as static shapes
	for _, path := range level.Paths {
		if path.Loops && len(path.Points) == 4 {
			shapes = append(shapes, CreateBoxFromPath(space, path, component.BoxCollisionType))
		}
	}
	return space, shapes
}
