package img

import (
	"image"
	"testing"
)

func TestScannerInRow(t *testing.T) {
	s := NewScanner()
	x, y := 3, 7
	gi := NewEmptyGrayImage(image.Rect(0, 0, x, y))
	bi := NewBinaryImage(gi, 124)
	bi[0] = []bool{true, true, false, true, true, false, true}
	os := s.scanRow(bi, 0)
	if len(os) != 3 {
		t.Log("Objects in row 0: \n", os)
		t.FailNow()
	}

}

func TestScannerInMultiplelines(t *testing.T) {
	data := [][]bool{{true, false, false, false, true},
		{false, true, false, true}}
	bi := BinaryImage(data)
	s := NewScanner()
	s.SearchObjects(bi)
	for id, obj := range s.set {
		t.Log(id)
		for _, row := range obj.points {
			for _, point := range row {
				t.Log(point)
			}
		}
	}
}
