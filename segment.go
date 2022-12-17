package main

import (
	"strconv"
	"strings"

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
	posX               float64
	posY               float64
	legFrame           int
	cx                 int
	cy                 int
	health             int
	fast               bool
	head               bool
	inEdge             Direction
	outEdge            Direction
	disallowDirection  Direction
	previousXDirection Direction
	direction          Direction
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
	s.update()
	s.sprite.MoveTo(s.posX, s.posY)

	imageName := &strings.Builder{}
	imageName.WriteString("seg")
	if s.fast {
		imageName.WriteByte('1')
	} else {
		imageName.WriteByte('0')
	}
	if s.health == 2 {
		imageName.WriteByte('1')
	} else {
		imageName.WriteByte('0')
	}
	if s.head {
		imageName.WriteByte('1')
	} else {
		imageName.WriteByte('0')
	}
	imageName.WriteString(strconv.Itoa(int(s.direction)))
	imageName.WriteString(strconv.Itoa(s.legFrame))
	s.sprite.SetImage(images[imageName.String()])
	s.sprite.Update()
}

func (s *Segment) Draw(screen *ebiten.Image) {
	s.sprite.Draw(screen)
}

func (s *Segment) update() {
	// Segments take either 16 or 8 frames to pass through each grid cell, depending on the amount by which
	// game.time is updated each frame. phase will be a number between 0 and 15 indicating where we're at
	// in that cycle.
	phase := s.game.time % 16

	if phase == 0 {
		// At this point, the segment is entering a new grid cell. We first update our current grid cell coordinates.
		s.cx += DX[s.outEdge]
		s.cy += DY[s.outEdge]

		// We then need to update in_edge. If, for example, we left the previous cell via its right edge, that means
		// we're entering the new cell via its left edge.
		s.inEdge = s.outEdge.Inverse()

		// During normal gameplay, once a segment reaches the bottom of the screen, it starts moving up again.
		// Once it reaches row 18, it starts moving down again, so that it remains a threat to the player.
		// During the title screen, we allow segments to go all the way back up to the top of the screen.
		tempY := 0
		if s.game.player != nil {
			tempY = 18
		}
		if s.cy == tempY {
			s.disallowDirection = DirectionUp
		}
		if s.cy == NumGridRows-1 {
			s.disallowDirection = DirectionDown
		}

	} else if phase == 4 {
		// At this point we decide which new cell we're going to go into (and therefore, which edge of the current
		// cell we will leave via - to be stored in out_edge)
		directions := make([]int, 4)
		for i := 0; i <= 3; i++ {
			directions[i] = s.rank(Direction(i))
		}
		min := 128
		minDirection := 0
		for i := 0; i <= 3; i++ {
			if directions[i] < min {
				min = directions[i]
				minDirection = i
			}
		}
		s.outEdge = Direction(minDirection)

		if s.outEdge.IsHorizontal() {
			s.previousXDirection = s.outEdge
		}

		newCellX := s.cx + DX[s.outEdge]
		newCellY := s.cy + DY[s.outEdge]

		// Destroy any rock that might be in the new cell
		if newCellX >= 0 && newCellX < NumGridCols {
			s.game.Damage(newCellX, newCellY, 5, false)
		}

		// Set new cell as occupied. It's a case of whichever segment is processed first, gets first dibs on a cell
		// The second line deals with the case where two segments are moving towards each other and are in
		// neighbouring cells. It allows a segment to tell if another segment trying to enter its cell from
		// the opposite direction
		s.game.occupied = append(s.game.occupied,
			Cell{X: newCellX, Y: newCellY},
			Cell{X: newCellX, Y: newCellY, Edge: s.outEdge.Inverse()},
		)
	}
	// turnIdx tells us whether the segment is going to be making a 90 degree turn in the current cell, or moving
	// in a straight line. 1 = anti-clockwise turn, 2 = straight ahead, 3 = clockwise turn, 0 = leaving through same
	// edge from which we entered (unlikely to ever happen in practice)
	turnIdx := (s.outEdge - s.inEdge) % 4

	// Calculate segment offset in the cell, measured from the cell's centre
	// We start off assuming that the segment is starting from the top of the cell - i.e. s.inEdge being DIRECTION_UP,
	// corresponding to zero. The primary and secondary axes, as described under "SEGMENT MOVEMENT" above, are Y and X.
	// We then apply a calculation to rotate these X and Y offsets, based on the actual direction the segment is coming from.
	// Let's take as an example the case where the segment is moving in a straight line from top to bottom.
	// We calculate offsetX by multiplying SECONDARY_AXIS_POSITIONS[phase] by 2-turn_idx. In this case, turn_idx
	// will be 2.  So 2 - turn_idx will be zero. Multiplying anything by zero gives zero, so we end up with no
	// movement on the X axis - which is what we want in this case.
	// The starting point for the offset_y calculation is that the segment starts at an offset of -16 and must cover
	// 32 pixels over the 16 phases - therefore we must multiply phase by 2. We then subtract the result of the
	// previous line, in which stolen_y_movement was calculated by multiplying SECONDARY_AXIS_POSITIONS[phase] by
	// turn_idx % 2.  mod 2 gives either zero (if turn_idx is 0 or 2), or 1 if it's 1 or 3. In the case we're looking
	// at, turn_idx is 2, so stolen_y_movement is zero.
	// The end result of all this is that in the case where the segment is moving in a straight line through a cell,
	// it just moves at 2 pixels per frame along the primary axis. If it's turning, it starts out moving at 2px
	// per frame on the primary axis, but then starts moving along the secondary axis based on the values in
	// SECONDARY_AXIS_POSITIONS. In this case we don't want it to continue moving along the primary axis - it should
	// initially slow to moving at 1px per phase, and then stop moving completely. Effectively, the secondary axis
	// is stealing movement from the primary axis - hence the name 'stolen_y_movement'
	offsetX := SecondaryAxisPositions[phase] * (2 - int(turnIdx))
	stolen_y_movement := (int(turnIdx) % 2) * SecondaryAxisPositions[phase]
	offsetY := -16 + (phase * 2) - stolen_y_movement

	// A rotation matrix is a set of numbers which, when multiplied by a set of coordinates, result in those
	// coordinates being rotated. Recall that the code above  makes the assumption that segment is starting from the
	// top edge of the cell and moving down. The code below chooses the appropriate rotation matrix based on the
	// actual edge the segment started from, and then modifies offset_x and offset_y based on this rotation matrix.
	rotation_matrix := RotationData[s.inEdge]
	offsetX = offsetX*rotation_matrix[0] + offsetY*rotation_matrix[1]
	offsetY = offsetX*rotation_matrix[2] + offsetY*rotation_matrix[3]

	// Finally, we can calculate the segment's position on the screen. See cell2pos function above.
	s.posX, s.posY = CellToPos(s.cx, s.cy, offsetX, offsetY)

	// We now need to decide which image the segment should use as its sprite.
	// Images for segment sprites follow the format 'segABCDE' where A is 0 or 1 depending on whether this is a
	// fast-moving segment, B is 0 or 1 depending on whether we currently have 1 or 2 health, C is whether this
	// is the head segment of a myriapod, D represents the direction we're facing (0 = up, 1 = top right,
	// up to 7 = top left) and E is how far we are through the walking animation (0 to 3)

	// Three variables go into the calculation of the direction. turn_idx tells us if we're making a turn in this
	// cell - and if so, whether we're turning clockwise or anti-clockwise. s.inEdge tells us which side of the
	// grid cell we entered from. And we can use SECONDARY_AXIS_SPEED[phase] to find out whether we should be facing
	// along the primary axis, secondary axis or diagonally between them.
	// (turn_idx - 2) gives 0 if straight, -1 if turning anti-clockwise, 1 if turning clockwise
	// Multiplying this by SECONDARY_AXIS_SPEED[phase] gives 0 if we're not doing a turn in this cell, or if
	// we are going to be turning but have not yet begun to turn. If we are doing a turn in this cell, and we're
	// at a phase where we should be showing a sprite with a new rotation, the result will be -1 or 1 if we're
	// currently in the first (45°) part of a turn, or -2 or 2 if we have turned 90°.
	// The next part of the calculation multiplies in_edge by 2 and then adds the result to the result of the previous
	// part. in_edge will be a number from 0 to 3, representing all possible directions in 90° increments.
	// It must be multiplied by two because the direction value we're calculating will be a number between 0 and 7,
	// representing all possible directions in 45° increments.
	// In the sprite filenames, the penultimate number represents the direction the sprite is facing, where a value
	// of zero means it's facing up. But in this code, if, for example, in_edge were zero, this means the segment is
	// coming from the top edge of its cell, and therefore should be facing down. So we add 4 to account for this.
	// After all this, we may have ended up with a number outside the desired range of 0 to 7. So the final step
	// is to MOD by 8.
	s.direction = Direction(((SecondaryAxisSpeed[phase] * (int(turnIdx) - 2)) + (int(s.inEdge) * 2) + 4) % 8)

	s.legFrame = phase / 4 // 16 phase cycle, 4 frames of animation
}

func (s *Segment) rank(proposedOutEdge Direction) int {
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
	// rock will either be the Rock object at the new grid cell, or nil.
	// It will be set to nil if there is no Rock object is at the new location, or if the new location is
	// outside the grid. We also have to account for the special case where the segment is off the left-hand
	// side of the screen on the first row, where it is initially created. We mustn't try to access that grid
	// cell
	var rock *Rock
	if out || (newCellY == 0 && newCellX < 0) {
		rock = nil
	} else {
		rock = s.game.grid[newCellY][newCellX]
	}

	rockPresent := rock != nil

	// Is new cell already occupied by another segment, or is another segment trying to enter my cell from
	// the opposite direction?
	occupiedBySegment := s.game.IsOccupied(newCellX, newCellY) || s.game.IsCellOccupied(Cell{s.cx, s.cy, proposedOutEdge})

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
	total := 0
	if out {
		total += 64
	}
	if turningBackOnSelf {
		total += 32
	}
	if directionDisallowed {
		total += 16
	}
	if occupiedBySegment {
		total += 8
	}
	if rockPresent {
		total += 4
	}
	if horizontalBlocked {
		total += 2
	}
	if sameAsPreviousXDirection {
		total += 1
	}
	return total
}
