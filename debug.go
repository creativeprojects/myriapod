//go:build !prod

package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Debug = true

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := " TPS: %0.2f - Rocks: %d - Bullets: %d - Explosions: %d\n%s\n%s"
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
		g.RockCount(),
		len(g.bullets),
		len(g.explosions),
		g.player,
		g.enemy,
	)
	ebitenutil.DebugPrint(screen, msg)
}

// String returns a debug string
func (p *Player) String() string {
	return fmt.Sprintf(" Player lives %d \n Player coordinates: %s",
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
