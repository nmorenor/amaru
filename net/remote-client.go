package net

import (
	"amaru/assets"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/maniartech/signals"
	"github.com/nmorenor/chezmoi-net/client"
)

const (
	ConnectionURL        = "wss://nmorenor.com/ws"
	AvailableSessionsURL = "https://nmorenor.com/hub-sessions"
)

func NewRemoteClient(currentClient *client.Client, userName string, hostMode bool) *RemoteClient {
	remoteClient := &RemoteClient{
		Client:                    currentClient,
		Participants:              nil,
		outmutex:                  &sync.Mutex{},
		inmutex:                   &sync.Mutex{},
		locationMutex:             &sync.Mutex{},
		Host:                      hostMode,
		Username:                  userName,
		RemoteUpdate:              signals.New[RemoteUpdateMessage](),
		SessionJoin:               signals.New[SessionJoinMessage](),
		SessionLeave:              signals.New[SessionLeaveMessage](),
		RemoteChat:                signals.New[RemoteChatMessage](),
		RemoteGameData:            signals.New[RemoteGameDataMessage](),
		RemoteInitialPositionData: signals.New[RemoteInitialPositionMessage](),
		SessionEnd:                signals.New[int](),
		ctx:                       context.Background(),
		GameData: &GameData{
			WasteLocations:      make(map[string]*WasteLocation),
			SessionParticipants: make(map[string]*SessionParticipant),
		},
	}
	remoteClient.Client.OnConnect = remoteClient.onReady
	remoteClient.Client.OnSessionChange = remoteClient.onSessionChange
	return remoteClient
}

type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Time.Format("2006-01-02T15:04:05.999Z07:00"))
	return []byte(formatted), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	tt, err := time.Parse("\"2006-01-02T15:04:05.999Z07:00\"", string(b))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

type WasteLocation struct {
	Id        string
	Location  assets.Path
	Collected bool
}

type SessionParticipant struct {
	Id        string
	Name      *string
	Position  *Point
	Anim      *string
	HasPlayer bool
	Score     int
}

type GameData struct {
	WasteLocations      map[string]*WasteLocation
	SessionParticipants map[string]*SessionParticipant
	LevelIndex          int
	Counter             int
	Frames              int
	OnGameState         bool
}
type Point struct {
	X float64
	Y float64
}
type ChatMessage struct {
	Source  string
	Message string
}
type GameDataMessage struct {
	Source   string
	GameData GameData
}

type Message struct {
	Source    string
	Position  Point
	Point     Point
	Animation string
}

type GetGameDataMessage struct {
	Id string
}

type GetGameDataResponse struct {
	GameData GameData
}

type PositionMessage struct {
	Id string
}

type PositionResponseMessage struct {
	Position Point
	Anim     string
}

type RemoteUpdateMessage struct {
	Client *RemoteClient
	From   *string
	Msg    Message
}

type RemoteChatMessage struct {
	Client *RemoteClient
	From   *string
	Msg    ChatMessage
}

type RemoteGameDataMessage struct {
	Client *RemoteClient
	From   *string
	Msg    *GameData
}

type SessionJoinMessage struct {
	Client   *RemoteClient
	Target   string
	Position *Point
	Anim     *string
}

type SessionLeaveMessage struct {
	Client *RemoteClient
	Target *string
}

type RemoteInitialPositionMessage struct {
	From     string
	Position Point
}

type RemoteClient struct {
	Host                      bool
	initialized               bool
	Client                    *client.Client
	Ready                     bool
	InvalidSession            bool
	InitAndReady              bool
	Participants              map[string]*string
	HostParticipant           *string
	outmutex                  *sync.Mutex
	inmutex                   *sync.Mutex
	locationMutex             *sync.Mutex
	Username                  string
	Session                   *string
	LocalPosition             *Point
	LocalAnimation            *string
	GameData                  *GameData
	RemoteUpdate              signals.Signal[RemoteUpdateMessage]
	SessionJoin               signals.Signal[SessionJoinMessage]
	SessionLeave              signals.Signal[SessionLeaveMessage]
	RemoteChat                signals.Signal[RemoteChatMessage]
	RemoteGameData            signals.Signal[RemoteGameDataMessage]
	RemoteInitialPositionData signals.Signal[RemoteInitialPositionMessage]
	SessionEnd                signals.Signal[int]
	ctx                       context.Context
}

// This will be called when web socket is connected
func (remoteClient *RemoteClient) onReady() {
	// Register this (RemoteClient) instance to receive rcp calls
	client.RegisterService(remoteClient, remoteClient.Client)

	if remoteClient.Host {
		remoteClient.Client.StartHosting(remoteClient.Username)
		//clipboard.Write(clipboard.FmtText, []byte(*remoteClient.Client.Session))
		fmt.Println("Session: " + *remoteClient.Client.Session)
	} else {
		response := remoteClient.Client.JoinSession(remoteClient.Username, *remoteClient.Session)
		if response == "" {
			remoteClient.InvalidSession = true
			return
		}
	}

	response := remoteClient.Client.SessionMembers()
	remoteClient.Participants = response.Members
	remoteClient.HostParticipant = &response.Host
	remoteClient.Ready = true
}

func (remoteClient *RemoteClient) Initialize() {
	if remoteClient.initialized {
		return
	}
	remoteClient.initialized = true
	if !remoteClient.Host {
		remoteClient.outmutex.Lock()
		defer remoteClient.outmutex.Unlock()
		for id := range remoteClient.Participants {
			if id != *remoteClient.Client.Id {
				rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)
				if !remoteClient.Host && remoteClient.HostParticipant != nil && *remoteClient.HostParticipant == id {
					sname := remoteClient.Client.GetServiceName(*remoteClient, "GetGameData", &id)
					var getGameDataResponse GetGameDataResponse
					if rpcClient != nil {
						msg := GetGameDataMessage{Id: id}
						rpcClient.Call(sname, msg, &getGameDataResponse)
						remoteClient.GameData = &getGameDataResponse.GameData
						remoteClient.RemoteGameData.Emit(remoteClient.ctx, RemoteGameDataMessage{
							Client: remoteClient,
							From:   &id,
							Msg:    &getGameDataResponse.GameData,
						})
					}
				}
				sname := remoteClient.Client.GetServiceName(*remoteClient, "GetPosition", &id)
				var position PositionResponseMessage
				if rpcClient != nil {
					msg := PositionMessage{Id: id}
					rpcClient.Call(sname, msg, &position)
					sessionParticipant := SessionParticipant{
						Id:        id,
						Name:      remoteClient.Participants[id],
						Position:  &position.Position,
						Anim:      &position.Anim,
						HasPlayer: false,
					}
					remoteClient.GameData.SessionParticipants[id] = &sessionParticipant
				}
			}
		}
	}
	remoteClient.InitAndReady = true
}

func (remoteClient *RemoteClient) SetLocalPosition(position *Point, anim *string) {
	remoteClient.locationMutex.Lock()
	defer remoteClient.locationMutex.Unlock()
	remoteClient.LocalPosition = position
	if anim != nil {
		remoteClient.LocalAnimation = anim
	}
}

func (remoteClient *RemoteClient) SendMessage(vector Point, position Point, animation string) {
	remoteClient.outmutex.Lock()
	defer remoteClient.outmutex.Unlock()
	rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)

	msg := &Message{
		Source:    *remoteClient.Client.Id,
		Point:     vector,
		Position:  position,
		Animation: animation,
	}

	sname := remoteClient.Client.GetServiceName(*remoteClient, "OnMessage", nil)
	if rpcClient != nil {
		var reply string
		rpcClient.Call(sname, msg, &reply)
	}
}

func (remoteClient *RemoteClient) SendChatMessage(message string) {
	remoteClient.outmutex.Lock()
	defer remoteClient.outmutex.Unlock()
	rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)

	msg := &ChatMessage{
		Source:  *remoteClient.Client.Id,
		Message: message,
	}

	sname := remoteClient.Client.GetServiceName(*remoteClient, "OnChatMessage", nil)
	if rpcClient != nil {
		var reply string
		rpcClient.Call(sname, msg, &reply)
	}
}

func (remoteClient *RemoteClient) SendGameDataMessage(gameData GameData) {
	remoteClient.outmutex.Lock()
	defer remoteClient.outmutex.Unlock()
	rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)

	msg := &GameDataMessage{
		Source:   *remoteClient.Client.Id,
		GameData: gameData,
	}

	sname := remoteClient.Client.GetServiceName(*remoteClient, "OnGetGameData", nil)
	if rpcClient != nil {
		var reply string
		rpcClient.Call(sname, msg, &reply)
	}
}

func (remoteClient *RemoteClient) SendInitialPositionDataMessage(position Point) {
	remoteClient.outmutex.Lock()
	defer remoteClient.outmutex.Unlock()
	rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)

	msg := &RemoteInitialPositionMessage{
		From:     *remoteClient.Client.Id,
		Position: position,
	}

	sname := remoteClient.Client.GetServiceName(*remoteClient, "OnNotifyInitialPosition", nil)
	if rpcClient != nil {
		var reply string
		rpcClient.Call(sname, msg, &reply)
	}
}

func (remoteClient *RemoteClient) RequestGameData() {
	remoteClient.outmutex.Lock()
	defer remoteClient.outmutex.Unlock()
	rpcClient := remoteClient.Client.GetRpcClientForService(*remoteClient)

	sname := remoteClient.Client.GetServiceName(*remoteClient, "GetGameData", remoteClient.HostParticipant)
	var getGameDataResponse GetGameDataResponse
	if rpcClient != nil {
		msg := GetGameDataMessage{Id: *remoteClient.Client.Id}
		rpcClient.Call(sname, msg, &getGameDataResponse)
		remoteClient.GameData = &getGameDataResponse.GameData
		remoteClient.RemoteGameData.Emit(remoteClient.ctx, RemoteGameDataMessage{
			Client: remoteClient,
			From:   remoteClient.HostParticipant,
			Msg:    &getGameDataResponse.GameData,
		})
	}
	for id := range remoteClient.Participants {
		if id != *remoteClient.Client.Id {
			if remoteClient.GameData != nil && remoteClient.GameData.SessionParticipants[id] == nil || !remoteClient.GameData.SessionParticipants[id].HasPlayer {
				sname = remoteClient.Client.GetServiceName(*remoteClient, "GetPosition", &id)
				var position PositionResponseMessage
				if rpcClient != nil {
					msg := PositionMessage{Id: id}
					rpcClient.Call(sname, msg, &position)
					remoteClient.RemoteInitialPositionData.Emit(remoteClient.ctx, RemoteInitialPositionMessage{
						From:     id,
						Position: position.Position,
					})
				}
			}
		}
	}
}

func (remoteClient *RemoteClient) FindParticipantFromName(target string) *string {
	for id, name := range remoteClient.Participants {
		if *name == target {
			return &id
		}
	}
	return nil
}

/**
 * Message received from rcp call, RPC methods must follow the signature
 */
func (remoteClient *RemoteClient) OnMessage(message *Message, reply *string) error {
	remoteClient.inmutex.Lock()
	defer remoteClient.inmutex.Unlock()
	if remoteClient.Participants[message.Source] != nil {
		remoteClient.RemoteUpdate.Emit(remoteClient.ctx, RemoteUpdateMessage{
			Client: remoteClient,
			From:   &message.Source,
			Msg:    *message,
		})
	}
	*reply = "OK"
	return nil
}

func (remoteClient *RemoteClient) OnChatMessage(message *ChatMessage, reply *string) error {
	remoteClient.inmutex.Lock()
	defer remoteClient.inmutex.Unlock()
	if remoteClient.Participants[message.Source] != nil {
		remoteClient.RemoteChat.Emit(remoteClient.ctx, RemoteChatMessage{
			Client: remoteClient,
			From:   &message.Source,
			Msg:    *message,
		})
	}
	*reply = "OK"
	return nil
}

func (remoteClient *RemoteClient) OnGetGameData(message *GameDataMessage, reply *string) error {
	remoteClient.inmutex.Lock()
	defer remoteClient.inmutex.Unlock()
	if remoteClient.Participants[message.Source] != nil {
		remoteClient.RemoteGameData.Emit(remoteClient.ctx, RemoteGameDataMessage{
			Client: remoteClient,
			From:   &message.Source,
			Msg:    &message.GameData,
		})
	}
	*reply = "OK"
	return nil
}

func (remoteClient *RemoteClient) OnNotifyInitialPosition(message *RemoteInitialPositionMessage, reply *string) error {
	remoteClient.inmutex.Lock()
	defer remoteClient.inmutex.Unlock()
	if remoteClient.Participants[message.From] != nil {
		remoteClient.RemoteInitialPositionData.Emit(remoteClient.ctx, *message)
	}
	*reply = "OK"
	return nil
}

func (remoteClient *RemoteClient) GetPosition(message *PositionMessage, reply *PositionResponseMessage) error {
	remoteClient.locationMutex.Lock()
	defer remoteClient.locationMutex.Unlock()
	if remoteClient.LocalPosition != nil && remoteClient.LocalAnimation != nil {
		*reply = PositionResponseMessage{Position: Point{X: remoteClient.LocalPosition.X, Y: remoteClient.LocalPosition.Y}, Anim: *remoteClient.LocalAnimation}
	}
	return nil
}

func (remoteClient *RemoteClient) GetGameData(message *GetGameDataMessage, reply *GetGameDataResponse) error {
	if remoteClient.GameData != nil {
		*reply = GetGameDataResponse{GameData: *remoteClient.GameData}
	}
	return nil
}

func (remoteClient *RemoteClient) onSessionChange(event client.SessionChangeEvent) {
	remoteClient.inmutex.Lock()
	defer remoteClient.inmutex.Unlock()
	response := remoteClient.Client.SessionMembers()
	oldParticipants := remoteClient.Participants
	remoteClient.Participants = response.Members
	if event.EventType == client.SESSION_JOIN && remoteClient.Participants[event.EventSource] != nil {
		remoteClient.SessionJoin.Emit(remoteClient.ctx, SessionJoinMessage{
			Client:   remoteClient,
			Target:   event.EventSource,
			Position: nil,
			Anim:     nil,
		})
	}
	if event.EventType == client.SESSION_LEAVE && oldParticipants[event.EventSource] != nil {
		remoteClient.SessionLeave.Emit(remoteClient.ctx, SessionLeaveMessage{
			Client: remoteClient,
			Target: &event.EventSource,
		})
	}
	if event.EventType == client.SESSION_END {
		remoteClient.ResetListeners()
	}
}

func (remoteClient *RemoteClient) ResetListeners() {
	remoteClient.RemoteUpdate.Reset()
	remoteClient.SessionJoin.Reset()
	remoteClient.SessionLeave.Reset()
	remoteClient.SessionEnd.Reset()
	remoteClient.RemoteChat.Reset()
	remoteClient.RemoteGameData.Reset()
	remoteClient.RemoteInitialPositionData.Reset()
}
