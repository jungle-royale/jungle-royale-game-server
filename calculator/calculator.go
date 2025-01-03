package calculator

import (
	"jungle-royale/cons"
	"jungle-royale/object"
	"jungle-royale/state"
	"time"
)

type Calculator struct {
	chunk *Chunk
	state *state.State
}

func NewCalculator(state *state.State) *Calculator {
	return &Calculator{
		state: state,
	}
}

func (calculator *Calculator) GetChunk() *Chunk {
	return calculator.chunk
}

func (calculator *Calculator) CalcMover(mo object.Mover) {
	beforeIndexSet := calculator.chunk.getChunkIndexSet(*mo.GetPhysical())
	mo.CalcGameTick()
	AfterIndexSet := calculator.chunk.getChunkIndexSet(*mo.GetPhysical())
	deleteIndexSet := beforeIndexSet.Difference(AfterIndexSet)
	calculator.chunk.RemoveKey(mo.GetObjectId(), mo.GetObjectType(), deleteIndexSet)
	addIndexSet := AfterIndexSet.Difference(beforeIndexSet)
	calculator.chunk.AddKey(mo.GetObjectId(), mo.GetObjectType(), addIndexSet)
}

func (calculator *Calculator) SetLocation(obj object.Object, x, y float32) {
	(*obj.GetPhysical()).SetCoord(x, y)
	indexSet := calculator.chunk.getChunkIndexSet(*obj.GetPhysical())
	calculator.chunk.AddKey(obj.GetObjectId(), obj.GetObjectType(), indexSet)
}

func (calculator *Calculator) ReLocation(mo object.Mover, x, y float32) {
	calculator.CalcMover(mo)
	(*mo.GetPhysical()).SetCoord(x, y)
	indexSet := calculator.chunk.getChunkIndexSet(*mo.GetPhysical())
	calculator.chunk.AddKey(mo.GetObjectId(), mo.GetObjectType(), indexSet)
}

func (calculator *Calculator) ConfigureCalculator(chunkNum int) {
	calculator.chunk = NewChunk(chunkNum)
}

func (calculator *Calculator) IsCollider(colliderA object.Collider, colliderB object.Collider) bool {
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
		calculator.CalcMover(player)

		chunkIndexSet := calculator.chunk.getChunkIndexSet(*player.GetPhysical())
		chunkIndexSet.Range(func(ci ChunkIndex) bool {
			objectSet := calculator.chunk.chunkTable[ci.X][ci.Y]

			// player - healPack
			objectSet[object.OBJECT_HEALPACK].Range(func(s string) bool {
				if healPack, ok := calculator.state.HealPacks.Get(s); ok {
					if calculator.IsCollider(player, *healPack) {
						calculator.state.HealPacks.Delete((*healPack).Id)
						player.GetHealPack()
						return false
					}
				}
				return true
			})

			// player - magicItem
			objectSet[object.OBJECT_MAGICITEM].Range(func(s string) bool {
				if magicItem, ok := calculator.state.MagicItems.Get(s); ok {
					if calculator.IsCollider(player, *magicItem) {
						calculator.state.MagicItems.Delete((*magicItem).ItemId)
						return false
					}
				}
				return true
			})

			return true
		})

		// check player is on tile
		if calculator.state.GameState == state.Playing {
			onTile := false
			calculator.state.Tiles.Range(func(key string, tile *object.Tile) bool {
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

		calculator.CalcMover(bullet)

		chunkIndexSet := calculator.chunk.getChunkIndexSet(*bullet.GetPhysical())
		chunkIndexSet.Range(func(ci ChunkIndex) bool {
			objectSet := calculator.chunk.chunkTable[ci.X][ci.Y]

			// bullet - player
			objectSet[object.OBJECT_PLAYER].Range(func(s string) bool {
				if player, ok := calculator.state.Players.Get(s); ok {
					if calculator.IsCollider(bullet, *player) {
						if calculator.state.GameState == state.Playing {
							if (*player).HeatedBullet(bullet) {
								calculator.state.Bullets.Delete(bulletId)
								return false
							}
						}
					}
				}
				return true
			})
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
			time.AfterFunc(cons.TILE_FALL_ALERT_TIME*time.Second, func() {
				calculator.state.Tiles.Delete(tileId)
			})
		}
		calculator.state.LastGameTick--
	}
}

func (calculator *Calculator) SecLoop() {
}
