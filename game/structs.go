// File: structs.go (Raycaster)
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
