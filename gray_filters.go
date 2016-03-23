package img

import (
	"image"
	"math"
	"math/rand"
)

//FilterType defines a kind of filter applied to an image
type FilterType int

const (
	//SaltAndPepper is the Salt & Pepper filter
	SaltAndPepper FilterType = 1
	//Average filter applies an average filter for each pixel on the image
	Average FilterType = 2
	//ArithmeticMean filter
	ArithmeticMean FilterType = 3
	//HarmonicMean filter
	HarmonicMean FilterType = 4
	//ContraharmonicMean filter
	ContraharmonicMean FilterType = 5
	//Fashion filter
	Fashion FilterType = 6
)

//FilterOptions defines options for a filter to apply
type FilterOptions struct {
	Percentage float64
	Dimensions image.Point
	Exponent   int
}

//Filter applies a filter to an image, with the options received
func (gi GrayImage) Filter(ft FilterType, fo *FilterOptions) (i GrayImage) {
	switch ft {
	case SaltAndPepper:
		i = gi.Clone()
		i.saltAndPepper(fo.Percentage)
	case Average:
		i = gi.avg()
	case ArithmeticMean:
		i = gi.arithMean()
	case HarmonicMean:
		i = gi.harmonicMean()
	case ContraharmonicMean:
		i = gi.contraharmonicMean(float64(fo.Exponent))
	case Fashion:
		i = gi.fashion()
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

func (gi GrayImage) forEach(mutator func(GrayImage, GrayImage, int, int), copy GrayImage) {
	maxPoint := gi.Bounds().Max
	for i := 0; i < maxPoint.X; i++ {
		for j := 0; j < maxPoint.Y; j++ {
			mutator(gi, copy, i, j)
		}
	}
}

func (gi GrayImage) avg() GrayImage {
	copy := NewEmptyGrayImage(gi.Bounds())
	mutator := func(o, g GrayImage, x, y int) {
		vecinity := o.getVecinity(x, y)
		avg := 0
		for _, v := range vecinity {
			avg += int(v)
		}
		avg /= len(vecinity)
		g.SetGray(x, y, uint8(avg))
	}
	gi.forEach(mutator, copy)
	return copy
}

func (gi GrayImage) arithMean() GrayImage {
	bounds := gi.Bounds()
	copy := NewEmptyGrayImage(bounds)
	sz := bounds.Max
	k := 1 / (sz.X * sz.Y)
	mutator := func(o, g GrayImage, x, y int) {
		vecinity := o.getVecinity(x, y)
		sum := 0
		for _, v := range vecinity {
			sum += int(v)
		}
		g.SetGray(x, y, uint8(sum*k))
	}
	gi.forEach(mutator, copy)
	return copy
}

func (gi GrayImage) harmonicMean() GrayImage {
	bounds := gi.Bounds()
	copy := NewEmptyGrayImage(bounds)
	sz := bounds.Max
	k := 1 / (sz.X * sz.Y)
	mutator := func(o, g GrayImage, x, y int) {
		vecinity := o.getVecinity(x, y)
		sum := 0
		for _, v := range vecinity {
			sum += int(v)
		}
		g.SetGray(x, y, uint8(k/sum))
	}
	gi.forEach(mutator, copy)
	return copy
}

func (gi GrayImage) contraharmonicMean(r float64) GrayImage {
	copy := NewEmptyGrayImage(gi.Bounds())
	r1 := r + 1.0
	mutator := func(o, g GrayImage, x, y int) {
		vecinity := g.getVecinity(x, y)
		up, down := 0.0, 0.0
		for _, v := range vecinity {
			up += math.Pow(float64(v), r1)
			down += math.Pow(float64(v), r)
		}
		g.SetGray(x, y, uint8(up/down))
	}
	gi.forEach(mutator, copy)
	return copy
}

func (gi GrayImage) fashion() GrayImage {
	copy := NewEmptyGrayImage(gi.Bounds())
	mutator := func(o, g GrayImage, x, y int) {
		vecinity := o.getVecinity(x, y)
		repeated := make(map[uint8]int)
		for _, v := range vecinity {
			repeated[v]++
		}
		max := 0
		var maxRep uint8
		for val, rep := range repeated {
			if rep > max {
				maxRep, max = val, rep
			}
		}
		g.SetGray(x, y, maxRep)
	}
	gi.forEach(mutator, copy)
	return copy
}
