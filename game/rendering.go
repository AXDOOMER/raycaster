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
	if x >= 0 && x < 320 {
		if y >= 0 && y < 200 {
			index := (y*320 + x) * 4
			screenbuffer[index+0] = uint8(color & 0xFF)
			screenbuffer[index+1] = uint8((color >> 8) & 0xFF)
			screenbuffer[index+2] = uint8((color >> 16) & 0xFF)
			screenbuffer[index+3] = uint8((color >> 24) & 0xFF)
		}
	}
}

func renderSky(player *Player) {
	// Do cylindrical projection?
	for x := 0; x < 320; x++ {
		for y := 0; y < 200; y++ {
			slide := x + int(player.Angle*205)

			offset := slide % 640	/* 640 is the sky's horizontal resolution*/
			if offset < 0 {
				// This accounts for Go's modulo behavior
				offset += 640
			}

			var color uint32 = sky_texture[offset+y*640]
			putPixel(int32(x), int32(y+int(player.LookY))-100, color)
		}
	}
}

func renderMinimap(player *Player) {
	for y := 0; y < 24; y++ {
		for x := 0; x < 24; x++ {
			if worldmap[y][x] > 0 {
				var color uint32 = 0x00FF00FF
				putPixel(int32(x), int32(y), color)
			}
		}
	}

	var color uint32 = 0xFF0000FF
	putPixel(int32(player.PosY), int32(player.PosX), color)
}

func renderFloors(player *Player) {
	for y := 100 + int(player.LookY); y < 200; y++ {
		rayDirX0 := player.DirX - player.PlaneX
		rayDirY0 := player.DirY - player.PlaneY
		rayDirX1 := player.DirX + player.PlaneX
		rayDirY1 := player.DirY + player.PlaneY

		// current pos compared to screen center
		p := y - 200/2 - int(player.LookY) + 1
		posZ := 0.5 * 200
		rowDistance := posZ / float64(p)

		// step vector on floor texture
		floorStepX := rowDistance * (rayDirX1 - rayDirX0) / 320
		floorStepY := rowDistance * (rayDirY1 - rayDirY0) / 320

		floorX := player.PosX + rowDistance*rayDirX0
		floorY := player.PosY + rowDistance*rayDirY0

		for x := 0; x < 320; x++ {
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

		// screen texture to pixel stuff
		var step float64 = 1.0 * 64 / float64(lineHeight)
		var texPos float64 = float64((drawStart-player.LookY)-200/2+lineHeight/2) * step

		for y := drawStart; y < drawEnd; y++ {
			var texY int32 = int32(texPos) & (64 - 1)
			texPos += step
			var color uint32 = wall_texture[texX+texY*64]
			if side == 1 {
				color = (color >> 1) & 0x7F7F7F7F
			}

			putPixel(int32(x), int32(y), color)
		}
	}
}
