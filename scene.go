package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg   *sdl.Texture
	bird *sdl.Texture
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 43}
	err = r.Copy(s.bird, nil, rect)
	if err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
}
func newScene(r *sdl.Renderer) (*scene, error) {
	texture, err := img.LoadTexture(r, "res/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load background image: %v", err)
	}

	bird, err := img.LoadTexture(r, "res/images/bird_frame_1.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load bird image: %v", err)
	}

	return &scene{bg: texture, bird: bird}, nil
}
