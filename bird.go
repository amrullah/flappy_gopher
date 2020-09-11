package main

import (
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"

	"fmt"
)

const (
	gravity   = 1.0
	jumpSpeed = 10
)

type Bird struct {
	mu       sync.RWMutex
	textures []*sdl.Texture // 4 different frames to show, to give illusion of wing flapping
	time     int
	dead     bool

	x, y  int32
	w, h  int32
	speed float64
}

func (b *Bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *Bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	rect := &sdl.Rect{X: 10, Y: (600 - b.y), W: b.w, H: b.h}

	birdToShow := b.time % len(b.textures)
	err := r.Copy(b.textures[birdToShow], nil, rect)
	if err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	return nil
}

func (b *Bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = -jumpSpeed
}

func (b *Bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}
func (b *Bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, texture := range b.textures {
		texture.Destroy()
	}
}

func (b *Bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *Bird) touch(p *Pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w { // too far right
		return
	}
	if p.x+p.w < b.x { // too far left
		return
	}
	if !p.inverted && p.h < b.y-b.h/2 { // pipe below the bird
		return
	}
	if p.inverted && (600-p.h) > b.y+b.h/2 { // pipe above the bird
		return
	}

	b.dead = true
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

	return &Bird{textures: textures, x: 10, y: 300, w: 50, h: 43, speed: 0}, nil
}
