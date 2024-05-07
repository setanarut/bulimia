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
	cam *engine.Camera
	// spawnTimer          *engine.Timer
}

func NewEntitySpawnSystem() *EntitySpawnSystem {

	return &EntitySpawnSystem{
		// spawnTimer: engine.NewTimer(time.Second * 2),
	}
}

func (sys *EntitySpawnSystem) Init() {
	e := arche.SpawnCamera(res.ScreenBox.Center(), res.ScreenBox.R, res.ScreenBox.T)
	sys.cam = comp.Camera.Get(e)
	res.CurrentRoom = res.ScreenBox

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

	arche.SpawnPlayer(0.1, 0.3, 0, 20, res.CurrentRoom.Center().Add(cm.Vec2{0, -120}))
	arche.SpawnWall(res.CurrentRoom.Center(), 100, 100)

	// top room
	for i := 5; i < 8; i++ {
		arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.Rooms[1], 20))
		// arche.DefaultKeyCollectible(i,  engine.RandomPointInBB(resources.Rooms[1], 20))
	}
	// center room
	for i := 1; i < 5; i++ {
		arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.Rooms[0], 20))
		arche.SpawnDefaultKeyCollectible(i, engine.RandomPointInBB(res.Rooms[0], 20))
	}
	// bottom room
	for i := 8; i < 11; i++ {
		arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.Rooms[2], 20))
		// arche.DefaultKeyCollectible(i,  engine.RandomPointInBB(resources.Rooms[2], 20))
	}

	res.World.OnRemove(func(world donburi.World, entity donburi.Entity) {
		e := world.Entry(entity)

		// adds trauma to the camera when the bomb is removed
		if e.HasComponent(comp.BombTag) {
			sys.cam.AddTrauma(0.2)
		}
		// adds trauma to the camera when the bomb is removed
		if e.HasComponent(comp.EnemyTag) {
			p := comp.Body.Get(e).Position()
			arche.SpawnCollectible(comp.Food, 10, -1, 10, p)
		}

	})

}

func (sys *EntitySpawnSystem) Update() {

	worldPos := sys.cam.ScreenToWorld(ebiten.CursorPosition())
	cursor := engine.InvPosVectY(worldPos, res.CurrentRoom.T)

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		arche.SpawnDefaultBomb(cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		arche.SpawnDefaultEnemy(cursor)
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

		for range 10 {
			arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.CurrentRoom, 64))
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
