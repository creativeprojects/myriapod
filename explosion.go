package main

import (
	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Explosion struct {
	sprite *lib.Sprite
	images [][]*ebiten.Image
	done   bool
}

func NewExplosion() *Explosion {
	return &Explosion{
		sprite: lib.NewSprite(lib.XCentre, lib.YCentre),
		images: [][]*ebiten.Image{
			{images["exp00"], images["exp01"], images["exp02"], images["exp03"], images["exp04"], images["exp05"], images["exp06"], images["exp07"]},
			{images["exp10"], images["exp11"], images["exp12"], images["exp13"], images["exp14"], images["exp15"], images["exp16"], images["exp17"]},
			{images["exp20"], images["exp21"], images["exp22"], images["exp23"], images["exp24"], images["exp25"], images["exp26"], images["exp27"]},
		},
	}
}

func (e *Explosion) Start(x, y float64, expType int) {
	e.sprite.MoveTo(x, y)
	e.sprite.Animate(e.images[expType], nil, 4, false)
	e.done = false
}

func (e *Explosion) Update() {
	if e.done {
		return
	}
	e.sprite.Update()
	if e.sprite.IsFinished() {
		e.done = true
	}
}

func (e *Explosion) Draw(screen *ebiten.Image) {
	if e.done {
		return
	}
	e.sprite.Draw(screen)
}

func (e *Explosion) Y() float64 {
	return e.sprite.RawY()
}

func (e *Explosion) IsDone() bool {
	return e.done
}
