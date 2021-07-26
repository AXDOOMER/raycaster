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
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// go get -v github.com/veandco/go-sdl2/sdl

type Player struct {
	PosX   float64
	PosY   float64
	DirX   float64
	DirY   float64
	PlaneX float64
	PlaneY float64
	Angle  float64
	LookY  int32
	Speed  float64
}

type Keyboard struct {
	KeyUp       int32
	KeyDown     int32
	KeyLeft     int32
	KeyRight    int32
	KeyAction   int32
	KeyLookUp   int32
	KeyLookDown int32
}

var worldmap = [24][24]int32{
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 7, 7, 7, 7, 7, 7, 7, 7},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 7, 7, 0, 7, 7, 7, 7, 7},
	{4, 0, 5, 0, 0, 0, 0, 5, 0, 5, 0, 5, 0, 5, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 6, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 8, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 0, 5, 5, 5, 5, 7, 7, 7, 7, 7, 7, 7, 1},
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 6, 0, 6, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2},
	{4, 0, 0, 5, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3},
}

func main() {

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
	//renderer.SetIntegerScale(true)
	defer renderer.Destroy()

	// Virtual screen
	//virtual := sdl.Rect{0, 0, 320, 200}

	// Define player and keyboard
	player := Player{22, 11.5, -1, 0, 0, 0.66, 0, 0, 0.05}
	keyboard := Keyboard{0, 0, 0, 0, 0, 0, 0}

	drawsky(&player, renderer)

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

				case sdl.K_DOWN:
					println("down key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyDown += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyDown = 0
					}

				case sdl.K_LEFT:
					println("left key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyLeft += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyLeft = 0
					}

				case sdl.K_RIGHT:
					println("Right key")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyRight += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyRight = 0
					}

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

				case sdl.K_PAGEUP:
					println("page up")
					if keyPressed == sdl.PRESSED {
						keyboard.KeyLookUp += 1
					} else if keyPressed == sdl.RELEASED {
						keyboard.KeyLookUp = 0
					}

				case sdl.K_PAGEDOWN:
					println("page down")
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
			//player.PosY -= player.Speed

			if worldmap[int32(player.PosX+player.DirX*player.Speed)][int32(player.PosY)] == 0 {
				player.PosX += player.DirX * player.Speed
			}

			if worldmap[int32(player.PosX)][int32(player.PosY+player.DirY*player.Speed)] == 0 {
				player.PosY += player.DirY * player.Speed
			}
		}

		if keyboard.KeyDown > 0 {
			//player.PosY += player.Speed

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
			player.DirX = player.DirX*math.Cos(-rotSpeed) - player.DirY*math.Sin(-rotSpeed)
			player.DirY = oldDirX*math.Sin(-rotSpeed) + player.DirY*math.Cos(-rotSpeed)
			var oldPlaneX float64 = player.PlaneX
			player.PlaneX = player.PlaneX*math.Cos(-rotSpeed) - player.PlaneY*math.Sin(-rotSpeed)
			player.PlaneY = oldPlaneX*math.Sin(-rotSpeed) + player.PlaneY*math.Cos(-rotSpeed)
			//player.PosX += player.Speed
		}

		if keyboard.KeyLeft > 0 {
			//player.PosX -= player.Speed
			var oldDirX float64 = player.DirX
			var rotSpeed float64 = player.Speed
			player.DirX = player.DirX*math.Cos(rotSpeed) - player.DirY*math.Sin(rotSpeed)
			player.DirY = oldDirX*math.Sin(rotSpeed) + player.DirY*math.Cos(rotSpeed)
			var oldPlaneX float64 = player.PlaneX
			player.PlaneX = player.PlaneX*math.Cos(rotSpeed) - player.PlaneY*math.Sin(rotSpeed)
			player.PlaneY = oldPlaneX*math.Sin(rotSpeed) + player.PlaneY*math.Cos(rotSpeed)
		}

		if keyboard.KeyLookUp > 0 {
			if player.LookY < 145 {
				player.LookY += 10
			}
		}
		if keyboard.KeyLookDown > 0 {
			if player.LookY > -145 {
				player.LookY -= 10
			}
		}

		////////////////////////////////////////////////////////////////////////////
		// UPDATE SCREEN
		////////////////////////////////////////////////////////////////////////////

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		background := sdl.Rect{0, 0, 320, 200}
		renderer.SetDrawColor(0, 0, 128, 255)
		renderer.FillRect(&background)

		raycast(&player, renderer)
		drawmap(&player, renderer)

		//rect := sdl.Rect{int32(player.PosX), int32(player.PosY), 100, 100}
		//renderer.SetDrawColor(255, 0, 0, 255)
		//renderer.FillRect(&rect)

		////////////////////////////////////////////////////////////////////////////
		// UPDATE SDL WINDOW
		////////////////////////////////////////////////////////////////////////////

		renderer.Present()
		sdl.Delay(16)
	}
}

func drawsky(player *Player, renderer *sdl.Renderer) {
	var previous float64 = 0
	for i := 0; i < 320; i++ {
		x := math.Sin(float64(i) * math.Pi / 320.0)
		previous += x * (320.0 / 256.0)
		fmt.Println(i, " ", uint32(previous))
	}
}

func drawmap(player *Player, renderer *sdl.Renderer) {
	for y := 0; y < 24; y++ {
		for x := 0; x < 24; x++ {
			if worldmap[y][x] > 0 {
				renderer.SetDrawColor(0, 255, 0, 255)
				renderer.DrawPoint(int32(x), int32(y))
			}
		}
	}

	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawPoint(int32(player.PosY), int32(player.PosX))
}

func raycast(player *Player, renderer *sdl.Renderer) {
	var w int32 = 320
	for x := 0; x < 320; x++ {
		var cameraX float64 = 2.0*float64(x)/float64(w) - 1
		var rayDirX = player.DirX + player.PlaneX*cameraX
		var rayDirY = player.DirY + player.PlaneY*cameraX

		var mapX int32 = int32(player.PosX)
		var mapY int32 = int32(player.PosY)

		var sideDistX float64 = 0
		var sideDistY float64 = 0

		var deltaDistX = math.Abs(1.0 / rayDirX)
		var deltaDistY = math.Abs(1.0 / rayDirY)
		var perpWallDist float64 = 0

		var stepX int32
		var stepY int32

		var hit bool = false
		var side int8

		if rayDirX < 0 {
			stepX = -1
			sideDistX = (player.PosX - float64(mapX)) * deltaDistX
		} else {
			stepX = 1
			sideDistX = (float64(mapX) + 1.0 - player.PosX) * deltaDistX
		}

		if rayDirY < 0 {
			stepY = -1
			sideDistY = (player.PosY - float64(mapY)) * deltaDistY
		} else {
			stepY = 1
			sideDistY = (float64(mapY) + 1.0 - player.PosY) * deltaDistY
		}

		// DDA
		for hit == false {
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				mapX += stepX
				side = 0
			} else {
				sideDistY += deltaDistY
				mapY += stepY
				side = 1
			}

			// ray hit a wall?
			if worldmap[mapX][mapY] > 0 {
				hit = true
			}
		}

		// distance
		if side == 0 {
			perpWallDist = (float64(mapX) - player.PosX + (1-float64(stepX))/2) / rayDirX
		} else {
			perpWallDist = (float64(mapY) - player.PosY + (1-float64(stepY))/2) / rayDirY
		}

		lineHeight := int32(200 / perpWallDist)

		var drawStart int32 = -lineHeight/2 + 200/2
		drawStart += player.LookY
		if drawStart < 0 {
			drawStart = 0
		}
		var drawEnd int32 = lineHeight/2 + 200/2
		drawEnd += player.LookY
		if drawEnd >= 200 {
			drawEnd = 200 /*- 1*/
		}

		// texture calculations

		// calculate value of WallX
		var wallX float64 // where exactly the wall was hit
		if side == 0 {
			wallX = player.PosY + perpWallDist*rayDirY
		} else {
			wallX = player.PosX + perpWallDist*rayDirX
		}
		wallX -= math.Floor(wallX)

		// x coordinate of the texture
		var texX = int32(wallX * 64)
		if side == 0 && rayDirX > 0 {
			texX = 64 - texX - 1
		}
		if side == 1 && rayDirY < 0 {
			texX = 64 - texX - 1
		}

		// screen texture->pixel stuff
		//var step float64 = 1.0 * 64 / float64(lineHeight)
		//var texPos float64 = float64(drawStart-200/2+lineHeight/2) * step

		/*for y := drawStart; y < drawEnd; y++ {
			var color uint32 = 0xFFFF00FF
			if side == 1 {
				color = 0xAAAA00FF
			}

			renderer.SetDrawColor(uint8(color&0xFF000000>>24), uint8(color&0x00FF0000>>16), uint8(color&0x0000FF00>>8), uint8(color&0x000000FF))
			renderer.DrawPoint(int32(x), y)
		}*/

		var color uint32 = 0xFFFF00FF
		if side == 1 {
			color = 0xAAAA00FF
		}
		pixel := sdl.Rect{int32(x), drawStart, 1, drawEnd - drawStart}
		renderer.SetDrawColor(uint8(color&0xFF000000>>24), uint8(color&0x00FF0000>>16), uint8(color&0x0000FF00>>8), uint8(color&0x000000FF))
		renderer.FillRect(&pixel)
	}
}
