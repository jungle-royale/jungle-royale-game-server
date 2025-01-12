package game

import (
	"fmt"
	"jungle-royale/calculator"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/serverlog"
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

const (
	END_GAME_MAX_TICK_COUNT = 60 * 60 * 1 // 60초(1분)
)

type LoopState struct {
	lastCalcLoopCheck      int
	calcLoopCheck          int
	lastBroadcastLoopCheck int
	broadcastLoopCheck     int
}

type ClientIdAllocator struct {
	ClientId   int
	ClientIdMu sync.Mutex
}

func NewClientIdAllocator() *ClientIdAllocator {
	return &ClientIdAllocator{
		0,
		sync.Mutex{},
	}
}

func (cia *ClientIdAllocator) AllocateClientId() int {
	cia.ClientIdMu.Lock()
	id := cia.ClientId
	cia.ClientId++
	cia.ClientIdMu.Unlock()
	return id
}

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
	gameRecorder      *statistic.Recorder
	debug             bool
	endTickCountMu    sync.Mutex
	endTickCount      int
	gameLogger        *serverlog.GameLogger
	loopState         LoopState
	loopWaitGroup     sync.WaitGroup
	ClientIdAllocator *ClientIdAllocator
	ObjectIdAllocator *object.ObjectIdAllocator
	broadcastDataSize int
}

// playing time - second
func NewGame(debug bool) *Game {
	objectIdAllocator := object.NewObjectIdAllocator()
	gameState := state.NewState(objectIdAllocator)
	gameRecorder := statistic.NewGameLogger()
	game := &Game{
		playerNum:         0,
		state:             gameState,
		calculator:        nil,
		clients:           util.NewSyncMap[ClientId, *Client](),
		serverClientTable: util.NewSyncMap[string, ClientId](),
		gameRecorder:      gameRecorder,
		debug:             debug,
		endTickCountMu:    sync.Mutex{},
		endTickCount:      0,
		loopState:         LoopState{0, 0, 0, 0},
		loopWaitGroup:     sync.WaitGroup{},
		ClientIdAllocator: NewClientIdAllocator(),
		ObjectIdAllocator: objectIdAllocator,
		broadcastDataSize: 0,
	}
	game.calculator = calculator.NewCalculator(
		gameState,
		func(clientId, rank, kill int) {
			game.MakeLog(ClientId(clientId), rank, kill)
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
	gameLogger *serverlog.GameLogger,
) *Game {
	game.minPlayerNum = minPlayerNum
	game.playingTime = playingTime
	game.alertGameStart = startHandler
	game.alertPlayerLeavae = leaveHandler
	game.alertGameEnd = endHandler
	game.gameLogger = gameLogger
	gameLogger.Log("Set Game")
	return game
}

func (game *Game) MakeLog(clientId ClientId, rank, kill int) {
	if ci, ok := game.clients.Get(clientId); ok {
		game.gameRecorder.AddRecord((*ci).serverClientId, rank, kill)
	}
}

func (game *Game) SetReadyStatus() *Game {
	game.state.ConfigureState(cons.WAITING_MAP_CHUNK_NUM, int(math.MaxInt))
	game.state.GameState = state.Waiting
	game.calculator.ConfigureCalculator(cons.WAITING_MAP_CHUNK_NUM)
	game.gameLogger.Log("Set Ready Status")
	return game
}

func (game *Game) SetPlayingStatus(length int) *Game {

	// map setting
	game.state.ConfigureState(length, game.playingTime)
	game.calculator.ConfigureCalculator(length)

	// player relocation
	game.state.Players.Range(func(key int, player *object.Player) bool {
		x := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		game.calculator.ReLocation(player, x, y)
		return true
	})

	// healpack setting
	for i := 0; i < length*length; i++ {
		x := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newHealPack := object.NewHealPack(x, y, game.ObjectIdAllocator.AllocateHealPackId())
		game.calculator.SetLocation(newHealPack, x, y)
		game.state.HealPacks.Store(newHealPack.Id, newHealPack)
	}

	// magic item setting
	for i := 0; i < length*length; i++ {
		x := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y := float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newStoneItem := object.NewMagicItem(object.STONE_MAGIC, x, y, game.ObjectIdAllocator.AllocateMagicId())
		game.calculator.SetLocation(newStoneItem, x, y)
		game.state.MagicItems.Store(newStoneItem.ItemId, newStoneItem)
		x = float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		y = float64(rand.Intn(int(game.state.MaxCoord-1))) + 0.5
		newFireItem := object.NewMagicItem(object.FIRE_MAGIC, x, y, game.ObjectIdAllocator.AllocateMagicId())
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

	game.gameLogger.Log("Set Playing Status")

	return game
}

func (game *Game) StartGame() *Game {

	game.gameLogger.Log("init game")

	game.loopWaitGroup.Add(4)

	go game.CalcGameTickLoop() // start main loop
	go game.BroadcastLoop()    // broadcast to client
	go game.CalcSecLoop()

	go game.GameStateCheckLoop()

	game.loopWaitGroup.Wait()

	game.gameLogger.Log("Game End")

	game.clients.Range(func(ci ClientId, c *Client) bool {
		c.close()
		return true
	})

	go game.alertGameEnd()

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
		func() {
			defer func() {
				if r := recover(); r != nil {
					game.gameLogger.Log("Recover CalcLoop Panic: " + fmt.Sprintf("%v", r))
				}
			}()

			game.state.ConfigMu.Lock()
			game.calculator.CalcGameTickState()
			game.state.ConfigMu.Unlock()
			if !game.debug {
				game.PlusEndCount()
				game.CheckEndGame()
			}
		}()
		game.loopState.calcLoopCheck++
		if game.IsEndState() {
			game.loopWaitGroup.Done()
			break
		}
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
		func() {
			defer func() {
				if r := recover(); r != nil {
					game.gameLogger.Log("Recover SecLoop Panic: " + fmt.Sprintf("%v", r))
				}
			}()
			if game.state.GameState == state.Waiting &&
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
						TotalPlayerNum: int32(game.state.Players.Length()),
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
					go game.alertGameStart()
					game.broadcast(gameStart)
					game.SetPlayingStatus(mapLength)
					game.gameLogger.Log("game start")
				}
			}
			game.calculator.SecLoop()
		}()

		if game.IsEndState() {
			game.loopWaitGroup.Done()
			break
		}
	}
}

func (game *Game) BroadcastLoop() {
	ticker := time.NewTicker(cons.BroadCastLoopInterval * time.Millisecond)
	// ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C { // broadcast loop

		func() {
			defer func() {
				if r := recover(); r != nil {
					game.gameLogger.Log("Recover BroadcastLoop Panic: " + fmt.Sprintf("%v", r))
				}
			}()

			game.state.ConfigMu.Lock()

			playerList := make([]*message.PlayerState, 0)
			game.state.Players.Range(func(key int, player *object.Player) bool {
				playerList = append(playerList, player.MakeSendingData())
				return true
			})

			bulletList := make([]*message.BulletState, 0)
			game.state.Bullets.Range(func(key int, bullet *object.Bullet) bool {
				bulletList = append(bulletList, bullet.MakeSendingData())
				return true
			})

			healPackList := make([]*message.HealPackState, 0)
			game.state.HealPacks.Range(func(key int, healPack *object.HealPack) bool {
				healPackList = append(healPackList, healPack.MakeSendingData())
				return true
			})

			magicItemList := make([]*message.MagicItemState, 0)
			game.state.MagicItems.Range(func(key int, magicItem *object.Magic) bool {
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
				LastSec:        int32((game.state.LastGameTick*16)/1000 + cons.WAITING_MAP_CHUNK_NUM),
				PlayerState:    playerList,
				BulletState:    bulletList,
				HealPackState:  healPackList,
				MagicItemState: magicItemList,
				TileState:      tileStateList,
				ChangingState:  game.state.ChangingState.MakeSendingData(),
			}

			// log.Println(gameState)
			// log.Printf("\n\n")

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

		}()
		game.loopState.broadcastLoopCheck++
		if game.IsEndState() {
			game.loopWaitGroup.Done()
			break
		}
	}
}

func (game *Game) OnMessage(data []byte, id int) {
	game.handleMessage(id, data)
}

func (game *Game) OnClient(client *Client) {

	if game.IsEndState() {
		client.close()
		return
	}
	lastClientId, exists := game.serverClientTable.Get(client.serverClientId)
	if exists {
		log.Printf("reconnection: %s", client.serverClientId)
		game.SetReconnectionPlayer(client, *lastClientId)
		client.ID = ClientId(game.ClientIdAllocator.AllocateClientId())
		game.clients.Store(client.ID, client)
		// MARK: reconnection 처리 후에 store 해야 함!
		// - broad cast 하기 전에, client한테 reconnection 메시지를 보내야 함
	} else {
		// 해당 serverClientId가 존재 하지 않는 경우,
		// wating 상태면 추가, 아니면 에러
		if game.state.GameState == state.Waiting {
			log.Printf("new game client's server client id: %s", client.serverClientId)
			game.playerNum++
			client.ID = ClientId(game.ClientIdAllocator.AllocateClientId())
			game.SetPlayer(client)
			game.clients.Store(client.ID, client)
			game.serverClientTable.Store(client.serverClientId, ClientId(client.ID))
		} else {
			log.Printf("fail to reconnect: %s", client.serverClientId)
			client.close()
			return
		}
	}

	if game.clients.Length() > 0 {
		game.ResetEndCount()
	}

	currentPlayerList := make([]*message.CurrentPlayers, 0)
	game.clients.Range(func(ci ClientId, c *Client) bool {
		currentPlayerList = append(currentPlayerList, &message.CurrentPlayers{
			PlayerId:   int32(c.ID),
			PlayerName: c.userName,
		})
		return true
	})

	data, err := proto.Marshal(&message.Wrapper{
		MessageType: &message.Wrapper_NewUser{
			NewUser: &message.NewUser{
				CurrentPlayers: currentPlayerList,
			},
		},
	})

	if err != nil {
		log.Printf("Failed to marshal GameState: %v", err)
		return
	}

	game.broadcast(data)

	game.gameLogger.Log("On Client " + client.serverClientId)
}

func (game *Game) OnObserver(client *Client) {

	if game.IsEndState() {
		client.close()
		return
	}

	log.Printf("new game client's server client id: %s", client.serverClientId)
	client.ID = ClientId(game.ClientIdAllocator.AllocateClientId())
	game.clients.Store(client.ID, client)

	game.gameLogger.Log("On Observer")
}

func (game *Game) OnClose(client *Client) {
	game.clients.Delete(ClientId(client.ID))
	if game.state.GameState == state.Waiting {
		go game.alertPlayerLeavae(client)
	}

	if game.clients.Length() == 0 {
		log.Println("clinet's count zero")

		game.PlusEndCount()
	}

	game.gameLogger.Log("On Close " + client.serverClientId)
}

func (game *Game) handleMessage(clientId int, data []byte) {
	var wrapper message.Wrapper
	if err := proto.Unmarshal(data, &wrapper); err != nil {
		log.Printf("Failed to unmarshal message from client %d: %v", clientId, err)
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

	clientId := client.ID

	x := rand.Intn(int(game.state.MaxCoord))
	y := rand.Intn(int(game.state.MaxCoord))

	game.state.AddPlayer(int(clientId), float64(x), float64(y))

	// send GameInit message
	gameInit := &message.GameInit{
		Id:           int32(clientId),
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
	client.write(data)
}

func (game *Game) SetReconnectionPlayer(client *Client, lastClientId ClientId) {
	client.ID = lastClientId
	gameReconnect := &message.GameReconnect{
		Id:             int32(client.ID),
		MinPlayerNum:   int32(game.minPlayerNum),
		TotalPlayerNum: int32(game.state.Players.Length()),
	}
	data, err := proto.Marshal(&message.Wrapper{
		MessageType: &message.Wrapper_GameReconnect{
			GameReconnect: gameReconnect,
		},
	})
	if err != nil {
		// log.Printf("Failed to marshal GameInit: %v", err)
		game.gameLogger.Log("Failed to marshal GameInit: " + err.Error())
		return
	}
	client.write(data)
}

func (game *Game) broadcast(data []byte) {
	game.broadcastDataSize = len(data)
	game.clients.Range(func(id ClientId, client *Client) bool {
		// log.Println(client)
		client.write(data)
		return true
	})
	// log.Printf("end broadcast")
}

func (game *Game) IsPlaying() bool {
	if game.debug {
		return true
	}
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

	// 한 번만 end 처리! end 요청이 여러번 가지 않도록!
	if game.endTickCount == END_GAME_MAX_TICK_COUNT {
		// log.Print("alert game end")
		game.gameLogger.Log("alert game end")
		game.state.GameState = state.End
	} else {
		return
	}
}

// client가 zero가 되거나, game이 끝나면 plus 시작
func (game *Game) PlusEndCount() {
	game.endTickCountMu.Lock()
	defer game.endTickCountMu.Unlock()

	IDEL_TICK := 60 * 5
	endThreashold := END_GAME_MAX_TICK_COUNT - IDEL_TICK

	// 이미 게임이 끝났다면, 5초 정도 여유를 두고 종료 메시지 전송
	if game.IsEndState() {
		if game.endTickCount < endThreashold {
			game.endTickCount = END_GAME_MAX_TICK_COUNT - IDEL_TICK
		} else {
			game.endTickCount += 1
		}
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
	if game.debug {
		return false
	}
	return game.state.GameState == state.End
}

func (game *Game) GameStateCheckLoop() {

	time.Sleep(5 * time.Second)

	ticker := time.NewTicker(cons.CheckLoopStateInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {

		// if game.loopState.calcLoopCheck == game.loopState.lastCalcLoopCheck {
		// 	game.gameLogger.Log("Calc loop restart")
		// 	go game.CalcGameTickLoop()
		// }
		// game.loopState.lastCalcLoopCheck = game.loopState.calcLoopCheck

		// if game.loopState.broadcastLoopCheck == game.loopState.lastBroadcastLoopCheck {
		// 	game.gameLogger.Log("Broadcast loop restart")
		// 	go game.BroadcastLoop()
		// }
		// game.loopState.lastBroadcastLoopCheck = game.loopState.broadcastLoopCheck

		if game.IsEndState() {
			game.loopWaitGroup.Done()
			game.gameLogger.Log("🚀 break")
			break
		}

		game.loopState.lastCalcLoopCheck = game.loopState.calcLoopCheck
		// calcLoopInSec := game.loopState.calcLoopCheck - game.loopState.lastCalcLoopCheck
		// broadCastLoopInSec := game.loopState.broadcastLoopCheck - game.loopState.lastBroadcastLoopCheck
		game.loopState.lastBroadcastLoopCheck = game.loopState.broadcastLoopCheck
		// game.gameLogger.Log("calc, broad, player, last, size: " + fmt.Sprint(calcLoopInSec) + " " + fmt.Sprint(broadCastLoopInSec) + " " + fmt.Sprint(game.state.Players.Length()) + " " + fmt.Sprint(game.state.LastGameTick) + " " + fmt.Sprint(game.broadcastDataSize))
	}
}
