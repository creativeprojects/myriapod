package main

import (
	"math/rand"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	game   *Game
	sprite *lib.Sprite
	done   bool
}

func NewBullet(game *Game) *Bullet {
	return &Bullet{
		game:   game,
		sprite: lib.NewSprite(lib.XCentre, lib.YCentre).SetImage(images["bullet"]),
		done:   false,
	}
}

func (b *Bullet) Start(x, y float64) {
	b.done = false
	b.sprite.MoveTo(x, y)
}

func (b *Bullet) IsDone() bool {
	return b.done
}

func (b *Bullet) Update() {
	if b.done {
		return
	}

	b.sprite.Move(0, -24)

	x := b.sprite.X(lib.XCentre)
	y := b.sprite.Y(lib.YCentre)

	if y <= 0 {
		b.done = true
	}
	cellX, cellY := PosToCell(x, y)
	if b.game.Damage(cellX, cellY, 1, true) {
		// Hit a rock - destroy self
		b.done = true
		return
	}
	if b.game.enemy.Collision(x, y) {
		b.game.score += 20
		b.game.SoundEffect("meanie_explode0")
		b.game.Explosion(x, y, 2)
		b.done = true
		return
	}

	for i := 0; i < len(b.game.segments); i++ {
		if b.game.segments[i].Collision(x, y) {
			b.game.score += 10
			b.game.SoundEffect("segment_explode0")
			b.game.Explosion(x, y, 2)
			b.done = true
			if b.game.segments[i].health == 0 {
				if b.game.grid[cellY][cellX] == nil && b.game.AllowPlayerMovement2(b.game.player.sprite.X(lib.XCentre), b.game.player.sprite.Y(lib.YCentre), cellX, cellY) {
					// Create new rock - 20% chance of being a totem
					b.game.grid[cellY][cellX] = NewRock(b.game, cellX, cellY, rand.Float64() < .2)
				}
				b.game.segments[i] = nil
				b.game.segments = append(b.game.segments[:i], b.game.segments[i+1:]...)
			}
			return
		}
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.done {
		return
	}
	b.sprite.Draw(screen)
}
