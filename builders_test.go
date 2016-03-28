package img

import (
	"image"
	"image/color"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

var (
	randGen *rand.Rand
)

func init() {
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestCloneGrayImg(t *testing.T) {
	r := image.Rect(0, 0, 50, 50)
	grayImg := NewEmptyGrayImage(r)
	for i := range grayImg {
		for j := range grayImg[i] {
			grayImg[i][j] = color.Gray{uint8(randGen.Intn(256))}
		}
	}
	clone := grayImg.Clone()
	if !reflect.DeepEqual(grayImg, clone) {
		t.FailNow()
	}
}

func TestNewGrayFromImg(t *testing.T) {
	sz := image.Rect(0, 0, randGen.Intn(100), randGen.Intn(100))
	colorImg := NewEmptyImage(sz)
	grayImg := NewGrayImage(colorImg)
	emptyGrayImg := NewEmptyGrayImage(sz)
	if !reflect.DeepEqual(emptyGrayImg, grayImg) {
		t.FailNow()
	}
}
