package main

import (
	"math"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	game      *Game
	sprite    *lib.Sprite
	images    [][]*ebiten.Image
	direction int
	frame     int
	lives     int
	alive     bool
	timer     int
	fireTimer int
}

func NewPlayer(game *Game) *Player {
	return &Player{
		game:   game,
		sprite: lib.NewSprite(lib.XCentre, lib.YCentre).MoveTo(PlayerSpawnX, PlayerSpawnY),
		images: [][]*ebiten.Image{
			{images["player00"], images["player01"], images["player02"]},
			{images["player10"], images["player11"], images["player12"]},
			{images["player20"], images["player21"], images["player22"]},
			{images["player30"], images["player31"], images["player32"]},
		},
		direction: 0,
		frame:     0,
		lives:     3,
		alive:     true,
		timer:     0,
		fireTimer: 0,
	}
}

// Move the player sprite
// dx and dy are either 0, -1 or 1. speed is an integer indicating
// how many pixels we should move in the specified direction.
func (p *Player) Move(dx, dy float64, speed int) {
	for i := 0; i < speed; i++ {
		if p.game.AllowPlayerMovement(p.sprite.X(lib.XCentre)+dx, p.sprite.Y(lib.YCentre)+dy) {
			p.sprite.Move(float64(dx), float64(dy))
		}
	}
}

func (p *Player) Update() {
	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy = 1
	}
	if dx != 0 || dy != 0 {
		// Move in the relevant directions by the specified number of pixels. The purpose of 3 - abs(dy) is to
		// generate vectors which look either like (3,0) (which is 3 units long) or (2, 2) (which is sqrt(8) long)
		// so we move roughly the same distance regardless of whether we're travelling straight along the x or y axis.
		// or at 45 degrees. Without this, we would move noticeably faster when travelling diagonally.
		p.Move(dx, 0, int(3-math.Abs(dy)))
		p.Move(0, dy, int(3-math.Abs(dx)))
	}

	p.fireTimer--
	// Fire cannon (or allow firing animation to finish)
	if p.fireTimer < 0 && (p.frame > 0 || ebiten.IsKeyPressed(ebiten.KeySpace)) {
		if p.frame == 0 {
			// Create a bullet
			p.game.SoundEffect("laser0")
			p.game.Fire(p.sprite.X(lib.XCentre), p.sprite.Y(lib.YCentre)-8)
		}
		p.frame = (p.frame + 1) % 3
		p.fireTimer = ReloadTime
	}

	p.sprite.SetImage(p.images[p.direction][p.frame])
	p.sprite.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.sprite.Draw(screen)
}
