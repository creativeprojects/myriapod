package main

import (
	"math/rand"
	"sort"
	"time"

	"github.com/cavern/creativeprojects/myriapod/lib"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	healthTable = [][]int{{1, 1}, {1, 2}, {2, 2}, {1, 1}}
)

type Drawable interface {
	Draw(screen *ebiten.Image)
	Y() float64
}

type Game struct {
	audioContext *audio.Context
	musicPlayer  *AudioPlayer
	background   []*ebiten.Image
	state        GameState
	space        *lib.Sprite
	grid         [][]*Rock
	occupation   []Cell
	player       *Player
	enemy        *FlyingEnemy
	segments     []*Segment
	bullets      []*Bullet
	explosions   []*Explosion
	wave         int
	time         int
	score        int
	slow         bool
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
		background:   []*ebiten.Image{images["bg0"], images["bg1"], images["bg2"]},
		state:        StateMenu,
		space: lib.NewSprite(lib.XLeft, lib.YTop).MoveTo(0, 420).Animate([]*ebiten.Image{
			images["space0"], images["space1"], images["space2"], images["space3"], images["space4"],
			images["space5"], images["space6"], images["space7"], images["space8"], images["space9"],
			images["space10"], images["space11"], images["space12"], images["space13"],
		}, nil, 4, true),
	}

	return g.Initialize(), nil
}

// Initialize a new game
func (g *Game) Initialize() *Game {
	g.state = StateMenu
	g.wave = -1
	g.time = 0
	g.score = 0
	g.segments = make([]*Segment, 0, 20)
	g.bullets = make([]*Bullet, 0, 10)
	g.explosions = make([]*Explosion, 0, 10)
	// g.occupation = make([]Cell, StartSegments)
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

func (g *Game) AddScore(score int) {
	g.score += score
}

// Layout defines the size of the game in pixels
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func (g *Game) AllowPlayerMovement(x, y float64) bool {
	if x < PlayerMinX || x > PlayerMaxX || y < PlayerMinY || y > PlayerMaxY {
		return false
	}

	// get coordinates of corners of player sprite's collision rectangle
	x0, y0 := PosToCell(x-18, y-10)
	x1, y1 := PosToCell(x+18, y+10)

	// check each corner against grid
	for yi := y0; yi <= y1; yi++ {
		for xi := x0; xi <= x1; xi++ {
			if g.grid[yi][xi] != nil {
				return false
			}
		}
	}

	return true
}

func (g *Game) AllowPlayerMovement2(x, y float64, ax, ay int) bool {
	if x < PlayerMinX || x > PlayerMaxX || y < PlayerMinY || y > PlayerMaxY {
		return false
	}

	// get coordinates of corners of player sprite's collision rectangle
	x0, y0 := PosToCell(x-18, y-10)
	x1, y1 := PosToCell(x+18, y+10)

	// check each corner against grid
	for yi := y0; yi <= y1; yi++ {
		for xi := x0; xi <= x1; xi++ {
			if g.grid[yi][xi] != nil || xi == ax && yi == ay {
				return false
			}
		}
	}

	return true
}

// Damage returns whether or not there was a rock at this position
func (g *Game) Damage(cellX, cellY, amount int, fromBullet bool) bool {
	if cellY < 0 || cellX < 0 {
		return false
	}
	// Find the rock at this grid cell
	rock := g.grid[cellY][cellX]

	if rock == nil {
		return false
	}

	// rock.damage returns False if the rock has lost all its health
	// in this case, the grid cell will be set to nil
	if rock.Damage(amount, fromBullet) {
		g.grid[cellY][cellX] = nil
	}

	return true
}

func (g *Game) ClearRocksForRespawn(x, y float64) {
	// Destroy any rocks that might be overlapping with the player when they respawn
	// Could be more than one rock, hence the loop
	x0, y0 := PosToCell(x-18, y-10)
	x1, y1 := PosToCell(x+18, y+10)

	for yi := y0; yi <= y1; yi++ {
		for xi := x0; xi <= x1; xi++ {
			g.Damage(xi, yi, 5, false)
		}
	}
}

func (g *Game) IsOccupied(x, y int) bool {
	for _, cell := range g.occupation {
		if cell.X == x && cell.Y == y {
			return true
		}
	}
	return false
}

func (g *Game) IsCellOccupied(cell Cell) bool {
	for _, existingCell := range g.occupation {
		if existingCell.Equal(cell) {
			return true
		}
	}
	return false
}

// Update game events
func (g *Game) Update() error {
	g.time++
	if g.wave%4 == 3 {
		g.time++
	}

	// At the start of each frame, we reset occupied to be an empty set. As each individual myriapod segment is
	// updated, it will create entries in the occupied set to indicate that other segments should not attempt to
	// enter its current grid cell. There are two types of entries that are created in the occupied set. One is a
	// tuple consisting of a pair of numbers, representing grid cell coordinates. The other is a tuple consisting of
	// three numbers - the first two being grid cell coordinates, the third representing an edge through which a
	// segment is trying to enter a cell.
	// It is only used for myriapod segments - not rocks.
	g.occupation = make([]Cell, 0, StartSegments*20)

	if g.state == StateMenu {
		g.space.Update()
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Start()
		}
		return nil
	}

	if g.state == StatePlaying {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			Debug = !Debug
		}
		// toggle between slow and normal speed mode
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.slow = !g.slow
			if g.slow {
				ebiten.SetTPS(GameSlowSpeed)
			} else {
				ebiten.SetTPS(GameNormalSpeed)
			}
		}
		if g.enemy.IsInactive() {
			if rand.Float64() < .01 {
				g.enemy.Start(g.player.sprite.X(lib.XCentre))
			}
		}
		if len(g.segments) == 0 {
			if g.RockCount() <= InitialRockCount+g.wave {
				g.newRock()
			} else {
				// New wave and enough rocks - create a new myriapod
				g.SoundEffect("wave0")
				g.wave++
				g.time = 0
				numSegments := StartSegments + g.wave/4*2 // On the first four waves there are 8 segments - then 10, and so on
				for i := 0; i < numSegments; i++ {
					cellX, cellY := -1-i, 0
					// Determines whether segments take one or two hits to kill, based on the wave number.
					// e.g. on wave 0 all segments take one hit; on wave 1 they alternate between one and two hits
					health := healthTable[g.wave%4][i%2]
					fast := g.wave%4 == 3 // Every fourth myriapod moves faster than usual
					head := i == 0        // The first segment of each myriapod is the head
					segment := NewSegment(g, cellX, cellY, health, fast, head)
					g.segments = append(g.segments, segment)
				}
			}
		}
		g.updateSegments()
		g.updateBullets()
		g.updateExplosions()
		g.updateGrid()
		g.player.Update()
		g.enemy.Update()

		if g.player.lives == 0 && g.player.timer == 100 {
			g.SoundEffect("gameover")
			g.state = StateGameOver
		}
		return nil
	}

	if g.state == StateGameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Initialize()
		}
		return nil
	}
	return nil
}

// Draw game events
func (g *Game) Draw(screen *ebiten.Image) {
	if g.wave < 0 {
		screen.DrawImage(g.background[0], nil)
	} else {
		screen.DrawImage(g.background[g.wave%3], nil)
	}

	if g.state == StateMenu {
		screen.DrawImage(images["title"], nil)
		g.space.Draw(screen)
		return
	}

	if g.state == StatePlaying {
		g.drawObjects(screen)
		g.enemy.Draw(screen)
		if Debug {
			g.displayDebug(screen)
		}
		return
	}

	if g.state == StateGameOver {
		screen.DrawImage(images["over"], nil)
		return
	}
}

// drawObjects from top to bottom
func (g *Game) drawObjects(screen *ebiten.Image) {
	objects := make([]Drawable, 0, NumGridCols*NumGridRows)
	for _, row := range g.grid {
		for _, rock := range row {
			if rock != nil {
				objects = append(objects, rock)
			}
		}
	}
	for _, segment := range g.segments {
		if segment != nil {
			objects = append(objects, segment)
		}
	}
	for _, bullet := range g.bullets {
		if bullet != nil {
			objects = append(objects, bullet)
		}
	}
	for _, explosion := range g.explosions {
		if explosion != nil {
			objects = append(objects, explosion)
		}
	}
	objects = append(objects, g.player)
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Y() < objects[j].Y()
	})
	for _, object := range objects {
		object.Draw(screen)
	}
}

func (g *Game) Fire(x, y float64) {
	bullet := g.findAvailableBullet()
	if bullet == nil {
		bullet = NewBullet(g)
		g.bullets = append(g.bullets, bullet)
	}
	bullet.Start(x, y)
}

func (g *Game) findAvailableBullet() *Bullet {
	for _, bullet := range g.bullets {
		if bullet == nil {
			continue
		}
		if bullet.IsDone() {
			return bullet
		}
	}
	return nil
}

func (g *Game) Explosion(x, y float64, expType int) {
	explosion := g.findAvailableExplosion()
	if explosion == nil {
		explosion = NewExplosion()
		g.explosions = append(g.explosions, explosion)
	}
	explosion.Start(x, y, expType)
}

func (g *Game) findAvailableExplosion() *Explosion {
	for _, explosion := range g.explosions {
		if explosion == nil {
			continue
		}
		if explosion.IsDone() {
			return explosion
		}
	}
	return nil
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

func (g *Game) RockCount() int {
	count := 0
	for _, row := range g.grid {
		for _, element := range row {
			if element != nil {
				count++
			}
		}
	}
	return count
}

func (g *Game) updateGrid() {
	for _, row := range g.grid {
		for _, element := range row {
			if element != nil {
				element.Update()
			}
		}
	}
}
func (g *Game) updateBullets() {
	for _, bullet := range g.bullets {
		if bullet != nil {
			bullet.Update()
		}
	}
}

func (g *Game) updateExplosions() {
	for _, explosion := range g.explosions {
		if explosion != nil {
			explosion.Update()
		}
	}
}

func (g *Game) updateSegments() {
	for _, segment := range g.segments {
		if segment != nil {
			segment.Update()
		}
	}
}

func (g *Game) newRock() {
	// retry every time we pick coordinates that already contain a rock
	for {
		x := rand.Intn(NumGridCols)
		y := rand.Intn(NumGridRows-3) + 1 // Leave last 2 rows rock-free
		if rock := g.grid[y][x]; rock == nil {
			g.grid[y][x] = NewRock(g, x, y, false)
			return
		}
	}
}

func (g *Game) AddOccupation(cell1, cell2 Cell) {
	g.occupation = append(g.occupation, cell1, cell2)
}

// Convert a position in pixel units to a position in grid units.
// In this game, a grid square is 32 pixels.
func PosToCell(x, y float64) (int, int) {
	return (int(x) - 16) / 32, int(y) / 32
}

// Convert grid cell position to pixel coordinates, with a given offset
func CellToPos(cellX, cellY, XOffset, YOffset int) (float64, float64) {
	// If the requested offset is zero, returns the centre of the requested cell, hence the +16.
	// In the case of the X axis, there's a 16 pixel border at the
	// left and right of the screen, hence +16 becomes +32.
	return float64((cellX * 32) + 32 + XOffset), float64((cellY * 32) + 16 + YOffset)
}
