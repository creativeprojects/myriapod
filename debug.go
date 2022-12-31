//go:build !prod

package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	DebugBuild = true
	Debug      = false
)

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := "\n\n\n TPS: %0.2f - time: %d \n Rocks: %d - Segments: %d - Bullets: %d - Explosions: %d - Occupation: %d\n%s\n%s"
	msg := fmt.Sprintf(template,
		ebiten.ActualTPS(),
		g.time,
		g.RockCount(),
		len(g.segments),
		len(g.bullets),
		len(g.explosions),
		len(g.occupation),
		g.player,
		g.enemy,
	)
	ebitenutil.DebugPrint(screen, msg)
}

// String returns a debug string
func (p *Player) String() string {
	return fmt.Sprintf(" Player lives: %d \n Player coordinates: %s",
		p.lives,
		p.sprite.String(),
	)
}

// String returns a debug string
func (e *FlyingEnemy) String() string {
	return fmt.Sprintf(" Enemy coordinates: %s",
		e.sprite.String(),
	)
}
