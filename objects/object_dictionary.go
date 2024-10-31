package objects

import (
	"fmt"
	"strings"
)

type Object_Dictionary map[string]*Object

func New_Object_Dictionary(obj Object_Dictionary) *Object {
	return &Object{
		Type:       ObjectType_Dictionary,
		dictionary: obj,
	}
}
func (od *Object_Dictionary) GetElement(key string) (*Object, bool) {
	element, hasElement := map[string]*Object(*od)[key]
	return element, hasElement
}

func (od *Object_Dictionary) String() string {
	strArr := make([]string, 0, len(*od))
	for key, val := range *od {
		strArr = append(strArr, fmt.Sprintf("%v : %v", key, val.String()))
	}
	return "{" + strings.Join(strArr, ", ") + "}"
}

func (od *Object_Dictionary) PdfString() string {
	strArr := make([]string, 0, len(*od))
	for key, val := range *od {
		strArr = append(strArr, fmt.Sprintf("/%v %v", key, val.PdfString()))
	}
	return "<<\n" + strings.Join(strArr, "\n") + "\n>>"
}
