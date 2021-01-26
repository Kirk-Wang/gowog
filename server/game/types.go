package game

import (
	"gowog-cloud/game/ws"
)

type Game interface {
	ProcessInput(message []byte)
	NewPlayerConnect(client ws.Client)
	RemovePlayer(playerID int32, clientID int32)
	GetQuitChannel() chan bool
	Update()
}
