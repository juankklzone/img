package img

import "image"

func (gi GrayImage) getVecinity(x, y int) []uint8 {
	var vecinity []uint8
	positions := []image.Point{
		image.Point{X: x - 1, Y: y - 1},
		image.Point{X: x, Y: y - 1},
		image.Point{X: x + 1, Y: y - 1},
		image.Point{X: x - 1, Y: y},
		image.Point{X: x + 1, Y: y},
		image.Point{X: x - 1, Y: y + 1},
		image.Point{X: x, Y: y + 1},
		image.Point{X: x + 1, Y: y + 1},
	}
	sz := gi.Bounds().Max
	for _, pos := range positions {
		if pos.X < sz.X && pos.Y < sz.Y && pos.X >= 0 && pos.Y >= 0 {
			vecinity = append(vecinity, gi.AtGray(pos.X, pos.Y))
		}
	}
	return vecinity
}
