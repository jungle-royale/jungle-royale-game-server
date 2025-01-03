package state

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/util"
	"math"
	"math/rand"

	"github.com/google/uuid"
)

// gamestate
const (
	Waiting = iota
	Counting
	Playing
)

type State struct {
	GameState    int
	ChunkNum     int
	Tiles        *util.Map[string, *object.Tile]
	Players      *util.Map[string, *object.Player]
	PlayerDead   *util.Map[string, *object.PlayerDead]
	Bullets      *util.Map[string, *object.Bullet]
	HealPacks    *util.Map[string, *object.HealPack]
	MagicItems   *util.Map[string, *object.Magic]
	FallenTime   int
	MaxCoord     float32
	LastGameTick int
}

func NewState() *State {
	return &State{
		Tiles:        util.NewSyncMap[string, *object.Tile](),
		Players:      util.NewSyncMap[string, *object.Player](),
		PlayerDead:   util.NewSyncMap[string, *object.PlayerDead](),
		Bullets:      util.NewSyncMap[string, *object.Bullet](),
		HealPacks:    util.NewSyncMap[string, *object.HealPack](),
		MagicItems:   util.NewSyncMap[string, *object.Magic](),
		FallenTime:   int(math.MaxInt),
		LastGameTick: -1,
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

	state.ChunkNum = chunkNum
	state.MaxCoord = float32(chunkNum * cons.CHUNK_LENGTH)

	state.LastGameTick = playingTime * 1000 / 16

	// map tile setting
	for i := 0; i < chunkNum; i++ {
		for j := 0; j < chunkNum; j++ {
			// log.Println(i, j)
			tildId := uuid.New().String()
			state.Tiles.Store(tildId, object.NewTile(
				tildId,
				float32(i*cons.CHUNK_LENGTH),
				float32(j*cons.CHUNK_LENGTH),
			))
		}
	}

	// nonfallen tile setting
	state.FallenTime = (state.LastGameTick - (cons.TILE_FALL_ALERT_TIME)) / (chunkNum * chunkNum)
}

func (state *State) AddPlayer(id string, x float32, y float32) {
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
