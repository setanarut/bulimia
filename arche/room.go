package arche

import (
	"bulimia/engine/cm"

	"github.com/yohamta/donburi"
)

func MakeRoom(world donburi.World, space *cm.Space, roomBB cm.BB, opts RoomOptions) {

	topDoorLength := roomBB.Width() / 5
	leftDoorLength := roomBB.Height() / 5

	topDoorCenter := roomBB.LT().Lerp(roomBB.RT(), 0.5)
	bottomDoorCenter := roomBB.LB().Lerp(roomBB.RB(), 0.5)

	leftDoorCenter := roomBB.LT().Lerp(roomBB.LB(), 0.5)
	rightDoorCenter := roomBB.RT().Lerp(roomBB.RB(), 0.5)

	topLeftWallCenter := cm.Vec2{topDoorLength, roomBB.T}
	topRightWallCenter := cm.Vec2{roomBB.R - topDoorLength, roomBB.T}
	bottomLeftWallCenter := cm.Vec2{topDoorLength, roomBB.B}
	bottomRightWallCenter := cm.Vec2{roomBB.R - topDoorLength, roomBB.B}

	leftDoorBottom := cm.Vec2{roomBB.L, roomBB.B + leftDoorLength}
	leftDoorTop := cm.Vec2{roomBB.L, roomBB.T - leftDoorLength}

	rightDoorBottom := cm.Vec2{roomBB.R, roomBB.B + leftDoorLength}
	rightDoorTop := cm.Vec2{roomBB.R, roomBB.T - leftDoorLength}

	// Top Wall
	if opts.TopWall {
		NewWallEntity(world, space, topLeftWallCenter, topDoorLength*2, 10)
		NewDoorEntity(world, space, topDoorCenter, topDoorLength, 10, opts.TopDoorKeyNumber)
		NewWallEntity(world, space, topRightWallCenter, topDoorLength*2, 10)
	}

	// Bottom Wall
	if opts.BottomWall {
		NewWallEntity(world, space, bottomLeftWallCenter, topDoorLength*2, 10)
		NewDoorEntity(world, space, bottomDoorCenter, topDoorLength, 10, opts.BottomDoorKeyNumber)
		NewWallEntity(world, space, bottomRightWallCenter, topDoorLength*2, 10)
	}

	// Left Wall
	if opts.LeftWall {
		NewWallEntity(world, space, leftDoorTop, 10, leftDoorLength*2)
		NewDoorEntity(world, space, leftDoorCenter, 10, leftDoorLength, opts.LeftDoorKeyNumber)
		NewWallEntity(world, space, leftDoorBottom, 10, leftDoorLength*2)
	}

	// Right Wall
	if opts.RightWall {
		NewWallEntity(world, space, rightDoorTop, 10, leftDoorLength*2)
		NewDoorEntity(world, space, rightDoorCenter, 10, leftDoorLength, opts.RightDoorKeyNumber)
		NewWallEntity(world, space, rightDoorBottom, 10, leftDoorLength*2)
	}
}

type RoomOptions struct {
	TopWall, BottomWall, LeftWall, RightWall                                     bool
	TopDoorKeyNumber, BottomDoorKeyNumber, LeftDoorKeyNumber, RightDoorKeyNumber int
}
