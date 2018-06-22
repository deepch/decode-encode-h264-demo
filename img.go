/*
 * Copyright (c) 2017 Simon Schmidt
 *
 * This software is provided 'as-is', without any express or implied
 * warranty. In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 *    claim that you wrote the original software. If you use this software
 *    in a product, an acknowledgment in the product documentation would be
 *    appreciated but is not required.
 * 2. Altered source versions must be plainly marked as such, and must not be
 *    misrepresented as being the original software.
 * 3. This notice may not be removed or altered from any source distribution.
 */

/*
A Component for working with Diffuse maps aka. Textures aka. Pictures.
*/
package main

import "image"
import "image/draw"
import "image/color"

type YuvOutput interface {
	SetYCbCr(x, y int, c color.YCbCr)
}
type nilYuvOutput struct{}

func (n nilYuvOutput) SetYCbCr(x, y int, c color.YCbCr) {}

type YuvInput interface {
	YCbCrAt(x, y int) color.YCbCr
}
type yuvConversion struct {
	image.Image
}

func (self yuvConversion) YCbCrAt(x, y int) color.YCbCr {
	return color.YCbCrModel.Convert(self.At(x, y)).(color.YCbCr)
}
func NewYuvInput(i image.Image) YuvInput {
	// image.YCbCr and image.NYCbCrA
	if y, ok := i.(YuvInput); ok {
		return y
	}
	return yuvConversion{i}
}

type yuvAdapter struct {
	draw.Image
}

func (self yuvAdapter) SetYCbCr(x, y int, c color.YCbCr) {
	self.Set(x, y, c)
}

type yuvStore struct {
	*image.YCbCr
}

func (self yuvStore) SetYCbCr(x, y int, c color.YCbCr) {
	yi := self.YOffset(x, y)
	ci := self.COffset(x, y)
	self.Y[yi] = c.Y
	self.Cb[ci] = c.Cb
	self.Cr[ci] = c.Cr
}

// Supports *image.YCbCr and anything that implements draw.Image!
func ToYuvOutput(i image.Image) YuvOutput {
	if y, ok := i.(*image.YCbCr); ok {
		return yuvStore{y}
	}
	//if y,ok := i.(*image.NYCbCrA); ok { return yuvStore{&(y.YCbCr)} }
	if d, ok := i.(draw.Image); ok {
		return yuvAdapter{d}
	}
	panic("unsupported")
}

type yuvUniform struct {
	color.YCbCr
}

func (yu yuvUniform) YCbCrAt(x, y int) color.YCbCr {
	return yu.YCbCr
}
