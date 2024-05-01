package component

import (
	"bulimia/engine"
	"bulimia/engine/cm"

	"github.com/yohamta/donburi"
)

type ItemType int

const (
	Food ItemType = iota
	Bomb
	Key
)

type InventoryData struct {
	Foods int
	Bombs int
	Keys  []int
}

type CollectibleData struct {
	Type      ItemType
	ItemCount int
	KeyNumber int
}

type AIData struct {
	Follow         bool
	FollowSpeed    float64
	FollowDistance float64
}

type DoorData struct {
	LockNumber   int
	Open         bool
	PlayerHasKey bool
}

var CameraComp = donburi.NewComponentType[engine.Camera]()
var InventoryComp = donburi.NewComponentType[InventoryData]()
var DoorComp = donburi.NewComponentType[DoorData]()
var CollectibleComp = donburi.NewComponentType[CollectibleData]()
var AnimPlayerComp = donburi.NewComponentType[engine.AnimationPlayer]()
var BodyComp = donburi.NewComponentType[cm.Body]()
var AIComp = donburi.NewComponentType[AIData](AIData{Follow: false, FollowSpeed: 500, FollowDistance: 300})

// Primitives
var AccelComp = donburi.NewComponentType[float64](40.0)
var SpeedComp = donburi.NewComponentType[float64](400.0)
var HealthComp = donburi.NewComponentType[int](100)
var DamageComp = donburi.NewComponentType[int](1)

// Tags
var PlayerTagComp = donburi.NewTag()
var WallTagComp = donburi.NewTag()
var FoodTagComp = donburi.NewTag()
var BombTagComp = donburi.NewTag()
var EnemyTagComp = donburi.NewTag()
