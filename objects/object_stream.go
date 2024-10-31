package objects

import (
	"fmt"

	"github.com/Wafl97/wpdf/filters"
)

type Object_Stream []byte

func New_Object_Stream(obj Object_Stream) *Object {
	return &Object{
		Type:   ObjectType_Stream,
		stream: obj,
	}
}

func (os *Object_Stream) Decode(filter string) ([]byte, error) {
	switch filterFunc, contains := filters.Filters[filter]; contains {
	case true:
		decoded, err := filterFunc(*os)
		if err != nil {
			return nil, fmt.Errorf("stream.Decode: %w", err)
		}
		return decoded, nil
	default:
		return nil, fmt.Errorf("stream.Decode: unknown/unimplemented filter [%s]", filter)
	}
}

func (os *Object_Stream) String() string {
	return fmt.Sprintf("STREAM %d bytes", len(*os))
}

func (os *Object_Stream) PdfString() string {
	return fmt.Sprintf("stream\n%s\nendstream", string(*os))
}
