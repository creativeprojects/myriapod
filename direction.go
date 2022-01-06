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
