package main

type Cell struct {
	X    int
	Y    int
	Edge Direction
}

func (c Cell) Equal(cell Cell) bool {
	return c.X == cell.X && c.Y == cell.Y && c.Edge == cell.Edge
}
