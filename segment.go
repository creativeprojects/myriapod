package main

import (
	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// X and Y directions indexed into by in_edge and out_edge in Segment
	// The indices correspond to the direction numbers above, i.e. 0 = up, 1 = right, 2 = down, 3 = left
	DX = []int{0, 1, 0, -1}
	DY = []int{-1, 0, 1, 0}
)

type Segment struct {
	game               *Game
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

func NewSegment(game *Game, cx, cy, health int, fast, head bool) *Segment {
	return &Segment{
		game:   game,
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

func (s *Segment) Collision(x, y float64) bool {
	if s.sprite.CollidePoint(x, y) {
		s.health--
		return true
	}
	return false
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

func (s *Segment) rank(proposedOutEdge Direction) (bool, bool, bool, bool, bool, bool, bool) {
	// proposed_out_edge is a number between 0 and 3, representing a possible direction to move - see DIRECTION_UP etc and DX/DY above
	// This function returns a tuple consisting of a series of factors determining which grid cell the segment should try to move into next.
	// These are not absolute rules - rather they are used to rank the four directions in order of preference,
	// i.e. which direction is the best (or at least, least bad) to move in. The factors are boolean (True or False)
	// values. A value of False is preferable to a value of True.
	// The order of the factors in the returned tuple determines their importance in deciding which way to go,
	// with the most important factor coming first.
	newCellX := s.cx + DX[proposedOutEdge]
	newCellY := s.cy + DY[proposedOutEdge]

	// Does this direction take us to a cell which is outside the grid?
	// Note: when the segments start, they are all outside the grid so this would be True, except for the case of
	// walking onto the top-left cell of the grid. But the end result of this and the following factors is that
	// it will still be allowed to continue walking forwards onto the screen.
	out := newCellX < 0 || newCellX > NumGridCols-1 || newCellY < 0 || newCellY > NumGridRows-1

	// We don't want it to to turn back on itself..
	turningBackOnSelf := proposedOutEdge == s.inEdge

	// ..or go in a direction that's disallowed (see comments in update method)
	directionDisallowed := proposedOutEdge == s.disallowDirection

	// Check to see if there's a rock at the proposed new grid cell.
	// rock will either be the Rock object at the new grid cell, or None.
	// It will be set to None if there is no Rock object is at the new location, or if the new location is
	// outside the grid. We also have to account for the special case where the segment is off the left-hand
	// side of the screen on the first row, where it is initially created. We mustn't try to access that grid
	// cell (unlike most languages, in Python trying to access a list index with negative value won't necessarily
	// result in a crash, but it's still not a good idea)
	var rock *Rock
	if out || (newCellY == 0 && newCellX < 0) {
		rock = nil
	} else {
		rock = s.game.grid[newCellY][newCellX]
	}

	rockPresent := rock != nil

	// Is new cell already occupied by another segment, or is another segment trying to enter my cell from
	// the opposite direction?
	// occupiedBySegment := (new_cell_x, new_cell_y) in game.occupied or (s.cx, s.cy, proposed_out_edge) in game.occupied
	occupiedBySegment := false

	// Prefer to move horizontally, unless there's a rock in the way.
	// If there are rocks both horizontally and vertically, prefer to move vertically
	var horizontalBlocked bool
	if rockPresent {
		horizontalBlocked = proposedOutEdge.IsHorizontal()
	} else {
		horizontalBlocked = !proposedOutEdge.IsHorizontal()
	}

	// Prefer not to go in the previous horizontal direction after we move up/down
	sameAsPreviousXDirection := proposedOutEdge == s.previousXDirection

	// Finally we create and return a tuple of factors determining which cell segment should try to move into next.
	// Most important first - e.g. we shouldn't enter a new cell if if's outside the grid
	return out, turningBackOnSelf, directionDisallowed, occupiedBySegment, rockPresent, horizontalBlocked, sameAsPreviousXDirection
}
