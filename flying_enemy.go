package main

import (
	"math"
	"math/rand"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type FlyingEnemy struct {
	sprite  *lib.Sprite
	images  [][]*ebiten.Image
	movingX float64
	dx      float64
	dy      float64
	color   int
	health  int
	timer   int
}

func NewFlyingEnemy() *FlyingEnemy {
	return &FlyingEnemy{
		sprite: lib.NewSprite(lib.XCentre, lib.YCentre),
		images: [][]*ebiten.Image{
			{images["meanie00"], images["meanie01"], images["meanie02"]},
			{images["meanie10"], images["meanie11"], images["meanie12"]},
			{images["meanie20"], images["meanie21"], images["meanie22"]},
		},
		movingX: 1,
		health:  1,
		timer:   0,
	}
}

func (e *FlyingEnemy) Start(playerX float64) {
	// Choose which side of the screen we start from.
	// Don't start right next to the player as that would be unfair
	// if not near player, start on a random side
	var side float64
	if playerX < 160 {
		side = 1
	} else if playerX > 320 {
		side = 0
	} else {
		side = math.Round(rand.Float64() * 2)
	}

	e.sprite.MoveTo(550*side-35, 688)

	// Always moves in the same X direction, but randomly pauses to just fly straight up or down
	e.movingX = 1                   // 0 if we're currently moving only vertically, 1 if moving along x axis (as well as y axis)
	e.dx = 1 - 2*side               // Move left or right depending on which side of the screen we're on
	e.dy = choice([]float64{-1, 1}) // Start moving either up or down
	e.color = rand.Intn(3)          // 3 different colours

	e.health = 1
	e.timer = 0
	e.sprite.Animate(e.images[e.color], []int{0, 2, 1, 2}, 4, true)
}

func (e *FlyingEnemy) IsInactive() bool {
	x := e.sprite.X(lib.XCentre)
	return e.health <= 0 || x < -35 || x > 515
}

func (e *FlyingEnemy) Collision(x, y float64) bool {
	if e.IsInactive() {
		return false
	}
	if e.sprite.CollidePoint(x, y) {
		e.health--
		return true
	}
	return false
}

func (e *FlyingEnemy) Update() {
	if e.IsInactive() {
		return
	}
	e.timer++

	// Move
	x := e.sprite.X(lib.XCentre) + e.dx*e.movingX*(3-math.Abs(e.dy))
	y := e.sprite.Y(lib.YCentre) + e.dy*(3-math.Abs(e.dx*e.movingX))

	e.sprite.MoveTo(x, y)

	if y < PlayerMinY || y > PlayerMaxY {
		// Gone too high or low - reverse y direction
		e.movingX = math.Round(rand.Float64())
		e.dy = -e.dy
	}

	e.sprite.Update()
}

func (e *FlyingEnemy) Draw(screen *ebiten.Image) {
	if e.IsInactive() {
		return
	}
	e.sprite.Draw(screen)
}

func choice(choices []float64) float64 {
	i := rand.Intn(len(choices))
	return choices[i]
}
