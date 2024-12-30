package state

import (
	"container/list"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/object/physical"
	"jungle-royale/util"
	"math"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type Tile struct {
	TileBoundary *physical.Rectangle
}

func NewTile(x, y float32) *Tile {
	return &Tile{
		&physical.Rectangle{
			X:      x,
			Y:      y,
			Width:  float32(cons.CHUNK_LENGTH),
			Length: float32(cons.CHUNK_LENGTH),
		},
	}
}

func (tile *Tile) MakeSendingData() *message.FallenReadyTileState {
	return &message.FallenReadyTileState{
		X: tile.TileBoundary.X,
		Y: tile.TileBoundary.Y,
	}
}

// gamestate
const (
	Waiting = iota
	Counting
	Playing
)

type State struct {
	GameState       int
	Players         *util.Map[string, *object.Player]
	PlayerDead      *util.Map[string, *object.PlayerDead]
	Bullets         *util.Map[string, *object.Bullet]
	HealPacks       *util.Map[string, *object.HealPack]
	MagicItems      *util.Map[string, *object.Magic]
	MapBoundary     *physical.Rectangle
	NonFallenTile   *list.List
	FallenReadyTile *list.List
	FallenTile      *list.List
	TileMu          sync.Mutex // only for FallenReadyTile
	FallenTime      int
	MaxCoord        float32
	LastGameSec     int
}

func NewState() *State {
	return &State{
		Players:         util.NewSyncMap[string, *object.Player](),
		PlayerDead:      util.NewSyncMap[string, *object.PlayerDead](),
		Bullets:         util.NewSyncMap[string, *object.Bullet](),
		HealPacks:       util.NewSyncMap[string, *object.HealPack](),
		MagicItems:      util.NewSyncMap[string, *object.Magic](),
		NonFallenTile:   list.New(),
		FallenReadyTile: list.New(),
		FallenTile:      list.New(),
		TileMu:          sync.Mutex{},
		FallenTime:      int(math.MaxInt),
		LastGameSec:     -1,
	}
}

func randomShuffle(n int) []int {
	numbers := make([]int, n)
	for i := 0; i < n; i++ {
		numbers[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}

	return numbers
}

func (state *State) ConfigureState(chunkNum int, playingTime int) {
	state.MaxCoord = float32(chunkNum * cons.CHUNK_LENGTH)
	state.MapBoundary = physical.NewRectangle(0, 0, state.MaxCoord, state.MaxCoord)

	// nonfallen tile setting
	state.FallenTime = playingTime / (chunkNum * chunkNum)

	tempTile := make(map[int]util.Pair[float32, float32])
	idx := 0
	for i := 0; i < chunkNum; i++ {
		for j := 0; j < chunkNum; j++ {
			tempTile[idx] = util.Pair[float32, float32]{V0: float32(i), V1: float32(j)}
			idx++
		}
	}
	randomNum := randomShuffle(chunkNum * chunkNum)
	for _, v := range randomNum {
		state.NonFallenTile.PushBack(NewTile(tempTile[v].V0, tempTile[v].V1))
	}
}

func (state *State) AddPlayer(id string, x float32, y float32) {
	newPlayer := object.NewPlayer(
		id,
		x,
		y,
	)
	state.Players.Store(id, newPlayer)
}

func (state *State) AddBullet(x float32, y float32, clientId string, BulletCreateMessage *message.CreateBullet) {
	bulletId := uuid.New().String()
	if player, exists := state.Players.Get(clientId); exists {
		newBullet := object.NewBullet(
			uuid.New().String(),
			clientId,
			(*player).MagicType,
			x,
			y,
			float64(BulletCreateMessage.Angle),
		)
		state.Bullets.Store(bulletId, newBullet)
	}
}

func (state *State) ChangeDirection(clientId string, msg *message.ChangeDir) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DirChange(float64(msg.GetAngle()), msg.IsMoved)
	}
}

func (state *State) CreateBullet(clientId string, msg *message.CreateBullet) {
	if player, exists := state.Players.Get(clientId); exists {
		state.AddBullet((*(*player).GetPhysical()).GetX(), (*(*player).GetPhysical()).GetY(), clientId, msg)
	}
}

func (state *State) DoDash(clientId string, msg *message.DoDash) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DoDash()
	}
}
