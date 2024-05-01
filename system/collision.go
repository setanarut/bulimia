package system

import (
	"bulimia/arche"
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type PhysicsSystem struct {
	world donburi.World
	space *cm.Space
	DT    float64
}

func NewPhysicsSystem(w donburi.World) *PhysicsSystem {
	return &PhysicsSystem{
		world: w,
		DT:    1.0 / 60.0,
	}
}

func (ps *PhysicsSystem) Init(world donburi.World, space *cm.Space, ScreenBox *cm.BB) {

	ps.space = space
	space.UseSpatialHash(64, 200)
	space.CollisionBias = math.Pow(0.3, 60)
	space.CollisionSlop = 0.5
	space.Damping = 0.01

	// Player
	space.NewCollisionHandler(arche.CollisionTypePlayer, arche.CollisionTypeDoor).BeginFunc = PlayerDoorEnter
	space.NewCollisionHandler(arche.CollisionTypePlayer, arche.CollisionTypeDoor).SeparateFunc = PlayerDoorExit
	space.NewCollisionHandler(arche.CollisionTypePlayer, arche.CollisionTypeCollectible).BeginFunc = ps.PlayerCollectibleCollisionBegin

	// Food
	space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeEnemy).BeginFunc = FoodEnemyCollisionBegin
	space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeBomb).BeginFunc = FoodBombCollisionBegin
	space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeWall).BeginFunc = FoodWallCollisionBegin
	space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeDoor).BeginFunc = FoodDoorCollisionBegin
	// space.NewCollisionHandler(arche.CollisionTypeFood, arche.CollisionTypeCollectible).BeginFunc = FoodCollectibleCollisionBegin

	space.Step(ps.DT)
}

func (ps *PhysicsSystem) Update(world donburi.World, space *cm.Space) {

	component.FoodTagComp.Each(world, func(e *donburi.Entry) {
		b := component.BodyComp.Get(e)
		if engine.IsMoving(b.Velocity(), 80) {
			DestroyBodyWithEntry(b)
		}
	})

	if pla, ok := component.PlayerTagComp.First(world); ok {

		component.EnemyTagComp.Each(world, func(e *donburi.Entry) {

			playerBody := component.BodyComp.Get(pla)
			ene := component.BodyComp.Get(e)
			ai := *component.AIComp.Get(e)

			if ai.Follow {
				dist := playerBody.Position().Distance(ene.Position())
				if dist < ai.FollowDistance {
					a := playerBody.Position().Sub(ene.Position()).Normalize().Mult(ai.FollowSpeed)
					ene.ApplyForceAtLocalPoint(a, ene.CenterOfGravity())
				}
			}

		})

	}
	space.Step(ps.DT)
}

func (ps *PhysicsSystem) Draw(world donburi.World, space *cm.Space, screen *ebiten.Image) {}

// Player <-> Collectible
func (ps *PhysicsSystem) PlayerCollectibleCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	playerBody, bodyCollectible := arb.Bodies()
	playerEntry := playerBody.UserData.(*donburi.Entry)
	collectibleEntry := bodyCollectible.UserData.(*donburi.Entry)

	inventory := component.InventoryComp.Get(playerEntry)
	collectibleComponent := component.CollectibleComp.Get(collectibleEntry)

	if collectibleComponent.Type == component.Food {
		inventory.Foods += collectibleComponent.ItemCount
	}
	if collectibleComponent.Type == component.Bomb {
		inventory.Bombs += collectibleComponent.ItemCount
	}

	if collectibleComponent.Type == component.Key {
		// oyuncu anahtara sahip değilse ekle
		keyNum := collectibleComponent.KeyNumber
		if !slices.Contains(inventory.Keys, keyNum) {
			inventory.Keys = append(inventory.Keys, keyNum)
		}
		component.DoorComp.Each(ps.world, func(e *donburi.Entry) {
			door := component.DoorComp.Get(e)
			if door.LockNumber == keyNum {
				door.PlayerHasKey = true
			}

		})
	}

	DestroyBodyWithEntry(bodyCollectible)

	return false
}

func Explode(bomb *donburi.Entry) {
	body := component.BodyComp.Get(bomb)

	space := body.FirstShape().Space()
	pos := body.Position()
	component.EnemyTagComp.Each(bomb.World, func(enemy *donburi.Entry) {
		body := component.BodyComp.Get(enemy)
		queryInfo := space.SegmentQueryFirst(pos, body.Position(), 0, arche.FilterBombRaycast)
		contactShape := queryInfo.Shape
		if contactShape != nil {
			if contactShape.Body() == body {
				ApplyRaycastImpulse(queryInfo, 1000)
				damage := int(engine.MapRange(queryInfo.Alpha, 0.5, 1, 200, 0))
				*component.HealthComp.Get(enemy) -= damage
				if *component.HealthComp.Get(enemy) < 0 {
					DestroyEntryWithBody(enemy)
				}

			}
		}

	})

	DestroyEntryWithBody(bomb)
}

// Player <-> Door (enter)
func PlayerDoorEnter(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	playerBody, doorBody := arb.Bodies()

	doorEntry := doorBody.UserData.(*donburi.Entry)
	playerEntry := playerBody.UserData.(*donburi.Entry)
	door := component.DoorComp.Get(doorEntry)
	inv := component.InventoryComp.Get(playerEntry)

	if slices.Contains(inv.Keys, door.LockNumber) {
		door.Open = true
		doorBody.FirstShape().SetSensor(true)
	}
	return true
}

// Player <-> Door (exit)
func PlayerDoorExit(arb *cm.Arbiter, space *cm.Space, userData interface{}) {
	playerBody, doorBody := arb.Bodies()
	doorEntry := doorBody.UserData.(*donburi.Entry)
	d := component.DoorComp.Get(doorEntry)
	d.Open = false
	doorBody.FirstShape().SetSensor(false)

	for _, room := range Rooms {
		if room.ContainsVect(playerBody.Position()) {
			CurrentRoom = room
		}
	}

}

// Food <-> Enemy
func FoodEnemyCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	bulletBody, enemyBody := arb.Bodies()
	bulletEntry := bulletBody.UserData.(*donburi.Entry)
	enemyEntry := enemyBody.UserData.(*donburi.Entry)

	if enemyEntry.Valid() {

		if enemyEntry.HasComponent(component.HealthComp) {
			health := component.HealthComp.Get(enemyEntry)

			if bulletEntry.Valid() {
				*health -= *component.DamageComp.Get(bulletEntry)
			}

			if *health < 0 {
				DestroyBodyWithEntry(enemyBody)
			}
		}
	}

	// çarpan bulletı yok et
	DestroyEntryWithBody(bulletEntry)
	return true
}

// Food <-> Wall
func FoodWallCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	Food, _ := arb.Bodies()
	DestroyBodyWithEntry(Food)
	return false
}

// Food <-> Bomb
func FoodBombCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	food, bomb := arb.Bodies()
	DestroyBodyWithEntry(food)
	Explode(bomb.UserData.(*donburi.Entry))
	return false
}

// Food <-> Collectible
func FoodCollectibleCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	arb.Ignore()
	return false
}

// Food <-> Door
func FoodDoorCollisionBegin(arb *cm.Arbiter, space *cm.Space, userData interface{}) bool {
	bodyA, _ := arb.Bodies()
	bulletEntry := bodyA.UserData.(*donburi.Entry)
	DestroyEntryWithBody(bulletEntry)
	return true
}

func ApplyRaycastImpulse(sqi cm.SegmentQueryInfo, power float64) {
	impulseVec2 := sqi.Normal.Neg().Mult(power * engine.MapRange(sqi.Alpha, 0.5, 1, 1, 0))
	sqi.Shape.Body().ApplyImpulseAtWorldPoint(impulseVec2, sqi.Point)
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
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
		if entry.HasComponent(component.BodyComp) {
			body := component.BodyComp.Get(entry)
			DestroyBodyWithEntry(body)
		}
	}
}
