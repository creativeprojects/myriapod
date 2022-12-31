package main

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

func (d Direction) IsHorizontal() bool {
	return d == DirectionLeft || d == DirectionRight
}

func (d Direction) Inverse() Direction {
	switch d {
	case DirectionUp:
		return DirectionDown
	case DirectionRight:
		return DirectionLeft
	case DirectionDown:
		return DirectionUp
	case DirectionLeft:
		return DirectionRight
	}
	return 0
}

var (
	// X and Y directions indexed by in_edge and out_edge in Segment
	// The indices correspond to the direction numbers above, i.e. 0 = up, 1 = right, 2 = down, 3 = left
	DX = []int{0, 1, 0, -1}
	DY = []int{-1, 0, 1, 0}
)
