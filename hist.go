package img

import "image"

//Histogram represents the ocurrence of a level of color into an image.
type Histogram map[uint8]int

//GetHistogram creates an histogram per channel, which counts the ocurrences of a color
func GetHistogram(img image.Image) (h []Histogram) {
	h = make([]Histogram, 4)
	h[0], h[1], h[2], h[3] = make(Histogram), make(Histogram), make(Histogram), make(Histogram)
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			h[0][uint8(r)]++
			h[1][uint8(g)]++
			h[2][uint8(b)]++
			h[3][uint8(a)]++
		}
	}
	return
}
