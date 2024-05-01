/* package system

import (
	"bulimia/comp"
	"bulimia/engine/cm"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type sys struct {
	query     *donburi.Query
	screenBox *cm.BB
}

func NewTemplate(screenBox *cm.BB) *sys {
	return &sys{
		screenBox: screenBox,
		// query: query.NewQuery(filter.Or(
		// 	filter.Contains(comp.EnemyTagComp), filter.Contains(comp.PlayerTagComp))),
		query: query.NewQuery(filter.Contains(comp.PlayerTagComp)),
	}
}
func (s *sys) Init(world donburi.World, space *cm.Space, screenBox *cm.BB) {
	world.OnRemove(func(world donburi.World, entity donburi.Entity) {

	})
}

func (s *sys) Update(world donburi.World, space *cm.Space) {
	// s.query.Each(world, func(e *donburi.Entry) {
	// 	if *comp.IsDeadComp.Get(e) {
	// 		fmt.Println(e)
	// 	}
	// })
}
func (s *sys) Draw(world donburi.World, space *cm.Space, screen *ebiten.Image) {

}
*/