package img

import (
	"bytes"
	"fmt"
)

//object identifies a set of positions into the image that together create an image. map[row]points
type object struct {
	points map[int][]int
	id     int
}

func newObject() *object {
	return &object{points: make(map[int][]int)}
}

func (o *object) append(row, y int) {
	o.points[row] = append(o.points[row], y)
}

func (o *object) maxInRow(row int) (y int) {
	idx := len(o.points[row])
	if idx > 0 {
		y = o.points[row][idx-1]
	} else {
		y = -1
	}
	return
}

func (o *object) minInRow(row int) (y int) {
	if o.points[row] != nil {
		y = o.points[row][0]
	} else {
		y = -1
	}
	return
}

func (o *object) hasPointsInRow(row, minY, maxY int) bool {
	if points := o.points[row]; points != nil {
		minReachable, maxReachable := minY-2, maxY+2
		for _, p := range points {
			if p > minReachable && p < maxReachable {
				return true
			}
		}
	}
	return false
}

func (o *object) String() string {
	buffer := new(bytes.Buffer)
	for row, points := range o.points {
		for _, p := range points {
			fmt.Fprintf(buffer, "[%d,%d], ", row, p)
		}
	}
	return buffer.String()
}

func (o *object) len() int {
	count := 0
	for _, rowpoints := range o.points {
		count += len(rowpoints)
	}
	return count
}
