package main

import (
	//"strings"
	//"bytes"
	//"net/http"
	//"io/ioutil"
	//"math/rand"
	//"time"
	//"encoding/base64"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// go get -v github.com/veandco/go-sdl2/sdl

type Player struct {
	X     int32
	Y     int32
	Angle float32
	Speed int32
}

type Keyboard struct {
	KeyUp     int32
	KeyDown   int32
	KeyLeft   int32
	KeyRight  int32
	KeyAction int32
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

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}
	defer renderer.Destroy()

	// Define player and keyboard
	player := Player{0, 0, 0, 4}
	keyboard := Keyboard{0, 0, 0, 0, 0}

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
				keyPressed := t.State

				switch keyCode {

				case sdl.K_ESCAPE:
					running = false
					println("Esc key")

				case sdl.K_UP:
					println("Up key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyUp += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyUp = 0
					}

					//player.Y -= player.Speed

				case sdl.K_DOWN:
					println("down key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyDown += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyDown = 0
					}
					//player.Y += player.Speed

				case sdl.K_LEFT:
					println("left key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyLeft += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyLeft = 0
					}
					//player.X -= player.Speed

				case sdl.K_RIGHT:
					println("Right key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyRight += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyRight = 0
					}
					//player.X += player.Speed

				case sdl.K_RCTRL, sdl.K_LCTRL:
					println("ctrl key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyAction += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyAction = 0
					}

				case sdl.K_RSHIFT, sdl.K_LSHIFT:
					println("Modify speed")
					if keyPressed == sdl.PRESSED {
						player.Speed *= 2
					} else if keyPressed == sdl.RELEASED {
						player.Speed /= 2
					}

				} // END SWITCH
			}
		}

		////////////////////////////////////////////////////////////////////////////
		// UPDATE PLAYER
		////////////////////////////////////////////////////////////////////////////

		if keyboard.KeyUp > 0 {
			player.Y -= player.Speed
		}

		if keyboard.KeyDown > 0 {
			player.Y += player.Speed
		}

		if keyboard.KeyLeft > 0 {
			player.X -= player.Speed
		}

		if keyboard.KeyRight > 0 {
			player.X += player.Speed
		}

		////////////////////////////////////////////////////////////////////////////
		// UPDATE SCREEN
		////////////////////////////////////////////////////////////////////////////

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		rect := sdl.Rect{player.X, player.Y, 200, 200}
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.FillRect(&rect)

		////////////////////////////////////////////////////////////////////////////
		// UPDATE SDL WINDOW
		////////////////////////////////////////////////////////////////////////////

		renderer.Present()
		sdl.Delay(16)
	}

	defer window.Destroy()
}
