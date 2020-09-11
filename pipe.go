package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Pipes struct {
	mu sync.RWMutex

	speed   int32
	texture *sdl.Texture
	pipes   []*Pipe
}

func newPipes(r *sdl.Renderer) (*Pipes, error) {
	texture, err := img.LoadTexture(r, "res/images/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load Pipe image: %v", err)
	}

	ps := &Pipes{
		speed:   8,
		texture: texture,
	}

	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(2 * time.Second)
		}
	}()
	return ps, nil
}

func (ps *Pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {

		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}
	// rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	// flip := sdl.FLIP_NONE
	// if p.inverted {
	// 	rect.Y = 0
	// 	flip = sdl.FLIP_VERTICAL
	// }
	// err := r.CopyEx(p.texture, nil, rect, 0, nil, flip)
	// if err != nil {
	// 	return fmt.Errorf("Could not copy pipe: %v", err)
	// }
	return nil
}

func (ps *Pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *Pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var remaining []*Pipe
	for _, p := range ps.pipes {
		p.mu.Lock()
		p.x -= ps.speed
		p.mu.Unlock()
		if p.x+p.w > 0 {
			remaining = append(remaining, p)
		}
	}
	ps.pipes = remaining
}

func (ps *Pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}

func (ps *Pipes) touch(bird *Bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.touch(bird)
	}
}

type Pipe struct {
	mu sync.RWMutex

	h        int32
	w        int32
	x        int32
	inverted bool
}

func newPipe() *Pipe {
	// texture, err := img.LoadTexture(r, "res/images/pipe.png")
	// if err != nil {
	// return nil, fmt.Errorf("Could not load Pipe image: %v", err)
	// }

	return &Pipe{
		x:        800,
		h:        int32(100 + rand.Intn(300)),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}
}

// func (p *Pipe) restart() {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	p.x = 400
// }

func (p *Pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}
	err := r.CopyEx(texture, nil, rect, 0, nil, flip)
	if err != nil {
		return fmt.Errorf("Could not copy pipe: %v", err)
	}
	return nil
}

func (p *Pipe) touch(bird *Bird) {
	p.mu.RLock()
	p.mu.RUnlock()

	bird.touch(p)
}

// func (p *Pipe) destroy() {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// }
