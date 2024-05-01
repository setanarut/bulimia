package arche

import (
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"
	"bulimia/resources"

	db "github.com/yohamta/donburi"
)

func NewBodyEntity(m, e, f, r float64, world db.World, space *cm.Space, pos cm.Vec2) *db.Entry {
	body := cm.NewBody(m, cm.MomentForCircle(m, 0, r*2, cm.Vec2{}))
	shape := cm.NewCircle(body, r, cm.Vec2{})
	shape.SetElasticity(e)
	shape.SetFriction(f)

	space.AddShape(shape)
	space.AddBody(shape.Body())

	body.SetPosition(pos)

	bodyEntry := world.Entry(world.Create(
		component.BodyComp,
	))
	body.UserData = bodyEntry

	component.BodyComp.Set(bodyEntry, body)
	return bodyEntry
}

func NewPlayerEntity(m, e, f, r float64, world db.World, space *cm.Space, pos cm.Vec2) *db.Entry {
	entry := NewBodyEntity(m, e, f, r, world, space, pos)
	body := component.BodyComp.Get(entry)
	body.FirstShape().SetCollisionType(CollisionTypePlayer)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskPlayer, cm.AllCategories&^BitmaskFood)
	body.SetMoment(cm.Intinity)
	entry.AddComponent(component.PlayerTagComp)
	entry.AddComponent(component.InventoryComp)
	entry.AddComponent(component.HealthComp)
	entry.AddComponent(component.SpeedComp)
	entry.AddComponent(component.AccelComp)

	// animasyon player
	entry.AddComponent(component.AnimPlayerComp)
	ap := engine.NewAnimationPlayer(resources.Pacman)
	w := 100
	ap.AddStateAnimation("shootR", 0, 0, w, w, 1)
	ap.AddStateAnimation("shootL", 0, w, w, w, 1)
	ap.AddStateAnimation("shootU", 0, w*2, w, w, 1)
	ap.AddStateAnimation("shootD", 0, w*3, w, w, 1)

	ap.AddStateAnimation("right", 0, 0, w, w, 4)
	ap.AddStateAnimation("left", 0, w, w, w, 4)
	ap.AddStateAnimation("up", 0, w*2, w, w, 4)
	ap.AddStateAnimation("down", 0, w*3, w, w, 4)

	ap.SetFPS(6)
	ap.DrawScaleX = 2 * r / float64(ap.CurrentFrame.Bounds().Dx())
	ap.DrawScaleY = ap.DrawScaleX

	component.AnimPlayerComp.Set(entry, ap)

	component.HealthComp.SetValue(entry, 100)
	component.InventoryComp.Set(entry, &component.InventoryData{Bombs: 100, Foods: 100})
	return entry
}

func NewEnemyEntity(m, e, f, r, viewRadius float64, world db.World, space *cm.Space, pos cm.Vec2) *db.Entry {
	entry := NewBodyEntity(m, e, f, r, world, space, pos)
	body := component.BodyComp.Get(entry)
	body.SetMoment(cm.Intinity)

	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskEnemy, cm.AllCategories)
	body.FirstShape().SetCollisionType(CollisionTypeEnemy)

	// body.SetType(cm.BODY_DYNAMIC)
	entry.AddComponent(component.EnemyTagComp)
	entry.AddComponent(component.HealthComp)
	entry.AddComponent(component.AIComp)

	// anim player
	entry.AddComponent(component.AnimPlayerComp)
	ap := engine.NewAnimationPlayer(resources.Enemy)
	w := resources.Enemy.Bounds().Dx()
	ap.AddStateAnimation("idle", 0, 0, w, w, 1)
	ap.DrawScaleX = 2 * r / float64(ap.CurrentFrame.Bounds().Dx())
	ap.DrawScaleY = ap.DrawScaleX
	ap.Paused = true
	component.AnimPlayerComp.Set(entry, ap)
	return entry
}

func NewBombEntity(m, e, f, r float64, world db.World, space *cm.Space, pos cm.Vec2) *db.Entry {
	entry := NewBodyEntity(m, e, f, r, world, space, pos)
	body := component.BodyComp.Get(entry)
	body.FirstShape().SetCollisionType(CollisionTypeBomb)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskBomb, cm.AllCategories)

	entry.AddComponent(component.BombTagComp)
	// anim player
	entry.AddComponent(component.AnimPlayerComp)
	ap := engine.NewAnimationPlayer(resources.Bomb)
	w := resources.Bomb.Bounds().Dx()
	ap.AddStateAnimation("idle", 0, 0, w, w, 1)
	ap.DrawScaleX = 2 * r / float64(ap.CurrentFrame.Bounds().Dx())
	ap.DrawScaleY = ap.DrawScaleX
	ap.Paused = true
	component.AnimPlayerComp.Set(entry, ap)

	return entry
}

// NewFoodEntity e: elastiklik f: friction
func NewFoodEntity(m, e, f, r float64, world db.World, space *cm.Space, pos cm.Vec2) *db.Entry {
	entry := NewBodyEntity(m, e, f, r, world, space, pos)
	body := component.BodyComp.Get(entry)
	body.FirstShape().SetCollisionType(CollisionTypeFood)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskFood, cm.AllCategories)

	entry.AddComponent(component.FoodTagComp)
	entry.AddComponent(component.DamageComp)
	entry.AddComponent(component.AnimPlayerComp)

	// Set Damage
	component.DamageComp.SetValue(entry, 40)

	// Set animation player
	ap := engine.NewAnimationPlayer(resources.Food)
	ap.AddStateAnimation("idle", 0, 0, 64, 64, 1)
	ap.DrawScaleX = 2 * r / float64(ap.CurrentFrame.Bounds().Dx())
	ap.DrawScaleY = ap.DrawScaleX
	ap.Paused = true
	component.AnimPlayerComp.Set(entry, ap)

	// body.FirstShape().SetShapeFilter(cm.NewShapeFilter(1,1,1))
	return entry
}

func NewWallEntity(world db.World, space *cm.Space, pos cm.Vec2, w, h float64) *db.Entry {

	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, w, h, 0)
	wallShape.Filter = cm.NewShapeFilter(0, BitmaskWall, cm.AllCategories)
	wallShape.CollisionType = CollisionTypeWall
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0)
	sbody.SetPosition(pos)
	space.AddShape(wallShape)
	space.AddBody(wallShape.Body())

	// components
	entry := world.Entry(world.Create(
		component.BodyComp,
		component.WallTagComp,
	))

	// anim player component
	entry.AddComponent(component.AnimPlayerComp)
	ap := engine.NewAnimationPlayer(resources.Wall)
	ap.AddStateAnimation("idle", 0, 0, 30, 30, 1)
	ap.DrawScaleX, ap.DrawScaleY = engine.ScaleFactor(30, 30, w, h)
	ap.Paused = true
	component.AnimPlayerComp.Set(entry, ap)

	wallShape.Body().UserData = entry
	component.BodyComp.Set(entry, wallShape.Body())
	return entry
}
func NewDoorEntity(world db.World,
	space *cm.Space,
	pos cm.Vec2,
	w,
	h float64,
	lockNumber int) *db.Entry {

	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, w, h, 0)
	wallShape.Filter = cm.NewShapeFilter(0, BitmaskDoor, cm.AllCategories)
	wallShape.SetSensor(false)
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0)
	wallShape.CollisionType = CollisionTypeDoor
	sbody.SetPosition(pos)
	space.AddShape(wallShape)
	space.AddBody(wallShape.Body())

	// components
	entry := world.Entry(world.Create(
		component.BodyComp,
		component.DoorComp,
	))
	component.DoorComp.SetValue(entry, component.DoorData{LockNumber: lockNumber})

	// anim player component
	entry.AddComponent(component.AnimPlayerComp)
	ap := engine.NewAnimationPlayer(resources.Wall)
	ap.AddStateAnimation("idle", 0, 0, 30, 30, 1)
	ap.DrawScaleX, ap.DrawScaleY = engine.ScaleFactor(30, 30, w, h)
	ap.Paused = true
	component.AnimPlayerComp.Set(entry, ap)
	wallShape.Body().UserData = entry
	component.BodyComp.Set(entry, wallShape.Body())

	return entry
}

func NewCollectibleEntity(itemType component.ItemType, count, keyNumber int,
	r float64, world db.World, s *cm.Space, pos cm.Vec2) *db.Entry {

	entry := NewBodyEntity(1, 0, 1, r, world, s, pos)
	body := component.BodyComp.Get(entry)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskCollectible, BitmaskPlayer|BitmaskWall|BitmaskCollectible)
	body.FirstShape().SetCollisionType(CollisionTypeCollectible)
	// anim player
	entry.AddComponent(component.AnimPlayerComp)
	var ap *engine.AnimationPlayer

	switch itemType {

	case component.Food:
		ap = engine.NewAnimationPlayer(resources.Food)

		ap.AddStateAnimation("idle", 0, 0, resources.Food.Bounds().Dx(), resources.Food.Bounds().Dy(), 1)

	case component.Bomb:
		ap = engine.NewAnimationPlayer(resources.Bomb)
		ap.AddStateAnimation("idle", 0, 0, resources.Bomb.Bounds().Dx(), resources.Bomb.Bounds().Dy(), 1)

	case component.Key:
		ap = engine.NewAnimationPlayer(resources.Key)
		ap.AddStateAnimation("idle", 0, 0, resources.Key.Bounds().Dx(), resources.Key.Bounds().Dy(), 1)

	default:
		ap = engine.NewAnimationPlayer(resources.Food)
		ap.AddStateAnimation("idle", 0, 0, resources.Food.Bounds().Dx(), resources.Food.Bounds().Dy(), 1)

	}

	ap.DrawScaleX = 2 * r / float64(ap.CurrentFrame.Bounds().Dx())
	ap.DrawScaleY = ap.DrawScaleX
	ap.Paused = true
	component.AnimPlayerComp.Set(entry, ap)

	entry.AddComponent(component.CollectibleComp)
	component.CollectibleComp.SetValue(entry, component.CollectibleData{Type: itemType, ItemCount: count, KeyNumber: keyNumber})
	return entry
}
