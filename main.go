package main

import (
	"bulimia/res"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	g := NewGame()
	g.Init()
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(res.Screen.Bounds().Dx(), res.Screen.Bounds().Dy())
	gOpt := &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryAuto, InitUnfocused: true}
	gOpt.DisableHiDPI = true
	// RUN
	if err := ebiten.RunGameWithOptions(g, gOpt); err != nil {
		log.Fatal(err)
	}

}
