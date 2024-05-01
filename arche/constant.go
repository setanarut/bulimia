package arche

import (
	"bulimia/engine/cm"
)

// Collision Bitmask Category
const (
	BitmaskPlayer      uint = 1
	BitmaskEnemy       uint = 2
	BitmaskBomb        uint = 4
	BitmaskFood        uint = 8
	BitmaskWall        uint = 16
	BitmaskDoor        uint = 32
	BitmaskCollectible uint = 64
	BitmaskBombRaycast uint = 128
)

// Collision type
const (
	CollisionTypePlayer cm.CollisionType = iota
	CollisionTypeEnemy
	CollisionTypeWall
	CollisionTypeFood
	CollisionTypeBomb
	CollisionTypeCollectible
	CollisionTypeDoor
)

var FilterBombRaycast cm.ShapeFilter = cm.NewShapeFilter(0, BitmaskBombRaycast, cm.AllCategories&^BitmaskBomb)

/* var (
	QueryEnemy  = donburi.NewQuery(filter.Contains(component.EnemyTagComp))
	QueryPlayer = donburi.NewQuery(filter.Contains(component.PlayerTagComp))
	QueryDoor   = donburi.NewQuery(filter.Contains(component.DoorComp))
	QueryFood   = donburi.NewQuery(filter.Contains(component.FoodTagComp))
	QueryBomb   = donburi.NewQuery(filter.Contains(component.BombTagComp))
	QueryAI     = donburi.NewQuery(filter.Contains(component.AIComp))
	QueryCamera = donburi.NewQuery(filter.Contains(component.CameraComp))
)
*/
