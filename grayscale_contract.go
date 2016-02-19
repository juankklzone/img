package img

import "image/color"

//Contract scales the gray scale image to a specified range
func (i GrayImage) Contract(scale float64, offset float64) (gi GrayImage) {
	gi = NewEmptyGrayImage(i.Bounds())
	for x := range i {
		for y := range i[x] {
			r, g, b, _ := i.At(x, y).RGBA()
			r8 := uint8(round(float64(uint8(r))*scale + offset))
			g8 := uint8(round(float64(uint8(g))*scale + offset))
			b8 := uint8(round(float64(uint8(b))*scale + offset))
			gi.Set(x, y, color.Gray{
				Y: uint8(0.299*float32(r8) + 0.587*float32(g8) + 0.114*float32(b8)),
			})
		}
	}
	return
}
