package img

import "sync"

//NewBinaryImage returns a BinaryImage with img's size
//if the color is greater than div, it's set up to true
func NewBinaryImage(img GrayImage, div uint8) BinaryImage {
	borders := img.Bounds()
	maxX, maxY := borders.Max.X, borders.Max.Y
	bi := make([][]bool, maxX)
	wg := new(sync.WaitGroup)
	wg.Add(maxX)
	for x := range bi {
		bi[x] = make([]bool, maxY)
		go readBinaryRow(wg, img, bi, x, maxY, div)
	}
	wg.Wait()
	return bi
}

func readBinaryRow(wg *sync.WaitGroup, from GrayImage, to BinaryImage, x, maxY int, div uint8) {
	for y := 0; y < maxY; y++ {
		to[x][y] = from.AtGray(x, y) > div
	}
	wg.Done()
}
