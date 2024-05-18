package comp

import (
	"bulimia/engine"
	"bulimia/engine/cm"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazznoer/colorgrad"
	"github.com/yohamta/donburi"
)

type ItemType int

const (
	Food ItemType = iota
	Bomb
	Key
	EmeticDrug
)

type InventoryData struct {
	Foods      int
	Bombs      int
	EmeticDrug int
	Keys       []int
}

type CollectibleData struct {
	Type      ItemType
	ItemCount int
	KeyNumber int
}

type AIData struct {
	Follow         bool
	FollowDistance float64
}

type DoorData struct {
	LockNumber   int
	Open         bool
	PlayerHasKey bool
}
type RenderData struct {
	Offset     cm.Vec2
	DrawScale  cm.Vec2
	DrawAngle  float64
	AnimPlayer *engine.AnimationPlayer
	DIO        *ebiten.DrawImageOptions
	ScaleColor color.Color
}
type CharacterData struct {
	Speed, Accel, Health float64
	VomitCooldownTimer   engine.Timer
	FoodPerCooldown      int
}
type DrugEffectData struct {
	AddMovementSpeed, Accel, Health float64
	AddVomitCooldownDuration        time.Duration
	AddFoodPerCooldown              int
	EffectTimer                     engine.Timer
}

var Inventory = donburi.NewComponentType[InventoryData](InventoryData{
	Bombs:      100,
	EmeticDrug: 20,
	Foods:      1000,
	Keys:       make([]int, 0),
})
var Door = donburi.NewComponentType[DoorData]()

var DrugEffect = donburi.NewComponentType[DrugEffectData](DrugEffectData{
	AddVomitCooldownDuration: -(time.Second / 10),
	AddFoodPerCooldown:       2,
	AddMovementSpeed:         -200,
	EffectTimer:              engine.NewTimer(time.Second * 6),
})

var Collectible = donburi.NewComponentType[CollectibleData]()

var Render = donburi.NewComponentType[RenderData](RenderData{
	Offset:     cm.Vec2{},
	DrawScale:  cm.Vec2{1, 1},
	DrawAngle:  0.0,
	DIO:        &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
	ScaleColor: color.White,
})

var Gradient = donburi.NewComponentType[colorgrad.Gradient](colorgrad.NewGradient().
	HtmlColors("rgb(255, 0, 179)", "rgb(255, 0, 0)", "rgb(255, 255, 255)").
	Domain(0, 100).
	Mode(colorgrad.BlendOklab).
	Interpolation(colorgrad.InterpolationBasis).
	Build())

var Body = donburi.NewComponentType[cm.Body]()
var AI = donburi.NewComponentType[AIData](AIData{Follow: true, FollowDistance: 300})

var Char = donburi.NewComponentType[CharacterData](CharacterData{
	Speed:              500,
	Accel:              100,
	Health:             100.,
	VomitCooldownTimer: engine.NewTimer(time.Second / 5),
	FoodPerCooldown:    1,
})

var Damage = donburi.NewComponentType[float64](20.0)

// Tags
var PlayerTag = donburi.NewTag()
var WallTag = donburi.NewTag()
var FoodTag = donburi.NewTag()
var BombTag = donburi.NewTag()
var EnemyTag = donburi.NewTag()
