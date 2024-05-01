package arche

import (
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"
	"math/rand/v2"

	"github.com/yohamta/donburi"
)

func DefaultEnemy(world donburi.World, space *cm.Space, pos cm.Vec2) {
	NewEnemyEntity(1, 0.5, 0, 20, 500, world, space, pos)
}

func DefaultBomb(world donburi.World, space *cm.Space, pos cm.Vec2) {
	NewBombEntity(1, 0.1, 0, 20, world, space, pos)
}
func DefaultFoodCollectible(world donburi.World, space *cm.Space, pos cm.Vec2) {
	NewCollectibleEntity(component.Food, 1, -1, 5, world, space, pos)
}

func DefaultKeyCollectible(keyNumber int, world donburi.World, space *cm.Space, pos cm.Vec2) {
	NewCollectibleEntity(component.Key, 1, keyNumber, 10, world, space, pos)
}

func RandomCollectible(world donburi.World, space *cm.Space, pos cm.Vec2) {
	randomType := component.ItemType(rand.IntN(3))
	NewCollectibleEntity(randomType, 1, engine.RandRangeInt(1, 10), 10, world, space, pos)
}

func RandomKeyCollectible(world donburi.World, space *cm.Space, pos cm.Vec2) {
	NewCollectibleEntity(component.Key, 1, engine.RandRangeInt(1, 10), 10, world, space, pos)
}
