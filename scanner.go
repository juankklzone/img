package img

import (
	"fmt"
	"sync"
	"time"
)

//Scanner finds objects in a binary image.
type Scanner struct {
	set *objectSet
	bi  BinaryImage
}

//NewScanner creates a scanner
func NewScanner(i BinaryImage) *Scanner {
	return &Scanner{set: newObjectset(), bi: i}
}

//SearchObjects count the objects found in the image
//and save them for later use
func (s *Scanner) SearchObjects() int {
	rows := s.bi.Bounds().Max.X
	wg := new(sync.WaitGroup)
	wg.Add(2)
	half := rows / 2
	ini := time.Now()
	go s.scanRows(wg, 0, half)
	go s.scanRows(wg, half, rows)
	wg.Wait()
	now := time.Now()
	s.joinSegments(half)
	fmt.Println("time req for joinSegments ", time.Since(now))
	fmt.Println("total time ", time.Since(ini))
	return len(s.set.objs)
}

func (s *Scanner) scanRows(wg *sync.WaitGroup, initial, final int) {
	defer wg.Done()
	for i := initial; i < final; i++ {
		lastRow := s.set.objectsInLastRow(i)
		currentRowObjects := s.scanRow(i)
		if len(lastRow.objs) == 0 {
			s.set.append(currentRowObjects)
			continue
		}
		for _, currentObject := range currentRowObjects.objs {
			var belongsTo []int
			min, max := currentObject.minInRow(i), currentObject.maxInRow(i)
			for _, upperObj := range lastRow.objs {
				if upperObj.hasPointsInRow(i-1, min, max) {
					belongsTo = append(belongsTo, upperObj.id)
				}
			}
			lenParents := len(belongsTo)
			if lenParents == 0 {
				s.set.add(currentObject)
			} else {
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

func (s *Scanner) joinSegments(row int) {
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

func (s *Scanner) scanRow(row int) *objectSet {
	os := newObjectset()
	currentObj := newObject()
	objInProgress := false
	id := 0
	for y := range s.bi[row] {
		if s.bi[row][y] {
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
func (s *Scanner) Filter(minSz int) int {
	s.set = s.set.filter(minSz)
	return len(s.set.objs)
}

//DrawObjects draw each object found in a random color inside ci
func (s *Scanner) DrawObjects() (ci ColorImage) {
	ci = NewColorFromImage(s.bi.Bounds())
	s.set.draw(ci)
	return
}
