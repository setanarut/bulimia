package system

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/res"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

// Chipmunk Space draw system
type DrawHUDSystem struct {
	textOptions *text.DrawOptions
	cam         *engine.Camera
	player      *donburi.Entry
}

func NewDrawHUDSystem() *DrawHUDSystem {
	return &DrawHUDSystem{
		textOptions: &text.DrawOptions{},
	}
}
func (hs *DrawHUDSystem) Init() {
	hs.textOptions.ColorScale.ScaleWithColor(color.White)
	hs.textOptions.LineSpacing = res.FontFace.Size * 1.2
	hs.textOptions.GeoM.Translate(30, 25)

	if camE, ok := comp.Camera.First(res.World); ok {
		hs.cam = comp.Camera.Get(camE)
	}

	if player, ok := comp.PlayerTag.First(res.World); ok {
		hs.player = player
	}
}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw() {

	// debug
	if false {
		text.Draw(res.Screen, fmt.Sprintf("%v", Input.LastPressedDirection), res.FontFace, hs.textOptions)
	}
	// debug
	if false {
		text.Draw(
			res.Screen,
			fmt.Sprintf(
				"bodies : %d \nentities : %d \nActualTPS : %v \nActualFPS : %v",
				len(res.Space.DynamicBodies),
				res.World.Len(),
				ebiten.ActualTPS(),
				ebiten.ActualFPS(),
			),
			res.FontFace,
			hs.textOptions)
	}

	// inventory
	if true {
		if hs.player.Valid() {
			playerInventory := *comp.Inventory.Get(hs.player)
			liv := *comp.Living.Get(hs.player)
			text.Draw(
				res.Screen,
				fmt.Sprintf(
					"Foods: %d\nBombs: %d\nKeys: %v\nHealth: %v",
					playerInventory.Foods,
					playerInventory.Bombs,
					playerInventory.Keys,
					liv.Health,
				),
				res.FontFace,
				hs.textOptions)
		} else {
			text.Draw(res.Screen, "You are dead", res.FontFace, hs.textOptions)
		}
	}
}
