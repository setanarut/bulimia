package comp

import (
	"bulimia/engine"
	"bulimia/engine/cm"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazznoer/colorgrad"
	"github.com/yohamta/donburi"
)

type ItemType int

const (
	Food ItemType = iota
	Bomb
	Key
)

type InventoryData struct {
	Foods int
	Bombs int
	Keys  []int
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
	AnimPlayer *engine.AnimationPlayer
	DIO        *ebiten.DrawImageOptions
	ScaleColor color.Color
}
type LivingData struct {
	Speed, Accel, Health, Damage float64
}

var Camera = donburi.NewComponentType[engine.Camera]()
var Inventory = donburi.NewComponentType[InventoryData](InventoryData{Bombs: 100, Foods: 100, Keys: make([]int, 0)})
var Door = donburi.NewComponentType[DoorData]()
var Collectible = donburi.NewComponentType[CollectibleData]()

var Render = donburi.NewComponentType[RenderData](RenderData{
	Offset:     cm.Vec2{},
	DrawScale:  cm.Vec2{1, 1},
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
var AI = donburi.NewComponentType[AIData](AIData{Follow: false, FollowDistance: 300})
var Living = donburi.NewComponentType[LivingData](LivingData{
	Speed:  400,
	Accel:  40,
	Health: 100.,
})

var Damage = donburi.NewComponentType[float64](20.0)

// Tags
var PlayerTag = donburi.NewTag()
var WallTag = donburi.NewTag()
var FoodTag = donburi.NewTag()
var BombTag = donburi.NewTag()
var EnemyTag = donburi.NewTag()
