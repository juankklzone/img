package img

import (
	"image/color"
	"image/draw"
)

//Red remove blue and green channels from the image
//and returns a new RGBA image with the result
func Red(img draw.Image) ColorImage {
	borders := img.Bounds()
	result := NewEmptyImage(borders)
	model := result.ColorModel()
	maxX, maxY := borders.Max.X, borders.Max.Y
	for xaxis := 0; xaxis < maxX; xaxis++ {
		for yaxis := 0; yaxis < maxY; yaxis++ {
			r, _, _, a := img.At(xaxis, yaxis).RGBA()
			result.Set(xaxis, yaxis, model.Convert(color.RGBA{uint8(r), 0, 0, uint8(a)}))
		}
	}
	return result
}

//Green remove red and blue channels from the image
//and returns a new RGBA image with the result
func Green(img draw.Image) ColorImage {
	borders := img.Bounds()
	result := NewEmptyImage(borders)
	model := result.ColorModel()
	maxX, maxY := borders.Max.X, borders.Max.Y
	for xaxis := 0; xaxis < maxX; xaxis++ {
		for yaxis := 0; yaxis < maxY; yaxis++ {
			_, g, _, _ := img.At(xaxis, yaxis).RGBA()
			result.Set(xaxis, yaxis, model.Convert(color.RGBA{0, uint8(g), 0, 0}))
		}
	}
	return result
}

//Blue remove red and green channels from the image
//and returns a new RGBA image with the result
func Blue(img draw.Image) ColorImage {
	borders := img.Bounds()
	result := NewEmptyImage(borders)
	model := result.ColorModel()
	maxX, maxY := borders.Max.X, borders.Max.Y
	for xaxis := 0; xaxis < maxX; xaxis++ {
		for yaxis := 0; yaxis < maxY; yaxis++ {
			_, _, b, _ := img.At(xaxis, yaxis).RGBA()
			result.Set(xaxis, yaxis, model.Convert(color.RGBA{0, 0, uint8(b), 0}))
		}
	}
	return result
}

//Add sums an addition into an base image and returns
//the result in other image.
func Add(base, addition draw.Image) (dest ColorImage) {
	borders := base.Bounds()
	dest = NewEmptyImage(borders)
	maxX, maxY := borders.Max.X, borders.Max.Y
	for xaxis := 0; xaxis < maxX; xaxis++ {
		for yaxis := 0; yaxis < maxY; yaxis++ {
			br, bg, bb, ba := base.At(xaxis, yaxis).RGBA()
			ar, ag, ab, aa := addition.At(xaxis, yaxis).RGBA()
			dest.Set(xaxis, yaxis, color.RGBA{
				uint8((br + ar)),
				uint8((bg + ag)),
				uint8((bb + ab)),
				uint8((ba + aa)),
			})
		}
	}
	return
}
