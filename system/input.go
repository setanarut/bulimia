package system

import (
	"bulimia/arche"
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"
	"fmt"
	"time"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

const (
	shootDown  eb.Key = 28
	shootLeft  eb.Key = 29
	shootRight eb.Key = 30
	shootUp    eb.Key = 31
)

var ArrowKeys []eb.Key = []eb.Key{28, 29, 30, 31}
var WASDDirection cm.Vec2
var WASDDirectionTEMP cm.Vec2
var (
	NoDirection    = cm.Vec2{0, 0}
	RightDirection = cm.Vec2{1, 0}
	LeftDirection  = cm.Vec2{-1, 0}
	UpDirection    = cm.Vec2{0, 1}
	DownDirection  = cm.Vec2{0, -1}
)

type PlayerControlSystem struct {
	distance, bulletRadius float64
	bulletTimer            *engine.Timer
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{
		bulletTimer:  engine.NewTimer(time.Second / 4),
		bulletRadius: 5,
	}
}

func (ss *PlayerControlSystem) Init(world donburi.World, space *cm.Space, ScreenBox *cm.BB) {
	// set SetVelocityUpdateFunc
	if playerEntry, ok := component.PlayerTagComp.First(world); ok {
		body := component.BodyComp.Get(playerEntry)
		ss.distance = body.FirstShape().Class.(*cm.Circle).Radius()
		body.SetVelocityUpdateFunc(playerVelocityFunc)
	}
}

func (ss *PlayerControlSystem) Update(world donburi.World, space *cm.Space) {

	UpdateWASDDirection()

	if playerEntry, ok := component.PlayerTagComp.First(world); ok {
		inventory := component.InventoryComp.Get(playerEntry)
		playerBody := component.BodyComp.Get(playerEntry)
		playerAnimPlayer := component.AnimPlayerComp.Get(playerEntry)
		playerPos := playerBody.Position()

		if IsPressedAndNotABC(shootRight, shootLeft, shootDown, shootUp) {
			playerAnimPlayer.SetState("shootR")

			if inventory.Foods > 0 {

				if ss.bulletTimer.IsReady() {
					ss.bulletTimer.Reset()
				}

				if ss.bulletTimer.IsStart() {
					bullet := arche.NewFoodEntity(0.1, 0, 0.5, ss.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{100, 0}, playerPos)
				}

			}
		}

		if IsPressedAndNotABC(shootLeft, shootRight, shootDown, shootUp) {
			playerAnimPlayer.SetState("shootL")

			if inventory.Foods > 0 {

				if ss.bulletTimer.IsReady() {
					ss.bulletTimer.Reset()
				}

				if ss.bulletTimer.IsStart() {

					bullet := arche.NewFoodEntity(0.1, 0, 0.5, ss.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{-100, 0}, playerPos)
				}
			}
		}

		if IsPressedAndNotABC(shootUp, shootLeft, shootRight, shootDown) {
			playerAnimPlayer.SetState("shootU")

			if inventory.Foods > 0 {
				if ss.bulletTimer.IsReady() {
					ss.bulletTimer.Reset()
				}
				if ss.bulletTimer.IsStart() {
					bullet := arche.NewFoodEntity(0.1, 0, 0.5, ss.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{0, 100}, playerPos)
				}

			}

		}

		if IsPressedAndNotABC(shootDown, shootUp, shootLeft, shootRight) {
			playerAnimPlayer.SetState("shootD")

			if inventory.Foods > 0 {

				if ss.bulletTimer.IsReady() {
					ss.bulletTimer.Reset()
				}
				if ss.bulletTimer.IsStart() {
					bullet := arche.NewFoodEntity(0.1, 0, 0.5, ss.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{0, -100}, playerPos)
				}

			}
		}

		if inpututil.IsKeyJustReleased(shootUp) {
			playerAnimPlayer.SetState("up")
		}
		if inpututil.IsKeyJustReleased(shootDown) {
			playerAnimPlayer.SetState("down")
		}
		if inpututil.IsKeyJustReleased(shootLeft) {
			playerAnimPlayer.SetState("left")
		}
		if inpututil.IsKeyJustReleased(shootRight) {
			playerAnimPlayer.SetState("right")
		}

		if !AnyKeyDown(ArrowKeys) {
			switch WASDDirection {
			case RightDirection:
				playerAnimPlayer.SetState("right")
			case LeftDirection:
				playerAnimPlayer.SetState("left")
			case UpDirection:
				playerAnimPlayer.SetState("up")
			case DownDirection:
				playerAnimPlayer.SetState("down")
			}
		}

		if inventory.Bombs > 0 {

			// Bomba bırak
			if inpututil.IsKeyJustPressed(eb.KeyE) {
				pos := playerPos.Add(cm.Vec2{ss.distance, 0})
				arche.DefaultBomb(world, space, pos)
				inventory.Bombs -= 1
			}

		}
	}

	// Explode all bombs
	if inpututil.IsKeyJustPressed(eb.KeyQ) {
		component.BombTagComp.Each(world, func(e *donburi.Entry) {
			Explode(e)
		})
	}

	// AI on/off
	if inpututil.IsKeyJustPressed(eb.KeyC) {
		component.AIComp.Each(world, func(e *donburi.Entry) {
			ai := component.AIComp.Get(e)
			ai.Follow = !ai.Follow
		})

	}

	// remove all enemies
	if inpututil.IsKeyJustPressed(eb.KeyBackspace) {
		component.EnemyTagComp.Each(world, func(e *donburi.Entry) {
			DestroyEntryWithBody(e)
		})
	}

	ss.bulletTimer.Update()
}

func (ds *PlayerControlSystem) Draw(world donburi.World, space *cm.Space, screen *eb.Image) {
}

func playerVelocityFunc(body *cm.Body, gravity cm.Vec2, damping, dt float64) {

	entry, ok := body.UserData.(*donburi.Entry)

	if ok {
		if entry.Valid() {
			speed := component.SpeedComp.Get(entry)
			accel := component.AccelComp.Get(entry)
			WASDAxisVector2 := WASDDirection.Normalize().Mult(*speed)
			body.SetVelocityVector(body.Velocity().LerpDistance(WASDAxisVector2, *accel))
		}
	}
}

func UpdateWASDDirection() {
	WASDDirection = cm.Vec2{}
	if inpututil.IsKeyJustPressed(eb.KeyW) {
		WASDDirectionTEMP.Y = 1
	}
	if inpututil.IsKeyJustPressed(eb.KeyS) {
		WASDDirectionTEMP.Y = -1
	}
	if inpututil.IsKeyJustPressed(eb.KeyA) {
		WASDDirectionTEMP.X = -1
	}
	if inpututil.IsKeyJustPressed(eb.KeyD) {
		WASDDirectionTEMP.X = 1
	}

	if inpututil.IsKeyJustReleased(eb.KeyW) && WASDDirectionTEMP.Y > 0 {
		fmt.Println("W")
		WASDDirectionTEMP.Y = 0
	}
	if inpututil.IsKeyJustReleased(eb.KeyS) && WASDDirectionTEMP.Y < 0 {
		WASDDirectionTEMP.Y = 0
	}
	if inpututil.IsKeyJustReleased(eb.KeyA) && WASDDirectionTEMP.X < 0 {
		WASDDirectionTEMP.X = 0
	}
	if inpututil.IsKeyJustReleased(eb.KeyD) && WASDDirectionTEMP.X > 0 {
		WASDDirectionTEMP.X = 0
	}

	WASDDirection = WASDDirectionTEMP
}

func AnyKeyDown(keys []eb.Key) bool {
	for _, key := range keys {
		if eb.IsKeyPressed(key) {
			return true
		}
	}
	return false
}

func IsPressedAndNotABC(onlyKey, a, b, c eb.Key) bool {
	return eb.IsKeyPressed(onlyKey) && !eb.IsKeyPressed(a) && !eb.IsKeyPressed(b) && !eb.IsKeyPressed(c)
}
