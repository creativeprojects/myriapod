package main

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Rock struct {
	game       *Game
	sprite     *lib.Sprite
	timer      int
	isTotem    bool
	rockType   int
	health     int
	showHealth int
	cellX      int
	cellY      int
	posX       float64
	posY       float64
}

func NewRock(game *Game, cellX, cellY int, isTotem bool) *Rock {
	health := 5
	showHealth := 5
	if !isTotem {
		health = rand.Intn(2) + 3
		showHealth = 1
	}
	posX, posY := CellToPos(cellX, cellY, 0, 0)
	return &Rock{
		game:       game,
		sprite:     lib.NewSprite(lib.XCentre, lib.YCentre).MoveTo(posX, posY),
		timer:      1,
		isTotem:    isTotem,
		rockType:   rand.Intn(4),
		health:     health,
		showHealth: showHealth,
		cellX:      cellX,
		cellY:      cellY,
		posX:       posX,
		posY:       posY,
	}
}

func (r *Rock) Damage(amount int, damagedByBullet bool) bool {
	// Damage can occur by being hit by bullets, or by being destroyed by a segment, or by being cleared from the
	// player's respawn location. Points can be earned by hitting special "totem" rocks, which have 5 health, but
	// this should only happen when they are hit by a bullet.
	if damagedByBullet && r.health == 5 {
		r.game.SoundEffect("totem_destroy0")
		r.game.AddScore(100)
	} else {
		if amount > r.health-1 {
			r.game.SoundEffect("rock_destroy0")
		} else {
			r.game.SoundEffect("hit" + strconv.Itoa(rand.Intn(4)))
		}
	}

	expType := 0
	if r.health == 5 {
		expType = 2
	}
	r.game.Explosion(r.sprite.X(lib.XCentre), r.sprite.Y(lib.YCentre), expType)
	r.health -= amount
	r.showHealth = r.health

	// Return false if we've lost all our health
	return r.health < 1
}

func (r *Rock) Update() {
	r.timer++
	// Every other frame, update the growing animation
	if r.timer%2 == 1 && r.showHealth < r.health {
		r.showHealth++
	}
	colour := max(r.game.wave, 0) % 3
	health := max(r.showHealth-1, 0)
	image := "rock" +
		strconv.Itoa(colour) +
		strconv.Itoa(r.rockType) +
		strconv.Itoa(health)
	r.sprite.SetImage(images[image])
}

func (r *Rock) Draw(screen *ebiten.Image) {
	if !r.sprite.HasImage() {
		log.Printf("No image for rock: totem=%v, type=%d, health=%d", r.isTotem, r.rockType, r.health)
		return
	}
	r.sprite.Draw(screen)
}

func (r *Rock) Y() float64 {
	return r.posY
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
