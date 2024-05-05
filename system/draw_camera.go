package system

import (
	"bulimia/comp"
	"bulimia/engine"
	"bulimia/engine/cm"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// DrawCameraSystem
type DrawCameraSystem struct {
	screenBox cm.BB
	cam       *engine.Camera
}

func NewDrawCameraSystem(screenBox cm.BB) *DrawCameraSystem {

	dbs := &DrawCameraSystem{
		screenBox: screenBox,
	}

	return dbs
}

func (ds *DrawCameraSystem) Init(world donburi.World, space *cm.Space, screenBox cm.BB) {
	if cam, ok := comp.Camera.First(world); ok {
		ds.cam = comp.Camera.Get(cam)
	}
}

func (ds *DrawCameraSystem) Update(world donburi.World, space *cm.Space) {

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		ds.cam.ZoomFactor -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		ds.cam.ZoomFactor += 1
	}

	ds.cam.LookAt(engine.InvPosVectY(CurrentRoom.Center(), ds.screenBox.T))

	comp.Render.Each(world, func(e *donburi.Entry) {
		comp.Render.Get(e).AnimPlayer.Update()
	})

}

func (ds *DrawCameraSystem) Draw(world donburi.World, space *cm.Space, screen *ebiten.Image) {

	// arka plan
	screen.Fill(color.Gray{0})

	comp.WallTag.Each(world, func(e *donburi.Entry) {
		ds.DrawEntry(e, screen)
	})

	comp.Door.Each(world, func(e *donburi.Entry) {
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

		ds.DrawEntry(e, screen)

	})

	comp.Collectible.Each(world, func(e *donburi.Entry) {
		ds.DrawEntry(e, screen)

	})

	comp.EnemyTag.Each(world, func(e *donburi.Entry) {
		r := comp.Render.Get(e)
		r.ScaleColor = comp.Gradient.Get(e).At(comp.Living.Get(e).Health)
		ds.DrawEntry(e, screen)

	})
	comp.BombTag.Each(world, func(e *donburi.Entry) {
		ds.DrawEntry(e, screen)

	})
	comp.FoodTag.Each(world, func(e *donburi.Entry) {

		ds.DrawEntry(e, screen)

	})

	if e, ok := comp.PlayerTag.First(world); ok {

		ds.DrawEntry(e, screen)
	}

}

func (ds *DrawCameraSystem) DrawEntry(e *donburi.Entry, screen *ebiten.Image) {

	body := comp.Body.Get(e)
	render := comp.Render.Get(e)
	pos := engine.InvPosVectY(body.Position(), ds.screenBox.T)

	render.DIO.GeoM.Reset()
	render.DIO.GeoM.Translate(render.Offset.X, render.Offset.Y)
	render.DIO.GeoM.Scale(render.DrawScale.X, render.DrawScale.Y)
	// render.DIO.GeoM.Rotate(engine.InvertAngle(body.Angle()))
	render.DIO.GeoM.Translate(pos.X, pos.Y)

	render.DIO.ColorScale.ScaleWithColor(render.ScaleColor)
	ds.cam.Draw(render.AnimPlayer.CurrentFrame, render.DIO, screen)
	render.DIO.ColorScale.Reset()
}
