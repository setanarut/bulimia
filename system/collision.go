package system

import (
	"bulimia/arche"
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"bulimia/res"
	"math"
	"slices"

	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

type CollisionSystem struct {
	DT float64
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		DT: 1.0 / 60.0,
	}
}

func (ps *CollisionSystem) Init() {
	res.Space.UseSpatialHash(50, 1000)
	res.Space.CollisionBias = math.Pow(0.3, 60)
	res.Space.CollisionSlop = 0.5
	res.Space.Damping = 0.03
	res.Space.Iterations = 1

	// Player
	res.Space.NewCollisionHandler(arche.CollisionTypePlayer, arche.CollisionTypeDoor).BeginFunc = playerDoorEnter
	res.Space.NewCollisionHandler(arche.CollisionTypePlayer, arche.CollisionTypeDoor).SeparateFunc = playerDoorExit
	res.Space.NewCollisionHandler(arche.CollisionTypePlayer, arche.CollisionTypeCollectible).BeginFunc = playerCollectibleCollisionBegin

	// Enemy
	res.Space.NewCollisionHandler(arche.CollisionTypeEnemy, arche.CollisionTypePlayer).PostSolveFunc = enemyPlayerPostSolve
	res.Space.NewCollisionHandler(arche.CollisionTypeEnemy, arche.CollisionTypePlayer).SeparateFunc = enemyPlayerSep

	// Food
	res.Space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeEnemy).BeginFunc = foodEnemyCollisionBegin
	res.Space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeBomb).BeginFunc = foodBombCollisionBegin
	res.Space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeWall).BeginFunc = foodWallCollisionBegin
	res.Space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeDoor).BeginFunc = foodDoorCollisionBegin

	res.Space.Step(ps.DT)

}

func (ps *CollisionSystem) Update() {

	comp.FoodTag.Each(res.World, func(e *donburi.Entry) {
		b := comp.Body.Get(e)
		if engine.IsMoving(b.Velocity(), 80) {
			DestroyBodyWithEntry(b)
		}
	})

	if pla, ok := comp.PlayerTag.First(res.World); ok {
		playerBody := comp.Body.Get(pla)

		comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {

			ene := comp.Body.Get(e)
			ai := *comp.AI.Get(e)
			livingData := comp.Char.Get(e)

			if ai.Follow {
				dist := playerBody.Position().Distance(ene.Position())
				if dist < ai.FollowDistance {
					speed := ene.Mass() * (livingData.Speed * 4)
					a := playerBody.Position().Sub(ene.Position()).Normalize().Mult(speed)
					ene.ApplyForceAtLocalPoint(a, ene.CenterOfGravity())
				}
			}

		})
		comp.Collectible.Each(res.World, func(e *donburi.Entry) {

			ene := comp.Body.Get(e)
			dist := playerBody.Position().Distance(ene.Position())

			if dist < 80 {
				speed := engine.MapRange(dist, 500, 0, 0, 1000)
				a := playerBody.Position().Sub(ene.Position()).Normalize().Mult(speed)
				ene.ApplyForceAtLocalPoint(a, ene.CenterOfGravity())
			}

		})

	}
	res.Space.Step(ps.DT)
}

func (ps *CollisionSystem) Draw() {}

// Player <-> Collectible
func playerCollectibleCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	playerBody, bodyCollectible := arb.Bodies()
	playerEntry, pok := playerBody.UserData.(*donburi.Entry)
	collectibleEntry, cok := bodyCollectible.UserData.(*donburi.Entry)

	if pok && cok {

		if playerEntry.Valid() &&
			collectibleEntry.Valid() &&
			collectibleEntry.HasComponent(comp.Collectible) &&
			playerEntry.HasComponent(comp.Inventory) {

			inventory := comp.Inventory.Get(playerEntry)
			collectibleComponent := comp.Collectible.Get(collectibleEntry)

			if collectibleComponent.Type == comp.Food {
				inventory.Foods += collectibleComponent.ItemCount
			}
			if collectibleComponent.Type == comp.Bomb {
				inventory.Bombs += collectibleComponent.ItemCount
			}
			if collectibleComponent.Type == comp.EmeticDrug {
				inventory.EmeticDrug += collectibleComponent.ItemCount

			}

			if collectibleComponent.Type == comp.Key {
				// oyuncu anahtara sahip değilse ekle
				keyNum := collectibleComponent.KeyNumber
				if !slices.Contains(inventory.Keys, keyNum) {
					inventory.Keys = append(inventory.Keys, keyNum)
				}

				comp.Door.Each(res.World, func(e *donburi.Entry) {
					door := comp.Door.Get(e)
					if door.LockNumber == keyNum {
						door.PlayerHasKey = true
					}

				})
			}

			DestroyBodyWithEntry(bodyCollectible)
		}
	}

	return false
}

// Player <-> Door (enter)
func playerDoorEnter(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	playerBody, doorBody := arb.Bodies()

	doorEntry := doorBody.UserData.(*donburi.Entry)
	playerEntry := playerBody.UserData.(*donburi.Entry)
	door := comp.Door.Get(doorEntry)
	inv := comp.Inventory.Get(playerEntry)

	if slices.Contains(inv.Keys, door.LockNumber) {
		door.Open = true
		doorBody.FirstShape().SetSensor(true)
	}
	return true
}

// Player <-> Door (exit)
func playerDoorExit(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	playerBody, doorBody := arb.Bodies()
	doorEntry := doorBody.UserData.(*donburi.Entry)
	d := comp.Door.Get(doorEntry)
	d.Open = false
	doorBody.FirstShape().SetSensor(false)

	for _, room := range res.Rooms {
		if room.ContainsVect(playerBody.Position()) {
			res.CurrentRoom = room
		}
	}

}

// Enemy <-> Player
func enemyPlayerPostSolve(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	enemyBody, playerBody := arb.Bodies()
	enemyEntry, eok := enemyBody.UserData.(*donburi.Entry)
	playerEntry, pok := playerBody.UserData.(*donburi.Entry)
	var livingData *comp.CharacterData
	if eok && pok {

		if playerEntry.Valid() && enemyEntry.Valid() {
			if playerEntry.HasComponent(comp.Char) && enemyEntry.HasComponent(comp.Damage) && playerEntry.HasComponent(comp.Render) {
				livingData = comp.Char.Get(playerEntry)
				comp.Render.Get(playerEntry).ScaleColor = colornames.Red
				livingData.Health -= *comp.Damage.Get(enemyEntry)
				// livingData.Health -= donburi.GetValue[float64](enemyEntry, comp.Damage)
				if livingData.Health < 0 {
					DestroyBodyWithEntry(playerBody)
				}
			}
		}

	}

}

// Enemy <-> Player Sep
func enemyPlayerSep(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	_, playerBody := arb.Bodies()
	playerEntry := playerBody.UserData.(*donburi.Entry)
	if playerEntry.Valid() {
		comp.Render.Get(playerEntry).ScaleColor = colornames.Yellow
	}

}

// Food <-> Enemy
func foodEnemyCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	bulletBody, enemyBody := arb.Bodies()
	bulletEntry := bulletBody.UserData.(*donburi.Entry)
	enemyEntry := enemyBody.UserData.(*donburi.Entry)

	if enemyEntry.Valid() {

		if enemyEntry.HasComponent(comp.Char) {
			livingData := comp.Char.Get(enemyEntry)

			if bulletEntry.Valid() {
				livingData.Health -= *comp.Damage.Get(bulletEntry)
			}

			if livingData.Health < 0 {
				DestroyBodyWithEntry(enemyBody)
			}
		}
	}

	// çarpan bulletı yok et
	DestroyEntryWithBody(bulletEntry)
	return true
}

// Food <-> Wall
func foodWallCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	Food, _ := arb.Bodies()
	DestroyBodyWithEntry(Food)
	return false
}

// Food <-> Bomb
func foodBombCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	food, bomb := arb.Bodies()
	DestroyBodyWithEntry(food)
	Explode(bomb.UserData.(*donburi.Entry))
	return false
}

// Food <-> Door
func foodDoorCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	bodyA, _ := arb.Bodies()
	bulletEntry := bodyA.UserData.(*donburi.Entry)
	DestroyEntryWithBody(bulletEntry)
	return true
}

func DestroyBodyWithEntry(b *cm.Body) {
	s := b.FirstShape().Space()
	if s.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		s.AddPostStepCallback(removeBodyPostStep, b, false)
	}
}
func DestroyEntryWithBody(entry *donburi.Entry) {
	if entry.Valid() {
		if entry.HasComponent(comp.Body) {
			body := comp.Body.Get(entry)
			DestroyBodyWithEntry(body)
		}
	}
}

func Explode(bomb *donburi.Entry) {
	bombBody := comp.Body.Get(bomb)
	space := bombBody.FirstShape().Space()

	comp.EnemyTag.Each(bomb.World, func(enemy *donburi.Entry) {

		livingData := comp.Char.Get(enemy)
		enemyBody := comp.Body.Get(enemy)

		queryInfo := space.SegmentQueryFirst(bombBody.Position(), enemyBody.Position(), 0, arche.FilterBombRaycast)
		contactShape := queryInfo.Shape

		if contactShape != nil {
			if contactShape.Body() == enemyBody {
				ApplyRaycastImpulse(queryInfo, 1000)
				damage := engine.MapRange(queryInfo.Alpha, 0.5, 1, 200, 0)
				livingData.Health -= damage
				if livingData.Health < 0 {
					DestroyEntryWithBody(enemy)
				}

			}
		}

	})
	res.Camera.AddTrauma(0.2)
	DestroyEntryWithBody(bomb)
}

func ApplyRaycastImpulse(sqi cm.SegmentQueryInfo, power float64) {
	impulseVec2 := sqi.Normal.Neg().Mult(power * engine.MapRange(sqi.Alpha, 0.5, 1, 1, 0))
	sqi.Shape.Body().ApplyImpulseAtWorldPoint(impulseVec2, sqi.Point)
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}
