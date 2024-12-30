package calculator

import (
	"jungle-royale/object"
	"jungle-royale/object/physical"
	"jungle-royale/state"
)

type Calculator struct {
	state *state.State
}

type Collider interface {
	GetPhysical() *physical.Physical
}

func NewCalculator(state *state.State) *Calculator {
	return &Calculator{state}
}

func (calculator *Calculator) IsCollider(colliderA Collider, colliderB Collider) bool {
	return (*colliderA.GetPhysical()).IsCollide(colliderB.GetPhysical())
}

func (calculator *Calculator) CalcGameTickState() {

	// player
	calculator.state.Players.Range(func(playerId string, player *object.Player) bool {
		if !player.IsValid() {
			calculator.state.Players.Delete(playerId)
			return true
		}
		player.CalcGameTick()

		calculator.state.HealPacks.Range(func(key string, healPack *object.HealPack) bool {
			if calculator.IsCollider(player, healPack) {
				calculator.state.Bullets.Delete(playerId)
			}
			player.GetHealPack()
			calculator.state.HealPacks.Delete(healPack.Id)
			return true
		})

		calculator.state.MagicItems.Range(func(key string, magicItem *object.Magic) bool {
			if calculator.IsCollider(player, magicItem) {
				calculator.state.Bullets.Delete(playerId) // collide 되면 bullet을 왜 삭제 하는지?
			}
			return true
		})

		return true
	})

	// bullet
	calculator.state.Bullets.Range(func(bulletId string, bullet *object.Bullet) bool {
		bullet.CalcGameTick()

		calculator.state.Players.Range(func(key string, player *object.Player) bool {
			if calculator.IsCollider(bullet, player) {
				calculator.state.Bullets.Delete(bulletId)
				player.HeatedBullet(bullet.BulletType)
			}
			return true
		})

		if !bullet.IsValid() {
			calculator.state.Bullets.Delete(bulletId)
		}
		return true
	})

}

func (calculator *Calculator) SecLoop() {

}
