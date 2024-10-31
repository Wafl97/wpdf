package wpdf

import (
	"bytes"
	"fmt"
	. "github.com/Wafl97/wpdf/objects"
	"strconv"
	"strings"
)

type XRefEntry struct {
	offset int
	gen    int
	object *Object
}

func newXRefEntry(buffer []byte) (*XRefEntry, error) {
	bufferSplit := bytes.Split(buffer, []byte(" "))
	if len(bufferSplit) < 2 {
		return nil, fmt.Errorf("newXRefEntry: failed to get offset and generation number: %q", string(buffer))
	}
	offset, err := strconv.Atoi(string(bufferSplit[0]))
	if err != nil {
		return nil, fmt.Errorf("newXRefEntry: failed to get offset: %q", string(buffer))
	}
	gen, err := strconv.Atoi(string(bufferSplit[1]))
	if err != nil {
		return nil, fmt.Errorf("newXRefEntry: failed to get generation number: %q", string(buffer))
	}
	return &XRefEntry{offset: offset, gen: gen, object: nil}, nil
}

type XRefTable struct {
	firstObj int
	table    []XRefEntry
	reader   *DocumentReader
}

func NewXRefTable(firstObj, size int) *XRefTable {
	return &XRefTable{
		firstObj: firstObj,
		table:    make([]XRefEntry, size),
	}
}

func (ot *XRefTable) GetObject(id int) (*Object, error) {
	index := id - ot.firstObj
	if ot.table[index].object != nil {
		return ot.table[index].object, nil
	}
	obj, err := ot.reader.readObjectFromOffset(ot.table[index].offset)
	if err != nil {
		return nil, err
	}
	ot.table[index].object = obj
	return obj, nil
}

type Trailer struct {
	xRefTable XRefTable
	summary   Object_Dictionary
	prev      *Trailer
}

func (t *Trailer) GetObject(id int) (*Object, error) {
	if t.xRefTable.firstObj <= id && len(t.xRefTable.table) > id {
		//fmt.Printf("reading object %d\n", id)
		return t.xRefTable.GetObject(id)
	} else if t.prev != nil {
		//fmt.Printf("looking for object %d\n", id)
		return t.prev.GetObject(id)
	}
	return nil, fmt.Errorf("xref: failed to find offset for object %d", id)
}

func (t *Trailer) String() string {
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("first obj = %d, size = %d", t.xRefTable.firstObj, len(t.xRefTable.table)))
	if t.prev != nil {
		sb.WriteString(fmt.Sprintf(" prev: %s", t.prev.String()))
	}
	return fmt.Sprintf("Trailer:\n%s\nXRefTables: %s", t.summary.String(), sb.String())
}
