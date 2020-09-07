package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"

	"fmt"
)

const gravity = 0.2

type Bird struct {
	textures []*sdl.Texture // 4 different frames to show, to give illusion of wing flapping
	time     int

	y, speed float64
}

func (b *Bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y -= b.speed
	if b.y < 0 {
		b.speed = -b.speed
		b.y = 0
	}
	b.speed += gravity

	rect := &sdl.Rect{X: 10, Y: (600 - int32(b.y)), W: 50, H: 43}

	birdToShow := b.time % len(b.textures)
	err := r.Copy(b.textures[birdToShow], nil, rect)
	if err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	return nil
}

func (b *Bird) destroy() {
	for _, texture := range b.textures {
		texture.Destroy()
	}
}
func newBird(r *sdl.Renderer) (*Bird, error) {
	var textures []*sdl.Texture

	for i := 1; i <= 4; i++ {
		filePath := fmt.Sprintf("res/images/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(r, filePath)
		if err != nil {
			return nil, fmt.Errorf("Could not load bird image: %v", err)
		}
		textures = append(textures, texture)
	}

	return &Bird{textures: textures, y: 300, speed: 0}, nil
}
