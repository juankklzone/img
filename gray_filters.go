package img

import (
	"image"
	"math/rand"
)

type FilterType int

const (
	GrayFilterSaltAndPepper FilterType = 1
)

type FilterOptions struct {
	Percentage float64
	Dimensions image.Point
}

func (gi GrayImage) Filter(ft FilterType, fo *FilterOptions) (i GrayImage) {
	i = gi.Clone()
	switch ft {
	case GrayFilterSaltAndPepper:
		i.saltAndPepper(fo.Percentage)
	}
	return
}

func (gi GrayImage) saltAndPepper(affected float64) {
	max := gi.Bounds().Max
	maxx, maxy := max.X, max.Y
	noAffected := int(float64(max.X*max.Y) * affected)
	for i := 0; i < noAffected; i++ {
		x, y := rand.Intn(maxx), rand.Intn(maxy)
		saltOrPepper := rand.Intn(2) * 255
		gi[x][y].Y = uint8(saltOrPepper)
	}
}
