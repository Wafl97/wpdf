package objects

import "fmt"

type Object_RealNumber struct {
	first int // -> nn.nn
	last  int //    nn.nn <-
}

func New_Object_RealNumber(first, last int) *Object {
	return &Object{Type: ObjectType_NumberReal, realNumber: Object_RealNumber{first, last}}
}

func (or *Object_RealNumber) String() string {
	return fmt.Sprintf("%d.%d", or.first, or.last)
}

func (or *Object_RealNumber) PdfString() string {
	return fmt.Sprintf("%d.%d", or.first, or.last)
}
