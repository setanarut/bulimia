package system

import (
	"bulimia/arche"
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var Input *engine.InputManager = &engine.InputManager{}

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

func (sys *PlayerControlSystem) Init(world donburi.World, space *cm.Space, ScreenBox *cm.BB) {
	// set SetVelocityUpdateFunc
	if playerEntry, ok := component.PlayerTagComp.First(world); ok {
		body := component.BodyComp.Get(playerEntry)
		sys.distance = body.FirstShape().Class.(*cm.Circle).Radius()
		body.SetVelocityUpdateFunc(sys.playerVelocityFunc)
	}
}

func (sys *PlayerControlSystem) Update(world donburi.World, space *cm.Space) {

	// Input.UpdateJustArrowDirection()
	Input.UpdateArrowDirection()
	Input.UpdateWASDDirection()

	if playerEntry, ok := component.PlayerTagComp.First(world); ok {
		inventory := component.InventoryComp.Get(playerEntry)
		playerBody := component.BodyComp.Get(playerEntry)
		playerAnimPlayer := component.AnimPlayerComp.Get(playerEntry)
		playerPos := playerBody.Position()

		if Input.ArrowDirection.Equal(engine.RightDirection) {
			playerAnimPlayer.SetState("shootR")

			if inventory.Foods > 0 {

				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}

				if sys.bulletTimer.IsStart() {
					bullet := arche.NewFoodEntity(0.1, 0, 0.5, sys.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{100, 0}, playerPos)
				}

			}
		}

		if Input.ArrowDirection.Equal(engine.LeftDirection) {
			playerAnimPlayer.SetState("shootL")

			if inventory.Foods > 0 {

				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}

				if sys.bulletTimer.IsStart() {

					bullet := arche.NewFoodEntity(0.1, 0, 0.5, sys.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{-100, 0}, playerPos)
				}
			}
		}

		if Input.ArrowDirection.Equal(engine.UpDirection) {
			playerAnimPlayer.SetState("shootU")

			if inventory.Foods > 0 {
				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}
				if sys.bulletTimer.IsStart() {
					bullet := arche.NewFoodEntity(0.1, 0, 0.5, sys.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{0, 100}, playerPos)
				}

			}

		}

		if Input.ArrowDirection == engine.DownDirection {
			playerAnimPlayer.SetState("shootD")

			if inventory.Foods > 0 {

				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}
				if sys.bulletTimer.IsStart() {
					bullet := arche.NewFoodEntity(0.1, 0, 0.5, sys.bulletRadius, world, space, playerPos)
					inventory.Foods -= 1
					bulletBody := component.BodyComp.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{0, -100}, playerPos)
				}

			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
			playerAnimPlayer.SetState("up")
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
			playerAnimPlayer.SetState("down")
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
			playerAnimPlayer.SetState("left")
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
			playerAnimPlayer.SetState("right")
		}

		if Input.ArrowDirection.Equal(engine.NoDirection) {

			switch Input.WASDDirection {

			case engine.RightDirection:
				playerAnimPlayer.SetState("right")

			case engine.LeftDirection:
				playerAnimPlayer.SetState("left")

			case engine.UpDirection:
				playerAnimPlayer.SetState("up")

			case engine.DownDirection:
				playerAnimPlayer.SetState("down")

			}

		}

		if inventory.Bombs > 0 {

			// Bomba bırak
			if inpututil.IsKeyJustPressed(ebiten.KeyE) {
				pos := playerPos.Add(cm.Vec2{sys.distance, 0})
				arche.DefaultBomb(world, space, pos)
				inventory.Bombs -= 1
			}

		}

	}

	// Explode all bombs
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		component.BombTagComp.Each(world, func(e *donburi.Entry) {
			Explode(e)
		})
	}

	// AI on/off
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		component.AIComp.Each(world, func(e *donburi.Entry) {
			ai := component.AIComp.Get(e)
			ai.Follow = !ai.Follow
		})

	}

	// remove all enemies
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		component.EnemyTagComp.Each(world, func(e *donburi.Entry) {
			DestroyEntryWithBody(e)
		})
	}

	sys.bulletTimer.Update()
}

func (sys *PlayerControlSystem) Draw(world donburi.World, space *cm.Space, screen *ebiten.Image) {
}

func (sys *PlayerControlSystem) playerVelocityFunc(body *cm.Body, gravity cm.Vec2, damping, dt float64) {

	entry, ok := body.UserData.(*donburi.Entry)

	if ok {
		if entry.Valid() {
			speed := component.SpeedComp.Get(entry)
			accel := component.AccelComp.Get(entry)
			WASDAxisVector2 := Input.WASDDirection.Normalize().Mult(*speed)
			body.SetVelocityVector(body.Velocity().LerpDistance(WASDAxisVector2, *accel))
		}
	}
}
