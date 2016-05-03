package img

import "image/color"

//WithRows slices the rows of an image and return the segment into another image
func (i ColorImage) WithRows(start, end int) ColorImage {
	diffs := end - start
	newImage := make([][]color.Color, i.Bounds().Max.X)
	for x := range newImage {
		newImage[x] = make([]color.Color, diffs)
		copy(newImage[x][:], i[x][start:end])
	}
	return newImage
}

//WithCols slices the columns of an image and return the segment into another image
func (i ColorImage) WithCols(start, end int) ColorImage {
	newDiff := end - start
	newImage := make([][]color.Color, newDiff)
	y := i.Bounds().Max.Y
	counter := start
	for x := range newImage {
		newImage[x] = make([]color.Color, y)
		copy(newImage[x][:], i[counter][:])
		counter++
	}
	return newImage
}
