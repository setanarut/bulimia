package system

import (
	"bulimia/comp"
	"bulimia/res"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
)

// Chipmunk Space draw system
type DrawHUDSystem struct {
	statsTextOptions  *text.DrawOptions
	centerTextOptions *text.DrawOptions
}

func NewDrawHUDSystem() *DrawHUDSystem {
	return &DrawHUDSystem{
		statsTextOptions: &text.DrawOptions{},
		centerTextOptions: &text.DrawOptions{

			LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
		},
	}
}
func (hs *DrawHUDSystem) Init() {
	hs.statsTextOptions.ColorScale.ScaleWithColor(colornames.White)
	hs.statsTextOptions.LineSpacing = res.Iosevka.Size * 1.2
	hs.statsTextOptions.GeoM.Translate(30, 25)
	hs.statsTextOptions.Filter = ebiten.FilterLinear

	hs.centerTextOptions.LayoutOptions.PrimaryAlign = text.AlignCenter
	hs.centerTextOptions.LayoutOptions.SecondaryAlign = text.AlignCenter
	hs.centerTextOptions.Filter = ebiten.FilterLinear
	hs.centerTextOptions.LineSpacing = res.IosevkaBig.Size * 1.2
	center := res.ScreenBox.Center()
	hs.centerTextOptions.GeoM.Translate(center.X, center.Y)

}

func (hs *DrawHUDSystem) Update() {

}
func (hs *DrawHUDSystem) Draw() {

	if ebiten.IsFocused() {
		// inventory
		if true {
			p, ok := comp.PlayerTag.First(res.World)
			if ok {
				playerInventory := *comp.Inventory.Get(p)

				liv := *comp.Living.Get(p)
				text.Draw(
					res.Screen,
					fmt.Sprintf(
						"Foods: %d\nBombs: %d\nKeys: %v\nHealth: %v",
						playerInventory.Foods,
						playerInventory.Bombs,
						playerInventory.Keys,
						liv.Health,
					),
					res.Iosevka,
					hs.statsTextOptions)
			} else {
				text.Draw(res.Screen, "You are dead \n Press Backspace key to restart", res.IosevkaBig, hs.centerTextOptions)
			}
		}
	} else {

		// unfocused
		if true {
			text.Draw(res.Screen, "PAUSED\n Click to resume", res.IosevkaBig, hs.centerTextOptions)
		}

	}

	// FPS/TPS Debug text
	if false {
		text.Draw(
			res.Screen,
			fmt.Sprintf(
				"DynamicBodies : %d\nStaticBodies : %dEntities : %d\nActualTPS : %v\nActualFPS : %v",
				len(res.Space.DynamicBodies),
				len(res.Space.StaticBodies),
				res.World.Len(),
				ebiten.ActualTPS(),
				ebiten.ActualFPS(),
			),
			res.Iosevka,
			hs.statsTextOptions)
	}

}
