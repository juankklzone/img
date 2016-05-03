package img

//Contract scales the gray scale image to a specified range
func (i GrayImage) Contract(scale float64, offset float64) (gi GrayImage) {
	gi = i.Clone()
	for x := range i {
		for y := range i[x] {
			gi[x][y].Y = uint8(float64(i[x][y].Y)*scale + offset)
		}
	}
	return
}
