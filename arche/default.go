package arche

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"math/rand/v2"

	"github.com/yohamta/donburi"
)

func SpawnDefaultPlayer(pos cm.Vec2) *donburi.Entry {
	return SpawnPlayer(1, 0.3, 0.5, 20, pos)

}
func SpawnDefaultEnemy(pos cm.Vec2) *donburi.Entry {
	return SpawnEnemy(0.3, 0.3, 0.5, 20, pos)
}

func SpawnDefaultFood(pos cm.Vec2) *donburi.Entry {
	return SpawnFood(1, 0.3, 0.5, 5, pos)
}

func SpawnDefaultBomb(pos cm.Vec2) {
	SpawnBomb(1, 0.1, 0, 20, pos)
}

// collectible
func SpawnDefaultFoodCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.Food, 1, -1, 5, pos)
}
func SpawnDefaultEmeticCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.EmeticDrug, 1, -1, 13, pos)
}

func SpawnDefaultKeyCollectible(keyNumber int, pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, keyNumber, 13, pos)
}

// Random

func SpawnRandomCollectible(pos cm.Vec2) {
	randomType := comp.ItemType(rand.IntN(4))
	SpawnCollectible(randomType, 1, engine.RandRangeInt(1, 10), 10, pos)
}

func SpawnRandomKeyCollectible(pos cm.Vec2) {
	SpawnCollectible(comp.Key, 1, engine.RandRangeInt(1, 10), 10, pos)
}

func SpawnRandomEnemy(pos cm.Vec2) {
	SpawnEnemy(1, 0.3, 0.5, engine.RandRange(5, 30), pos)
}
