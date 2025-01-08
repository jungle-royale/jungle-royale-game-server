package state

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/util"
	"math"
	"math/rand/v2"
	"sync"

	"github.com/google/uuid"
)

// gamestate
const (
	Empty = iota
	Waiting
	Counting
	Playing
	End
)

type State struct {
	GameState     int
	ChunkNum      int
	Tiles         [][]*object.Tile
	Players       *util.Map[string, *object.Player]
	Bullets       *util.Map[string, *object.Bullet]
	HealPacks     *util.Map[string, *object.HealPack]
	MagicItems    *util.Map[string, *object.Magic]
	FallenTime    int
	MaxCoord      float32
	LastGameTick  int
	ChangingState *object.ChangingState
	ConfigMu      sync.Mutex
}

func NewState() *State {
	return &State{
		GameState:     Empty,
		Players:       util.NewSyncMap[string, *object.Player](),
		Bullets:       util.NewSyncMap[string, *object.Bullet](),
		HealPacks:     util.NewSyncMap[string, *object.HealPack](),
		MagicItems:    util.NewSyncMap[string, *object.Magic](),
		FallenTime:    int(math.MaxInt),
		LastGameTick:  -1,
		ChangingState: object.NewChangingState(),
		ConfigMu:      sync.Mutex{},
	}
}

func (state *State) ConfigureState(chunkNum int, playingTime int) {

	state.ConfigMu.Lock()

	state.ChunkNum = chunkNum
	state.MaxCoord = float32(chunkNum * cons.CHUNK_LENGTH)

	state.LastGameTick = playingTime * 1000 / 16

	// map tile setting
	state.Tiles = make([][]*object.Tile, chunkNum)
	for i := 0; i < chunkNum; i++ {
		state.Tiles[i] = make([]*object.Tile, chunkNum)
		for j := 0; j < chunkNum; j++ {
			tildId := uuid.New().String()
			tileType := rand.IntN(object.TILE_TYPE_NUM)
			newTile := object.NewTile(
				tildId,
				float64(i*cons.CHUNK_LENGTH),
				float64(j*cons.CHUNK_LENGTH),
				i,
				j,
			).SetTileState(object.TILE_TYPE_NUM).SetTileType(tileType)
			state.Tiles[i][j] = newTile
		}
	}

	// nonfallen tile setting
	state.FallenTime = (state.LastGameTick - (cons.TILE_FALL_ALERT_TIME)) / (chunkNum * chunkNum)

	state.ConfigMu.Unlock()
}

func (state *State) AddPlayer(id string, x, y float64) {
	newPlayer := object.NewPlayer(
		id,
		x,
		y,
	)
	state.Players.Store(id, newPlayer)
}

func (state *State) ChangeDirection(clientId string, msg *message.ChangeDir) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DirChange(float64(msg.GetAngle()), msg.IsMoved)
	}
}

func (state *State) ChangeAngle(clientId string, msg *message.ChangeAngle) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).AngleChange(float64(msg.GetAngle()))
	}
}

func (state *State) DoDash(clientId string, msg *message.DoDash) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DoDash()
	}
}

func (state *State) ChangeBulletState(clientId string, msg *message.ChangeBulletState) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).IsShooting = msg.IsShooting
	}
}
