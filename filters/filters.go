package filters

type Filter func(encoded []byte) ([]byte, error)

var Filters map[string]Filter = map[string]Filter{
	"FlateDecode": FlateDecode_Decode,
}
