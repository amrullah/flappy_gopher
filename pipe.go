package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Pipe struct {
	mu       sync.RWMutex
	x        int32
	h        int32
	w        int32
	speed    int32
	inverted bool
	texture  *sdl.Texture
}

func newPipe(r *sdl.Renderer) (*Pipe, error) {
	texture, err := img.LoadTexture(r, "res/images/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load Pipe image: %v", err)
	}

	return &Pipe{
		x:        400,
		h:        300,
		w:        50,
		speed:    10,
		inverted: false,
		texture:  texture,
	}, nil
}

func (p *Pipe) update() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x -= p.speed
}

func (p *Pipe) restart() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x = 400
}

func (p *Pipe) paint(r *sdl.Renderer) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}

	err := r.Copy(p.texture, nil, rect)
	if err != nil {
		return fmt.Errorf("Could not copy pipe: %v", err)
	}
	return nil
}

func (p *Pipe) destroy() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.texture.Destroy()
}
