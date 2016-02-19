package img

import (
	"image/color"
	"math"
)

//Contract scales the color of the actual image to a specified range
func (i ColorImage) Contract(scale float64, offset float64) (ci ColorImage) {
	ci = NewEmptyImage(i.Bounds())
	for x := range i {
		for y := range i[x] {
			r, g, b, a := i.At(x, y).RGBA()
			r8 := uint8(round(float64(uint8(r))*scale + offset))
			g8 := uint8(round(float64(uint8(g))*scale + offset))
			b8 := uint8(round(float64(uint8(b))*scale + offset))
			a8 := uint8(round(float64(uint8(a))*scale + offset))
			ci.Set(x, y, color.RGBA{r8, g8, b8, a8})
		}
	}
	return
}

func round(x float64) float64 {
	return math.Floor(x + .5)
}
