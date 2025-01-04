package object

import (
	"jungle-royale/message"
	"jungle-royale/util"
)

type ChangingState struct {
	HeatBulletStateList *util.Set[HeatBulletState]
	GetItemStateList    *util.Set[GetItemState]
	PlayerDeadStateList *util.Set[PlayerDeadState]
}

func NewChangingState() *ChangingState {
	return &ChangingState{
		HeatBulletStateList: util.NewSyncSet[HeatBulletState](),
		GetItemStateList:    util.NewSyncSet[GetItemState](),
		PlayerDeadStateList: util.NewSyncSet[PlayerDeadState](),
	}
}

func (cs *ChangingState) MakeSendingData() *message.ChangingState {

	HeatBulletStateList := make([]*message.HeatBulletState, 0)
	GetItemStateList := make([]*message.GetItemState, 0)
	playerDeadStateList := make([]*message.PlayerDeadState, 0)

	cs.HeatBulletStateList.Range(func(hbs HeatBulletState) bool {
		HeatBulletStateList = append(HeatBulletStateList, hbs.MakeSendingData())
		cs.HeatBulletStateList.Remove(hbs)
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
		HeatBulletState: HeatBulletStateList,
		GetItemState:    GetItemStateList,
		PlayerDeadState: playerDeadStateList,
	}
}

type HeatBulletState struct {
	bulletId string
	playerId string
}

func NewHeatBulletState(bulletId, playerId string) HeatBulletState {
	return HeatBulletState{bulletId, playerId}
}

func (hbs *HeatBulletState) MakeSendingData() *message.HeatBulletState {
	return &message.HeatBulletState{
		BulletId: hbs.bulletId,
		PlayerId: hbs.playerId,
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
