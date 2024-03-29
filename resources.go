package main

import (
	"embed"
	"fmt"
	"image"
	"io/fs"
	"path"
	"strings"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed images sounds music
var embeddedFiles embed.FS

func loadImages() (map[string]*ebiten.Image, error) {
	imageNames, err := fs.Glob(embeddedFiles, "images/*.png")
	if err != nil {
		return nil, err
	}
	imagesMap := make(map[string]*ebiten.Image, len(imageNames))
	for _, imageName := range imageNames {
		file, err := embeddedFiles.Open(imageName)
		if err != nil {
			return imagesMap, fmt.Errorf("%s: %w", imageName, err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			return imagesMap, fmt.Errorf("%s: %w", imageName, err)
		}
		img2 := ebiten.NewImageFromImage(img)
		if err != nil {
			return imagesMap, fmt.Errorf("%s: %w", imageName, err)
		}
		imageName = path.Base(imageName)
		imageName = strings.TrimSuffix(imageName, path.Ext(imageName))
		imagesMap[imageName] = img2
	}
	return imagesMap, nil
}

func loadSounds() (map[string][]byte, error) {
	soundNames, err := fs.Glob(embeddedFiles, "sounds/*.ogg")
	if err != nil {
		return nil, err
	}
	soundsMap := make(map[string][]byte, len(soundNames))
	for _, soundName := range soundNames {
		file, err := embeddedFiles.Open(soundName)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		snd, err := vorbis.DecodeWithSampleRate(SampleRate, file)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		buf := make([]byte, snd.Length())
		_, err = snd.Read(buf)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		soundName = path.Base(soundName)
		soundName = strings.TrimSuffix(soundName, path.Ext(soundName))
		soundsMap[soundName] = buf
	}
	return soundsMap, nil
}
