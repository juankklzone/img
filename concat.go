package img

import "image/color"

//ConcatBottom expects a same-width image to concat at the bottom of i
//the result of the concatenation is returned
func (i ColorImage) ConcatBottom(ci ColorImage) ColorImage {
	maxX := ci.Bounds().Max.X
	maxYa := i.Bounds().Max.Y
	maxYb := ci.Bounds().Max.Y
	newWidth := maxYa + maxYb
	if r := ci.Bounds().Max.X; r < i.Bounds().Max.X {
		maxX = r
	} else {
		maxX = i.Bounds().Max.X
	}
	newImage := make([][]color.Color, maxX)
	for x := 0; x < maxX; x++ {
		newImage[x] = make([]color.Color, newWidth)
		copy(newImage[x][:maxYa], i[x][:])
		copy(newImage[x][maxYa:], ci[x][:])
	}
	return newImage
}

//ConcatRight expects a same-heigth image to concat at the left of i
//the result of the concatenation is returned
func (i ColorImage) ConcatRight(ci ColorImage) ColorImage {
	maxXa := i.Bounds().Max.X
	maxXb := ci.Bounds().Max.X
	newWidth := maxXa + maxXb
	newImage := make([][]color.Color, newWidth)
	ySize := i.Bounds().Max.Y
	for x := 0; x < maxXa; x++ {
		newImage[x] = make([]color.Color, ySize)
		copy(newImage[x][:], i[x][:])
	}
	secondIdx := 0
	for x := maxXa; x < newWidth; x++ {
		newImage[x] = make([]color.Color, ySize)
		copy(newImage[x][:], ci[secondIdx][:])
		secondIdx++
	}
	return newImage
}
