// Adapted from github.com/RobCherry/vibrant
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Rob Cherry
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package basiccolor

import (
	"image/color"
	"math"
)

var _ color.Color = HSL{}

// HSL represents the HSL value for a color.
type HSL struct {
	H, S, L float64
	A       uint8
}

func (c HSL) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := hslToNRGB(c.H, c.S, c.L)
	return color.NRGBA{r, g, b, c.A}.RGBA()
}

// HSLModel is the color.Model for the HSL type.
var HSLModel = color.ModelFunc(func(c color.Color) color.Color {
	if _, ok := c.(HSL); ok {
		return c
	}
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	h, s, l := nrgbToHSL(nrgba.R, nrgba.G, nrgba.B)
	return HSL{h, s, l, nrgba.A}
})

// Returns the Hue [0..360], Saturation and Lightness [0..1] of the color.
func nrgbToHSL(r, g, b uint8) (h, s, l float64) {
	fr := float64(r) / 255.0
	fg := float64(g) / 255.0
	fb := float64(b) / 255.0

	min := math.Min(math.Min(fr, fg), fb)
	max := math.Max(math.Max(fr, fg), fb)

	l = (max + min) / 2
	if min == max {
		s = 0
		h = 0
	} else {
		if l < 0.5 {
			s = (max - min) / (max + min)
		} else {
			s = (max - min) / (2.0 - max - min)
		}
		if max == fr {
			h = (fg - fb) / (max - min)
		} else if max == fg {
			h = 2.0 + (fb-fr)/(max-min)
		} else {
			h = 4.0 + (fr-fg)/(max-min)
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	return h, s, l
}

// Returns the RGB [0..255] values given a Hue [0..360], Saturation and Lightness [0..1]
func hslToNRGB(h, s, l float64) (uint8, uint8, uint8) {
	if s == 0 {
		return clampUint8(uint8(roundFloat64(l*255.0)), 0, 255), clampUint8(uint8(roundFloat64(l*255.0)), 0, 255), clampUint8(uint8(roundFloat64(l*255.0)), 0, 255)
	}

	var (
		r, g, b    float64
		t1, t2     float64
		tr, tg, tb float64
	)

	if l < 0.5 {
		t1 = l * (1.0 + s)
	} else {
		t1 = l + s - l*s
	}

	t2 = 2*l - t1
	h = h / 360
	tr = h + 1.0/3.0
	tg = h
	tb = h - 1.0/3.0

	if tr < 0 {
		tr++
	}
	if tr > 1 {
		tr--
	}
	if tg < 0 {
		tg++
	}
	if tg > 1 {
		tg--
	}
	if tb < 0 {
		tb++
	}
	if tb > 1 {
		tb--
	}

	// Red
	if 6*tr < 1 {
		r = t2 + (t1-t2)*6*tr
	} else if 2*tr < 1 {
		r = t1
	} else if 3*tr < 2 {
		r = t2 + (t1-t2)*(2.0/3.0-tr)*6
	} else {
		r = t2
	}

	// Green
	if 6*tg < 1 {
		g = t2 + (t1-t2)*6*tg
	} else if 2*tg < 1 {
		g = t1
	} else if 3*tg < 2 {
		g = t2 + (t1-t2)*(2.0/3.0-tg)*6
	} else {
		g = t2
	}

	// Blue
	if 6*tb < 1 {
		b = t2 + (t1-t2)*6*tb
	} else if 2*tb < 1 {
		b = t1
	} else if 3*tb < 2 {
		b = t2 + (t1-t2)*(2.0/3.0-tb)*6
	} else {
		b = t2
	}

	return clampUint8(uint8(roundFloat64(r*255.0)), 0, 255), clampUint8(uint8(roundFloat64(g*255.0)), 0, 255), clampUint8(uint8(roundFloat64(b*255.0)), 0, 255)
}

func clampUint8(value, min, max uint8) uint8 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Utility function for rounding the value of a number.
func roundFloat64(value float64) float64 {
	if value < 0.0 {
		return value - 0.5
	}
	return value + 0.5
}
