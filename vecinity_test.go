package img

import (
	"fmt"
	"image"
	"testing"
)

func TestVecinity(t *testing.T) {
	i := NewEmptyGrayImage(image.Rect(0, 0, 4, 4))
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			k := uint8(x*y + y)
			i.SetGray(x, y, k)
			fmt.Printf("%d\t", k)
		}
		fmt.Println()
	}
	v := i.getVecinity(2, 2)
	if len(v) < 8 {
		t.Log("Vecinos de 0,0:")
		for _, vecino := range v {
			t.Log(vecino)
		}
	} else {
		t.Fail()
	}
}
