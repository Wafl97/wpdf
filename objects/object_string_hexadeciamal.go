package objects

type Object_HexadeciamalString []byte

func New_Object_HexadeciamalString(obj Object_HexadeciamalString) *Object {
	return &Object{
		Type:              ObjectType_StringHexadeciamal,
		hexadecimalString: obj,
	}
}

func (os *Object_HexadeciamalString) String() string {
	return "0x" + string(*os)
}

func (os *Object_HexadeciamalString) PdfString() string {
	return "<" + string(*os) + ">"
}
