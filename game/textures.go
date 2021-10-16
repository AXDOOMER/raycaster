// File: texture.go (Raycaster)
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
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/nfnt/resize"
)

func decodeBase64(data string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
}

func imageFormatDecode(data string, dataType string) (image.Image, error) {
	if dataType == "jpg" {
		return jpeg.Decode(decodeBase64(data))
	} else if dataType == "png" {
		return png.Decode(decodeBase64(data))
	}
	return nil, errors.New("invalid data type for image to be decoded")
}

func textureDecoder(data string, dataType string, rescale bool, dest []uint32) {
	// Decode image from base64 to image binary according to type
	img, err := imageFormatDecode(data, dataType)
	if err != nil {
		panic(err)
	}

	if rescale {
		img = resize.Resize(64, 64, img, resize.Bilinear)
	}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			color := img.At(x, y)
			rr, gg, bb, aa := color.RGBA()
			// pack values into pixels
			var pixel uint32 = ((rr >> 8) << 24) | ((gg >> 8) << 16) | ((bb >> 8) << 8) | aa>>8
			dest[x+y*img.Bounds().Max.X] = pixel
		}
	}
}
