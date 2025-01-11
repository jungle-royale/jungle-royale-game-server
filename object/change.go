package object

import (
	"jungle-royale/message"
	"jungle-royale/util"
)

type ChangingState struct {
	HitBulletStateList  *util.Set[HitBulletState]
	GetItemStateList    *util.Set[GetItemState]
	PlayerDeadStateList *util.Set[PlayerDeadState]
}

func NewChangingState() *ChangingState {
	return &ChangingState{
		HitBulletStateList:  util.NewSyncSet[HitBulletState](),
		GetItemStateList:    util.NewSyncSet[GetItemState](),
		PlayerDeadStateList: util.NewSyncSet[PlayerDeadState](),
	}
}

func (cs *ChangingState) MakeSendingData() *message.ChangingState {

	HitBulletStateList := make([]*message.HitBulletState, 0)
	GetItemStateList := make([]*message.GetItemState, 0)
	playerDeadStateList := make([]*message.PlayerDeadState, 0)

	cs.HitBulletStateList.Range(func(hbs HitBulletState) bool {
		HitBulletStateList = append(HitBulletStateList, hbs.MakeSendingData())
		cs.HitBulletStateList.Remove(hbs)
		return true
	})

	cs.GetItemStateList.Range(func(gis GetItemState) bool {
		GetItemStateList = append(GetItemStateList, gis.MakeSendingData())
		cs.GetItemStateList.Remove(gis)
		return true
	})

	cs.PlayerDeadStateList.Range(func(pds PlayerDeadState) bool {
		playerDeadStateList = append(playerDeadStateList, pds.MakeSendingData())
		cs.PlayerDeadStateList.Remove(pds)
		return true
	})

	return &message.ChangingState{
		HitBulletState:  HitBulletStateList,
		GetItemState:    GetItemStateList,
		PlayerDeadState: playerDeadStateList,
	}
}

type HitBulletState struct {
	bulletId   int
	objectType int
	objectId   int
	X          float64
	Y          float64
	bulletType int
}

func NewHitBulletState(bulletId, objectId, objectType int, x, y float64, bulletType int) HitBulletState {
	return HitBulletState{bulletId, objectType, objectId, x, y, bulletType}
}

func (hbs *HitBulletState) MakeSendingData() *message.HitBulletState {
	return &message.HitBulletState{
		ObjectType: int32(hbs.objectType),
		BulletId:   int32(hbs.bulletId),
		ObjectId:   int32(hbs.objectId),
		X:          float32(hbs.X),
		Y:          float32(hbs.Y),
		BulletType: int32(hbs.bulletType),
	}
}

// item type
const (
	ITEM_HEALPACK = iota
	ITEM_STONE_MAGIC
	ITEM_FIRE_MAGIC
)

type GetItemState struct {
	itemId   int
	playerId int
	itemType int
	X        float64
	Y        float64
}

func NewGetItemState(itemId, playerId, itemType int, x, y float64) GetItemState {
	return GetItemState{itemId, playerId, itemType, x, y}
}

func (gis *GetItemState) MakeSendingData() *message.GetItemState {
	return &message.GetItemState{
		ItemId:   int32(gis.itemId),
		PlayerId: int32(gis.playerId),
		ItemType: int32(gis.itemType),
		X:        float32(gis.X),
		Y:        float32(gis.Y),
	}
}

type PlayerDeadState struct {
	Killer      int
	Dead        int
	DyingStatus int
	KillNum     int
	Placement   int
}

func NewPlayerDeadState(killer, dead, ds int) *PlayerDeadState {
	return &PlayerDeadState{
		killer,
		dead,
		ds,
		0,
		-1,
	}
}

func (pd *PlayerDeadState) MakeSendingData() *message.PlayerDeadState {
	return &message.PlayerDeadState{
		KillerId:    int32(pd.Killer),
		DeadId:      int32(pd.Dead),
		DyingStatus: int32(pd.DyingStatus),
		KillNum:     int32(pd.KillNum),
		Placement:   int32(pd.Placement),
	}
}
