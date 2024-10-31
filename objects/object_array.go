package objects

import "strings"

type Object_Array []*Object

func New_Object_Array(obj Object_Array) *Object {
	return &Object{
		Type:  ObjectType_Array,
		array: obj,
	}
}

func (oa *Object_Array) String() string {
	strArr := make([]string, len(*oa))
	for idx, obj := range *oa {
		strArr[idx] = obj.String()
	}
	return "[" + strings.Join(strArr, ", ") + "]"
}

func (oa *Object_Array) PdfString() string {
	strArr := make([]string, len(*oa))
	for idx, obj := range *oa {
		strArr[idx] = obj.PdfString()
	}
	return "[" + strings.Join(strArr, " ") + "]"
}
