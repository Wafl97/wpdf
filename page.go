package wpdf

import (
	"bytes"
	"strings"
)

type Line struct {
}

type Page struct {
	pageNr        int
	decodedStream []byte
	lines         []string
}

// TODO:
func (p *Page) Text() string {
	if p.lines != nil || len(p.lines) == 0 {
		p.lines = newPageReader(p.decodedStream).readLines()
	}
	return strings.Join(p.lines, "\n")
}

func (p *Page) PageNr() int {
	return p.pageNr
}

func (p *Page) DecodedStream() []byte {
	return p.decodedStream
}

type PageReader struct {
	decodedStream []byte
	index         int
	textBuffer    *bytes.Buffer
	lines         []string
}

func newPageReader(decodedStream []byte) *PageReader {
	return &PageReader{decodedStream: decodedStream, index: 0, textBuffer: new(bytes.Buffer)}
}

func (p *PageReader) readLines() []string {
	for p.index < len(p.decodedStream) {
		if bytes.Equal(p.decodedStream[p.index:p.index+2], []byte("BT")) {
			p.index += 2
			p.extractLine()
		}
		p.index++
	}
	return p.lines
}

func (p *PageReader) extractLine() {
	for p.index < len(p.decodedStream) {
		if bytes.Equal(p.decodedStream[p.index:p.index+2], []byte("ET")) {
			p.lines = append(p.lines, p.textBuffer.String())
			p.textBuffer.Reset()
			p.index += 2
			return
		}
		//if p.decodedStream[p.index] == '[' {
		//	p.index++
		//	p.readLine()
		//}
		p.textBuffer.WriteByte(p.decodedStream[p.index])

		p.index++
	}
}

func (p *PageReader) readLine() {
	for p.index < len(p.decodedStream) {
		char := p.decodedStream[p.index]
		switch char {
		case ']':
			p.lines = append(p.lines, p.textBuffer.String())
			p.textBuffer.Reset()
			return
		default:
			p.textBuffer.WriteByte(p.decodedStream[p.index])
			p.index++
		}
	}
}
