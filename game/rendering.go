// File: rendering.go (Raycaster)
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
	"math"
)

func putPixel(x int32, y int32, color uint32) {
	// ignore values that are out of range
	if x >= 0 && x < 320*screen_scaling {
		if y >= 0 && y < 200*screen_scaling {
			index := (y*320*screen_scaling + x) * 4
			screenbuffer[index+0] = uint8(color & 0xFF)
			screenbuffer[index+1] = uint8((color >> 8) & 0xFF)
			screenbuffer[index+2] = uint8((color >> 16) & 0xFF)
			screenbuffer[index+3] = uint8((color >> 24) & 0xFF)
		}
	}
}

func putPixelScaled(x int32, y int32, color uint32, multiplier int32) {
	xstop := x*multiplier + multiplier
	ystop := y*multiplier + multiplier
	for j := y * multiplier; j < ystop; j++ {
		for i := x * multiplier; i < xstop; i++ {
			putPixel(i, j, color)
		}
	}
}

func renderSky(player *Player) {
	// Do cylindrical projection?
	for y := 0; y < 200; y++ {
		height := int32(y+int(player.LookY)) - 100
		for x := 0; x < 320; x++ {
			slide := x + int(player.Angle*205)

			offset := slide % 640 /* 640 is the sky's horizontal resolution*/
			if offset < 0 {
				// This accounts for Go's modulo behavior
				offset += 640
			}

			var color uint32 = sky_texture[offset+y*640]
			putPixelScaled(int32(x), height, color, screen_scaling)
		}
	}
}

func renderMinimap(player *Player) {
	for y := 0; y < 24; y++ {
		for x := 0; x < 24; x++ {
			if worldmap[y][x] > 0 {
				var color uint32 = 0x00FF00FF
				putPixelScaled(int32(x), int32(y), color, screen_scaling)
			}
		}
	}

	var color uint32 = 0xFF0000FF
	putPixelScaled(int32(player.PosY), int32(player.PosX), color, screen_scaling)
}

func renderFloors(player *Player) {
	for y := int(100*screen_scaling) + int(player.LookY*screen_scaling); y < int(200*screen_scaling); y++ {
		rayDirX0 := player.DirX - player.PlaneX
		rayDirY0 := player.DirY - player.PlaneY
		rayDirX1 := player.DirX + player.PlaneX
		rayDirY1 := player.DirY + player.PlaneY

		// current pos compared to screen center
		p := y - int(200*screen_scaling)/2 - int(player.LookY*screen_scaling) + 1
		posZ := 0.5 * float64(200*screen_scaling)
		rowDistance := posZ / float64(p)

		// step vector on floor texture
		floorStepX := rowDistance * (rayDirX1 - rayDirX0) / float64(320*screen_scaling)
		floorStepY := rowDistance * (rayDirY1 - rayDirY0) / float64(320*screen_scaling)

		floorX := player.PosX + rowDistance*rayDirX0
		floorY := player.PosY + rowDistance*rayDirY0

		for x := 0; x < int(320*screen_scaling); x++ {
			cellX := int32(floorX)
			cellY := int32(floorY)

			tx := int32(64*(floorX-float64(cellX))) & (64 - 1)
			ty := int32(64*(floorY-float64(cellY))) & (64 - 1)

			floorX += floorStepX
			floorY += floorStepY

			var color uint32 = floor_texture[tx+ty*64]
			color = (color >> 1) & 0x7F7F7F7F
			putPixel(int32(x), int32(y), color)
		}
	}
}

func renderWalls(player *Player) {
	for x := 0; x < int(320*screen_scaling); x++ {
		var cameraX float64 = 2.0*float64(x)/float64(320*screen_scaling) - 1
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
		for !hit {
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

		lineHeight := int32(float64(200*screen_scaling) / perpWallDist)

		var drawStart int32 = -lineHeight/2 + 200*screen_scaling/2
		drawStart += player.LookY * screen_scaling
		if drawStart < 0 {
			drawStart = 0
		}
		var drawEnd int32 = lineHeight/2 + 200*screen_scaling/2
		drawEnd += player.LookY * screen_scaling
		if drawEnd >= 200*screen_scaling {
			drawEnd = 200 * screen_scaling /*- 1*/
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

		// screen texture to pixel stuff
		var step float64 = 1.0 * 64 / float64(lineHeight)
		var texPos float64 = float64((drawStart-player.LookY*screen_scaling)-200*screen_scaling/2+lineHeight/2) * step

		for y := drawStart; y < drawEnd; y++ {
			var texY int32 = int32(texPos) & (64 - 1)
			texPos += step
			var color uint32 = wall_texture[texX+texY*64]
			if side == 1 {
				color = (color >> 1) & 0x7F7F7F7F
			}

			putPixel(int32(x), int32(y), color)
		}

		// for each drawn slipe of wall, save its distance to the 1D depth buffer
		zbuffer[x] = perpWallDist
	}
}

func renderSprites(player *Player) {
	// Extra float precision is used to avoid graphic glitches
	var spriteX float64 = 5.0 - player.PosX
	var spriteY float64 = 5.0 - player.PosY
	var h float64 = float64(200 * screen_scaling)
	var w float64 = float64(320 * screen_scaling)

	// Compute the needed value for correct matrix multiplication
	// It's used to compute the transformed sprite via inverse camera matrix
	var invDet float64 = 1.0 / (player.PlaneX*player.DirY - player.DirX*player.PlaneY)

	// Compute the actual depth inside the screen, that's our Z in the 3D. Used to avoid fisheye effect.
	var transformX float64 = invDet * (player.DirY*spriteX - player.DirX*spriteY)
	var transformY float64 = invDet * (-player.PlaneY*spriteX + player.PlaneX*spriteY)

	var spriteScreenX float64 = float64((w / 2) * (1 + float64(transformX/transformY)))

	// Height of the sprite on screen, using transformY as real distance to prevent fisheye.
	var spriteHeight float64 = float64(math.Abs(float64(float64(h) / transformY)))

	// Calculate the lowest and highest pixel to fill in current stripe
	var drawStartY float64 = float64(-spriteHeight/2 + h/2)
	drawStartY += float64(player.LookY * screen_scaling)
	if drawStartY < 0 {
		drawStartY = 0
	}

	var drawEndY float64 = float64(spriteHeight/2 + h/2)
	drawEndY += float64(player.LookY * screen_scaling)
	if drawEndY >= h {
		drawEndY = h - 1
	}

	// Calculate width of the sprite at its distance
	var spriteWidth float64 = float64(math.Abs(float64(h) / (transformY)))
	var drawStartX float64 = -spriteWidth/2 + spriteScreenX

	if drawStartX < 0 {
		drawStartX = 0
	}

	var drawEndX float64 = spriteWidth/2 + spriteScreenX
	if drawEndX >= w {
		drawEndX = w - 1
	}

	// Draw sprite vertically, just like walls
	for column := float64(drawStartX); column < float64(drawEndX); column++ {
		var texX int32 = int32((column - (-spriteWidth/2 + spriteScreenX)) * 64 / spriteWidth)
		// The conditions:
		// 1     - It's in front of the plane (the player doesn't see things behind him)
		// 2 & 3 - It's on the screen
		// 4     - It's in front of the value in the z-buffer, meaning it's not hidden behind a wall
		if transformY > 0 && column > 0 && column < w && transformY < zbuffer[int32(column)] {
			// For every pixel in the column
			for y := drawStartY; y < drawEndY; y++ {
				var d float64 = y - float64(player.LookY*screen_scaling) - h*0.5 + spriteHeight*0.5
				var texY int32 = int32((d * 64) / spriteHeight) // 64 is sprite height
				var color uint32 = donald_texture[64*texY+texX]
				// Don't draw 100% invisible colors
				// TODO: putPixel could be improved to handle this by itself
				/*
					putPixelAlpha implementation
					i.e. compute a ratio between the colors according to the alpha value and that's the new color
					if (color & 0x000000FF) != 255 {
						newcolor = (3 * colorbuffer[y][column] / 4 + color / 4) // 25% transparency
						// OR
						newcolor = (colorbuffer[y][column] / 2 + color / 2) // 50% transparency
					} else normal rendering. Could also check if 255 because can skip the entire pixel.
				*/
				if (color & 0x000000FF) == 255 {
					putPixel(int32(column), int32(y), color)
				}
			}
		}
	}
}
