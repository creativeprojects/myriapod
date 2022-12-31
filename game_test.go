package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPixelPosition(t *testing.T) {
	for x := 0; x < NumGridCols; x++ {
		for y := 0; y < NumGridRows; y++ {
			for offsetX := -16; offsetX < 16; offsetX++ {
				for offsetY := -16; offsetY < 16; offsetY++ {
					t.Run(fmt.Sprintf("X=%d,Y=%d,offset:X=%d,Y=%d", x, y, offsetX, offsetY), func(t *testing.T) {
						pixelX, pixelY := CellToPos(x, y, offsetX, offsetY)
						assert.LessOrEqual(t, pixelX, WindowWidth)
						assert.LessOrEqual(t, pixelY, WindowHeight)
						posX, posY := PosToCell(pixelX, pixelY)
						assert.Equal(t, x, posX)
						assert.Equal(t, y, posY)
					})
				}
			}
		}
	}
}

func BenchmarkPosToCell(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x, y := PosToCell(10.0, 11.0)
		assert.Equal(b, 0, x)
		assert.Equal(b, 0, y)
	}
}

func BenchmarkCellToPos(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x, y := CellToPos(10, 11, 12, 13)
		assert.EqualValues(b, 364, x)
		assert.EqualValues(b, 381, y)
	}
}
