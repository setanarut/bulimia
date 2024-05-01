package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	w := 800.0
	g := NewGame(w, w*0.618)
	g.Init()

	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(int(g.screenBox.R), int(g.screenBox.T))
	gOpt := &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryAuto}
	// RUN
	if err := ebiten.RunGameWithOptions(g, gOpt); err != nil {
		log.Fatal(err)
	}

}
