package objects

type Object_LiteralString string

func New_Object_LiteralString(obj Object_LiteralString) *Object {
	return &Object{
		Type:          ObjectType_StringLiteral,
		literalString: obj,
	}
}

func (ol *Object_LiteralString) String() string {
	return string(*ol)
}

func (ol *Object_LiteralString) PdfString() string {
	return "(" + string(*ol) + ")"
}
