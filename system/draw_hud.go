package system

import (
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"
	"bulimia/resources"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

// Chipmunk Space draw system
type DrawHUDSystem struct {
	screenBox   *cm.BB
	textOptions *text.DrawOptions
	cam         *engine.Camera
	player      *donburi.Entry
}

func NewDrawHUDSystem(screenBox *cm.BB) *DrawHUDSystem {
	return &DrawHUDSystem{
		screenBox:   screenBox,
		textOptions: &text.DrawOptions{},
	}
}
func (hs *DrawHUDSystem) Init(world donburi.World, space *cm.Space, screenBox *cm.BB) {
	hs.textOptions.ColorScale.ScaleWithColor(color.White)
	hs.textOptions.LineSpacing = resources.FontFace.Size * 1.2
	hs.textOptions.GeoM.Translate(30, 25)

	if camE, ok := component.CameraComp.First(world); ok {
		hs.cam = component.CameraComp.Get(camE)
	}

	if player, ok := component.PlayerTagComp.First(world); ok {
		hs.player = player
	}
}

func (hs *DrawHUDSystem) Update(world donburi.World, space *cm.Space) {

}
func (hs *DrawHUDSystem) Draw(world donburi.World, space *cm.Space, scr *ebiten.Image) {

	// debug
	if false {
		text.Draw(scr, fmt.Sprintf("%v", AnyKeyDown(ArrowKeys)), resources.FontFace, hs.textOptions)
	}
	// debug
	if false {
		text.Draw(
			scr,
			fmt.Sprintf(
				"bodies : %d \nentities : %d \nActualTPS : %v \nActualFPS : %v",
				len(space.DynamicBodies),
				world.Len(),
				ebiten.ActualTPS(),
				ebiten.ActualFPS(),
			),
			resources.FontFace,
			hs.textOptions)
	}

	// inventory
	if true {
		if hs.player.Valid() {
			playerInventory := *component.InventoryComp.Get(hs.player)
			playerHealth := *component.HealthComp.Get(hs.player)
			text.Draw(
				scr,
				fmt.Sprintf(
					"Foods: %d\nBombs: %d\nKeys: %v\nHealth: %v",
					playerInventory.Foods,
					playerInventory.Bombs,
					playerInventory.Keys,
					playerHealth,
				),
				resources.FontFace,
				hs.textOptions)
		} else {
			text.Draw(scr, "You are dead", resources.FontFace, hs.textOptions)
		}
	}
}
