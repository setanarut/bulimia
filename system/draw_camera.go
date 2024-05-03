package system

import (
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazznoer/colorgrad"
	"github.com/yohamta/donburi"
)

// DrawCameraSystem
type DrawCameraSystem struct {
	entityDIO *ebiten.DrawImageOptions
	screenBox *cm.BB
	grad      colorgrad.Gradient
	cam       *engine.Camera
}

func NewDrawCameraSystem(screenBox *cm.BB) *DrawCameraSystem {
	gr, _ := colorgrad.NewGradient().
		HtmlColors("rgb(255, 0, 179)", "rgb(255, 0, 0)", "rgb(255, 255, 255)").
		Domain(0, 100).
		Mode(colorgrad.BlendOklab).
		Interpolation(colorgrad.InterpolationBasis).
		Build()

	dbs := &DrawCameraSystem{
		entityDIO: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		},
		screenBox: screenBox,
		grad:      gr,
	}

	return dbs
}

func (ds *DrawCameraSystem) Init(world donburi.World, space *cm.Space, screenBox *cm.BB) {
	if cam, ok := component.CameraComp.First(world); ok {
		ds.cam = component.CameraComp.Get(cam)
	}
}

func (ds *DrawCameraSystem) Update(world donburi.World, space *cm.Space) {
	ds.cam.LookAt(engine.InvPosVectY(CurrentRoom.Center(), ds.screenBox.T))

	component.AnimPlayerComp.Each(world, func(e *donburi.Entry) {
		component.AnimPlayerComp.Get(e).Update()
	})
}

func (ds *DrawCameraSystem) Draw(world donburi.World, space *cm.Space, screen *ebiten.Image) {

	// arka plan
	screen.Fill(color.Gray{0})

	component.WallTagComp.Each(world, func(e *donburi.Entry) {
		body := component.BodyComp.Get(e)
		animComp := component.AnimPlayerComp.Get(e)
		ds.updateGeoM(body, animComp.DrawOffset, animComp.DrawScaleX, animComp.DrawScaleY)
		ds.entityDIO.ColorScale.ScaleWithColor(color.RGBA{0, 0, 255, 0})
		ds.cam.Draw(animComp.CurrentFrame, ds.entityDIO, screen)
		ds.entityDIO.ColorScale.Reset()

	})

	component.DoorComp.Each(world, func(e *donburi.Entry) {
		body := component.BodyComp.Get(e)
		animComp := component.AnimPlayerComp.Get(e)
		doorData := component.DoorComp.Get(e)
		ds.updateGeoM(body, animComp.DrawOffset, animComp.DrawScaleX, animComp.DrawScaleY)
		if doorData.Open {
			ds.entityDIO.ColorScale.ScaleWithColor(color.RGBA{0, 0, 0, 0})
		} else {
			if doorData.PlayerHasKey {
				ds.entityDIO.ColorScale.ScaleWithColor(color.RGBA{0, 255, 0, 255})
			} else {
				ds.entityDIO.ColorScale.ScaleWithColor(color.RGBA{0, 0, 200, 255})
			}
		}

		ds.cam.Draw(animComp.CurrentFrame, ds.entityDIO, screen)
		ds.entityDIO.ColorScale.Reset()

	})

	component.CollectibleComp.Each(world, func(entry *donburi.Entry) {
		body := component.BodyComp.Get(entry)
		animComp := component.AnimPlayerComp.Get(entry)
		ds.updateGeoM(body, animComp.DrawOffset, animComp.DrawScaleX, animComp.DrawScaleY)
		ds.entityDIO.ColorScale.ScaleWithColor(color.RGBA{0, 255, 100, 255})
		ds.cam.Draw(animComp.CurrentFrame, ds.entityDIO, screen)
		ds.entityDIO.ColorScale.Reset()

	})

	component.EnemyTagComp.Each(world, func(e *donburi.Entry) {
		body := component.BodyComp.Get(e)
		animComp := component.AnimPlayerComp.Get(e)
		h := *component.HealthComp.Get(e)
		ds.entityDIO.ColorScale.ScaleWithColor(ds.grad.At(float64(h)))
		ds.updateGeoM(body, animComp.DrawOffset, animComp.DrawScaleX, animComp.DrawScaleY)

		ds.cam.Draw(animComp.CurrentFrame, ds.entityDIO, screen)
		ds.entityDIO.ColorScale.Reset()

	})
	component.BombTagComp.Each(world, func(e *donburi.Entry) {
		body := component.BodyComp.Get(e)
		animComp := component.AnimPlayerComp.Get(e)
		ds.updateGeoM(body, animComp.DrawOffset, animComp.DrawScaleX, animComp.DrawScaleY)

		ds.cam.Draw(animComp.CurrentFrame, ds.entityDIO, screen)
		ds.entityDIO.ColorScale.Reset()

	})
	component.FoodTagComp.Each(world, func(e *donburi.Entry) {
		body := component.BodyComp.Get(e)
		animComp := component.AnimPlayerComp.Get(e)
		ds.entityDIO.ColorScale.ScaleWithColor(color.RGBA{0, 255, 100, 255})
		ds.updateGeoM(body, animComp.DrawOffset, animComp.DrawScaleX, animComp.DrawScaleY)
		ds.cam.Draw(animComp.CurrentFrame, ds.entityDIO, screen)
		ds.entityDIO.ColorScale.Reset()

	})

	if e, ok := component.PlayerTagComp.First(world); ok {
		body := component.BodyComp.Get(e)
		animComp := component.AnimPlayerComp.Get(e)
		ds.updateGeoM(body, animComp.DrawOffset, animComp.DrawScaleX, animComp.DrawScaleY)
		ds.entityDIO.ColorScale.ScaleWithColor(color.RGBA{255, 220, 0, 255})
		ds.cam.Draw(animComp.CurrentFrame, ds.entityDIO, screen)
		ds.entityDIO.ColorScale.Reset()
	}

}

func (ds *DrawCameraSystem) updateGeoM(body *cm.Body, centerOffset cm.Vec2, scaleX, scaleY float64) {
	ds.entityDIO.GeoM.Reset()
	pos := engine.InvPosVectY(body.Position(), ds.screenBox.T)
	ds.entityDIO.GeoM.Translate(centerOffset.X, centerOffset.Y)
	ds.entityDIO.GeoM.Scale(scaleX, scaleY)
	ds.entityDIO.GeoM.Rotate(engine.InvertAngle(body.Angle()))
	ds.entityDIO.GeoM.Translate(pos.X, pos.Y)
}
