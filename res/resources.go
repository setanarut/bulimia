package res

import (
	"bulimia/engine"
	"bulimia/engine/cm"
	"embed"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

//go:embed assets/*
var assets embed.FS

var (
	Screen    *ebiten.Image
	ScreenBox cm.BB
	World     donburi.World = donburi.NewWorld()
	Space     *cm.Space     = cm.NewSpace()

	Rooms       []cm.BB = make([]cm.BB, 0)
	CurrentRoom cm.BB
)

var (
	Wall     = ebiten.NewImage(30, 30)
	Pacman   = engine.LoadImage("assets/pac.png", assets)
	Items    = engine.LoadImage("assets/items.png", assets)
	Enemy    = engine.LoadImage("assets/enemy.png", assets)
	FontFace = engine.LoadTextFace("assets/iosevka.ttf", 20, assets)
)

func init() {
	Wall.Fill(color.White)
}
