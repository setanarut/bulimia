package main

import (
	"bulimia/engine/cm"
	"bulimia/system"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var Deneme int = 666

type System interface {
	Init(world donburi.World, space *cm.Space, ScreenBox *cm.BB)
	Update(world donburi.World, space *cm.Space)
	Draw(world donburi.World, space *cm.Space, screen *ebiten.Image)
}

type Game struct {
	world     donburi.World
	space     *cm.Space
	systems   []System
	screenBox *cm.BB
}

func NewGame(w, h float64) *Game {
	return &Game{
		space:     cm.NewSpace(),
		world:     donburi.NewWorld(),
		screenBox: &cm.BB{0, 0, w, h},
	}
}

// var Start = time.Now()

func (g *Game) Init() {

	g.systems = []System{
		system.NewEntitySpawnSystem(),
		system.NewPlayerControlSystem(),
		system.NewPhysicsSystem(g.world),
		// system.NewTemplate(g.screenBox),
		system.NewDrawCameraSystem(g.screenBox),
		system.NewDrawHUDSystem(g.screenBox),
	}

	// Initalize systems
	for _, s := range g.systems {
		s.Init(g.world, g.space, g.screenBox)
	}

}

func (g *Game) Update() error {
	for _, s := range g.systems {
		s.Update(g.world, g.space)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, s := range g.systems {
		s.Draw(g.world, g.space, screen)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return int(g.screenBox.R), int(g.screenBox.T)
}
