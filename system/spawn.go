package system

import (
	"bulimia/arche"
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"bulimia/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type EntitySpawnSystem struct {
}

func NewEntitySpawnSystem() *EntitySpawnSystem {

	return &EntitySpawnSystem{}
}

func (sys *EntitySpawnSystem) Init() {
	res.CurrentRoom = res.ScreenRect

	res.Rooms = make([]cm.BB, 0)
	res.Rooms = append(res.Rooms, res.CurrentRoom)                                        // middle 0
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{0, res.CurrentRoom.T}))  // top 1
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{0, -res.CurrentRoom.T})) // bottom 2
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{-res.CurrentRoom.R, 0})) // left 3
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{res.CurrentRoom.R, 0}))  // right 4

	arche.SpawnRoom(res.Rooms[0], arche.RoomOptions{true, true, true, true, 1, 2, 3, 4})
	arche.SpawnRoom(res.Rooms[1], arche.RoomOptions{true, false, true, true, 5, -1, 6, 7})
	arche.SpawnRoom(res.Rooms[2], arche.RoomOptions{false, true, true, true, -1, 8, 9, 10})
	arche.SpawnRoom(res.Rooms[3], arche.RoomOptions{true, true, true, false, 11, 12, 13, -1})
	arche.SpawnRoom(res.Rooms[4], arche.RoomOptions{true, true, false, true, 14, 15, -1, 16})

	ResetLevel()

	res.World.OnRemove(func(world donburi.World, entity donburi.Entity) {
		e := world.Entry(entity)
		if e.HasComponent(comp.EnemyTag) {
			p := comp.Body.Get(e).Position()
			i := comp.Inventory.Get(e)
			for _, v := range i.Keys {
				arche.SpawnDefaultKeyCollectible(v, p)
			}
			for range i.Bombs {
				arche.SpawnDefaultBomb(p)
			}
			for range i.Foods {
				arche.SpawnDefaultFood(p)
			}
		}

	})

}

func (sys *EntitySpawnSystem) Update() {

	// Reset Level
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		ResetLevel()
	}

	worldPos := res.Camera.ScreenToWorld(ebiten.CursorPosition())
	cursor := engine.InvPosVectY(worldPos, res.CurrentRoom.T)

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		arche.SpawnDefaultBomb(cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		arche.SpawnDefaultEmeticCollectible(cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		arche.SpawnRandomCollectible(cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		arche.SpawnWall(cursor, 200, 20)

	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		arche.SpawnWall(cursor, 20, 200)

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {

		for range 500 {
			// arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.CurrentRoom, 64))
			arche.SpawnEnemy(0.3, 0.3, 0.5, 8, engine.RandomPointInBB(res.CurrentRoom, 64))
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {

		for range 10 {
			arche.SpawnDefaultBomb(engine.RandomPointInBB(res.CurrentRoom, 64))
		}

	}

}

func (sys *EntitySpawnSystem) Draw() {
}

func ResetLevel() {
	res.CurrentRoom = res.ScreenRect
	res.Camera.LookAt(res.CurrentRoom.Center())

	player, ok := comp.PlayerTag.First(res.World)
	if ok {
		DestroyEntryWithBody(player)
		arche.SpawnDefaultPlayer(res.CurrentRoom.Center().Add(cm.Vec2{0, -100}))
	} else {
		arche.SpawnDefaultPlayer(res.CurrentRoom.Center().Add(cm.Vec2{0, -100}))
	}

	comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {
		DestroyEntryWithBody(e)
	})
	comp.Collectible.Each(res.World, func(e *donburi.Entry) {
		DestroyEntryWithBody(e)
	})
	comp.BombTag.Each(res.World, func(e *donburi.Entry) {
		DestroyEntryWithBody(e)
	})

	// reset doors
	comp.Door.Each(res.World, func(e *donburi.Entry) {
		comp.Door.Get(e).PlayerHasKey = false
		comp.Door.Get(e).Open = false
	})

	// top room
	for i := 5; i < 8; i++ {
		arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.Rooms[1], 20))
		// arche.DefaultKeyCollectible(i,  engine.RandomPointInBB(resources.Rooms[1], 20))
	}
	// center room
	arche.SpawnDefaultEmeticCollectible(engine.RandomPointInBB(res.Rooms[0], 20))

	for i := 1; i < 5; i++ {
		e := arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.Rooms[0], 20))
		inv := comp.Inventory.Get(e)
		inv.Keys = append(inv.Keys, i)

		// arche.SpawnDefaultKeyCollectible(i, engine.RandomPointInBB(res.Rooms[0], 20))
	}
	// bottom room
	for i := 8; i < 11; i++ {
		arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.Rooms[2], 20))
		// arche.DefaultKeyCollectible(i,  engine.RandomPointInBB(resources.Rooms[2], 20))
	}

}
