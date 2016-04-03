package img

type Scanner struct {
	set objectSet
}

func NewScanner() *Scanner {
	return &Scanner{set: make(objectSet)}
}

func (s *Scanner) SearchObjects(bi BinaryImage) {
	rows := bi.Bounds().Max.X
	for i := 0; i < rows; i++ {
		lastRow := s.set.objectsInLastRow(i)
		currentRowObjects := s.scanRow(bi, i)
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
					lastRow.add(currentObject)
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
			s.set.update(lastRow)
		}
	}
}

func (s *Scanner) scanRow(bi BinaryImage, row int) objectSet {
	os := make(objectSet)
	currentObj := newObject()
	objInProgress := false
	id := 0
	for y := range bi[row] {
		if bi[row][y] {
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
