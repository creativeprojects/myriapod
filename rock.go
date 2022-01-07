package main

import (
	"math/rand"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Rock struct {
	sprite     *lib.Sprite
	timer      int
	isTotem    bool
	rockType   int
	health     int
	showHealth int
	x          int
	y          int
}

func NewRock(x, y int, isTotem bool) *Rock {
	health := 5
	showHealth := 5
	if !isTotem {
		health = rand.Intn(2) + 3
		showHealth = 1
	}
	posX, posY := CellToPos(x, y, 0, 0)
	return &Rock{
		sprite:     lib.NewSprite(lib.XCentre, lib.YCentre).MoveTo(posX, posY),
		timer:      1,
		isTotem:    isTotem,
		rockType:   rand.Intn(4),
		health:     health,
		showHealth: showHealth,
		x:          x,
		y:          y,
	}
}

func (r *Rock) Update() {
	r.timer++
	r.sprite.SetImage(images["rock000"])
}

func (r *Rock) Draw(screen *ebiten.Image) {
	r.sprite.Draw(screen)
}
