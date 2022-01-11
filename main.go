package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

var (
	images map[string]*ebiten.Image
	sounds map[string][]byte
)

func main() {
	var err error

	if DebugBuild {
		flag.BoolVar(&Debug, "d", false, "Debug mode")
		flag.Parse()
	}

	images, err = loadImages()
	if err != nil {
		log.Fatal(err)
	}

	audioContext := audio.NewContext(SampleRate)

	sounds, err = loadSounds()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle(WindowTitle)
	game, err := NewGame(audioContext)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
