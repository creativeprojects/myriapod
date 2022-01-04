package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Rock struct {
	timer      int
	isTotem    bool
	rockType   int
	health     int
	showHealth int
}

func NewRock(isTotem bool) *Rock {
	health := 5
	showHealth := 5
	if !isTotem {
		health = rand.Intn(1) + 3
		showHealth = 1
	}
	return &Rock{
		timer:      1,
		isTotem:    isTotem,
		rockType:   rand.Intn(3),
		health:     health,
		showHealth: showHealth,
	}
}

func (r *Rock) Update() {
	r.timer++
}

func (r *Rock) Draw(screen *ebiten.Image) {
}
