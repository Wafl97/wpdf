package objects

type Object_Name string

func New_Object_Name(obj Object_Name) *Object {
	return &Object{
		Type: ObjectType_Name,
		name: obj,
	}
}

func (on *Object_Name) String() string {
	return string(*on)
}

func (on *Object_Name) PdfString() string {
	return "/" + string(*on)
}
