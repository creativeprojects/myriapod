package main

import "github.com/cavern/creativeprojects/myriapod/lib"

type Bullet struct {
	sprite *lib.Sprite
}

func NewBullet() *Bullet {
	return &Bullet{
		sprite: lib.NewSprite(lib.XCentre, lib.YCentre),
	}
}

func (b *Bullet) Update() {
	b.sprite.Move(0, -24)
}
