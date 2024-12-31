package state

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/object/physical"
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
	Tiles        *util.Map[int, *object.Tile]
	Players      *util.Map[string, *object.Player]
	PlayerDead   *util.Map[string, *object.PlayerDead]
	Bullets      *util.Map[string, *object.Bullet]
	HealPacks    *util.Map[string, *object.HealPack]
	MagicItems   *util.Map[string, *object.Magic]
	MapBoundary  *physical.Rectangle
	FallenTime   int
	MaxCoord     float32
	LastGameTick int
}

func NewState() *State {
	return &State{
		Tiles:        util.NewSyncMap[int, *object.Tile](),
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
	state.MaxCoord = float32(chunkNum * cons.CHUNK_LENGTH)
	state.MapBoundary = physical.NewRectangle(0, 0, state.MaxCoord, state.MaxCoord)

	state.LastGameTick = playingTime * 60

	// map tile setting
	tileIdx := 0
	for i := 0; i < chunkNum; i++ {
		for j := 0; j < chunkNum; j++ {
			state.Tiles.Store(tileIdx, object.NewTile(
				tileIdx, float32(i*cons.CHUNK_LENGTH),
				float32(j*cons.CHUNK_LENGTH),
			))
			tileIdx++
		}
	}

	// nonfallen tile setting
	state.FallenTime = (playingTime * 60) / (chunkNum * chunkNum)
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
