package img

import (
	"bytes"
	"fmt"
	"image/color"
	"image/draw"
	"sync"
)

//object identifies a set of positions into the image that together create an image. map[row]points
type object struct {
	points map[int][]int
	mtx    sync.RWMutex
	id     int
}

func newObject() *object {
	return &object{points: make(map[int][]int)}
}

func (o *object) append(row, y int) {
	o.mtx.Lock()
	o.points[row] = append(o.points[row], y)
	o.mtx.Unlock()
}

func (o *object) maxInRow(row int) (y int) {
	o.mtx.RLock()
	idx := len(o.points[row])
	if idx > 0 {
		y = o.points[row][idx-1]
	} else {
		y = -1
	}
	o.mtx.RUnlock()
	return
}

func (o *object) minInRow(row int) (y int) {
	o.mtx.RLock()
	if o.points[row] != nil {
		y = o.points[row][0]
	} else {
		y = -1
	}
	o.mtx.RUnlock()
	return
}

func (o *object) hasPointsInRow(row, minY, maxY int) bool {
	o.mtx.RLock()
	//defer, i thought you were my friend
	if points := o.points[row]; points != nil {
		minReachable, maxReachable := minY-2, maxY+2
		for _, p := range points {
			if p > minReachable && p < maxReachable {
				o.mtx.RUnlock()
				return true
			}
		}
	}
	o.mtx.RUnlock()
	return false
}

func (o *object) String() string {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	buffer := new(bytes.Buffer)
	for row, points := range o.points {
		for _, p := range points {
			fmt.Fprintf(buffer, "[%d,%d], ", row, p)
		}
	}
	return buffer.String()
}

func (o *object) len() int {
	o.mtx.RLock()
	count := 0
	for _, rowpoints := range o.points {
		count += len(rowpoints)
	}
	o.mtx.RUnlock()
	return count
}

func (o *object) draw(img draw.Image, color color.Color) {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	for row, points := range o.points {
		for _, y := range points {
			img.Set(row, y, color)
		}
	}
}

func (o *object) get(i int) []int {
	o.mtx.RLock()
	data := o.points[i]
	o.mtx.RUnlock()
	return data
}

func (o *object) set(i int, d []int) {
	o.mtx.Lock()
	o.points[i] = d
	o.mtx.Unlock()
}

func (o *object) isAdjacent(o2 *object, row int) bool {
	o.mtx.RLock()
	o2.mtx.RLock()
	defer o.mtx.RUnlock()
	defer o2.mtx.RUnlock()
	points := o.points[row]
	otherPoints := o2.points[row-1]
	for _, p := range points {
		for _, op := range otherPoints {
			k := p - op
			if k < 2 && k > -2 {
				return true
			}
		}
	}
	return false
}
