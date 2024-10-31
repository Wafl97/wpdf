package filters

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

func FlateDecode_Decode(encoded []byte) ([]byte, error) {
	buffer := bytes.NewReader(encoded)
	reader, err := zlib.NewReader(buffer)
	if err != nil {
		return nil, fmt.Errorf("FlateDecode: %w", err)
	}
	defer reader.Close()
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("FlateDecode: %w", err)
	}
	return decoded, nil
}
