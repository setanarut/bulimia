package system

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/res"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// DrawCameraSystem
type DrawCameraSystem struct {
}

func NewDrawCameraSystem() *DrawCameraSystem {
	return &DrawCameraSystem{}
}

func (ds *DrawCameraSystem) Init() {
}

func (ds *DrawCameraSystem) Update() {

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		res.Camera.ZoomFactor -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		res.Camera.ZoomFactor += 1
	}

	res.Camera.LookAt(engine.InvPosVectY(res.CurrentRoom.Center(), res.ScreenRect.T))

	comp.Render.Each(res.World, func(e *donburi.Entry) {
		comp.Render.Get(e).AnimPlayer.Update()
	})

}

func (ds *DrawCameraSystem) Draw() {

	// arka plan
	res.Screen.Fill(color.Gray{0})

	comp.WallTag.Each(res.World, ds.DrawEntry)

	comp.Door.Each(res.World, func(e *donburi.Entry) {
		render := comp.Render.Get(e)
		doorData := comp.Door.Get(e)

		if doorData.Open {
			render.ScaleColor = color.RGBA{0, 0, 0, 0}
		} else {
			if doorData.PlayerHasKey {
				render.ScaleColor = color.RGBA{0, 255, 0, 255}
			} else {
				render.ScaleColor = color.RGBA{0, 0, 200, 255}
			}
		}

		ds.DrawEntry(e)

	})

	comp.Collectible.Each(res.World, ds.DrawEntry)
	comp.BombTag.Each(res.World, ds.DrawEntry)
	comp.FoodTag.Each(res.World, ds.DrawEntry)

	comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {
		r := comp.Render.Get(e)
		r.ScaleColor = comp.Gradient.Get(e).At(comp.Char.Get(e).Health)
		ds.DrawEntry(e)

	})

	if e, ok := comp.PlayerTag.First(res.World); ok {
		ds.DrawEntry(e)
	}

}

func (ds *DrawCameraSystem) DrawEntry(e *donburi.Entry) {

	body := comp.Body.Get(e)
	render := comp.Render.Get(e)
	pos := engine.InvPosVectY(body.Position(), res.ScreenRect.T)

	render.DIO.GeoM.Reset()
	render.DIO.GeoM.Translate(render.Offset.X, render.Offset.Y)
	render.DIO.GeoM.Scale(render.DrawScale.X, render.DrawScale.Y)
	render.DIO.GeoM.Rotate(engine.InvertAngle(render.DrawAngle))
	render.DIO.GeoM.Translate(pos.X, pos.Y)

	render.DIO.ColorScale.ScaleWithColor(render.ScaleColor)
	res.Camera.Draw(render.AnimPlayer.CurrentFrame, render.DIO, res.Screen)
	render.DIO.ColorScale.Reset()
}
