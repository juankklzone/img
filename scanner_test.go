package img

import (
	"image"
	"testing"
)

func TestScannerInRow(t *testing.T) {
	x, y := 3, 7
	gi := NewEmptyGrayImage(image.Rect(0, 0, x, y))
	bi := NewBinaryImage(gi, 124)
	bi[0] = []bool{true, true, false, true, true, false, true}
	s := NewScanner(bi)
	os := s.scanRow(0)
	if len(os) != 3 {
		t.Log("Objects in row 0: \n", os)
		t.FailNow()
	}

}

func TestScannerInMultiplelines(t *testing.T) {
	data := [][]bool{
		{true, false, true, false, true},
		{false, true, false, true, false},
		{true, false, false, false, true},
		{true, true, false, false, false},
		{false, false, false, false, true},
		{true, true, true, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
	}
	bi := BinaryImage(data)
	s := NewScanner(bi)
	s.SearchObjects()
	t.Log(s.set)
}
