package arche

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"bulimia/res"
	"image/color"

	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

func SpawnBody(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {
	body := cm.NewBody(m, cm.MomentForCircle(m, 0, r*2, cm.Vec2{}))
	shape := cm.NewCircle(body, r, cm.Vec2{})
	shape.SetElasticity(e)
	shape.SetFriction(f)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())
	body.SetPosition(pos)
	bodyEntry := res.World.Entry(res.World.Create(comp.Body))
	body.UserData = bodyEntry
	comp.Body.Set(bodyEntry, body)
	return bodyEntry
}

func SpawnPlayer(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {
	entry := SpawnBody(m, e, f, r, pos)
	body := comp.Body.Get(entry)
	body.FirstShape().SetCollisionType(CollisionTypePlayer)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskPlayer, cm.AllCategories&^BitmaskFood)
	body.SetMoment(cm.Intinity)

	entry.AddComponent(comp.PlayerTag)
	entry.AddComponent(comp.Inventory)
	entry.AddComponent(comp.Living)
	entry.AddComponent(comp.Render)

	render := comp.Render.Get(entry)

	render.AnimPlayer = engine.NewAnimationPlayer(res.Pacman)
	w := 100
	render.AnimPlayer.AddStateAnimation("shootR", 0, 0, w, w, 1, true)
	render.AnimPlayer.AddStateAnimation("shootL", 0, w, w, w, 1, true)
	render.AnimPlayer.AddStateAnimation("shootU", 0, w*2, w, w, 1, true)
	render.AnimPlayer.AddStateAnimation("shootD", 0, w*3, w, w, 1, true)

	render.AnimPlayer.AddStateAnimation("right", 0, 0, w, w, 4, true)
	render.AnimPlayer.AddStateAnimation("left", 0, w, w, w, 4, true)
	render.AnimPlayer.AddStateAnimation("up", 0, w*2, w, w, 4, true)
	render.AnimPlayer.AddStateAnimation("down", 0, w*3, w, w, 4, true)

	render.AnimPlayer.SetFPS(8)
	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.ScaleColor = colornames.Yellow
	return entry
}

func SpawnEnemy(m, e, f, r, viewRadius float64, pos cm.Vec2) *donburi.Entry {
	entry := SpawnBody(m, e, f, r, pos)
	body := comp.Body.Get(entry)
	body.SetMoment(cm.Intinity)

	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskEnemy, cm.AllCategories)
	body.FirstShape().SetCollisionType(CollisionTypeEnemy)

	entry.AddComponent(comp.EnemyTag)
	entry.AddComponent(comp.AI)
	entry.AddComponent(comp.Living)
	entry.AddComponent(comp.Render)
	entry.AddComponent(comp.Gradient)

	render := comp.Render.Get(entry)
	liv := comp.Living.Get(entry)
	liv.Speed = 480
	render.AnimPlayer = engine.NewAnimationPlayer(res.Enemy)
	render.AnimPlayer.AddStateAnimation("idle", 200, 0, 100, 100, 1, false)
	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	return entry
}

func SpawnBomb(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {
	entry := SpawnBody(m, e, f, r, pos)
	body := comp.Body.Get(entry)
	body.FirstShape().SetCollisionType(CollisionTypeBomb)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskBomb, cm.AllCategories)

	entry.AddComponent(comp.BombTag)
	entry.AddComponent(comp.Render)

	render := comp.Render.Get(entry)

	render.AnimPlayer = engine.NewAnimationPlayer(res.Items)
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, 100, 100, 1, false)
	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	render.ScaleColor = color.RGBA{244, 117, 88, 255}
	return entry
}

// SpawnFood e: elastiklik f: friction
func SpawnFood(m, e, f, r float64, pos cm.Vec2) *donburi.Entry {
	entry := SpawnBody(m, e, f, r, pos)
	body := comp.Body.Get(entry)
	body.FirstShape().SetCollisionType(CollisionTypeFood)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskFood, cm.AllCategories)

	entry.AddComponent(comp.FoodTag)
	entry.AddComponent(comp.Render)
	entry.AddComponent(comp.Damage)

	render := comp.Render.Get(entry)

	render.AnimPlayer = engine.NewAnimationPlayer(res.Items)
	render.AnimPlayer.AddStateAnimation("idle", 200, 0, 100, 100, 1, false)
	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	render.ScaleColor = colornames.Orange
	return entry
}

func SpawnWall(boxCenter cm.Vec2, boxW, boxH float64) *donburi.Entry {

	sbody := cm.NewStaticBody()
	wallShape := cm.NewBox(sbody, boxW, boxH, 0)
	wallShape.Filter = cm.NewShapeFilter(0, BitmaskWall, cm.AllCategories)
	wallShape.CollisionType = CollisionTypeWall
	wallShape.SetElasticity(0)
	wallShape.SetFriction(0)
	sbody.SetPosition(boxCenter)
	res.Space.AddShape(wallShape)
	res.Space.AddBody(wallShape.Body())

	// components
	entry := res.World.Entry(res.World.Create(
		comp.Body,
		comp.WallTag,
		comp.Render,
	))
	wallShape.Body().UserData = entry
	comp.Body.Set(entry, wallShape.Body())

	render := comp.Render.Get(entry)

	render.AnimPlayer = engine.NewAnimationPlayer(res.Wall)
	imW := res.Wall.Bounds().Dx()
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, imW, imW, 1, false)
	render.DrawScale = engine.GetBoxScaleFactor(float64(imW), float64(imW), boxW, boxH)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	render.ScaleColor = colornames.Blue
	return entry
}
func SpawnDoor(
	boxCenter cm.Vec2,
	boxW,
	boxH float64,
	lockNumber int) *donburi.Entry {

	sbody := cm.NewStaticBody()
	shape := cm.NewBox(sbody, boxW, boxH, 0)
	shape.Filter = cm.NewShapeFilter(0, BitmaskDoor, cm.AllCategories)
	shape.SetSensor(false)
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.CollisionType = CollisionTypeDoor
	sbody.SetPosition(boxCenter)
	res.Space.AddShape(shape)
	res.Space.AddBody(shape.Body())

	// components
	entry := res.World.Entry(res.World.Create(
		comp.Body,
		comp.Door,
		comp.Render,
	))
	shape.Body().UserData = entry
	comp.Body.Set(entry, shape.Body())
	comp.Door.SetValue(entry, comp.DoorData{LockNumber: lockNumber})

	render := comp.Render.Get(entry)

	render.AnimPlayer = engine.NewAnimationPlayer(res.Wall)
	imW := res.Wall.Bounds().Dx()
	render.AnimPlayer.AddStateAnimation("idle", 0, 0, imW, imW, 1, false)
	render.DrawScale = engine.GetBoxScaleFactor(float64(imW), float64(imW), boxW, boxH)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.AnimPlayer.Paused = true
	return entry
}

func SpawnCollectible(itemType comp.ItemType, count, keyNumber int,
	r float64, pos cm.Vec2) *donburi.Entry {

	entry := SpawnBody(1, 0, 1, r, pos)
	body := comp.Body.Get(entry)
	body.FirstShape().Filter = cm.NewShapeFilter(0, BitmaskCollectible, BitmaskPlayer|BitmaskWall|BitmaskCollectible)
	body.FirstShape().SetCollisionType(CollisionTypeCollectible)

	entry.AddComponent(comp.Collectible)
	entry.AddComponent(comp.Render)

	comp.Collectible.SetValue(entry, comp.CollectibleData{
		Type:      itemType,
		ItemCount: count,
		KeyNumber: keyNumber})

	var ap *engine.AnimationPlayer

	switch itemType {

	case comp.Bomb:
		ap = engine.NewAnimationPlayer(res.Items)
		ap.AddStateAnimation("idle", 0, 0, 100, 100, 1, false)
	case comp.Key:
		ap = engine.NewAnimationPlayer(res.Items)
		ap.AddStateAnimation("idle", 100, 0, 100, 100, 1, false)
	case comp.Food:
		ap = engine.NewAnimationPlayer(res.Items)
		ap.AddStateAnimation("idle", 200, 0, 100, 100, 1, false)
	default:
		ap = engine.NewAnimationPlayer(res.Items)
		ap.AddStateAnimation("idle", 200, 0, 100, 100, 1, false)

	}

	render := comp.Render.Get(entry)
	render.AnimPlayer = ap

	render.AnimPlayer.Paused = true
	render.DrawScale = engine.GetCircleScaleFactor(r, render.AnimPlayer.CurrentFrame)
	render.Offset = engine.GetEbitenImageOffset(render.AnimPlayer.CurrentFrame)
	render.ScaleColor = colornames.Cyan

	return entry
}

func SpawnCamera(lookAt cm.Vec2, width, height float64) *donburi.Entry {
	e := res.World.Entry(res.World.Create(comp.Camera))
	cam := engine.NewCamera(lookAt, width, height)
	cam.Lerp = true
	comp.Camera.Set(e, cam)
	return e
}

func SpawnRoom(roomBB cm.BB, opts RoomOptions) {

	topDoorLength := roomBB.Width() / 5
	leftDoorLength := roomBB.Height() / 5

	topDoorCenter := roomBB.LT().Lerp(roomBB.RT(), 0.5)
	bottomDoorCenter := roomBB.LB().Lerp(roomBB.RB(), 0.5)

	leftDoorCenter := roomBB.LT().Lerp(roomBB.LB(), 0.5)
	rightDoorCenter := roomBB.RT().Lerp(roomBB.RB(), 0.5)

	topLeftWallCenter := cm.Vec2{roomBB.L + topDoorLength, roomBB.T}
	topRightWallCenter := cm.Vec2{roomBB.R - topDoorLength, roomBB.T}
	bottomLeftWallCenter := cm.Vec2{roomBB.L + topDoorLength, roomBB.B}
	bottomRightWallCenter := cm.Vec2{roomBB.R - topDoorLength, roomBB.B}

	leftDoorBottom := cm.Vec2{roomBB.L, roomBB.B + leftDoorLength}
	leftDoorTop := cm.Vec2{roomBB.L, roomBB.T - leftDoorLength}

	rightDoorBottom := cm.Vec2{roomBB.R, roomBB.B + leftDoorLength}
	rightDoorTop := cm.Vec2{roomBB.R, roomBB.T - leftDoorLength}

	// Top Wall
	if opts.TopWall {
		SpawnWall(topLeftWallCenter, topDoorLength*2, 10)
		SpawnDoor(topDoorCenter, topDoorLength, 10, opts.TopDoorKeyNumber)
		SpawnWall(topRightWallCenter, topDoorLength*2, 10)
	}

	// Bottom Wall
	if opts.BottomWall {
		SpawnWall(bottomLeftWallCenter, topDoorLength*2, 10)
		SpawnDoor(bottomDoorCenter, topDoorLength, 10, opts.BottomDoorKeyNumber)
		SpawnWall(bottomRightWallCenter, topDoorLength*2, 10)
	}

	// Left Wall
	if opts.LeftWall {
		SpawnWall(leftDoorTop, 10, leftDoorLength*2)
		SpawnDoor(leftDoorCenter, 10, leftDoorLength, opts.LeftDoorKeyNumber)
		SpawnWall(leftDoorBottom, 10, leftDoorLength*2)
	}

	// Right Wall
	if opts.RightWall {
		SpawnWall(rightDoorTop, 10, leftDoorLength*2)
		SpawnDoor(rightDoorCenter, 10, leftDoorLength, opts.RightDoorKeyNumber)
		SpawnWall(rightDoorBottom, 10, leftDoorLength*2)
	}
}

type RoomOptions struct {
	TopWall, BottomWall, LeftWall, RightWall                                     bool
	TopDoorKeyNumber, BottomDoorKeyNumber, LeftDoorKeyNumber, RightDoorKeyNumber int
}
