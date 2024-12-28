package game

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/network"
	"jungle-royale/object"
	"jungle-royale/state"
	"log"
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	waiting = iota
	counting
	playing
)

type Game struct {
	gameState    int
	minPlayerNum int
	playerNum    int
	state        *state.State
	socket       *network.Socket
}

func NewGame(socket *network.Socket, minPlayerNum int) *Game {
	game := &Game{
		minPlayerNum: minPlayerNum,
		playerNum:    0,
		state:        state.NewState(),
		socket:       socket,
	}

	return game
}

func (game *Game) SetReadyStatus() *Game {
	game.gameState = waiting
	game.state.SetState(cons.WAITING_MAP_CHUNK_NUM)
	return game
}

func (game *Game) SetPlayingStatus() *Game {
	game.gameState = playing
	game.state.SetState(game.playerNum * game.playerNum)
	game.state.MoverList.GetPlayers().Range(func(key, value any) bool {
		player := value.(*object.Player)
		x := float32(rand.Intn(int(game.state.MaxCoord)))
		y := float32(rand.Intn(int(game.state.MaxCoord)))
		player.SetLocation(x, y)
		return true
	})
	return game
}

func (game *Game) StartGame() *Game {
	go game.CalcGameTickLoop() // start main loop
	go game.BroadcastLoop()    // broadcast to client
	go game.CalcSecLoop()
	return game
}

func (game *Game) OnClient(clientId string) {
	if game.gameState == waiting {
		game.playerNum++
		game.SetPlayer(clientId)
	}
}

func (game *Game) SetPlayer(clientId string) {

	x := rand.Intn(int(game.state.MaxCoord))
	y := rand.Intn(int(game.state.MaxCoord))

	game.state.AddPlayer(clientId, float32(x), float32(y))

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

func (game *Game) CalcGameTickLoop() {
	ticker := time.NewTicker(cons.CalcLoopInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // calculation loop
		game.state.CalcGameTickState()
	}
}

func (game *Game) CalcSecLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	gameStartCount := 10
	for range ticker.C {
		if game.gameState != playing &&
			game.playerNum >= game.minPlayerNum &&
			gameStartCount >= 0 {
			Count := &message.GameCount{
				Count: int32(gameStartCount),
			}
			data, err := proto.Marshal(&message.Wrapper{
				MessageType: &message.Wrapper_GameCount{
					GameCount: Count,
				},
			})
			if err != nil {
				log.Printf("Failed to marshal GameState: %v", err)
				return
			}
			log.Printf("game play in %d", gameStartCount)
			(*game.socket).Broadcast(data)
			gameStartCount--
			if gameStartCount == -1 {
				game.SetPlayingStatus()
			}
		}
		game.state.SecLoop()
	}
}

func (game *Game) BroadcastLoop() {
	ticker := time.NewTicker(cons.BroadCastLoopInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // broadcast loop
		playerList := make([]*message.PlayerState, 0)
		game.state.MoverList.GetPlayers().Range(func(key, value any) bool {
			player := value.(*object.Player)
			playerList = append(playerList, player.MakeSendingData())
			return true
		})

		bulletList := make([]*message.BulletState, 0)
		game.state.MoverList.GetBullets().Range(func(key, value any) bool {
			bullet := value.(*object.Bullet)
			bulletList = append(bulletList, bullet.MakeSendingData())
			return true
		})

		gameState := &message.GameState{
			PlayerState: playerList,
			BulletState: bulletList,
		}

		data, err := proto.Marshal(&message.Wrapper{
			MessageType: &message.Wrapper_GameState{
				GameState: gameState,
			},
		})
		if err != nil {
			log.Printf("Failed to marshal GameState: %v", err)
			return
		}

		(*game.socket).Broadcast(data)
	}
}
