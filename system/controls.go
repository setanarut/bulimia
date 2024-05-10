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
	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {
		body := comp.Body.Get(playerEntry)
		sys.distance = body.FirstShape().Class.(*cm.Circle).Radius()

	}
}

func (sys *PlayerControlSystem) Update() {

	// Input.UpdateJustArrowDirection()
	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {

		inventory := comp.Inventory.Get(playerEntry)
		playerBody := comp.Body.Get(playerEntry)
		playerRenderData := comp.Render.Get(playerEntry)
		playerPos := playerBody.Position()

		if inventory.Foods > 0 {

			switch res.Input.ArrowDirection {

			case engine.RightDirection:
				playerRenderData.AnimPlayer.SetState("shootR")
				sys.shoot(playerPos)

			case engine.LeftDirection:
				playerRenderData.AnimPlayer.SetState("shootL")
				sys.shoot(playerPos)

			case engine.UpDirection:
				playerRenderData.AnimPlayer.SetState("shootU")
				sys.shoot(playerPos)

			case engine.DownDirection:
				playerRenderData.AnimPlayer.SetState("shootD")
				sys.shoot(playerPos)

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

		if res.Input.ArrowDirection.Equal(engine.NoDirection) {

			switch res.Input.WASDDirection {

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

		if inventory.Bombs > 0 {

			// Bomba bırak
			if inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) {
				bombPos := res.Input.LastPressedDirection.Neg().Mult(sys.distance)
				arche.SpawnDefaultBomb(playerPos.Add(bombPos))
				inventory.Bombs -= 1
			}

		}

	}

	// Explode all bombs
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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

	sys.bulletTimer.Update()
}

func (sys *PlayerControlSystem) Draw() {
}

func (sys *PlayerControlSystem) shoot(pos cm.Vec2) {
	if sys.bulletTimer.IsReady() {
		sys.bulletTimer.Reset()
	}
	if sys.bulletTimer.IsStart() {
		bullet := arche.SpawnFood(0.1, 0, 0.5, sys.bulletRadius, pos)
		bulletBody := comp.Body.Get(bullet)
		bulletBody.ApplyImpulseAtWorldPoint(res.Input.ArrowDirection.Mult(100), pos)
	}
}
