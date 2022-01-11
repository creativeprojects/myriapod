//go:build prod

package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	DebugBuild = false
	Debug      = false
)

func (g *Game) displayDebug(screen *ebiten.Image) {}

// String returns a debug string
func (p *Player) String() string {}

// String returns a debug string
func (e *FlyingEnemy) String() string {}
