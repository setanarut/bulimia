package main

import (
	"bulimia/res"
	"bulimia/system"

	"github.com/hajimehoshi/ebiten/v2"
)

type System interface {
	Init()
	Update()
	Draw()
}

type Game struct {
	systems []System
}

func NewGame() *Game {
	return &Game{}
}

// var Start = time.Now()

func (g *Game) Init() {

	g.systems = []System{
		system.NewEntitySpawnSystem(),
		system.NewPlayerControlSystem(),
		system.NewPhysicsSystem(),
		// system.NewTemplate(g.screenBox),
		system.NewDrawCameraSystem(),
		system.NewDrawHUDSystem(),
	}

	// Initalize systems
	for _, s := range g.systems {
		s.Init()
	}

}

func (g *Game) Update() error {
	for _, s := range g.systems {
		s.Update()
	}
	return nil
}

func (g *Game) Draw(s *ebiten.Image) {
	for _, s := range g.systems {
		s.Draw()
	}
	s.DrawImage(res.Screen, nil)
}

func (g *Game) Layout(w, h int) (int, int) {
	return res.Screen.Bounds().Dx(), res.Screen.Bounds().Dy()
}
