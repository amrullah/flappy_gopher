package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time int
	bg   *sdl.Texture
	bird *Bird
	pipes *Pipes
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(100 * time.Millisecond)

		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(1 * time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}
func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseMotionEvent, *sdl.WindowEvent:
	case *sdl.KeyboardEvent:
		s.bird.jump()
	default:
		log.Printf("Unknown Event: %T", event)
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
	// s.bird.touch(s.pipe)
}
func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}

	if err := s.pipes.paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}
func newScene(r *sdl.Renderer) (*scene, error) {
	texture, err := img.LoadTexture(r, "res/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load background image: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, err
	}
	pipes, err := newPipes(r)
	if err != nil {
		return nil, err
	}
	return &scene{bg: texture, bird: bird, pipes: pipes}, nil
}
