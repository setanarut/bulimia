package system

import (
	"bulimia/arche"
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

var Rooms []cm.BB
var CurrentRoom cm.BB

type EntitySpawnSystem struct {
	queryBody, queryCam *donburi.Query
	scr                 *cm.BB
	cam                 *engine.Camera
	// spawnTimer          *engine.Timer
}

func NewEntitySpawnSystem() *EntitySpawnSystem {

	return &EntitySpawnSystem{
		// spawnTimer: engine.NewTimer(time.Second * 2),
		queryBody: query.NewQuery(filter.Contains(component.BodyComp)),
		queryCam:  query.NewQuery(filter.Contains(component.CameraComp)),
	}
}

func (sys *EntitySpawnSystem) Init(world donburi.World, space *cm.Space, scr *cm.BB) {
	e := arche.NewCameraEntity(scr.Center(), scr.R, scr.T, world)
	sys.cam = component.CameraComp.Get(e)
	sys.scr = scr
	Rooms = make([]cm.BB, 0)
	Rooms = append(Rooms, *scr)                           // middle 0
	Rooms = append(Rooms, scr.Offset(cm.Vec2{0, scr.T}))  // top 1
	Rooms = append(Rooms, scr.Offset(cm.Vec2{0, -scr.T})) // bottom 2

	CurrentRoom = *scr
	arche.MakeRoom(world, space, Rooms[1], arche.RoomOptions{true, false, true, true, 5, -1, 6, 7})
	arche.MakeRoom(world, space, Rooms[0], arche.RoomOptions{true, true, true, true, 1, 2, 3, 4})
	arche.MakeRoom(world, space, Rooms[2], arche.RoomOptions{false, true, true, true, -1, 8, 9, 10})

	arche.NewPlayerEntity(0.1, 0.3, 0, 20, world, space, scr.Center().Add(cm.Vec2{0, -120}))

	arche.NewWallEntity(world, space, CurrentRoom.Center(), 100, 100)

	// top room
	for i := 5; i < 8; i++ {
		arche.DefaultEnemy(world, space, engine.RandomPointInBB(Rooms[1], 20))
		// arche.DefaultKeyCollectible(i, world, space, engine.RandomPointInBB(Rooms[1], 20))
	}
	// center room
	for i := 1; i < 5; i++ {
		arche.DefaultEnemy(world, space, engine.RandomPointInBB(Rooms[0], 20))
		arche.DefaultKeyCollectible(i, world, space, engine.RandomPointInBB(Rooms[0], 20))
	}
	// bottom room
	for i := 8; i < 11; i++ {
		arche.DefaultEnemy(world, space, engine.RandomPointInBB(Rooms[2], 20))
		// arche.DefaultKeyCollectible(i, world, space, engine.RandomPointInBB(Rooms[2], 20))
	}

	world.OnRemove(func(world donburi.World, entity donburi.Entity) {
		e := world.Entry(entity)

		// adds trauma to the camera when the bomb is removed
		if e.HasComponent(component.BombTagComp) {
			sys.cam.AddTrauma(0.2)
		}
		// adds trauma to the camera when the bomb is removed
		if e.HasComponent(component.EnemyTagComp) {
			p := component.BodyComp.Get(e).Position()
			arche.NewCollectibleEntity(component.Food, 10, -1, 10, world, space, p)
		}

	})

}

func (sys *EntitySpawnSystem) Update(world donburi.World, space *cm.Space) {

	worldPos := sys.cam.ScreenToWorld(ebiten.CursorPosition())
	cursor := engine.InvPosVectY(worldPos, CurrentRoom.T)

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		arche.DefaultBomb(world, space, cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		arche.DefaultEnemy(world, space, cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		arche.RandomCollectible(world, space, cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		arche.NewWallEntity(world, space, cursor, 200, 20)

	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		arche.NewWallEntity(world, space, cursor, 20, 200)

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {

		for range 10 {
			arche.DefaultEnemy(world, space, engine.RandomPointInBB(CurrentRoom, 64))
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {

		for range 10 {
			arche.DefaultBomb(world, space, engine.RandomPointInBB(CurrentRoom, 64))
		}

	}

}

func (sys *EntitySpawnSystem) Draw(world donburi.World, space *cm.Space, screen *ebiten.Image) {
}
