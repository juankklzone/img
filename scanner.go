package img

type Scanner struct {
	set objectSet
	bi  BinaryImage
}

func NewScanner(i BinaryImage) *Scanner {
	return &Scanner{set: make(objectSet), bi: i}
}

func (s *Scanner) SearchObjects() int {
	rows := s.bi.Bounds().Max.X
	for i := 0; i < rows; i++ {
		lastRow := s.set.objectsInLastRow(i)
		currentRowObjects := s.scanRow(i)
		if len(lastRow) == 0 {
			s.set.append(currentRowObjects)
		} else {
			for _, currentObject := range currentRowObjects {
				var belongsTo []int
				min, max := currentObject.minInRow(i), currentObject.maxInRow(i)
				for _, upperObj := range lastRow {
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
						lastRow.groupObjects(initial, lastRow[belongsTo[i]])
					}
					lastRow.groupObjects(initial, currentObject)
					s.set.drop(belongsTo[1:]...)
					lastRow.drop(belongsTo[1:]...)
				}
			}
		}
	}
	return len(s.set)
}

func (s *Scanner) scanRow(row int) objectSet {
	os := make(objectSet)
	currentObj := newObject()
	objInProgress := false
	id := 0
	for y := range s.bi[row] {
		if s.bi[row][y] {
			currentObj.append(row, y)
			objInProgress = true
		} else if objInProgress {
			currentObj.id = id
			os[id] = currentObj
			id++
			objInProgress = false
			currentObj = newObject()
		}
	}
	if objInProgress {
		currentObj.id = id
		os[id] = currentObj
	}
	return os
}

func (s *Scanner) Filter(minSz int) int {
	s.set = s.set.filter(minSz)
	return len(s.set)
}

func (s *Scanner) DrawObjects() (ci ColorImage) {
	ci = NewColorFromImage(s.bi.Bounds())
	s.set.draw(ci)
	return
}
