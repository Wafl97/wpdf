package objects

import "fmt"

type Object_Ref struct {
	id  int
	gen int
}

func New_Object_Ref(id, gen int) *Object {
	return &Object{Type: ObjectType_Ref, ref: Object_Ref{id, gen}}
}

func (or *Object_Ref) ID() int {
	return or.id
}

func (or *Object_Ref) String() string {
	return fmt.Sprintf("REF (%d %d)", or.id, or.gen)
}

func (or *Object_Ref) PdfString() string {
	return fmt.Sprintf("%d %d R", or.id, or.gen)
}
