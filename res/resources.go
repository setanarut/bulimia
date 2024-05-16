package res

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"embed"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"golang.org/x/text/language"
)

//go:embed assets/*
var assets embed.FS

var (
	Screen    *ebiten.Image
	Camera    *engine.Camera
	ScreenBox cm.BB
	World     donburi.World = donburi.NewWorld()
	Space     *cm.Space     = cm.NewSpace()

	Rooms       []cm.BB = make([]cm.BB, 0)
	CurrentRoom cm.BB
	Input       *engine.InputManager = &engine.InputManager{}
)

var (
	Wall      = ebiten.NewImage(30, 30)
	Pacman    = engine.LoadImage("assets/pac.png", assets)
	Items     = engine.LoadImage("assets/items.png", assets)
	Enemy     = engine.LoadImage("assets/enemy.png", assets)
	Futura    = engine.LoadTextFace("assets/futura.ttf", 20, assets)
	FuturaBig = &text.GoTextFace{
		Source:   Futura.Source,
		Size:     35,
		Language: language.English,
	}
)

func init() {
	Wall.Fill(color.White)

}

func PlayerVelocityFunc(body *cm.Body, gravity cm.Vec2, damping float64, dt float64) {

	entry, ok := body.UserData.(*donburi.Entry)

	if ok {
		if entry.Valid() {
			livingData := comp.Char.Get(entry)
			WASDAxisVector := Input.WASDDirection.Normalize().Mult(livingData.Speed)
			body.SetVelocityVector(body.Velocity().LerpDistance(WASDAxisVector, livingData.Accel))
		}
	}
}

// func ItemVelocityFunc(body *cm.Body, gravity cm.Vec2, damping float64, dt float64) {

// 	entry, ok := body.UserData.(*donburi.Entry)

// 	if ok {
// 		if entry.Valid() {

// 			WASDAxisVector := Input.WASDDirection.Normalize().Mult(livingData.Speed)
// 			body.SetVelocityVector(body.Velocity().LerpDistance(WASDAxisVector, livingData.Accel))
// 		}
// 	}
// }
