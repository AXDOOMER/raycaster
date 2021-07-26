package main

import (
	//"strings"
	//"bytes"
	//"net/http"
	//"io/ioutil"
	//"math/rand"
	//"time"
	//"encoding/base64"
	"github.com/veandco/go-sdl2/sdl"
)

// go get -v github.com/veandco/go-sdl2/sdl

type Player struct {
	X     int32
	Y     int32
	Angle float32
	Speed int32
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Raycaster", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	// Define player
	player := Player{0, 0, 0, 4}

	running := true
	for running {
		////////////////////////////////////////////////////////////////////////////
		// GET KEYS AND EVENTS
		////////////////////////////////////////////////////////////////////////////
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false

			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym
				switch keyCode {

				case sdl.K_ESCAPE:
					running = false
					println("Esc key")

				case sdl.K_UP:
					println("Up key")
					player.Y -= player.Speed

				case sdl.K_DOWN:
					println("down key")
					player.Y += player.Speed

				case sdl.K_LEFT:
					println("left key")
					player.X -= player.Speed

				case sdl.K_RIGHT:
					println("Right key")
					player.X += player.Speed

				case sdl.K_RCTRL, sdl.K_LCTRL:
					println("ctrl key")
				}
			}
		}

		////////////////////////////////////////////////////////////////////////////
		// UPDATE SCREEN
		////////////////////////////////////////////////////////////////////////////

		rect := sdl.Rect{player.X, player.Y, 200, 200}
		surface.FillRect(&rect, 0xffff0000)

		////////////////////////////////////////////////////////////////////////////
		// UPDATE SDL WINDOW
		////////////////////////////////////////////////////////////////////////////

		window.UpdateSurface()
		surface.FillRect(nil, 0)
	}

	defer window.Destroy()
}
