package img

import (
	"image/color"
	"math/rand"
	"time"
)

var (
	randGen *rand.Rand
)

func init() {
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randColors(sz int) (colors []color.RGBA) {
	colors = make([]color.RGBA, sz)
	for i := range colors {
		colors[i].A = 255
		colors[i].R = uint8(randGen.Intn(256))
		colors[i].G = uint8(randGen.Intn(256))
		colors[i].B = uint8(randGen.Intn(256))
	}
	return
}
