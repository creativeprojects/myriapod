package main

import (
	"math/rand"
	"time"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	audioContext *audio.Context
	musicPlayer  *AudioPlayer
	state        GameState
	space        *lib.Sprite
	grid         [][]*Rock
	player       *Player
	enemy        *FlyingEnemy
	bullets      []*Bullet
	wave         int
	time         int
	score        int
}

// NewGame creates a new game instance and prepares a demo AI game
func NewGame(audioContext *audio.Context) (*Game, error) {
	m, err := NewAudioPlayer(audioContext)
	if err != nil {
		return nil, err
	}

	g := &Game{
		audioContext: audioContext,
		musicPlayer:  m,
		state:        StateMenu,
		space: lib.NewSprite(lib.XLeft, lib.YTop).MoveTo(0, 420).Animate([]*ebiten.Image{
			images["space0"], images["space1"], images["space2"], images["space3"], images["space4"],
			images["space5"], images["space6"], images["space7"], images["space8"], images["space9"],
			images["space10"], images["space11"], images["space12"], images["space13"],
		}, nil, 4, true),
		bullets: make([]*Bullet, 0, 10),
	}

	return g.Initialize(), nil
}

// Initialize a new game
func (g *Game) Initialize() *Game {
	g.state = StateMenu
	g.wave = -1
	g.time = 0
	g.score = 0
	g.space.Start()
	return g
}

func (g *Game) Start() {
	rand.Seed(time.Now().UnixNano())
	g.newGrid()
	g.player = NewPlayer(g)
	g.enemy = NewFlyingEnemy()
	g.enemy.Start(g.player.sprite.X(lib.XCentre))
	g.state = StatePlaying
}

// Layout defines the size of the game in pixels
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func (g *Game) AllowPlayerMovement(x, y float64) bool {
	if x < PlayerMinX || x > PlayerMaxX || y < PlayerMinY || y > PlayerMaxY {
		return false
	}
	return true
}

// Update game events
func (g *Game) Update() error {
	if g.state == StateMenu {
		g.space.Update()
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Start()
		}
		return nil
	}

	if g.state == StatePlaying {
		if g.enemy.IsInactive() {
			if rand.Float64() < .01 {
				g.enemy.Start(g.player.sprite.X(lib.XCentre))
			}
		}
		g.player.Update()
		g.enemy.Update()
	}
	return nil
}

// Draw game events
func (g *Game) Draw(screen *ebiten.Image) {
	if g.wave < 0 {
		screen.DrawImage(images["bg0"], nil)
	}
	if g.state == StateMenu {
		screen.DrawImage(images["title"], nil)
		g.space.Draw(screen)
	}

	if g.state == StatePlaying {
		g.player.Draw(screen)
		g.enemy.Draw(screen)
		g.displayDebug(screen)
	}
}

func (g *Game) Fire(x, y float64) {
	//
}

func (g *Game) SoundEffect(name string) {
	PlaySE(g.audioContext, sounds[name])
}

// newGrid creates a new empty grid
func (g *Game) newGrid() {
	g.grid = make([][]*Rock, NumGridRows)
	for i := range g.grid {
		g.grid[i] = make([]*Rock, NumGridCols)
	}
}

// Convert a position in pixel units to a position in grid units.
// In this game, a grid square is 32 pixels.
func PosToCell(x, y int) (int, int) {
	return (x - 16) / 32, y / 32
}

// Convert grid cell position to pixel coordinates, with a given offset
func CellToPos(cellX, cellY, XOffset, YOffset int) (int, int) {
	// If the requested offset is zero, returns the centre of the requested cell, hence the +16.
	// In the case of the X axis, there's a 16 pixel border at the
	// left and right of the screen, hence +16 becomes +32.
	return (cellX * 32) + 32 + XOffset, (cellY * 32) + 16 + YOffset
}
