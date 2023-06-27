package engine

import "image/color"

func Ptr[T any](t T) *T {
	return &t
}

func ConvertToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()

	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
