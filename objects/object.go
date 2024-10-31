package objects

type ObjectType uint

const (
	ObjectType_Array ObjectType = iota
	ObjectType_Dictionary
	ObjectType_Name
	ObjectType_Stream
	ObjectType_StringLiteral
	ObjectType_StringHexadeciamal
	ObjectType_NumberReal
	ObjectType_IntergerNumber
	ObjectType_Indirect
	ObjectType_Ref
	ObjectType_Boolean
	ObjectType_Null
)

type Object struct {
	Type              ObjectType
	indirect          Object_Indirect
	dictionary        Object_Dictionary
	array             Object_Array
	stream            Object_Stream
	name              Object_Name
	literalString     Object_LiteralString
	hexadecimalString Object_HexadeciamalString
	realNumber        Object_RealNumber
	integerNumber     Object_IntegerNumber
	ref               Object_Ref
	boolean           Object_Boolean
	null              Object_Null
}

func (o *Object) AddIndirectObject(obj *Object) {
	o.indirect.direct[obj.Type] = obj
}

func (o *Object) AddStream(readStream func(length int) *Object) {
	lengthObject, _ := o.indirect.direct[ObjectType_Dictionary].dictionary.GetElement("Length")
	o.indirect.direct[ObjectType_Stream] = readStream(int(lengthObject.integerNumber))
}

func (o *Object) AsIndirect() *Object_Indirect {
	return &o.indirect
}

func (o *Object) IsIndirect() (*Object_Indirect, bool) {
	return &o.indirect, o.Type == ObjectType_Indirect
}

func (o *Object) AsDictionary() *Object_Dictionary {
	return &o.dictionary
}

func (o *Object) IsDictionary() (*Object_Dictionary, bool) {
	return &o.dictionary, o.Type == ObjectType_Dictionary
}

func (o *Object) AsArray() *Object_Array {
	return &o.array
}

func (o *Object) IsArray() (*Object_Array, bool) {
	return &o.array, o.Type == ObjectType_Array
}

func (o *Object) AsStream() *Object_Stream {
	return &o.stream
}

func (o *Object) IsStream() (*Object_Stream, bool) {
	return &o.stream, o.Type == ObjectType_Stream
}

func (o *Object) AsName() *Object_Name {
	return &o.name
}

func (o *Object) IsName() (*Object_Name, bool) {
	return &o.name, o.Type == ObjectType_Name
}

func (o *Object) IsLiteral() (*Object_LiteralString, bool) {
	return &o.literalString, o.Type == ObjectType_StringLiteral
}

func (o *Object) AsInteger() *Object_IntegerNumber {
	return &o.integerNumber
}

func (o *Object) IsInteger() (*Object_IntegerNumber, bool) {
	return &o.integerNumber, o.Type == ObjectType_IntergerNumber
}

func (o *Object) AsRef() *Object_Ref {
	return &o.ref
}

func (o *Object) IsRef() (*Object_Ref, bool) {
	return &o.ref, o.Type == ObjectType_Ref
}

func (o *Object) String() string {
	switch o.Type {
	case ObjectType_Indirect:
		return o.indirect.String()
	case ObjectType_Dictionary:
		return o.dictionary.String()
	case ObjectType_StringLiteral:
		return o.literalString.String()
	case ObjectType_Array:
		return o.array.String()
	case ObjectType_Name:
		return o.name.String()
	case ObjectType_Stream:
		return o.stream.String()
	case ObjectType_Ref:
		return o.ref.String()
	case ObjectType_StringHexadeciamal:
		return o.hexadecimalString.String()
	case ObjectType_NumberReal:
		return o.realNumber.String()
	case ObjectType_IntergerNumber:
		return o.integerNumber.String()
	case ObjectType_Boolean:
		return o.boolean.String()
	case ObjectType_Null:
		return o.null.String()
	}
	return "NO-OBJ"
}

func (o *Object) PdfString() string {
	switch o.Type {
	case ObjectType_Indirect:
		return o.indirect.PdfString()
	case ObjectType_Dictionary:
		return o.dictionary.PdfString()
	case ObjectType_StringLiteral:
		return o.literalString.PdfString()
	case ObjectType_Array:
		return o.array.PdfString()
	case ObjectType_Name:
		return o.name.PdfString()
	case ObjectType_Stream:
		return o.stream.PdfString()
	case ObjectType_Ref:
		return o.ref.PdfString()
	case ObjectType_StringHexadeciamal:
		return o.hexadecimalString.PdfString()
	case ObjectType_NumberReal:
		return o.realNumber.PdfString()
	case ObjectType_IntergerNumber:
		return o.integerNumber.PdfString()
	case ObjectType_Boolean:
		return o.boolean.PdfString()
	case ObjectType_Null:
		return o.null.PdfString()
	}
	return "NO-OBJ"
}
