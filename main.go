package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/img"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Could not initialize SDL: %v", err)
	}

	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("Could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Could not create Window: %v", err)
	}
	defer w.Destroy()

	err = drawTitle(r)
	if err != nil {
		return fmt.Errorf("Could not draw title: %v", err)
	}

	time.Sleep(5 * time.Second)

	err = drawBackground(r)
	if err != nil {
		return fmt.Errorf("Could not draw background: %v", err)
	}

	time.Sleep(5 * time.Second)

	return nil
}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()

	texture, err := img.LoadTexture(r, "res/images/background.png")
	if err != nil {
		return fmt.Errorf("Could not load background image: %v", err)
	}
	defer texture.Destroy()

	err = r.Copy(texture, nil, nil)
	if err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}
	r.Present()
	return nil
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	font, err := ttf.OpenFont("res/fonts/flappy-font.ttf", 20)
	if err != nil {
		return fmt.Errorf("Could not open font: %v", err)
	}
	defer font.Close()
	surface, err := font.RenderUTF8Solid("Flappy Gopher", sdl.Color{R: 255, G: 100, B: 0, A: 255})
	if err != nil {
		return fmt.Errorf("Could not render title: %v", err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer texture.Destroy()

	err = r.Copy(texture, nil, nil)
	if err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}
	r.Present()
	return nil
}
