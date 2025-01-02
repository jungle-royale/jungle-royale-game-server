package game

import (
	"jungle-royale/calculator"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/state"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type Game struct {
	minPlayerNum   int
	playingTime    int
	playerNum      int
	state          *state.State
	calculator     *calculator.Calculator
	clients        map[ClientId]*Client
	clientsMu      sync.Mutex
	alertGameStart func() // 게임 시작을 알림
	alertGameEnd   func() // 게임 종료를 알림
}

// playing time - second
func NewGame(
	minPlayerNum int,
	playingTime int,
	startHandler func(),
	endHandler func(),
) *Game {
	// playing time - sec
	gameState := state.NewState()
	game := &Game{
		minPlayerNum:   minPlayerNum,
		playingTime:    playingTime,
		playerNum:      0,
		state:          gameState,
		calculator:     calculator.NewCalculator(gameState),
		clients:        make(map[ClientId]*Client),
		clientsMu:      sync.Mutex{},
		alertGameStart: startHandler,
		alertGameEnd:   endHandler,
	}
	return game
}

func (game *Game) SetReadyStatus() *Game {
	game.state.ConfigureState(cons.WAITING_MAP_CHUNK_NUM, int(math.MaxInt))
	game.state.GameState = state.Waiting
	game.calculator.ConfigureCalculator(cons.WAITING_MAP_CHUNK_NUM)
	return game
}

func (game *Game) SetPlayingStatus(length int) *Game {
	game.state.GameState = state.Playing

	// map setting
	game.state.ConfigureState(length, game.playingTime)
	game.calculator.ConfigureCalculator(length)
	// game.state.ConfigureState(4)
	// game.calculator.ConfigureCalculator(4)

	// player relocation
	game.state.Players.Range(func(key string, player *object.Player) bool {
		x := float32(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float32(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		player.SetLocation(x, y)
		return true
	})

	// healpack setting
	for i := 0; i < length*length; i++ {
		x := float32(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float32(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newHealPack := object.NewHealPack(x, y)
		game.state.HealPacks.Store(newHealPack.Id, newHealPack)
	}

	// magic item setting
	for i := 0; i < length*length; i++ {
		x := float32(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float32(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newStoneItem := object.NewMagicItem(object.STONE_MAGIC, x, y)
		newFireItem := object.NewMagicItem(object.FIRE_MAGIC, x, y)
		game.state.MagicItems.Store(newStoneItem.ItemId, newStoneItem)
		game.state.MagicItems.Store(newFireItem.ItemId, newFireItem)
	}
	return game
}

func (game *Game) StartGame() *Game {
	go game.CalcGameTickLoop() // start main loop
	go game.BroadcastLoop()    // broadcast to client
	go game.CalcSecLoop()
	return game
}

// Room Interface
func (game *Game) OnClient(client *Client) {
	if game.state.GameState == state.Waiting {
		game.playerNum++
		game.SetPlayer(client)
		game.clientsMu.Lock()
		game.clients[client.ID] = client
		game.clientsMu.Unlock()
	}
}

func (game *Game) SetPlayer(client *Client) {

	clientId := string(client.ID)

	x := rand.Intn(int(game.state.MaxCoord))
	y := rand.Intn(int(game.state.MaxCoord))

	game.state.AddPlayer(clientId, float32(x), float32(y))

	// send GameInit message
	gameInit := &message.GameInit{
		Id: clientId,
	}
	data, err := proto.Marshal(&message.Wrapper{
		MessageType: &message.Wrapper_GameInit{
			GameInit: gameInit,
		},
	})
	if err != nil {
		log.Printf("Failed to marshal GameInit: %v", err)
		return
	}

	if err := client.write(data); err != nil {
		log.Printf("Failed to send GameInit message to client %s: %v", clientId, err)
		return
	}

	// log.Printf("보낸겨: %s", gameInit.String())
}

func (game *Game) CalcGameTickLoop() {
	ticker := time.NewTicker(cons.CalcLoopInterval * time.Millisecond)
	defer ticker.Stop()

	// currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	for range ticker.C { // calculation loop
		// tempTime := time.Now().UnixNano() / int64(time.Millisecond)
		// log.Printf("%d\n", tempTime-currentTime)
		// currentTime = tempTime
		game.calculator.CalcGameTickState()
	}
}

func (game *Game) CalcSecLoop() {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	gameStartCount := 3
	for range ticker.C {
		if game.state.GameState != state.Playing &&
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

			game.broadcast(data)

			gameStartCount--
			if gameStartCount == 0 {
				mapLength := game.playerNum / 10
				if mapLength == 0 {
					mapLength = 1
				}
				start := &message.GameStart{
					MapLength: int32(mapLength * cons.CHUNK_LENGTH),
				}
				gameStart, err := proto.Marshal(&message.Wrapper{
					MessageType: &message.Wrapper_GameStart{
						GameStart: start,
					},
				})
				if err != nil {
					log.Printf("Failed to marshal GameState: %v", err)
					return
				}

				game.alertGameStart()

				game.broadcast(gameStart)

				game.SetPlayingStatus(mapLength)
				log.Println("game start")
			}
		}
		game.calculator.SecLoop()
	}
}

func (game *Game) BroadcastLoop() {
	ticker := time.NewTicker(cons.BroadCastLoopInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // broadcast loop
		playerList := make([]*message.PlayerState, 0)
		game.state.Players.Range(func(key string, player *object.Player) bool {
			playerList = append(playerList, player.MakeSendingData())
			return true
		})

		bulletList := make([]*message.BulletState, 0)
		game.state.Bullets.Range(func(key string, bullet *object.Bullet) bool {
			bulletList = append(bulletList, bullet.MakeSendingData())
			return true
		})

		healPackList := make([]*message.HealPackState, 0)
		game.state.HealPacks.Range(func(key string, healPack *object.HealPack) bool {
			healPackList = append(healPackList, healPack.MakeSendingData())
			return true
		})

		magicItemList := make([]*message.MagicItemState, 0)
		game.state.MagicItems.Range(func(key string, magicItem *object.Magic) bool {
			magicItemList = append(magicItemList, magicItem.MakeSendingData())
			return true
		})

		playerDeadList := make([]*message.PlayerDeadState, 0)
		game.state.PlayerDead.Range(func(key string, status *object.PlayerDead) bool {
			playerDeadList = append(playerDeadList, status.MakeSendingData())
			game.state.PlayerDead.Delete(key)
			return true
		})

		tileStateList := make([]*message.TileState, 0)
		tileState := game.state.Tiles.ValueList()
		for _, tile := range tileState {
			tileStateList = append(tileStateList, tile.MakeSendingData())
		}

		gameState := &message.GameState{
			PlayerState:     playerList,
			BulletState:     bulletList,
			HealPackState:   healPackList,
			MagicItemState:  magicItemList,
			PlayerDeadState: playerDeadList,
			TileState:       tileStateList,
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
		// log.Print("healpack: ")
		// log.Println(healPackList)
		game.broadcast(data)
	}
}

// Room Interface
func (game *Game) OnMessage(data []byte, id string) {
	game.handleMessage(id, data)
}

func (game *Game) OnClose(client *Client) {
	game.clientsMu.Lock()
	delete(game.clients, client.ID)
	game.clientsMu.Unlock()
	if len(game.clients) == 0 {
		game.alertGameEnd()
	}
}

func (game *Game) handleMessage(clientId string, data []byte) {
	var wrapper message.Wrapper
	if err := proto.Unmarshal(data, &wrapper); err != nil {
		log.Printf("Failed to unmarshal message from client %s: %v", clientId, err)
		return
	}

	if changeAngle := wrapper.GetChangeAngle(); changeAngle != nil {
		game.state.ChangeAngle(clientId, changeAngle)
	}

	// dirChange message
	if changeDir := wrapper.GetChangeDir(); changeDir != nil {
		game.state.ChangeDirection(clientId, changeDir)
	}

	if doDash := wrapper.GetDoDash(); doDash != nil {
		game.state.DoDash(clientId, doDash)
	}

	// bulletCreate message
	if createBullet := wrapper.GetCreateBullet(); createBullet != nil {
		game.state.CreateBullet(clientId, createBullet)
	}
}

func (game *Game) broadcast(data []byte) {
	game.clientsMu.Lock()
	for id, client := range game.clients {
		if err := client.write(data); err != nil {
			log.Printf("Failed to send message to client %s: %v", id, err)
			client.close()
			delete(game.clients, id)
		}
	}
	game.clientsMu.Unlock()
}
