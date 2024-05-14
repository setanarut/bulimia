package system

import (
	"bulimia/arche"
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

var BombDistance float64 = 40
var HoldLivingData comp.LivingData

type PlayerControlSystem struct {
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{}
}

func (sys *PlayerControlSystem) Init() {

}

func (sys *PlayerControlSystem) Update() {

	// Input.UpdateJustArrowDirection()
	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {

		livingData := comp.Living.Get(playerEntry)
		inventory := comp.Inventory.Get(playerEntry)
		playerBody := comp.Body.Get(playerEntry)
		playerRenderData := comp.Render.Get(playerEntry)
		playerPos := playerBody.Position()

		if playerEntry.HasComponent(comp.DrugEffect) {

			de := comp.DrugEffect.Get(playerEntry)

			if de.EffectTimer.IsStart() {

				HoldLivingData = *livingData
				livingData.BulletPerCoolDown += de.ExtraBulletPerCoolDown
				livingData.ShootingCooldownTimer.SetDuration(de.ShootingCooldown)
				livingData.Speed *= de.SpeedScaleFactor

			}

			if de.EffectTimer.IsReady() {
				livingData.BulletPerCoolDown -= de.ExtraBulletPerCoolDown
				livingData.ShootingCooldownTimer.SetDuration(HoldLivingData.ShootingCooldownTimer.Duration())
				livingData.Speed = HoldLivingData.Speed
				playerEntry.RemoveComponent(comp.DrugEffect)
			}

			de.EffectTimer.Update()
		}

		if inventory.Foods > 0 {

			if !res.Input.ArrowDirection.Equal(engine.NoDirection) {

				playerRenderData.AnimPlayer.SetState("shootR")
				playerBody.SetAngle(res.Input.ArrowDirection.ToAngle())

				if livingData.ShootingCooldownTimer.IsReady() {
					livingData.ShootingCooldownTimer.Reset()
				}

				if livingData.ShootingCooldownTimer.IsStart() {
					for range livingData.BulletPerCoolDown {
						dir := engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
						bullet := arche.SpawnDefaultFood(playerPos)
						bulletBody := comp.Body.Get(bullet)
						bulletBody.ApplyImpulseAtWorldPoint(dir, playerPos)
					}

				}

			} else {
				playerBody.SetAngle(0)

			}

		}

		livingData.ShootingCooldownTimer.Update()

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
				bombPos := res.Input.LastPressedDirection.Neg().Mult(BombDistance)
				arche.SpawnDefaultBomb(playerPos.Add(bombPos))
				inventory.Bombs -= 1
			}

		}

		// ilac kullan
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

			if inventory.EmeticDrug > 0 {
				if !playerEntry.HasComponent(comp.DrugEffect) {
					playerEntry.AddComponent(comp.DrugEffect)
					// comp.DrugEffect.Get(playerEntry).EffectTimer.Reset()
					inventory.EmeticDrug -= 1
				}
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

}

func (sys *PlayerControlSystem) Draw() {
}
