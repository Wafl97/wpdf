package objects

import "fmt"

type Object_IntegerNumber int

func New_Object_IntegerNumber(integer int) *Object {
	return &Object{
		Type:          ObjectType_IntergerNumber,
		integerNumber: Object_IntegerNumber(integer),
	}
}

func (oi *Object_IntegerNumber) Value() int {
	return int(*oi)
}

func (oi *Object_IntegerNumber) String() string {
	return fmt.Sprintf("%d", *oi)
}

func (oi *Object_IntegerNumber) PdfString() string {
	return fmt.Sprintf("%d", *oi)
}
