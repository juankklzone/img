package canny

import (
	"image"
	"math"

	"github.com/juankklzone/img"
)

const (
	gaussianCutOff = 0.005
	magnitudeScale = 100.0
	magnitudeLimit = 1000.0
	magnitudeMax   = magnitudeScale * magnitudeLimit
	maxData        = 0xff000000
)

type Detector struct {
	image   *image.RGBA
	picsize int
	width   int
	heigth  int

	gaussKernelRad   float64
	gaussKernelWidth int

	lowThreshold  float64
	highThreshold float64

	data      []int
	magnitude []int
	xConv     []float64
	yConv     []float64
	xGradient []float64
	yGradient []float64
}

func NewDetector(g *image.RGBA) *Detector {
	d := &Detector{
		image:            g,
		gaussKernelRad:   2,
		gaussKernelWidth: 16,
		lowThreshold:     2.5,
		highThreshold:    7.5,
	}
	return d
}

func (d *Detector) Process() img.GrayImage {
	border := d.image.Bounds().Max
	d.width = border.X
	d.heigth = border.Y
	d.picsize = d.heigth * d.width
	d.initArrays()
	d.readLuminance()
	d.computeGradients()
	low, high := 0, 0
	low = int(math.Floor((d.lowThreshold * magnitudeScale) + .5))
	high = int(math.Floor((d.highThreshold * magnitudeScale) + .5))
	d.performHysterisis(low, high)
	d.thresholdEdges()
	borders := d.writeEdges()
	return borders
}

func castLuminance(r, g, b uint32) int {
	return int(math.Floor((.299*float64(r) + .587*float64(g) + .114*float64(b)) + .5))
}

func (d *Detector) readLuminance() {
	for y := 0; y < d.heigth; y++ {
		for x := 0; x < d.width; x++ {
			idx := y*d.width + x
			r, g, b, _ := d.image.At(x, y).RGBA()
			l := castLuminance(b&0xff, g&0xff, r&0xff)
			d.data[idx] = int(l)
		}
	}
}

func (d *Detector) initArrays() {
	if d.data == nil || d.picsize != len(d.data) {
		d.data = make([]int, d.picsize)
		d.magnitude = make([]int, d.picsize)
		d.xConv = make([]float64, d.picsize)
		d.yConv = make([]float64, d.picsize)
		d.xGradient = make([]float64, d.picsize)
		d.yGradient = make([]float64, d.picsize)
	}
}

func (d *Detector) computeGradients() {
	kernel := make([]float64, d.gaussKernelWidth)
	diffKernel := make([]float64, d.gaussKernelWidth)
	var kwidth int
	sumk := 0.0
	sumg := 0.0
	for kwidth = 0; kwidth < d.gaussKernelWidth; kwidth++ {
		g1 := gaussian(float64(kwidth), d.gaussKernelRad)
		if g1 <= gaussianCutOff && kwidth >= 2 {
			break
		}
		g2 := gaussian(float64(kwidth)-0.5, d.gaussKernelRad)
		g3 := gaussian(float64(kwidth)+0.5, d.gaussKernelRad)
		kernel[kwidth] = (g1 + g2 + g3) / 3 / (2 * math.Pi * d.gaussKernelRad * d.gaussKernelRad)
		sumk += kernel[kwidth]
		diffKernel[kwidth] = g3 - g2
		sumg += diffKernel[kwidth]
	}
	initX := kwidth - 1
	maxX := d.width - (kwidth - 1)
	initY := d.width * (kwidth - 1)
	maxY := d.width * (d.heigth - (kwidth - 1))
	for x := initX; x < maxX; x++ {
		for y := initY; y < maxY; y += d.width {
			index := x + y
			sumX := float64(d.data[index]) * kernel[0]
			sumY := sumX
			xOffset := 1
			yOffset := d.width
			for xOffset < kwidth {
				sumY += kernel[xOffset] * float64(d.data[index-yOffset]+d.data[index+yOffset])
				sumX += kernel[xOffset] * float64(d.data[index-xOffset]+d.data[index+xOffset])
				xOffset++
				yOffset += d.width
			}
			d.yConv[index] = sumY
			d.xConv[index] = sumX
		}
	}
	sumGradX := 0.0
	sumGradY := 0.0
	for y := initY; y < maxY; y += d.width {
		for x := initX; x < maxX; x++ {
			sum := 0.0
			index := x + y
			for i := 1; i < kwidth; i++ {
				sum += diffKernel[i] * (d.yConv[index-i] - d.yConv[index+i])
			}
			d.xGradient[index] = sum
			sumGradX += sum
		}
	}

	for y := initY; y < maxY; y += d.width {
		for x := kwidth; x < d.width-kwidth; x++ {
			sum := 0.0
			index := x + y
			yOffset := d.width
			for i := 1; i < kwidth; i++ {
				sum += diffKernel[i] * (d.xConv[index-yOffset] - d.xConv[index+yOffset])
				yOffset += d.width
			}
			d.yGradient[index] = sum
			sumGradY += sum
		}
	}
	initX = kwidth
	maxX = d.width - kwidth
	initY = d.width * kwidth
	maxY = d.width * (d.heigth - kwidth)
	for x := initX; x < maxX; x++ {
		for y := initY; y < maxY; y += d.width {
			index := x + y
			indexN, indexS := index-d.width, index+d.width
			indexW, indexE := index-1, index+1
			indexNW, indexNE := indexN-1, indexN+1
			indexSW, indexSE := indexS-1, indexS+1
			xGrad := d.xGradient[index]
			yGrad := d.yGradient[index]
			gradMag := hypot(xGrad, yGrad)

			nMag := hypot(d.xGradient[indexN], d.yGradient[indexN])
			sMag := hypot(d.xGradient[indexS], d.yGradient[indexS])
			eMag := hypot(d.xGradient[indexE], d.yGradient[indexE])
			wMag := hypot(d.xGradient[indexW], d.yGradient[indexW])
			neMag := hypot(d.xGradient[indexNE], d.yGradient[indexNE])
			seMag := hypot(d.xGradient[indexSE], d.yGradient[indexSE])
			swMag := hypot(d.xGradient[indexSW], d.yGradient[indexSW])
			nwMag := hypot(d.xGradient[indexNW], d.yGradient[indexNW])
			tmp := 0.0
			passed := false
			if xGrad*yGrad <= 0.0 {
				if math.Abs(xGrad) >= math.Abs(yGrad) {
					tmp = math.Abs(xGrad * gradMag)
					passed =
						tmp >= math.Abs(yGrad*neMag-(xGrad+yGrad)*eMag) &&
							tmp > math.Abs(yGrad*swMag-(xGrad+yGrad)*wMag)
				} else {
					tmp = math.Abs(yGrad * gradMag)
					passed =
						tmp >= math.Abs(xGrad*neMag-(yGrad+xGrad)*nMag) &&
							tmp > math.Abs(xGrad*swMag-(yGrad+xGrad)*sMag)
				}
			} else {
				if math.Abs(xGrad) >= math.Abs(yGrad) {
					tmp = math.Abs(xGrad * gradMag)
					passed = tmp >= math.Abs(yGrad*seMag+(xGrad-yGrad)*eMag) &&
						tmp > math.Abs(yGrad*nwMag+(xGrad-yGrad)*wMag)
				} else {
					tmp = math.Abs(yGrad * gradMag)
					passed = tmp >= math.Abs(xGrad*seMag+(yGrad-xGrad)*sMag) &&
						tmp > math.Abs(xGrad*nwMag+(yGrad-xGrad)*nMag)
				}
			}
			if passed {
				if gradMag > magnitudeLimit {
					d.magnitude[index] = magnitudeMax
				} else {
					d.magnitude[index] = int(magnitudeScale * gradMag)
				}
			} else {
				d.magnitude[index] = 0
			}
		}
	}
}

func (d *Detector) performHysterisis(low, high int) {
	d.data = make([]int, len(d.data))
	offset := 0
	for y := 0; y < d.heigth; y++ {
		for x := 0; x < d.width; x++ {
			if d.data[offset] == 0 && d.magnitude[offset] >= high {
				d.follow(x, y, offset, low)
			}
			offset++
		}
	}
}

func (d *Detector) follow(x1, y1, i1, threshold int) {
	x0, x2, y0, y2 := 0, x1, 0, y1
	if x1 != 0 {
		x0 = x1 - 1
	}
	if y1 != 0 {
		y0 = y1 - 1
	}
	if x1 != d.width-1 {
		x2++
	}
	if y1 != d.heigth-1 {
		y2++
	}
	d.data[i1] = d.magnitude[i1]
	for x := x0; x <= x2; x++ {
		for y := y0; y <= y2; y++ {
			i2 := x + y*d.width
			if ((y != y1) || (x != x1)) && d.data[i2] == 0 && d.magnitude[i2] >= threshold {
				d.follow(x, y, i2, threshold)
				return
			}
		}
	}
	return
}

func (d *Detector) thresholdEdges() {
	for i := 0; i < d.picsize; i++ {
		if d.data[i] > 0 {
			d.data[i] = -1
		} else {
			d.data[i] = maxData
		}
	}
}

func (d *Detector) writeEdges() img.GrayImage {
	edges := img.NewEmptyGrayImage(d.image.Bounds())
	index := 0
	for y := 0; y < d.heigth; y++ {
		for x := 0; x < d.width; x++ {
			index = y*d.width + x
			edges.SetGray(x, y, uint8(d.data[index]))
		}
	}
	return edges
}

func gaussian(x, sigma float64) float64 {
	return math.Exp(-(x * x) / (2 * sigma * sigma))
}

func hypot(x, y float64) float64 {
	return math.Hypot(x, y)
}

func (d *Detector) at(idx int) (x, y int) {
	y = idx / d.width
	x = idx % d.width
	return
}
