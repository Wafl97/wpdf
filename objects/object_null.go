package objects

type Object_Null struct{}

func New_Object_Null() *Object {
	return &Object{Type: ObjectType_Null, null: Object_Null{}}
}

func (on *Object_Null) String() string {
	return "null"
}

func (on *Object_Null) PdfString() string {
	return "null"
}
