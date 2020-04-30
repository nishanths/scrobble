package basiccolor

import (
	"fmt"
	"image/color"
)

//go:generate stringer -type=Color

// Color represents a basic color.
type Color int

const (
	Red Color = iota
	Orange
	Brown
	Yellow
	Green
	Blue
	Violet
	Pink
	Black
	Gray
	White
)

// Closest returns the closest basic color to c.
func Closest(c color.Color) Color {
	hsl := HSLModel.Convert(c).(HSL)
	hue, saturation, lightness := hsl.H, hsl.S, hsl.L

	if isWhite(saturation, lightness) {
		return White
	}
	if isBlack(saturation, lightness) {
		return Black
	}
	if isGray(saturation, lightness) {
		return Gray
	}
	if r, ok := determineColor(hue, saturation, lightness); ok {
		return r
	}

	// should not be reached
	panic(fmt.Sprintf("basiccolor: internal error: %+v", hsl))
}

func determineColor(hue, saturation, lightness float64) (Color, bool) {
	if hue >= 343 || hue < 17 {
		if lightness <= 0.34 {
			return Brown, true
		}
		if lightness <= 0.66 {
			return Red, true
		}
		return Orange, true
	}

	if hue >= 17 && hue < 46 {
		if lightness <= 0.40 {
			return Brown, true
		}
		return Orange, true
	}

	if hue >= 46 && hue < 69 {
		return Yellow, true
	}

	if hue >= 69 && hue < 172 {
		return Green, true
	}

	if hue >= 172 && hue < 185 {
		if lightness < 0.5 {
			return Green, true
		}
		return Blue, true
	}

	if hue >= 185 && hue < 244 {
		return Blue, true
	}

	if hue >= 244 && hue < 286 {
		return Violet, true
	}

	if hue >= 286 && hue < 307 {
		if lightness < 0.5 {
			return Violet, true
		}
		return Pink, true
	}

	if hue >= 307 && hue < 343 {
		return Pink, true
	}

	return Color(-1), false
}

func isGray(saturation, lightness float64) bool {
	// Based on these (saturation, lightness) points:
	//
	// (0.00, 0.20) // taken care of by isBlack()
	// (0.12, 0.16)
	// (0.12, 0.946)
	// (0.00, 0.93) // taken care of by isWhite()

	return saturation < 0.12
}

func isWhite(saturation, lightness float64) bool {
	// Based on these (saturation, lightness) points:
	//
	// (0.00, 0.93)
	// (0.10, 0.93)
	// (0.15, 0.97)
	// (0.30, 0.99)
	// (1.00, 0.99)

	if saturation < 0.10 {
		return lightness >= 0.93
	}
	if saturation < 0.15 {
		return lightness >= (4./5*saturation)+0.85
	}
	if saturation < 0.30 {
		return lightness >= (2./15*saturation)+0.95
	}
	return lightness >= 0.99
}

func isBlack(saturation, lightness float64) bool {
	// Based on these (saturation, lightness) points:
	//
	// (0.00, 0.20)
	// (0.08, 0.20)
	// (0.20, 0.08)
	// (0.50, 0.05)
	// (1.00, 0.04)

	if saturation < 0.08 {
		return lightness <= 0.20
	}
	if saturation < 0.20 {
		return lightness <= -saturation+0.28
	}
	if saturation < 0.50 {
		return lightness <= (-0.1*saturation)+0.1
	}
	return lightness <= (-0.02*saturation)+0.06
}