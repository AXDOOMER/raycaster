// File: gameloop.go (Raycaster)
// Copyright (C) 2021 Alexandre-Xavier Labonté-Lamoureux
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package game

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func Start() {
	////////////////////////////////////////////////////////////////////////////
	// INIT SDL, WINDOW, RENDERER, TEXTURE
	////////////////////////////////////////////////////////////////////////////
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Raycaster", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		960, 600, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	window.SetMinimumSize(320, 200)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}
	renderer.SetLogicalSize(320, 200)
	defer renderer.Destroy()

	// Create texture for intermediate rendering
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, 320, 200)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	////////////////////////////////////////////////////////////////////////////
	// INIT PLAYER STATE
	////////////////////////////////////////////////////////////////////////////

	// Define player and keyboard
	player := Player{22, 11.5, -1, 0, 0, 0.66, 0, 0, 0.05}
	keyboard := Keyboard{0, 0, 0, 0, 0, 0, 0}

	////////////////////////////////////////////////////////////////////////////
	// DECODE GAME TEXTURE
	////////////////////////////////////////////////////////////////////////////
	textureDecoder(rock_texture, "png", true, wall_texture[:])
	textureDecoder(clouds_texture, "jpg", false, sky_texture[:])
	textureDecoder(dirt_texture, "png", false, floor_texture[:])

	cycles := 0
	running := true
	for running {
		start := time.Now()

		////////////////////////////////////////////////////////////////////////////
		// GET KEYS AND EVENTS
		////////////////////////////////////////////////////////////////////////////
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym
				keyPressed := t.State

				switch keyCode {

				case sdl.K_ESCAPE:
					running = false

				case sdl.K_UP:
					if keyPressed == sdl.PRESSED {
						keyboard.KeyUp += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyUp = 0
					}

				case sdl.K_DOWN:
					if keyPressed == sdl.PRESSED {
						keyboard.KeyDown += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyDown = 0
					}

				case sdl.K_LEFT:
					if keyPressed == sdl.PRESSED {
						keyboard.KeyLeft += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyLeft = 0
					}

				case sdl.K_RIGHT:
					if keyPressed == sdl.PRESSED {
						keyboard.KeyRight += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyRight = 0
					}

				case sdl.K_RCTRL, sdl.K_LCTRL:
					if keyPressed == sdl.PRESSED {
						keyboard.KeyAction += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyAction = 0
					}

				case sdl.K_RSHIFT, sdl.K_LSHIFT:
					if keyPressed == sdl.PRESSED {
						player.Speed *= 2
					} else if keyPressed == sdl.RELEASED {
						player.Speed /= 2
					}

				case sdl.K_PAGEUP:
					if keyPressed == sdl.PRESSED {
						keyboard.KeyLookUp += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyLookUp = 0
					}

				case sdl.K_PAGEDOWN:
					if keyPressed == sdl.PRESSED {
						keyboard.KeyLookDown += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyLookDown = 0
					}

				} // END SWITCH
			}
		}

		////////////////////////////////////////////////////////////////////////////
		// UPDATE PLAYER
		////////////////////////////////////////////////////////////////////////////

		if keyboard.KeyUp > 0 {
			if worldmap[int32(player.PosX+player.DirX*player.Speed)][int32(player.PosY)] == 0 {
				player.PosX += player.DirX * player.Speed
			}

			if worldmap[int32(player.PosX)][int32(player.PosY+player.DirY*player.Speed)] == 0 {
				player.PosY += player.DirY * player.Speed
			}
		}

		if keyboard.KeyDown > 0 {
			if worldmap[int32(player.PosX-player.DirX*player.Speed)][int32(player.PosY)] == 0 {
				player.PosX -= player.DirX * player.Speed
			}

			if worldmap[int32(player.PosX)][int32(player.PosY-player.DirY*player.Speed)] == 0 {
				player.PosY -= player.DirY * player.Speed
			}
		}

		if keyboard.KeyRight > 0 {
			var oldDirX float64 = player.DirX
			var rotSpeed float64 = player.Speed
			player.Angle += rotSpeed
			player.DirX = player.DirX*math.Cos(-rotSpeed) - player.DirY*math.Sin(-rotSpeed)
			player.DirY = oldDirX*math.Sin(-rotSpeed) + player.DirY*math.Cos(-rotSpeed)
			var oldPlaneX float64 = player.PlaneX
			player.PlaneX = player.PlaneX*math.Cos(-rotSpeed) - player.PlaneY*math.Sin(-rotSpeed)
			player.PlaneY = oldPlaneX*math.Sin(-rotSpeed) + player.PlaneY*math.Cos(-rotSpeed)
		}

		if keyboard.KeyLeft > 0 {
			var oldDirX float64 = player.DirX
			var rotSpeed float64 = player.Speed
			player.Angle -= rotSpeed
			player.DirX = player.DirX*math.Cos(rotSpeed) - player.DirY*math.Sin(rotSpeed)
			player.DirY = oldDirX*math.Sin(rotSpeed) + player.DirY*math.Cos(rotSpeed)
			var oldPlaneX float64 = player.PlaneX
			player.PlaneX = player.PlaneX*math.Cos(rotSpeed) - player.PlaneY*math.Sin(rotSpeed)
			player.PlaneY = oldPlaneX*math.Sin(rotSpeed) + player.PlaneY*math.Cos(rotSpeed)
		}

		if keyboard.KeyLookUp > 0 {
			if player.LookY < 100 {
				player.LookY += 10
			}
		}
		if keyboard.KeyLookDown > 0 {
			if player.LookY > -100 {
				player.LookY -= 10
			}
		}

		////////////////////////////////////////////////////////////////////////////
		// CLEAR RENDERER
		////////////////////////////////////////////////////////////////////////////

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		////////////////////////////////////////////////////////////////////////////
		// RENDER GAME WORLD AND UPDATE RENDERER
		////////////////////////////////////////////////////////////////////////////

		renderSky(&player)
		renderFloors(&player)
		renderWalls(&player)
		renderMinimap(&player)

		texture.Update(nil, screenbuffer[:], 320*4)
		renderer.Copy(texture, nil, nil)

		////////////////////////////////////////////////////////////////////////////
		// UPDATE SDL WINDOW
		////////////////////////////////////////////////////////////////////////////

		renderer.Present()
		elapsed := int(time.Since(start).Milliseconds())

		if 16-elapsed < 0 {
			elapsed = 0
		} else {
			elapsed = 16 - elapsed
		}

		sdl.Delay(uint32(elapsed))
		cycles++
	}
}
