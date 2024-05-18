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

var bombDistance float64 = 40

type PlayerControlSystem struct {
}

func NewPlayerControlSystem() *PlayerControlSystem {
	return &PlayerControlSystem{}
}

func (sys *PlayerControlSystem) Init() {
	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {
		playerBody := comp.Body.Get(playerEntry)
		playerBody.SetVelocityUpdateFunc(PlayerVelocityFunc)
	}

}

func (sys *PlayerControlSystem) Update() {

	// Input.UpdateJustArrowDirection()
	res.Input.UpdateArrowDirection()
	res.Input.UpdateWASDDirection()

	if playerEntry, ok := comp.PlayerTag.First(res.World); ok {

		charData := comp.Char.Get(playerEntry)
		inventory := comp.Inventory.Get(playerEntry)
		playerBody := comp.Body.Get(playerEntry)
		playerRenderData := comp.Render.Get(playerEntry)
		playerPos := playerBody.Position()

		if playerEntry.HasComponent(comp.DrugEffect) {

			drugEffectData := comp.DrugEffect.Get(playerEntry)

			if drugEffectData.EffectTimer.IsStart() {
				AddDrugEffect(charData, drugEffectData)

			}

			if drugEffectData.EffectTimer.IsReady() {
				RemoveDrugEffect(charData, drugEffectData)
				playerEntry.RemoveComponent(comp.DrugEffect)
			}

			drugEffectData.EffectTimer.Update()
		}

		if !res.Input.ArrowDirection.Equal(engine.NoDirection) {
			playerRenderData.AnimPlayer.SetState("shoot")
			playerRenderData.DrawAngle = res.Input.ArrowDirection.ToAngle()

			// SHOOTING
			if inventory.Foods > 0 {
				if charData.VomitCooldownTimer.IsReady() {
					charData.VomitCooldownTimer.Reset()
				}

				if charData.VomitCooldownTimer.IsStart() {
					for range charData.FoodPerCooldown {
						if inventory.Foods > 0 {
							inventory.Foods -= 1
						}
						dir := engine.Rotate(res.Input.ArrowDirection.Mult(1000), engine.RandRange(0.2, -0.2))
						bullet := arche.SpawnDefaultFood(playerPos)
						bulletBody := comp.Body.Get(bullet)
						bulletBody.ApplyImpulseAtWorldPoint(dir, playerPos)
					}
				}
			}

		} else {
			playerRenderData.AnimPlayer.SetState("right")
			playerRenderData.DrawAngle = res.Input.LastPressedDirection.ToAngle()

		}

		charData.VomitCooldownTimer.Update()

		if inventory.Bombs > 0 {

			// Bomba bırak
			if inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) {
				bombPos := res.Input.LastPressedDirection.Neg().Mult(bombDistance)
				arche.SpawnDefaultBomb(playerPos.Add(bombPos))
				inventory.Bombs -= 1
			}

		}

		// ilac kullan
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

			if inventory.EmeticDrug > 0 {
				if !playerEntry.HasComponent(comp.DrugEffect) {
					playerEntry.AddComponent(comp.DrugEffect)
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

func AddDrugEffect(charData *comp.CharacterData, drugEffectData *comp.DrugEffectData) {
	charData.FoodPerCooldown += drugEffectData.ExtraVomit
	charData.VomitCooldownTimer.Target += drugEffectData.VomitCooldown
	charData.Speed += drugEffectData.AddMovementSpeed
}
func RemoveDrugEffect(charData *comp.CharacterData, drugEffectData *comp.DrugEffectData) {
	charData.FoodPerCooldown -= drugEffectData.ExtraVomit
	charData.VomitCooldownTimer.Target -= drugEffectData.VomitCooldown
	charData.Speed -= drugEffectData.AddMovementSpeed
}

func PlayerVelocityFunc(body *cm.Body, gravity cm.Vec2, damping float64, dt float64) {

	entry, ok := body.UserData.(*donburi.Entry)

	if ok {
		if entry.Valid() {
			livingData := comp.Char.Get(entry)
			WASDAxisVector := res.Input.WASDDirection.Normalize().Mult(livingData.Speed)
			body.SetVelocityVector(body.Velocity().LerpDistance(WASDAxisVector, livingData.Accel))
		}
	}
}
