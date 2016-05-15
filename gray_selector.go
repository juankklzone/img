package img

//Select takes all the values inside gi between min and max inclusive
//and returns the result
func (gi GrayImage) Select(min, max uint8) GrayImage {
	bounds := gi.Bounds()
	newGray := NewEmptyGrayImage(bounds)
	for x := range gi {
		for y := range gi[x] {
			v := gi[x][y]
			if v.Y >= min && v.Y <= max {
				newGray[x][y] = v
			}
		}
	}
	return newGray
}
