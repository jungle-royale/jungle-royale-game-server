package state

import (
	"jungle-royale/message"
	"jungle-royale/object"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type State struct {
	Players sync.Map // 동시성을 지원하는 sync.Map
	Bullets sync.Map // 동시성을 지원하는 sync.Map
}

func NewState() *State {
	return &State{
		Players: sync.Map{},
		Bullets: sync.Map{},
	}
}

func (state *State) AddPlayer(id string) {
	newPlayer := object.NewPlayer(
		id,
		float32(rand.Intn(1000)),
		float32(rand.Intn(1000)),
	)
	state.Players.Store(id, newPlayer) // sync.Map에 데이터 저장
}

func (state *State) AddBullet(BulletCreateMessage *message.BulletCreate) {
	bulletId := uuid.New().String()
	newBullet := object.NewBullet(
		bulletId,
		BulletCreateMessage.PlayerId,
		BulletCreateMessage.StartX,
		BulletCreateMessage.StartY,
		float64(BulletCreateMessage.Angle),
		state.GetPlayers(), // 플레이어 목록을 일반 맵으로 변환
	)
	state.Bullets.Store(bulletId, newBullet) // sync.Map에 데이터 저장
}

// GetPlayers는 sync.Map의 데이터를 일반 map[string]*object.Player로 변환
func (state *State) GetPlayers() map[string]*object.Player {
	players := make(map[string]*object.Player)
	state.Players.Range(func(key, value any) bool {
		players[key.(string)] = value.(*object.Player)
		return true
	})
	return players
}

func (state *State) CalcState() {
	// 플레이어 상태 계산
	state.Players.Range(func(key, value any) bool {
		player := value.(*object.Player)
		player.Move()
		return true
	})

	// 총알 상태 계산
	state.Bullets.Range(func(key, value any) bool {
		bulletId := key.(string)
		bullet := value.(*object.Bullet)
		isValid := bullet.Move()
		if !isValid {
			state.Bullets.Delete(bulletId) // 유효하지 않은 총알 삭제
		}
		return true
	})
}
