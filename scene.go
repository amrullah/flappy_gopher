package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time int
	bg   *sdl.Texture
	bird []*sdl.Texture // 4 different frames to show, to give illusion of wing flapping
}

func (s *scene) run(ctx context.Context, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		for range time.Tick(100 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}
func (s *scene) paint(r *sdl.Renderer) error {
	s.time++

	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 43}

	birdToShow := s.time % len(s.bird)
	err = r.Copy(s.bird[birdToShow], nil, rect)
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

	var bird []*sdl.Texture

	for i := 1; i <= 4; i++ {
		filePath := fmt.Sprintf("res/images/bird_frame_%d.png", i)
		birdFrame, err := img.LoadTexture(r, filePath)
		if err != nil {
			return nil, fmt.Errorf("Could not load bird image: %v", err)
		}
		bird = append(bird, birdFrame)
	}

	return &scene{bg: texture, bird: bird}, nil
}
