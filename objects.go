package img

import (
	"image"
)

//object identifies a set of positions into the image that together create an image. map[row]points
type object struct {
	points map[int][]image.Point
	id     int
}

func newObject() *object {
	return &object{points: make(map[int][]image.Point)}
}

func (o *object) append(row, y int) {
	o.points[row] = append(o.points[row], image.Point{X: row, Y: y})
}

func (o *object) maxInRow(row int) (y int) {
	idx := len(o.points[row])
	if idx > 0 {
		y = o.points[row][idx-1].Y
	} else {
		y = -1
	}
	return
}

func (o *object) minInRow(row int) (y int) {
	if o.points[row] != nil {
		y = o.points[row][0].Y
	} else {
		y = -1
	}
	return
}

func (o *object) hasPointsInRow(row, minY, maxY int) bool {
	if points := o.points[row]; points != nil {
		minReachable, maxReachable := minY-2, maxY+2
		for _, p := range points {
			if p.Y > minReachable && p.Y < maxReachable {
				return true
			}
		}
	}
	return false
}
