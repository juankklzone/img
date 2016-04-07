package img

import (
	"bytes"
	"fmt"
	"image/draw"
	"sync"
)

//objectSet is the list of objects inside an image map[idObj]obj
type objectSet struct {
	objs map[int]*object
	mtx  sync.RWMutex
}

func newObjectset() *objectSet {
	return &objectSet{objs: make(map[int]*object)}
}

//objectsInLastRow obtains the objectSet which contains a list of objects in the last row
func (os *objectSet) objectsInLastRow(row int) (o *objectSet) {
	o = newObjectset()
	os.mtx.RLock()
	defer os.mtx.RUnlock()
	for id, obj := range os.objs {
		points := obj.get(row - 1)
		if points != nil {
			o.set(id, obj)
		}
	}
	return
}

func (os *objectSet) add(toAdd *object) {
	os.mtx.Lock()
	defer os.mtx.Unlock()
	id := len(os.objs)
	for os.objs[id] != nil {
		id += 2
	}
	toAdd.id = id
	os.objs[id] = toAdd
}

func (os *objectSet) drop(toDrop ...int) {
	os.mtx.Lock()
	defer os.mtx.Unlock()
	for _, obj := range toDrop {
		delete(os.objs, obj)
	}
}

func (os *objectSet) append(toAppend *objectSet) {
	toAppend.mtx.RLock()
	defer toAppend.mtx.RUnlock()
	for _, obj := range toAppend.objs {
		os.add(obj)
	}
}

func (os *objectSet) update(changes *objectSet) {
	os.mtx.Lock()
	defer os.mtx.Unlock()
	for id, obj := range changes.objs {
		os.objs[id] = obj
	}
}

func (os *objectSet) groupObjects(dest int, toBeAppened *object) {
	os.mtx.Lock()
	defer os.mtx.Unlock()
	original := os.objs[dest]
	for row, points := range toBeAppened.points {
		newpoints := append(original.get(row), points...)
		original.set(row, newpoints)
	}
	os.objs[original.id] = original
}

func (os *objectSet) String() string {
	os.mtx.RLock()
	defer os.mtx.RUnlock()
	buf := new(bytes.Buffer)
	for id := range os.objs {
		fmt.Fprintf(buf, "%d\t", id)
		//fmt.Fprintf(buf, "%s \n\n", obj.String())
	}
	return buf.String()
}

func (os *objectSet) filter(minSize int) (s *objectSet) {
	os.mtx.RLock()
	defer os.mtx.RUnlock()
	s = newObjectset()
	for id, obj := range os.objs {
		if obj.len() >= minSize {
			s.objs[id] = obj
		}
	}
	return
}

func (os *objectSet) draw(img draw.Image) {
	os.mtx.RLock()
	defer os.mtx.RUnlock()
	colors := randColors(len(os.objs))
	idx := 0
	for _, o := range os.objs {
		o.draw(img, colors[idx])
		idx++
	}
}

func (os *objectSet) get(i int) *object {
	os.mtx.RLock()
	defer os.mtx.RUnlock()
	return os.objs[i]
}

func (os *objectSet) set(i int, o *object) {
	os.mtx.Lock()
	defer os.mtx.Unlock()
	os.objs[i] = o
}
