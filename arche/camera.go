package arche

import (
	"bulimia/component"
	"bulimia/engine"
	"bulimia/engine/cm"

	"github.com/yohamta/donburi"
)

func NewCameraEntity(lookAt cm.Vec2, width, height float64, w donburi.World) *donburi.Entry {
	e := w.Entry(w.Create(component.CameraComp))
	cam := engine.NewCamera(lookAt, width, height)
	cam.Lerp = true
	component.CameraComp.Set(e, cam)
	return e
}
