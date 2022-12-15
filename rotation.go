package main

var (
	RotationData = [][]int{
		{1, 0, 0, 1}, {0, -1, 1, 0}, {-1, 0, 0, -1}, {0, 1, -1, 0},
	}
	// This list represents how much the segment moves along the secondary axis, in situations where it makes two 45° turns
	// as described above. For the first four frames it doesn't move at all along the secondary axis. For the next eight
	// frames it moves at one pixel per frame, then for the last four frames it moves at two pixels per frame.
	SecondaryAxisSpeed = []int{
		0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
	}
	// This list stores the total secondary axis movement that will have occurred at each phase in the segment's movement
	// through the current grid cell (if the segment is turning)
	SecondaryAxisPositions = []int{
		0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 12, 14,
	}
)
