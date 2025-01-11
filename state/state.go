package state

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/util"
	"math"
	"math/rand/v2"
	"sync"
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
	GameState         int
	ChunkNum          int
	Tiles             [][]*object.Tile
	Players           *util.Map[int, *object.Player]
	Bullets           *util.Map[int, *object.Bullet]
	HealPacks         *util.Map[int, *object.HealPack]
	MagicItems        *util.Map[int, *object.Magic]
	FallenTime        int
	MaxCoord          float32
	LastGameTick      int
	ChangingState     *object.ChangingState
	ConfigMu          sync.Mutex
	ObjectIdAllocator *object.ObjectIdAllocator
}

func NewState(oia *object.ObjectIdAllocator) *State {
	return &State{
		GameState:         Empty,
		Tiles:             nil,
		Players:           util.NewSyncMap[int, *object.Player](),
		Bullets:           util.NewSyncMap[int, *object.Bullet](),
		HealPacks:         util.NewSyncMap[int, *object.HealPack](),
		MagicItems:        util.NewSyncMap[int, *object.Magic](),
		FallenTime:        int(math.MaxInt),
		LastGameTick:      -1,
		ChangingState:     object.NewChangingState(),
		ConfigMu:          sync.Mutex{},
		ObjectIdAllocator: oia,
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
			tildId := state.ObjectIdAllocator.AllocateTileId()
			tileType := rand.IntN(object.TILE_TYPE_NUM)
			newTile := object.NewTile(
				tildId,
				float64(i*cons.CHUNK_LENGTH),
				float64(j*cons.CHUNK_LENGTH),
				i,
				j,
				func() int {
					return state.ObjectIdAllocator.AllocateTileId()
				},
			).SetTileState(object.TILE_TYPE_NUM).SetTileType(tileType)
			state.Tiles[i][j] = newTile
		}
	}

	// nonfallen tile setting
	state.FallenTime = (state.LastGameTick - (cons.TILE_FALL_ALERT_TIME)) / (chunkNum * chunkNum)

	state.ConfigMu.Unlock()
}

func (state *State) AddPlayer(id int, x, y float64) {
	newPlayer := object.NewPlayer(
		id,
		x,
		y,
	)
	state.Players.Store(id, newPlayer)
}

func (state *State) ChangeDirection(clientId int, msg *message.ChangeDir) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DirChange(float64(msg.GetAngle()), msg.IsMoved)
	}
}

func (state *State) ChangeAngle(clientId int, msg *message.ChangeAngle) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).AngleChange(float64(msg.GetAngle()))
	}
}

func (state *State) DoDash(clientId int, msg *message.DoDash) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DoDash()
	}
}

func (state *State) ChangeBulletState(clientId int, msg *message.ChangeBulletState) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).IsShooting = msg.IsShooting
	}
}
