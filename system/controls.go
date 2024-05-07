package system

import (
	"bulimia/arche"
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"bulimia/res"
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

func (sys *PlayerControlSystem) Init() {
	// set SetVelocityUpdateFunc
	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {
		body := comp.Body.Get(playerEntry)
		sys.distance = body.FirstShape().Class.(*cm.Circle).Radius()
		body.SetVelocityUpdateFunc(sys.playerVelocityFunc)
	}
}

func (sys *PlayerControlSystem) Update() {

	// Input.UpdateJustArrowDirection()
	Input.UpdateArrowDirection()
	Input.UpdateWASDDirection()

	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {
		inventory := comp.Inventory.Get(playerEntry)
		playerBody := comp.Body.Get(playerEntry)
		playerRenderData := comp.Render.Get(playerEntry)

		playerPos := playerBody.Position()

		if Input.ArrowDirection.Equal(engine.RightDirection) && !ebiten.IsKeyPressed(ebiten.KeySpace) {
			playerRenderData.AnimPlayer.SetState("shootR")

			if inventory.Foods > 0 {

				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}

				if sys.bulletTimer.IsStart() {
					bullet := arche.SpawnFood(0.1, 0, 0.5, sys.bulletRadius, playerPos)
					inventory.Foods -= 1
					bulletBody := comp.Body.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{100, 0}, playerPos)
				}

			}
		}

		if Input.ArrowDirection.Equal(engine.LeftDirection) && !ebiten.IsKeyPressed(ebiten.KeySpace) {
			playerRenderData.AnimPlayer.SetState("shootL")

			if inventory.Foods > 0 {

				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}

				if sys.bulletTimer.IsStart() {

					bullet := arche.SpawnFood(0.1, 0, 0.5, sys.bulletRadius, playerPos)
					inventory.Foods -= 1
					bulletBody := comp.Body.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{-100, 0}, playerPos)
				}
			}
		}

		if Input.ArrowDirection.Equal(engine.UpDirection) && !ebiten.IsKeyPressed(ebiten.KeySpace) {
			playerRenderData.AnimPlayer.SetState("shootU")

			if inventory.Foods > 0 {
				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}
				if sys.bulletTimer.IsStart() {
					bullet := arche.SpawnFood(0.1, 0, 0.5, sys.bulletRadius, playerPos)
					inventory.Foods -= 1
					bulletBody := comp.Body.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{0, 100}, playerPos)
				}

			}

		}

		if Input.ArrowDirection == engine.DownDirection && !ebiten.IsKeyPressed(ebiten.KeySpace) {
			playerRenderData.AnimPlayer.SetState("shootD")

			if inventory.Foods > 0 {

				if sys.bulletTimer.IsReady() {
					sys.bulletTimer.Reset()
				}
				if sys.bulletTimer.IsStart() {
					bullet := arche.SpawnFood(0.1, 0, 0.5, sys.bulletRadius, playerPos)
					inventory.Foods -= 1
					bulletBody := comp.Body.Get(bullet)
					bulletBody.ApplyImpulseAtWorldPoint(cm.Vec2{0, -100}, playerPos)
				}

			}
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) && inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			playerRenderData.AnimPlayer.SetState("shootR")

			if inventory.Bombs > 0 {
				arche.SpawnDefaultBomb(playerPos.Add(cm.Vec2{50, 0}))
				inventory.Bombs -= 1

			}
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) && inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			playerRenderData.AnimPlayer.SetState("shootL")

			if inventory.Bombs > 0 {
				arche.SpawnDefaultBomb(playerPos.Add(cm.Vec2{-50, 0}))
				inventory.Bombs -= 1

			}
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) && inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			playerRenderData.AnimPlayer.SetState("shootU")

			if inventory.Bombs > 0 {
				arche.SpawnDefaultBomb(playerPos.Add(cm.Vec2{0, 50}))
				inventory.Bombs -= 1

			}
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) && inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			playerRenderData.AnimPlayer.SetState("shootD")

			if inventory.Bombs > 0 {
				arche.SpawnDefaultBomb(playerPos.Add(cm.Vec2{0, -50}))
				inventory.Bombs -= 1

			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
			playerRenderData.AnimPlayer.SetState("up")
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
			playerRenderData.AnimPlayer.SetState("down")
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
			playerRenderData.AnimPlayer.SetState("left")
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
			playerRenderData.AnimPlayer.SetState("right")
		}

		if Input.ArrowDirection.Equal(engine.NoDirection) {

			switch Input.WASDDirection {

			case engine.RightDirection:
				playerRenderData.AnimPlayer.SetState("right")

			case engine.LeftDirection:
				playerRenderData.AnimPlayer.SetState("left")

			case engine.UpDirection:
				playerRenderData.AnimPlayer.SetState("up")

			case engine.DownDirection:
				playerRenderData.AnimPlayer.SetState("down")

			}

		}

		// if inventory.Bombs > 0 {

		// 	// Bomba bırak
		// 	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		// 		bombPos := Input.LastPressedDirection.Mult(sys.distance)
		// 		arche.SpawnDefaultBomb(world, space, playerPos.Add(bombPos))
		// 		inventory.Bombs -= 1
		// 	}

		// }

	}

	// Explode all bombs
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		comp.BombTag.Each(res.World, func(e *donburi.Entry) {
			Explode(e)
		})
	}

	// AI on/off
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		comp.AI.Each(res.World, func(e *donburi.Entry) {
			ai := comp.AI.Get(e)
			ai.Follow = !ai.Follow
		})

	}

	// remove all enemies
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {
			DestroyEntryWithBody(e)
		})
	}

	sys.bulletTimer.Update()
}

func (sys *PlayerControlSystem) Draw() {
}

func (sys *PlayerControlSystem) playerVelocityFunc(body *cm.Body, gravity cm.Vec2, damping, dt float64) {

	entry, ok := body.UserData.(*donburi.Entry)

	if ok {
		if entry.Valid() {
			livingData := comp.Living.Get(entry)
			WASDAxisVector := Input.WASDDirection.Normalize().Mult(livingData.Speed)
			body.SetVelocityVector(body.Velocity().LerpDistance(WASDAxisVector, livingData.Accel))
		}
	}
}
