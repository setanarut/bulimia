package arche

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"math/rand/v2"
)

func SpawnDefaultEnemy(pos cm.Vec2) {
	SpawnEnemy(1, 0.5, 0, 20, pos)
}
func SpawnDefaultPlayer(pos cm.Vec2) {
	SpawnPlayer(0.1, 0.3, 0, 20, pos)

}

func SpawnDefaultBomb(pos cm.Vec2) {
	SpawnBomb(1, 0.1, 0, 20, pos)
}
func SpawnDefaultFoodCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.Food, 1, -1, 10, pos)
}

func SpawnDefaultKeyCollectible(keyNumber int, pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, keyNumber, 10, pos)
}

func SpawnRandomCollectible(pos cm.Vec2) {
	randomType := comp.ItemType(rand.IntN(3))
	SpawnCollectible(randomType, 1, engine.RandRangeInt(1, 10), 10, pos)
}

func SpawnRandomKeyCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, engine.RandRangeInt(1, 10), 10, pos)
}
