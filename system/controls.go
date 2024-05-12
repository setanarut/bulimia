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
	bombDistance float64
	bulletTimer  *engine.Timer
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{
		bulletTimer: engine.NewTimer(time.Second / 4),
	}
}

func (sys *PlayerControlSystem) Init() {
	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {
		body := comp.Body.Get(playerEntry)
		sys.bombDistance = body.FirstShape().Class.(*cm.Circle).Radius()

	}
}

func (sys *PlayerControlSystem) Update() {

	// Input.UpdateJustArrowDirection()
	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {

		inventory := comp.Inventory.Get(playerEntry)
		emetic := comp.Living.Get(playerEntry).Emetic
		playerBody := comp.Body.Get(playerEntry)
		playerRenderData := comp.Render.Get(playerEntry)
		playerPos := playerBody.Position()

		if inventory.Foods > 0 {

			switch res.Input.ArrowDirection {

			case engine.RightDirection:
				playerRenderData.AnimPlayer.SetState("shootR")
				sys.shoot(playerPos, emetic)
			case engine.LeftDirection:
				playerRenderData.AnimPlayer.SetState("shootL")
				sys.shoot(playerPos, emetic)
			case engine.UpDirection:
				playerRenderData.AnimPlayer.SetState("shootU")
				sys.shoot(playerPos, emetic)
			case engine.DownDirection:
				playerRenderData.AnimPlayer.SetState("shootD")
				sys.shoot(playerPos, emetic)

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
				bombPos := res.Input.LastPressedDirection.Neg().Mult(sys.bombDistance)
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
	// disable Emetic
	if inpututil.IsKeyJustPressed(ebiten.Key0) {
		player, ok := comp.PlayerTag.First(res.World)
		if ok {
			comp.Living.Get(player).Emetic = !comp.Living.Get(player).Emetic
			if comp.Living.Get(player).Emetic {
				sys.bulletTimer.SetDuration(time.Second / 30)
			} else {
				sys.bulletTimer.SetDuration(time.Second / 4)
			}

		}

	}

	sys.bulletTimer.Update()
}

func (sys *PlayerControlSystem) Draw() {
}

func (sys *PlayerControlSystem) shoot(pos cm.Vec2, emetic bool) {

	if sys.bulletTimer.IsReady() {
		sys.bulletTimer.Reset()
	}

	if sys.bulletTimer.IsStart() {
		var dir cm.Vec2

		if emetic {
			sys.bulletTimer.SetDuration(time.Second / 30)
			for range 10 {
				dir = engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
				bullet := arche.SpawnDefaultFood(pos)
				bulletBody := comp.Body.Get(bullet)
				bulletBody.ApplyImpulseAtWorldPoint(dir, pos)
			}
		} else {
			for range 1 {
				dir = engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.1, -0.1))
				bullet := arche.SpawnDefaultFood(pos)
				bulletBody := comp.Body.Get(bullet)
				bulletBody.ApplyImpulseAtWorldPoint(dir, pos)
			}
		}

	}

}
