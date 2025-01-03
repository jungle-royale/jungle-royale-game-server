package calculator

import (
	"jungle-royale/cons"
	"jungle-royale/object"
	"jungle-royale/state"
	"jungle-royale/util"
	"time"
)

type Calculator struct {
	chunk       *Chunk
	state       *state.State
	LeafTileSet *util.Set[*object.Tile]
}

func NewCalculator(state *state.State) *Calculator {
	return &Calculator{
		state:       state,
		LeafTileSet: util.NewSet[*object.Tile](),
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
						calculator.chunk.RemoveKey(
							(*healPack).GetObjectId(),
							object.OBJECT_HEALPACK,
							calculator.chunk.getChunkIndexSet(*(*healPack).GetPhysical()),
						)
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
						calculator.chunk.RemoveKey(
							(*magicItem).GetObjectId(),
							object.OBJECT_HEALPACK,
							calculator.chunk.getChunkIndexSet(*(*magicItem).GetPhysical()),
						)
						return false
					}
				}
				return true
			})
			return true
		})

		// check player is on tile
		if calculator.state.GameState == state.Playing {

			if playerChunkIdx, valid := calculator.GetChunk().getChunkIndex(
				(*player.GetPhysical()).GetX(),
				(*player.GetPhysical()).GetY(),
			); !valid ||
				calculator.state.Tiles[playerChunkIdx.X][playerChunkIdx.Y].TileState == object.TILE_FALL {
				player.Dead("", object.DYING_FALL, calculator.state.Players.Length())
				calculator.state.Players.Delete(playerId)
				calculator.state.PlayerDead.Store(playerId, player.DyingStatus)
			}

			// onTile := false
			// for i := 0; i < calculator.chunk.chunkNum; i++ {
			// 	for j := 0; j < calculator.chunk.chunkNum; j++ {
			// 		tile := calculator.state.Tiles[i][j]
			// 		if tile.TileState != object.TILE_FALL && tile.PhysicalObject.IsInRectangle(
			// 			(*player.GetPhysical()).GetX(),
			// 			(*player.GetPhysical()).GetY(),
			// 		) {
			// 			onTile = true
			// 		}
			// 		return true
			// 	}
			// }
			// if !onTile {
			// 	player.Dead("", object.DYING_FALL, calculator.state.Players.Length())
			// 	calculator.state.Players.Delete(playerId)
			// 	calculator.state.PlayerDead.Store(playerId, player.DyingStatus)
			// }
		}

		// create bullet
		if player.IsShooting && player.ShootingCoolTime <= 0 {
			newBullet := player.CreateBullet()
			calculator.state.Bullets.Store(newBullet.GetObjectId(), newBullet)
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
		if calculator.state.LastGameTick%calculator.state.FallenTime == cons.TILE_FALL_ALERT_TIME%calculator.state.FallenTime {
			if tile, ok := calculator.LeafTileSet.SelectRandom(func(t *object.Tile) bool { return true }); ok {
				calculator.LeafTileSet.Remove(tile)
				tile.SetTileState(object.TILE_DANGEROUS)
				tile.ParentTile.ChildTile.Remove(tile)
				if tile.ParentTile.ChildTile.Length() == 0 && tile != tile.ParentTile {
					calculator.LeafTileSet.Add(tile.ParentTile)
				}
				time.AfterFunc(cons.TILE_FALL_ALERT_TIME*time.Second, func() {
					calculator.state.Tiles[tile.IdxI][tile.IdxJ].SetTileState(object.TILE_FALL)
					onObjectList := calculator.chunk.GetChunkKeySet(tile.IdxI, tile.IdxJ)

					// healPack
					onObjectList[object.OBJECT_HEALPACK].Range(func(s string) bool {
						calculator.state.HealPacks.Delete(s)
						return true
					})

					// magicItem
					onObjectList[object.OBJECT_MAGICITEM].Range(func(s string) bool {
						calculator.state.MagicItems.Delete(s)
						return true
					})
				})
			}
		}
		calculator.state.LastGameTick--
	}
}

func (calculator *Calculator) SecLoop() {
}
