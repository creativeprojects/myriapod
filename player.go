package main

import (
	"math"
	"strconv"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	game      *Game
	sprite    *lib.Sprite
	life      *ebiten.Image
	images    [][]*ebiten.Image
	direction int
	frame     int
	lives     int
	alive     bool
	timer     int
	fireTimer int
	op        *ebiten.DrawImageOptions
}

func NewPlayer(game *Game) *Player {
	return &Player{
		game:   game,
		sprite: lib.NewSprite(lib.XCentre, lib.YCentre).MoveTo(PlayerSpawnX, PlayerSpawnY).SetImage(images["player00"]),
		life:   images["life"],
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
		op:        &ebiten.DrawImageOptions{},
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
	p.timer++
	if p.alive {
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

		x := p.sprite.X(lib.XCentre)
		y := p.sprite.Y(lib.YCentre)
		p.fireTimer--
		// Fire cannon (or allow firing animation to finish)
		if p.fireTimer < 0 && (p.frame > 0 || ebiten.IsKeyPressed(ebiten.KeySpace)) {
			if p.frame == 0 {
				// Create a bullet
				p.game.SoundEffect("laser0")
				p.game.Fire(x, y-8)
			}
			p.frame = (p.frame + 1) % 3
			p.fireTimer = ReloadTime
		}

		if p.game.enemy.Collision(x, y) {
			p.game.SoundEffect("player_explode0")
			p.game.Explosion(x, y, 1)
			p.alive = false
			p.timer = 0
			p.frame = 0
			p.lives--
			p.sprite.Animate([]*ebiten.Image{images["blank"], p.images[p.direction][p.frame]}, nil, 2, true)
		}
	} else {
		// player not alive
		if p.timer > RespawnTime {
			p.alive = true
			p.timer = 0
			p.sprite.MoveTo(240, 768)
			// Ensure there are no rocks at the player's respawn position
			p.game.ClearRocksForRespawn(240, 768)
		}
	}

	if p.timer > InvulnerabilityTime {
		p.sprite.Stop() // stop respawn animation (if started)
		p.sprite.SetImage(p.images[p.direction][p.frame])
	}
	p.sprite.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.sprite.Draw(screen)
	p.drawLives(screen)
	p.drawScore(screen)
}

func (p *Player) drawLives(screen *ebiten.Image) {
	// Display number of lives
	for i := 0; i < p.lives; i++ {
		p.op.GeoM.Reset()
		p.op.GeoM.Translate(float64(i)*40+8, 4)
		screen.DrawImage(p.life, p.op)
	}
}

func (p *Player) drawScore(screen *ebiten.Image) {
	// Display score
	score := strconv.Itoa(p.game.score)
	for i := 0; i < len(score); i++ {
		digit := string(score[len(score)-i-1])
		p.op.GeoM.Reset()
		p.op.GeoM.Translate(448-float64(i)*24, 5)
		screen.DrawImage(images["digit"+digit], p.op)
	}
}
