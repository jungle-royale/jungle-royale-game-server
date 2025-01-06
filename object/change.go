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
	bulletId   string
	objectType int
	objectId   string
}

func NewHitBulletState(bulletId, objectId string, objectType int) HitBulletState {
	return HitBulletState{bulletId, objectType, objectId}
}

func (hbs *HitBulletState) MakeSendingData() *message.HitBulletState {
	return &message.HitBulletState{
		BulletId: hbs.bulletId,
		ObjectId: hbs.objectId,
	}
}

// item type
const (
	ITEM_HEALPACK = iota
	ITEM_STONE_MAGIC
	ITEM_FIRE_MAGIC
)

type GetItemState struct {
	itemId   string
	playerId string
	itemType int
}

func NewGetItemState(itemId, playerId string, itemType int) GetItemState {
	return GetItemState{itemId, playerId, itemType}
}

func (gis *GetItemState) MakeSendingData() *message.GetItemState {
	return &message.GetItemState{
		ItemId:   gis.itemId,
		PlayerId: gis.playerId,
		ItemType: int32(gis.itemType),
	}
}

type PlayerDeadState struct {
	Killer      string
	dead        string
	DyingStatus int
	KillNum     int
	Placement   int
}

func NewPlayerDeadState(killer string, dead string, ds int) *PlayerDeadState {
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
		KillerId:    pd.Killer,
		DeadId:      pd.dead,
		DyingStatus: int32(pd.DyingStatus),
	}
}
