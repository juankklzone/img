package img

import (
	"image"
	"testing"
)

func TestVecinity(t *testing.T) {
	i := NewEmptyGrayImage(image.Rect(0, 0, 4, 4))
	/*for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			k := uint8(x*y + y)
			i.SetGray(x, y, k)
		}
	}*/
	v := i.getVecinity(2, 2)
	if l := len(v); l < 8 {
		t.Fail()
	}
}
