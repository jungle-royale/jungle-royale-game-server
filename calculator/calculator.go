package calculator

import (
	"jungle-royale/chunk"
	"jungle-royale/cons"
	"jungle-royale/object"
	"jungle-royale/object/physical"
	"jungle-royale/state"
	"time"
)

type Calculator struct {
	chunkNum  int
	chunkList [][]*chunk.Chunk
	state     *state.State
}

type Collider interface {
	GetPhysical() *physical.Physical
}

func NewCalculator(state *state.State) *Calculator {
	return &Calculator{
		state: state,
	}
}

func (calculator *Calculator) ConfigureCalculator(chunkNum int) {
	calculator.chunkNum = chunkNum
	calculator.chunkList = make([][]*chunk.Chunk, chunkNum)
	for i := 0; i < chunkNum; i++ {
		calculator.chunkList[i] = make([]*chunk.Chunk, chunkNum)
		for j := 0; j < chunkNum; j++ {
			calculator.chunkList[i][j] = chunk.NewChunk()
		}
	}
}

func (calculator *Calculator) IsCollider(colliderA Collider, colliderB Collider) bool {
	return (*colliderA.GetPhysical()).IsCollide(colliderB.GetPhysical())
}

func (calculator *Calculator) CalcGameTickState() {

	// player
	calculator.state.Players.Range(func(playerId string, player *object.Player) bool {
		if !player.IsValid() {
			player.DyingStatus.Placement = calculator.state.Players.Length()
			calculator.state.Players.Delete(playerId)
			calculator.state.PlayerDead.Store(playerId, player.DyingStatus)
			if killer, ok := calculator.state.Players.Get(player.DyingStatus.Killer); ok {
				(*killer).DyingStatus.Kill()
			}
			return true
		}
		player.CalcGameTick()

		calculator.state.HealPacks.Range(func(key string, healPack *object.HealPack) bool {
			if calculator.IsCollider(player, healPack) {
				calculator.state.HealPacks.Delete(healPack.Id)
				player.GetHealPack()
			}
			return true
		})

		calculator.state.MagicItems.Range(func(key string, magicItem *object.Magic) bool {
			if calculator.IsCollider(player, magicItem) {
				calculator.state.MagicItems.Delete(magicItem.ItemId)
			}
			return true
		})

		// check player is on tile
		if calculator.state.GameState == state.Playing {
			onTile := false
			calculator.state.Tiles.Range(func(i int, tile *object.Tile) bool {
				if tile.PhysicalObject.IsInRectangle(
					(*player.GetPhysical()).GetX(),
					(*player.GetPhysical()).GetY(),
				) {
					onTile = true
				}
				return true
			})
			if !onTile {
				player.Dead("", object.DYING_FALL, calculator.state.Players.Length())
				calculator.state.Players.Delete(playerId)
				calculator.state.PlayerDead.Store(playerId, player.DyingStatus)
			}
		}
		return true
	})

	// bullet
	calculator.state.Bullets.Range(func(bulletId string, bullet *object.Bullet) bool {
		bullet.CalcGameTick()

		calculator.state.Players.Range(func(key string, player *object.Player) bool {
			if calculator.IsCollider(bullet, player) {
				calculator.state.Bullets.Delete(bulletId)
				if calculator.state.GameState == state.Playing {
					player.HeatedBullet(bullet)
				}
			}
			return true
		})

		if !bullet.IsValid() {
			calculator.state.Bullets.Delete(bulletId)
		}
		return true
	})

	if calculator.state.GameState == state.Playing {
		// tile fall
		if calculator.state.LastGameTick%calculator.state.FallenTime == cons.TILE_FALL_ALERT_TIME {
			tileId, stile, _ := calculator.state.Tiles.SelectRandom(func(t *object.Tile) bool {
				return t.TileState == object.TILE_NORMAL
			})
			stile.SetTileState(object.TILE_DANGEROUS)
			ticker := time.NewTicker(cons.TILE_FALL_ALERT_TIME * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				calculator.state.Tiles.Delete(tileId)
			}
		}
		calculator.state.LastGameTick--
	}
}

func (calculator *Calculator) SecLoop() {
}
