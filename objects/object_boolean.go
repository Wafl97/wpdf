package objects

import "fmt"

type Object_Boolean bool

func New_Object_Boolean(boolean bool) *Object {
	return &Object{
		Type:    ObjectType_Boolean,
		boolean: Object_Boolean(boolean),
	}
}

func (ob *Object_Boolean) String() string {
	return fmt.Sprintf("%t", *ob)
}

func (ob *Object_Boolean) PdfString() string {
	return fmt.Sprintf("%t", *ob)
}
