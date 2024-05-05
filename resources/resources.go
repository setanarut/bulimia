package resources

import (
	"bulimia/engine"
	"embed"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

var Wall = ebiten.NewImage(30, 30)

var Pacman = engine.LoadImage("assets/pac.png", assets)
var Items = engine.LoadImage("assets/items.png", assets)
var Enemy = engine.LoadImage("assets/enemy.png", assets)
var FontFace = engine.LoadTextFace("assets/iosevka.ttf", 20, assets)

func init() {
	Wall.Fill(color.White)
}
