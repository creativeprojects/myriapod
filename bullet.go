package main

import (
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
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.done {
		return
	}
	b.sprite.Draw(screen)
}
