package img

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"
)

type RGBAScanner struct {
	set *objectSet
	img image.Image
	ev  Evaluation
}

type Evaluation func(c color.Color) bool

func NewRGBAScanner(i image.Image, ev Evaluation) *RGBAScanner {
	return &RGBAScanner{set: newObjectset(), img: i, ev: ev}
}

//SearchObjects count the objects found in the image
//and save them for later use
func (s *RGBAScanner) SearchObjects() int {
	rows := s.img.Bounds().Max.X
	wg := new(sync.WaitGroup)
	wg.Add(2)
	half := rows / 2
	ini := time.Now()
	maxy := s.img.Bounds().Max.Y
	go s.scanRows(wg, 0, half, maxy)
	go s.scanRows(wg, half, rows, maxy)
	wg.Wait()
	now := time.Now()
	s.joinSegments(half)
	fmt.Println("time req for joinSegments ", time.Since(now))
	fmt.Println("total time ", time.Since(ini))
	return len(s.set.objs)
}

func (s *RGBAScanner) scanRows(wg *sync.WaitGroup, initial, final, maxy int) {
	defer wg.Done()
	for i := initial; i < final; i++ {
		//search for objects in last row
		lastRow := s.set.objectsInLastRow(i)
		//search for objects in this row
		currentRowObjects := s.scanRow(i, maxy)
		//if there's no rows in the last row, they are new objects
		if len(lastRow.objs) == 0 {
			s.set.append(currentRowObjects)
			continue
		}
		//checks for every object in this row if it has some parent in the last row
		for _, currentObject := range currentRowObjects.objs {
			var belongsTo []int
			min, max := currentObject.minInRow(i), currentObject.maxInRow(i)
			for _, upperObj := range lastRow.objs {
				if upperObj.hasPointsInRow(i-1, min, max) {
					belongsTo = append(belongsTo, upperObj.id)
				}
			}
			lenParents := len(belongsTo)
			//if an object has no parents, it is a new object
			if lenParents == 0 {
				s.set.add(currentObject)
			} else {
				//an object can be a children of one or more, then we have to join them all
				initial := belongsTo[0]
				for i := 1; i < lenParents; i++ {
					lastRow.groupObjects(initial, lastRow.get(belongsTo[i]))
				}
				lastRow.groupObjects(initial, currentObject)
				s.set.drop(belongsTo[1:]...)
				lastRow.drop(belongsTo[1:]...)
			}
		}
	}
}

func (s *RGBAScanner) joinSegments(row int) {
	lastRow := s.set.objectsInLastRow(row + 1)
	if len(lastRow.objs) == 0 {
		return
	}
	for _, parent := range lastRow.objs {
		var children []int
		thisRow := s.set.objectsInLastRow(row)
		for _, child := range thisRow.objs {
			if parent.isAdjacent(child, row) {
				children = append(children, child.id)
			}
		}
		if k := len(children); k != 0 {
			for i := 0; i < k; i++ {
				s.set.groupObjects(parent.id, s.set.get(children[i]))
			}
			s.set.drop(children...)
		}
	}

}

func (s *RGBAScanner) scanRow(row, maxy int) *objectSet {
	os := newObjectset()
	currentObj := newObject()
	objInProgress := false
	id := 0
	for y := 0; y < maxy; y++ {
		if s.ev(s.img.At(row, y)) {
			currentObj.append(row, y)
			objInProgress = true
		} else if objInProgress {
			currentObj.id = id
			os.set(id, currentObj)
			id++
			objInProgress = false
			currentObj = newObject()
		}
	}
	if objInProgress {
		currentObj.id = id
		os.set(id, currentObj)
	}
	return os
}

//Filter discarts from the object set, those which size is less than minSz
func (s *RGBAScanner) Filter(minSz int) int {
	s.set = s.set.filter(minSz)
	return len(s.set.objs)
}

//DrawObjects draw each object found in a specified color inside ci
func (s *RGBAScanner) DrawObjects(c color.Color) (ci ColorImage) {
	ci = NewColorFromImage(s.img.Bounds())
	s.set.drawWithColor(ci, c)
	return
}
