package system

import (
	"amaru/archetype"
	"amaru/component"
	"amaru/engine"
	"amaru/net"
	"context"
	"time"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type RemoteSystem struct {
	game                    *component.GameData
	query                   *query.Query
	space                   *cp.Space
	debug                   *component.DebugData
	shouldEnd               bool
	shouldPlaceWaste        bool
	sessionLeaveMessages    *engine.Queue[net.SessionLeaveMessage]
	sessionJoinMessages     *engine.Queue[net.SessionJoinMessage]
	initialPositionMessages *engine.Queue[net.RemoteInitialPositionMessage]
	startTime               time.Time
}

func NewRemoteSystem() *RemoteSystem {
	return &RemoteSystem{
		query: query.NewQuery(filter.Contains(
			component.Player,
		)),
		sessionLeaveMessages:    engine.NewQueue[net.SessionLeaveMessage](),
		sessionJoinMessages:     engine.NewQueue[net.SessionJoinMessage](),
		initialPositionMessages: engine.NewQueue[net.RemoteInitialPositionMessage](),
		startTime:               time.Now(),
	}
}

// hook remote events
func (s *RemoteSystem) Initialize(game *component.GameData, world donburi.World) {
	s.game = game
	if s.space == nil {
		physics, _ := archetype.MustFindPhysics(world)
		if physics == nil {
			return
		}
		s.space = physics.Space
	}
	if s.debug == nil {
		debug, ok := query.NewQuery(filter.Contains(component.Debug)).First(world)
		if !ok {
			return
		}

		s.debug = component.Debug.Get(debug)
	}
	s.game.Session.RemoteClient.RemoteUpdate.AddListener(func(ctx context.Context, rum net.RemoteUpdateMessage) {
		s.game.Session.PlayerMessage[rum.Msg.Source] = &component.RemotePlayerMessage{
			Position:  &rum.Msg.Position,
			Animation: rum.Msg.Animation,
			Source:    rum.Msg.Source,
			Vector:    rum.Msg.Point,
		}
	})
	s.game.Session.RemoteClient.SessionJoin.AddListener(func(ctx context.Context, sjm net.SessionJoinMessage) {
		s.sessionJoinMessages.Add(&sjm)
	})
	s.game.Session.RemoteClient.SessionLeave.AddListener(func(ctx context.Context, slm net.SessionLeaveMessage) {
		s.sessionLeaveMessages.Add(&slm)
	})
	s.game.Session.RemoteClient.SessionEnd.AddListener(func(ctx context.Context, val int) {
		s.shouldEnd = true
	})
	s.game.Session.RemoteClient.RemoteChat.AddListener(func(ctx context.Context, rcm net.RemoteChatMessage) {
		s.game.ChatMessages.Add(&rcm.Msg)
	})
	s.game.Session.RemoteClient.RemoteGameData.AddListener(func(ctx context.Context, rwm net.RemoteGameDataMessage) {
		s.game.Session.RemoteClient.GameData = rwm.Msg
		s.shouldPlaceWaste = true
		if s.game.Session.Type == component.SessionTypeJoin && !s.game.Session.RemoteClient.GameData.OnGameState {
			s.game.GameOver = true
		}
	})
	s.game.Session.RemoteClient.RemoteInitialPositionData.AddListener(func(ctx context.Context, ripd net.RemoteInitialPositionMessage) {
		s.initialPositionMessages.Add(&ripd)
	})

	if !s.game.Session.RemoteClient.Host && s.game.Session.JustJoined {
		s.game.Session.JustJoined = false
		if s.space == nil {
			physics, _ := archetype.MustFindPhysics(world)
			if physics == nil {
				return
			}
			s.space = physics.Space
		}
		for id, participant := range s.game.Session.RemoteClient.GameData.SessionParticipants {
			anim := engine.Ptr(component.DefaultPlayerAnimation)
			if participant.Anim != nil && *participant.Anim != "" {
				anim = participant.Anim
			}
			pos := participant.Position
			if pos != nil && !s.game.Session.RemoteClient.GameData.SessionParticipants[id].HasPlayer {
				archetype.NewPlayer(world, s.space, math.NewVec2(pos.X, pos.Y), *anim, *participant.Name, id, false)
				s.game.Session.RemoteClient.GameData.SessionParticipants[id].HasPlayer = true
			}
		}
	}

	if !s.game.Session.RemoteClient.Host && s.game.Session.RemoteClient != nil && s.game.Session.RemoteClient.GameData != nil && s.game.Session.RemoteClient.GameData.WasteLocations != nil {
		for _, loc := range s.game.Session.RemoteClient.GameData.WasteLocations {
			s.game.Session.RemoteClient.GameData.WasteLocations[loc.Id] = loc
			if !loc.Collected {
				archetype.PlaceRemoteWasteFromPath(world, s.space, s.debug, loc.Id, loc.Location, loc.Collected)
			}
		}

	}
}

func (s *RemoteSystem) Update(w donburi.World) {
	if s.game == nil {
		s.game = component.MustFindGame(w)
		if s.game == nil {
			return
		}
	}
	if s.space == nil {
		physics, _ := archetype.MustFindPhysics(w)
		if physics == nil {
			return
		}
		s.space = physics.Space
	}
	if s.debug == nil {
		debug, ok := query.NewQuery(filter.Contains(component.Debug)).First(w)
		if !ok {
			return
		}

		s.debug = component.Debug.Get(debug)
	}
	playerCount := 0
	s.query.Each(w, func(entry *donburi.Entry) {
		playerCount++
	})

	if s.shouldPlaceWaste && s.game.Session.RemoteClient.GameData.WasteLocations != nil {
		s.shouldPlaceWaste = false
		for _, loc := range s.game.Session.RemoteClient.GameData.WasteLocations {
			archetype.PlaceRemoteWasteFromPath(w, s.space, s.debug, loc.Id, loc.Location, loc.Collected)
		}
	}
	if s.shouldEnd {
		s.game.Session.End = true
		s.removePlayers(w)
	}
	for s.initialPositionMessages.Length() > 0 {
		ripd := s.initialPositionMessages.Remove()
		var targetPlayer *component.PlayerData
		s.query.Each(w, func(entry *donburi.Entry) {
			player := component.Player.Get(entry)
			if player == nil || player.Local {
				return
			}
			if player.ID == ripd.From {
				targetPlayer = player
			}
		})

		if targetPlayer == nil {
			if s.game.Session.RemoteClient.GameData.SessionParticipants[ripd.From] == nil {
				s.game.Session.RemoteClient.GameData.SessionParticipants[ripd.From] = &net.SessionParticipant{
					Id:        ripd.From,
					Name:      s.game.Session.RemoteClient.Participants[ripd.From],
					Position:  &ripd.Position,
					Anim:      engine.Ptr(component.DefaultPlayerAnimation),
					HasPlayer: false,
				}
			}
			userName := s.game.Session.RemoteClient.Participants[ripd.From]
			archetype.NewPlayer(w, s.space, math.NewVec2(ripd.Position.X, ripd.Position.Y), component.DefaultPlayerAnimation, *userName, ripd.From, false)

			s.game.Session.RemoteClient.GameData.SessionParticipants[ripd.From].HasPlayer = true
		}
	}
	for s.sessionJoinMessages.Length() > 0 {
		joinMessage := s.sessionJoinMessages.Remove()
		s.addPlayer(w, joinMessage)
	}
	for s.sessionLeaveMessages.Length() > 0 {
		slm := s.sessionLeaveMessages.Remove()
		s.removePlayer(w, slm)
	}
}

func (s *RemoteSystem) addPlayer(w donburi.World, sjm *net.SessionJoinMessage) {
	if *s.game.Session.RemoteClient.Client.Id == sjm.Target {
		return
	}
	participant := net.SessionParticipant{
		Id:       sjm.Target,
		Name:     sjm.Client.Participants[sjm.Target],
		Position: sjm.Position,
		Anim:     sjm.Anim,
	}
	if s.game.Session.RemoteClient.GameData.SessionParticipants[participant.Id] != nil {
		return
	}
	if s.game.Session.RemoteClient.GameData.SessionParticipants == nil {
		s.game.Session.RemoteClient.GameData.SessionParticipants = map[string]*net.SessionParticipant{}
	}
	s.game.Session.RemoteClient.GameData.SessionParticipants[participant.Id] = &participant

}

func (s *RemoteSystem) removePlayers(w donburi.World) {
	s.query.Each(w, func(entry *donburi.Entry) {
		player := component.Player.Get(entry)
		if player == nil || player.Local {
			return
		}
		playerBody := player.Body
		s.space.RemoveBody(playerBody)
		entity := entry.Entity()
		w.Remove(entity)
	})
}

func (s *RemoteSystem) removePlayer(w donburi.World, slm *net.SessionLeaveMessage) {
	var targetPlayer *donburi.Entity
	var targetEntry *donburi.Entry
	var targetEntryLabel *donburi.Entity
	s.query.Each(w, func(entry *donburi.Entry) {
		player := component.Player.Get(entry)
		if player == nil || player.ID == *s.game.Session.RemoteClient.Client.Id {
			return
		}
		targetEntryLabel = engine.Ptr(player.Label.Entity())
		entity := entry.Entity()
		targetPlayer = &entity
		targetEntry = entry
	})

	if targetPlayer == nil {
		return
	}
	playerBody := component.Player.Get(targetEntry).Body
	if s.game.Session.RemoteClient.Participants[*slm.Target] != nil {
		delete(s.game.Session.RemoteClient.Participants, *slm.Target)
	}
	if s.game.Session.RemoteClient.GameData.SessionParticipants[*slm.Target] != nil {
		delete(s.game.Session.RemoteClient.GameData.SessionParticipants, *slm.Target)
	}
	s.space.RemoveBody(playerBody)
	w.Remove(*targetPlayer)
	w.Remove(*targetEntryLabel)
}
