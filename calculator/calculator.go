package calculator

import (
	"jungle-royale/object"
	"jungle-royale/state"
)

type Calculator struct {
	state *state.State
}

func NewCalculator(state *state.State) *Calculator {
	return &Calculator{state}
}

func (calculator *Calculator) CalcGameTickState() {

	// player
	players := calculator.state.ObjectList.GetPlayers()
	players.Range(func(key, value any) bool {
		playerId := key.(string)
		player := value.(*object.Player)
		if !player.IsValid() {
			calculator.state.ObjectList.GetPlayers().Delete(playerId)
			return true
		}
		player.CalcGameTick()
		return true
	})

	// bullet
	bullets := calculator.state.ObjectList.GetBullets()
	bullets.Range(func(key, value1 any) bool {
		bulletId := key.(string)
		bullet := value1.(*object.Bullet)
		bullet.CalcGameTick()
		if collider := bullet.CalcCollision(calculator.state.ObjectList); collider != nil {
			calculator.state.ObjectList.GetBullets().Delete(bulletId)
			if player, ok := (*collider).(*object.Player); ok {
				player.HeatedBullet()
			}
		}
		if !bullet.IsValid() {
			calculator.state.ObjectList.GetBullets().Delete(bulletId)
		}
		return true
	})
}

func (calculator *Calculator) SecLoop() {

}
