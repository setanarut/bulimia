package main

import (
	"bulimia/engine/cm"
	"bulimia/res"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	w, h := 800, 600
	res.Screen = ebiten.NewImage(w, h)
	res.ScreenBox = cm.NewBB(0, 0, float64(w), float64(h))

	g := NewGame()
	g.Init()

	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(res.Screen.Bounds().Dx(), res.Screen.Bounds().Dy())
	gOpt := &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryAuto}
	// RUN
	if err := ebiten.RunGameWithOptions(g, gOpt); err != nil {
		log.Fatal(err)
	}

}
