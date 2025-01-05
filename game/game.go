package game

import (
	"jungle-royale/calculator"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/state"
	"jungle-royale/util"
	"log"
	"math"
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
)

type Game struct {
	minPlayerNum   int
	playingTime    int
	playerNum      int
	state          *state.State
	calculator     *calculator.Calculator
	clients        *util.Map[ClientId, *Client]
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
		clients:        util.NewSyncMap[ClientId, *Client](),
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

	// map setting
	game.state.ConfigureState(length, game.playingTime)
	game.calculator.ConfigureCalculator(length)

	// player relocation
	game.state.Players.Range(func(key string, player *object.Player) bool {
		x := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		game.calculator.ReLocation(player, x, y)
		return true
	})

	// healpack setting
	for i := 0; i < length*length; i++ {
		x := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newHealPack := object.NewHealPack(x, y)
		game.calculator.SetLocation(newHealPack, x, y)
		game.state.HealPacks.Store(newHealPack.Id, newHealPack)
	}

	// magic item setting
	for i := 0; i < length*length; i++ {
		x := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newStoneItem := object.NewMagicItem(object.STONE_MAGIC, x, y)
		game.calculator.SetLocation(newStoneItem, x, y)
		game.state.MagicItems.Store(newStoneItem.ItemId, newStoneItem)
		x = float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y = float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newFireItem := object.NewMagicItem(object.FIRE_MAGIC, x, y)
		game.calculator.SetLocation(newFireItem, x, y)
		game.state.MagicItems.Store(newFireItem.ItemId, newFireItem)
	}

	// tile fall setting
	last_x_idx := rand.Intn(length)
	last_y_idx := rand.Intn(length)
	game.state.Tiles[last_x_idx][last_y_idx].ParentTile = game.state.Tiles[last_x_idx][last_y_idx]
	tileTreeSet := util.NewSet[calculator.ChunkIndex]()
	tileTreeSet.Add(calculator.ChunkIndex{X: last_x_idx, Y: last_y_idx})
	dir := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for tileTreeSet.Length() > 0 {
		currentTileIdx, _ := tileTreeSet.SelectRandom(func(t calculator.ChunkIndex) bool { return true })
		currentTile := game.state.Tiles[currentTileIdx.X][currentTileIdx.Y]
		tileTreeSet.Remove(currentTileIdx)
		hasChild := false
		for i := 0; i < 4; i++ {
			if 0 <= currentTileIdx.X+dir[i][0] &&
				currentTileIdx.X+dir[i][0] < length &&
				0 <= currentTileIdx.Y+dir[i][1] &&
				currentTileIdx.Y+dir[i][1] < length &&
				game.state.Tiles[currentTileIdx.X+dir[i][0]][currentTileIdx.Y+dir[i][1]].ParentTile == nil {
				childTile := game.state.Tiles[currentTileIdx.X+dir[i][0]][currentTileIdx.Y+dir[i][1]]
				tileTreeSet.Add(calculator.ChunkIndex{
					X: currentTileIdx.X + dir[i][0],
					Y: currentTileIdx.Y + dir[i][1],
				})
				currentTile.ChildTile.Add(childTile)
				childTile.ParentTile = currentTile
				hasChild = true
			}
		}
		if !hasChild {
			game.calculator.LeafTileSet.Add(currentTile)
		}
	}

	game.state.GameState = state.Playing

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
		game.clients.Store(client.ID, client)
	}
}

func (game *Game) SetPlayer(client *Client) {

	clientId := string(client.ID)

	x := float64(rand.Intn(int(game.state.MaxCoord)))
	y := float64(rand.Intn(int(game.state.MaxCoord)))

	game.state.AddPlayer(clientId, x, y)

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
				mapLength := int(math.Sqrt(float64(game.playerNum)))
				if mapLength < 2 {
					mapLength = 2
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

		tileStateList := make([]*message.TileState, 0)
		for i := 0; i < game.state.ChunkNum; i++ {
			for j := 0; j < game.state.ChunkNum; j++ {
				if game.state.Tiles[i][j].TileState != object.TILE_FALL {
					tileStateList = append(tileStateList, game.state.Tiles[i][j].MakeSendingData())
				}
			}
		}

		gameState := &message.GameState{
			PlayerState:    playerList,
			BulletState:    bulletList,
			HealPackState:  healPackList,
			MagicItemState: magicItemList,
			TileState:      tileStateList,
			ChangingState:  game.state.ChangingState.MakeSendingData(),
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

		game.broadcast(data)
	}
}

// Room Interface
func (game *Game) OnMessage(data []byte, id string) {
	game.handleMessage(id, data)
}

func (game *Game) OnClose(client *Client) {
	game.clients.Delete(client.ID)
	if game.clients.Length() == 0 {
		log.Print("clinets count zero")
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
	if changeBulletState := wrapper.GetChangeBulletState(); changeBulletState != nil {
		game.state.ChangeBulletState(clientId, changeBulletState)
	}
}

func (game *Game) broadcast(data []byte) {
	game.clients.Range(func(id ClientId, client *Client) bool {
		if err := client.write(data); err != nil {
			log.Printf("Failed to send message to client %s: %v", id, err)
			client.close()
			game.clients.Delete(id)
		}
		return true
	})
}
