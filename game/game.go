package game

import (
	"jungle-royale/calculator"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/state"
	"jungle-royale/statistic"
	"jungle-royale/util"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type Game struct {
	minPlayerNum      int
	playingTime       int
	playerNum         int
	state             *state.State
	calculator        *calculator.Calculator
	clients           *util.Map[ClientId, *Client]
	serverClientTable *util.Map[string, ClientId]
	alertGameStart    func() // 게임 시작을 알림
	alertGameEnd      func() // 게임 종료를 알림
	alertPlayerLeavae func(client *Client)
	gameLogger        *statistic.Logger
	debug             bool
	endTickCountMu    sync.Mutex
	endTickCount      int
}

// playing time - second
func NewGame(debug bool) *Game {
	// playing time - sec
	gameState := state.NewState()
	logger := statistic.NewGameLogger()
	game := &Game{
		playerNum:         0,
		state:             gameState,
		calculator:        nil,
		clients:           util.NewSyncMap[ClientId, *Client](),
		serverClientTable: util.NewSyncMap[string, ClientId](),
		gameLogger:        logger,
		debug:             debug,
		endTickCountMu:    sync.Mutex{},
		endTickCount:      0,
	}
	game.calculator = calculator.NewCalculator(
		gameState,
		func(clientId string, rank, kill int) {
			game.MakeLog(clientId, rank, kill)
		},
	)
	return game
}

func (game *Game) SetGame(
	minPlayerNum int,
	playingTime int,
	startHandler func(),
	leaveHandler func(client *Client),
	endHandler func(),
) *Game {
	game.minPlayerNum = minPlayerNum
	game.playingTime = playingTime
	game.alertGameStart = startHandler
	game.alertPlayerLeavae = leaveHandler
	game.alertGameEnd = endHandler
	return game
}

func (game *Game) MakeLog(clientId string, rank, kill int) {
	if ci, ok := game.clients.Get(ClientId(clientId)); ok {
		game.gameLogger.AddLog((*ci).serverClientId, rank, kill)
	}
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
	game.state.Tiles[last_x_idx][last_y_idx].Depth = 0
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
				childTile.Depth = currentTile.Depth + 1
				hasChild = true
			}
		}
		if !hasChild {
			game.calculator.LeafTileSet.Push(currentTile)
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

func (game *Game) CalcGameTickLoop() {
	ticker := time.NewTicker(cons.CalcLoopInterval * time.Millisecond)
	defer ticker.Stop()

	// currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	for range ticker.C { // calculation loop
		// tempTime := time.Now().UnixNano() / int64(time.Millisecond)
		// log.Printf("%d\n", tempTime-currentTime)
		// currentTime = tempTime
		game.calculator.CalcGameTickState()
		game.PlusEndCount()
		game.CheckEndGame()
	}
}

func (game *Game) CalcSecLoop() {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	gameStartCount := 15
	if game.debug {
		gameStartCount = 3
	}
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
					MapLength:      int32(mapLength * cons.CHUNK_LENGTH),
					TotalPlayerNum: int32(game.clients.Length()),
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

		game.state.ConfigMu.Lock()

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
			LastSec:        int32((game.state.LastGameTick * 16) / 1000),
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

		game.state.ConfigMu.Unlock()

		game.broadcast(data)
	}
}

func (game *Game) OnMessage(data []byte, id string) {
	game.handleMessage(id, data)
}

func (game *Game) OnClient(client *Client) {
	if game.state.GameState == state.Waiting {
		game.playerNum++
		game.SetPlayer(client)
		game.clients.Store(client.ID, client)
		game.serverClientTable.Store(client.serverClientId, client.ID)
	} else {
		_, exists := game.serverClientTable.Get(client.serverClientId)
		if exists {
			log.Printf("reconnection: %s", client.serverClientId)
			game.clients.Store(client.ID, client)
		} else {
			log.Printf("fail to reconnect: %s", client.serverClientId)
			client.close()
		}
	}
	if game.clients.Length() > 0 {
		game.ResetEndCount()
	}
}

func (game *Game) OnClose(client *Client) {
	game.clients.Delete(client.ID)
	if game.state.GameState == state.Waiting {
		game.alertPlayerLeavae(client)
	}

	if game.clients.Length() == 0 {
		log.Println("clinet's count zero")
		game.PlusEndCount()
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

func (game *Game) SetPlayer(client *Client) {

	clientId := string(client.ID)

	x := rand.Intn(int(game.state.MaxCoord))
	y := rand.Intn(int(game.state.MaxCoord))

	game.state.AddPlayer(clientId, float64(x), float64(y))

	// send GameInit message
	gameInit := &message.GameInit{
		Id:           clientId,
		MinPlayerNum: int32(game.minPlayerNum),
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

func (game *Game) IsPlaying() bool {
	if game.state.GameState == state.Empty {
		return false
	} else {
		return true
	}
}

func (game *Game) CheckEndGame() {
	// client가 들어올 때, close 될 때, tick 체크할 때 -> 동시 연산
	game.endTickCountMu.Lock()
	defer game.endTickCountMu.Unlock()

	var threeMinutesTick = 60 * 60 * 3
	if game.endTickCount > threeMinutesTick {
		game.alertGameEnd()
	}
}

// client가 zero가 되거나, game이 끝나면 plus 시작
func (game *Game) PlusEndCount() {
	game.endTickCountMu.Lock()
	defer game.endTickCountMu.Unlock()

	// 이미 게임이 끝났다면, 무조건 tick 증가
	if game.IsEndState() {
		game.endTickCount += 1
		return
	}

	if game.clients.Length() == 0 {
		game.endTickCount += 1
		return
	}
}

// client가 들어오면 end count ++
func (game *Game) ResetEndCount() {
	game.endTickCountMu.Lock()
	defer game.endTickCountMu.Unlock()

	// 이미 게임이 끝났다면, zero로 초기화 하지 않는다.
	if game.IsEndState() {
		return
	}

	game.endTickCount = 0
}

func (game *Game) IsEndState() bool {
	return game.state.GameState == state.End
}
