package main

import (
	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Segment struct {
	sprite             *lib.Sprite
	cx                 int
	cy                 int
	health             int
	fast               bool
	head               bool
	inEdge             Direction
	outEdge            Direction
	disallowDirection  Direction
	previousXDirection Direction
}

func NewSegment(cx, cy, health int, fast, head bool) *Segment {
	return &Segment{
		sprite: lib.NewSprite(lib.XCentre, lib.YCentre),
		cx:     cx,
		cy:     cy,
		health: health,
		fast:   fast,
		head:   head,
		// Each myriapod segment moves in a defined pattern within its current cell, before moving to the next one.
		// It will start at one of the edges - represented by a number, where 0=down,1=right,2=up,3=left
		// inEdge stores the edge through which it entered the cell.
		// Several frames after entering a cell, it chooses which edge to leave through - stored in outEdge
		// The path it follows is explained in the update and rank methods
		inEdge:  DirectionLeft,
		outEdge: DirectionRight,

		disallowDirection:  DirectionUp,    // Prevents segment from moving in a particular direction
		previousXDirection: DirectionRight, // Used to create winding/snaking motion
	}
}

func (s *Segment) Update() {
	x, y := CellToPos(s.cx, s.cy, 0, 0)
	s.sprite.MoveTo(x, y)
	s.sprite.SetImage(images["seg00000"])
	s.sprite.Update()
}

func (s *Segment) Draw(screen *ebiten.Image) {
	s.sprite.Draw(screen)
}
