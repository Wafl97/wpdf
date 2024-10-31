package objects

import (
	"fmt"
	"strings"
)

type Object_Direct map[ObjectType]*Object

func (od *Object_Direct) String() string {
	strArr := make([]string, 0, len(*od))
	for _, obj := range *od {
		strArr = append(strArr, obj.String())
	}
	return strings.Join(strArr, " ")
}

func (od *Object_Direct) PdfString() string {
	strArr := make([]string, 0, len(*od))
	for _, obj := range *od {
		strArr = append(strArr, obj.PdfString())
	}
	return strings.Join(strArr, "\n")
}

type Object_Indirect struct {
	id     int
	gen    int
	direct Object_Direct
}

func New_Object_Indirect(id, gen int) *Object {
	return &Object{
		Type:     ObjectType_Indirect,
		indirect: Object_Indirect{id, gen, make(Object_Direct, 4)},
	}
}

func (oi *Object_Indirect) Contains(objectType ObjectType) (*Object, bool) {
	obj, hasType := oi.direct[objectType]
	return obj, hasType
}

func (oi *Object_Indirect) ID() int {
	return oi.id
}

func (oi *Object_Indirect) String() string {
	return fmt.Sprintf("OBJ (%d %d) %s", oi.id, oi.gen, oi.direct.String())
}

func (oi *Object_Indirect) PdfString() string {
	return fmt.Sprintf("%d %d obj\n%s\nendobj\n", oi.id, oi.gen, oi.direct.PdfString())
}
