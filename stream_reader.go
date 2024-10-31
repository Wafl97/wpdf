package wpdf

type StreamReader struct {
	stream []byte
	index  int
}

func NewStreamReader(stream []byte) *StreamReader {
	return &StreamReader{stream: stream, index: 0}
}
