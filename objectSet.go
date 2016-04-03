package img

//objectSet is the list of objects inside an image map[idObj]obj
type objectSet map[int]*object

//objectsInLastRow obtains the objectSet which contains a list of objects in the last row
func (os objectSet) objectsInLastRow(row int) (o objectSet) {
	o = make(objectSet)
	for id, obj := range os {
		points := obj.points[row-1]
		if points != nil {
			o[id] = obj
		}
	}
	return
}

func (os objectSet) add(toAdd *object) {
	id := len(os)
	toAdd.id = id
	os[id] = toAdd
}

func (os objectSet) drop(toDrop ...int) {
	for _, obj := range toDrop {
		delete(os, obj)
	}
}

func (os objectSet) append(toAppend objectSet) {
	newID := len(os)
	for _, obj := range toAppend {
		obj.id = newID
		os[newID] = obj
		newID++
	}
}

func (os objectSet) update(changes objectSet) {
	for id, obj := range changes {
		os[id] = obj
	}
}

func (os objectSet) groupObjects(dest int, toBeAppened *object) {
	original := os[dest]
	for row, points := range toBeAppened.points {
		original.points[row] = append(original.points[row], points...)
	}
	os[original.id] = original
}
