package game

import (
	"jungle-royale/message"
	"jungle-royale/network"
	"jungle-royale/object"
	"jungle-royale/state"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	MaxClientCount        = 100
	CalcLoopInterval      = 16
	BroadCastLoopInterval = 16
)

type Game struct {
	state  *state.State
	socket *network.Socket
}

func NewGame(socket *network.Socket) *Game {
	game := &Game{
		state.NewState(),
		socket,
	}
	go game.CalcLoop()      // start main loop
	go game.BroadcastLoop() // broadcast to client
	return game
}

func (game *Game) OnClient(clientId string) {
	game.SetPlayer(clientId)
}

func (game *Game) OnMessage(data []byte, id string) {
	game.HandleMessage(id, data)
}

func (game *Game) SetPlayer(clientId string) {
	game.state.AddPlayer(clientId)

	// send GameInit message
	gameInit := &message.GameInit{Id: clientId}
	data, err := proto.Marshal(&message.Wrapper{
		MessageType: &message.Wrapper_GameInit{
			GameInit: gameInit,
		},
	})
	if err != nil {
		log.Printf("Failed to marshal GameInit: %v", err)
		return
	}

	if err := (*game.socket).Send(data, clientId); err != nil {
		log.Printf("Failed to send GameInit message to client %s: %v", clientId, err)
		return
	}

	log.Printf("보낸겨: %s", gameInit.String())
}

func (game *Game) HandleMessage(clientId string, data []byte) {
	var wrapper message.Wrapper
	if err := proto.Unmarshal(data, &wrapper); err != nil {
		log.Printf("Failed to unmarshal message from client %s: %v", clientId, err)
		return
	}

	// dirChange message
	if dirChange := wrapper.GetDirChange(); dirChange != nil {
		if value, exists := game.state.Players.Load(clientId); exists {
			player := value.(*object.Player)
			go player.DirChange(float64(dirChange.GetAngle()), dirChange.IsMoved)
		}
	}

	// bulletCreate message
	if bulletCreate := wrapper.GetBulletCreate(); bulletCreate != nil {
		game.state.AddBullet(bulletCreate)
	}
}

func (game *Game) CalcLoop() {
	ticker := time.NewTicker(CalcLoopInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // calculation loop
		game.state.CalcState()
	}
}

func (game *Game) BroadcastLoop() {
	ticker := time.NewTicker(BroadCastLoopInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // broadcast loop
		playerList := make([]*message.Player, 0)
		game.state.Players.Range(func(key, value any) bool {
			player := value.(*object.Player)
			playerList = append(playerList, player.MakeSendingData())
			return true
		})

		bulletList := make([]*message.BulletState, 0)
		game.state.Bullets.Range(func(key, value any) bool {
			bullet := value.(*object.Bullet)
			bulletList = append(bulletList, bullet.MakeSendingData())
			return true
		})

		gameState := &message.GameState{
			Players:     playerList,
			BulletState: bulletList,
		}

		data, err := proto.Marshal(&message.Wrapper{
			MessageType: &message.Wrapper_State{
				State: gameState,
			},
		})
		if err != nil {
			log.Printf("Failed to marshal GameState: %v", err)
			return
		}

		(*game.socket).Broadcast(data)
	}
}
