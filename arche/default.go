package arche

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"math/rand/v2"

	"github.com/yohamta/donburi"
)

func SpawnDefaultEnemy(world donburi.World, space *cm.Space, pos cm.Vec2) {
	SpawnEnemy(1, 0.5, 0, 20, 500, world, space, pos)
}

func SpawnDefaultBomb(world donburi.World, space *cm.Space, pos cm.Vec2) {
	SpawnBomb(1, 0.1, 0, 20, world, space, pos)
}
func SpawnDefaultFoodCollectible(world donburi.World, space *cm.Space, pos cm.Vec2) {
	SpawnCollectible(comp.Food, 1, -1, 10, world, space, pos)
}

func SpawnDefaultKeyCollectible(keyNumber int, world donburi.World, space *cm.Space, pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, keyNumber, 10, world, space, pos)
}

func SpawnRandomCollectible(world donburi.World, space *cm.Space, pos cm.Vec2) {
	randomType := comp.ItemType(rand.IntN(3))
	SpawnCollectible(randomType, 1, engine.RandRangeInt(1, 10), 10, world, space, pos)
}

func SpawnRandomKeyCollectible(world donburi.World, space *cm.Space, pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, engine.RandRangeInt(1, 10), 10, world, space, pos)
}
