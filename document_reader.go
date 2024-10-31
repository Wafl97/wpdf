package wpdf

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"

	. "github.com/Wafl97/wpdf/errors"
	. "github.com/Wafl97/wpdf/objects"
	. "github.com/Wafl97/wpdf/util"
)

type Token string

var (
	DictionaryBegin        Token = "<<"
	DictionaryEnd          Token = ">>"
	ArrayBegin             Token = "["
	ArrayEnd               Token = "]"
	NameDelimiter          Token = "/"
	ObjectBegin            Token = "obj"
	ObjectEnd              Token = "endobj"
	StreamBegin            Token = "stream"
	StreamEnd              Token = "endstream"
	LiteralStringBegin     Token = "("
	LiteralStringEnd       Token = ")"
	HexadecimalStringBegin Token = "<"
	HexadecimalStringEnd   Token = ">"
	ObjectRef              Token = "R"
	BooleanTrue            Token = "true"
	BooleanFalse           Token = "false"
	Null                   Token = "null"
	StartXref              Token = "startxref"
	Xref                   Token = "xref"
	TrailerBegin           Token = "trailer"
	EOF                    Token = "%%EOF"
	BeginTextObj           Token = "BT"
	EndTextObj             Token = "ET"
	TJ                     Token = "TJ"
	Tj                     Token = "Tj"
	Tf                     Token = "Tf"
	Td                     Token = "Td"
)

var eol = "\n"

type DocumentReader struct {
	//tokens []Token
	index int
	//revIndex int
	buffer  []byte
	objects []*Object
	trailer *Trailer
}

func NewDocumentReader(bytes []byte) *DocumentReader {
	return &DocumentReader{buffer: bytes}
}

func (r *DocumentReader) ReadPageCount() (int, error) {
	if r.trailer == nil {
		return 0, fmt.Errorf("reader.ReadPageCount: trailer is nil")
	}
	rootObjRef, hasRoot := r.trailer.summary.GetElement("Root")
	if !hasRoot {
		return 0, fmt.Errorf("reader.ReadPageCount: trailer has no Root object")
	}
	rootObj, err := r.trailer.GetObject(rootObjRef.AsRef().ID())
	if err != nil {
		return 0, err
	}
	rootDictionary, hasDictionary := rootObj.AsIndirect().Contains(ObjectType_Dictionary)
	if !hasDictionary {
		return 0, fmt.Errorf("reader.ReadPageCount: Root object does not contain a dictionary")
	}
	pagesObjRef, hasPagesRef := rootDictionary.AsDictionary().GetElement("Pages")
	if !hasPagesRef {
		return 0, fmt.Errorf("reader.ReadPageCount: Root object has to Pages ref")
	}
	pagesObj, err := r.trailer.GetObject(pagesObjRef.AsRef().ID())
	if err != nil {
		return 0, err
	}
	pagesDictionary, hasDictionary := pagesObj.AsIndirect().Contains(ObjectType_Dictionary)
	if !hasDictionary {
		return 0, fmt.Errorf("reader.ReadPageCount: Pages object does not contain a dictionary")
	}
	pageCount, hasPageCount := pagesDictionary.AsDictionary().GetElement("Count")
	if !hasPageCount {
		return 0, fmt.Errorf("reader.ReadPageCount: Pages dictionary does not contains a Count: %s", pagesObj.String())
	}
	return pageCount.AsInteger().Value(), nil
}

func (r *DocumentReader) ReadHeader() (*Header, error) {
	if r.buffer == nil {
		return nil, fmt.Errorf("document reader -> read trailer: no buffer provided")
	}
	headerBuffer := new(bytes.Buffer)
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			return nil, ErrInvalidHeader
		}
		if char == '\r' || char == '\n' {
			break
		}
		headerBuffer.WriteRune(rune(char))
	}
	version, err := VersionFromString(headerBuffer.String())
	if err != nil {
		return nil, err
	}
	return &Header{version: *version}, nil
}

/*
1. look for EOF
2. find EOL marker
3. find startxref offset
4. read xref table
5. read following trailer dictionary
6. if trailer dictionary has previous repeat from 3
*/
func (r *DocumentReader) ReadTrailer() (*Trailer, error) {
	r.index = len(r.buffer) - 1

	// trim end
	for r.index >= 0 {
		char := r.buffer[r.index]
		if char == '\r' || char == '\n' { // TODO: rewrite this logic flow
			r.buffer = r.buffer[:r.index]
			r.index--
		} else {
			r.index++
			break
		}
	}

	// look for %%EOF
	if !r.checkForRev(EOF) {
		fmt.Printf("%q", r.buffer[r.index-len(EOF):r.index])
		return nil, fmt.Errorf("reader->read trailer: invalid trailer, missing %%%%EOF")
	}

	// find eol marker
	if bytes.Equal(r.buffer[r.index-2:r.index], []byte("\r\n")) {
		eol = "\r\n" // default is just "\n"
	}
	// find 'startxref'
	for r.index >= 0 {
		if r.checkForRev(StartXref) {
			r.index += len(StartXref)
			break
		}
		r.index--
	}
	// read until number
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			return nil, err
		}
		if IsNumber(char) {
			r.unreadByte()
			break
		}
	}

	// read following offset
	offsetBuffer := new(bytes.Buffer)
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			return nil, err
		}
		if char == '\n' || char == '\r' {
			break
		}
		offsetBuffer.WriteByte(char)
	}
	offset, err := strconv.Atoi(offsetBuffer.String())
	if err != nil {
		return nil, fmt.Errorf("reader.ReadTrailer: failed to parse offset for xref table: %w", err)
	}

	r.trailer, err = r.readTrailer(offset)

	return r.trailer, err
}

// FIXME: expects eol which might be different
// example: eol = \r\n, line = xref\n
func (r *DocumentReader) readTrailer(xrefOffset int) (*Trailer, error) {

	// go to offset and past 'xref'
	r.index = xrefOffset

	fmt.Printf("%q\n", string(r.buffer[r.index:r.index+20]))

	// read until number
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			return nil, err
		}
		if IsNumber(char) {
			r.unreadByte()
			break
		}
	}

	// read the following 2 numbers (start id and size)
	xrefHeaderBuffer := new(bytes.Buffer)
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			return nil, err
		}
		xrefHeaderBuffer.WriteByte(char)
		if bytes.Equal(r.buffer[r.index:r.index+len(eol)], []byte(eol)) {
			r.index += len(eol)
			break
		}
	}
	xrefHeaderBufferSplit := strings.Split(xrefHeaderBuffer.String(), " ")
	if len(xrefHeaderBufferSplit) != 2 {
		return nil, fmt.Errorf("reader.readTrailer: failed to read xref table at offset %d", xrefOffset)
	}
	firstObjId, err := strconv.Atoi(xrefHeaderBufferSplit[0])
	if err != nil {
		return nil, fmt.Errorf("reader.readTrailer: failed to parse firstObjId for xref table: %w", err)
	}
	size, err := strconv.Atoi(xrefHeaderBufferSplit[1])
	if err != nil {
		return nil, fmt.Errorf("reader.readTrailer: failed to parse size for xref table: %w", err)
	}
	trailer := new(Trailer)
	trailer.xRefTable = *NewXRefTable(firstObjId, size)
	trailer.xRefTable.reader = r
	// read all the following offsets
	for row := 0; row < size; row++ {
		xrefEntry, err := newXRefEntry(r.buffer[r.index : r.index+20])
		if err != nil {
			return nil, err
		}
		trailer.xRefTable.table[row] = *xrefEntry
		r.index += 20
	}
	// go past 'trailer'
	r.index += len(TrailerBegin) + len(eol)
	// read dictionary
	trailer.summary = *r.readDictionary().AsDictionary()
	// if dictionary contains 'Prev', then go to that offset and repeat the process
	if prevOffsetObj, hasPrev := trailer.summary.GetElement("Prev"); hasPrev {
		prevTrailer, err := r.readTrailer(prevOffsetObj.AsInteger().Value())
		if err != nil {
			return nil, err
		}
		trailer.prev = prevTrailer
	}
	return trailer, nil
}

func (r *DocumentReader) BuildPageTree(pages *[]*Page) {
	// TODO: error handling
	rootObjRef, _ := r.trailer.summary.GetElement("Root")
	rootObj, _ := r.trailer.GetObject(rootObjRef.AsRef().ID())
	rootDictionary, _ := rootObj.AsIndirect().Contains(ObjectType_Dictionary)
	pageTreeRootRef, _ := rootDictionary.AsDictionary().GetElement("Pages")
	pageTreeRoot, _ := r.trailer.GetObject(pageTreeRootRef.AsRef().ID())
	if err := r.walkPageTree(pageTreeRoot, pages); err != nil {
		fmt.Println(err)
	}
}

func (r *DocumentReader) walkPageTree(treeNodeObject *Object, pages *[]*Page) error {
	treeNodeDictionary, _ := treeNodeObject.AsIndirect().Contains(ObjectType_Dictionary)
	nodeType, hasType := treeNodeDictionary.AsDictionary().GetElement("Type")
	if !hasType {
		return fmt.Errorf("wpdf: page tree node %d does not have a Type attribute", treeNodeObject.AsIndirect().ID())
	}
	if nodeType.AsName().String() == "Page" {
		page, err := r.readPage(treeNodeObject)
		if err != nil {
			return err
		}
		*pages = append(*pages, page)
		return nil
	}
	kids, _ := treeNodeDictionary.AsDictionary().GetElement("Kids")
	kidsArray := kids.AsArray()
	for kid := 0; kid < len(*kidsArray); kid++ {
		pageNode, _ := r.trailer.GetObject((*kidsArray)[kid].AsRef().ID())
		if err := r.walkPageTree(pageNode, pages); err != nil {
			return err
		}
	}
	return nil
}

func (r *DocumentReader) readPage(pageObj *Object) (*Page, error) {
	pageObjDictionary, _ := pageObj.AsIndirect().Contains(ObjectType_Dictionary)
	contentsObjRef, _ := pageObjDictionary.AsDictionary().GetElement("Contents")
	contentsObjRefArr, isArray := contentsObjRef.IsArray()
	if isArray {
		decodedStream := new(bytes.Buffer)
		for _, refObj := range *contentsObjRefArr {
			object, _ := r.trailer.GetObject(refObj.AsRef().ID())
			decoded, err := r.readPageContents(object)
			if err != nil {
				return nil, err
			}
			decodedStream.Write(decoded)
		}
		return &Page{decodedStream: decodedStream.Bytes()}, nil
	}

	contentsObj, _ := r.trailer.GetObject(contentsObjRef.AsRef().ID())
	decodedStream, err := r.readPageContents(contentsObj)
	if err != nil {
		return nil, err
	}
	return &Page{decodedStream: decodedStream}, nil
}

func (r *DocumentReader) readPageContents(contentsObj *Object) ([]byte, error) {
	contentsDictionary, _ := contentsObj.AsIndirect().Contains(ObjectType_Dictionary)
	filterObj, _ := contentsDictionary.AsDictionary().GetElement("Filter")
	filter := filterObj.AsName().String()
	streamObj, _ := contentsObj.AsIndirect().Contains(ObjectType_Stream)
	return streamObj.AsStream().Decode(filter)
}

func (r *DocumentReader) readNextByte() (byte, error) {
	if r.index >= len(r.buffer) {
		return byte(0), fmt.Errorf("next: index %d is out of range of buffer %d", r.index, len(r.buffer))
	}
	char := r.buffer[r.index]
	r.index++
	return char, nil
}

func (r *DocumentReader) checkFor(keyword Token) bool {
	if r.index+len(keyword) >= len(r.buffer) {
		return false
	}
	if bytes.Equal(r.buffer[r.index:r.index+len(keyword)], []byte(keyword)) {
		r.index += len(keyword)
		return true
	}
	return false
}

func (r *DocumentReader) checkForRev(keyword Token) bool {
	if r.index-len(keyword) < 0 {
		return false
	}
	if bytes.Equal(r.buffer[r.index-len(keyword):r.index], []byte(keyword)) {
		r.index -= len(keyword)
		return true
	}
	return false
}

func (r *DocumentReader) unreadByte() error {
	if r.index <= 0 {
		return fmt.Errorf("prev: index %d is out of range of buffer %d", r.index, len(r.buffer))
	}
	r.index--
	return nil
}

func (r *DocumentReader) checkForRef() bool {
	buffer := make([]byte, 0, 16)

	idx := r.index
	for idx < len(r.buffer) { // read until numerical char is found
		char := r.buffer[idx]
		if IsNumerical(char) { // numerical char found
			break
		}
		if !IsWhiteSpace(char) { // only white space is expected for ref not letters of delimiters
			//fmt.Println("failed on char", string(char))
			return false
		}
		idx++
	}
	// look for the object id
	for idx < len(r.buffer) { // read until not digit or '+' / '-'
		char := r.buffer[idx]
		if char == ' ' { // whitespace indicated separation between values
			idx++ // use the next index before breaking
			break
		}
		if IsNumerical(char) { // gather all the numerical chars
			buffer = append(buffer, char)
		} else if IsDelimiter(char) { //
			return false
		}
		idx++
	}

	// do the same again for the generation number
	for idx < len(r.buffer) { // read until not digit or '+' / '-'
		char := r.buffer[idx]
		if char == ' ' { // whitespace indicated separation between values
			idx++ // use the next index before breaking
			break
		}
		if IsNumerical(char) { // gather all the numerical chars
			buffer = append(buffer, char)
		} else if IsDelimiter(char) { //
			return false
		}
		idx++
	}
	if idx >= len(r.buffer) {
		return false
	}
	if Token(r.buffer[idx]) != ObjectRef { // check for ref keyword 'R'
		return false
	}
	r.index = idx + 1 // increment index to 1 more so the position is after the 'R'
	return true
}

func (r *DocumentReader) readNextToken() (Token, error) {
	for r.index <= len(r.buffer) {

		if r.checkFor(ObjectBegin) {
			return ObjectBegin, nil
		}
		if r.checkFor(ObjectEnd) {
			return ObjectEnd, nil
		}
		if r.checkFor(StreamBegin) {
			return StreamBegin, nil
		}
		if r.checkFor(StreamEnd) {
			return StreamEnd, nil
		}
		if r.checkFor(DictionaryBegin) {
			return DictionaryBegin, nil
		}
		if r.checkFor(DictionaryEnd) {
			return DictionaryEnd, nil
		}
		if r.checkFor(BooleanTrue) {
			return BooleanTrue, nil
		}
		if r.checkFor(BooleanFalse) {
			return BooleanFalse, nil
		}
		if r.checkFor(Null) {
			return Null, nil
		}
		if r.checkForRef() {
			return ObjectRef, nil
		}

		char, err := r.readNextByte()
		if err != nil {
			return "", err
		}

		switch char {
		case ' ', '\t', '\r', '\n':
			continue
		case '<':
			return HexadecimalStringBegin, nil
		case '>':
			return HexadecimalStringEnd, nil
		case '[':
			return ArrayBegin, nil
		case ']':
			return ArrayEnd, nil
		case '(':
			return LiteralStringBegin, nil
		case ')':
			return LiteralStringEnd, nil
		case '/':
			return NameDelimiter, nil
		//case 'R':
		//	return ObjectRef, nil
		default:
			return Token(char), nil
		}
	}
	return "NO-TOKEN", nil
}

func (r *DocumentReader) handleContentsRef(pageCount int, contentsRef *Object_Ref) (*Page, error) {
	actualContents, isIndirect := r.objects[contentsRef.ID()].IsIndirect()
	if !isIndirect {
		return nil, fmt.Errorf("read page contents: object is not indirect")
	}
	contentsDictionaryObj, hasDictionary := actualContents.Contains(ObjectType_Dictionary)
	if !hasDictionary {
		return nil, fmt.Errorf("read page contents: object is has no dictionary")
	}
	contentsDictionary := contentsDictionaryObj.AsDictionary()
	filterObj, hasFilter := contentsDictionary.GetElement("Filter")
	if !hasFilter {
		return nil, fmt.Errorf("read page contents: dictionary does not contain filter")
	}
	filter := filterObj.AsName()
	streamObj, hasStream := actualContents.Contains(ObjectType_Stream)
	if !hasStream {
		return nil, fmt.Errorf("read page contents: object does not contain stream")
	}
	stream := streamObj.AsStream()
	decoded, err := stream.Decode(filter.String())
	if err != nil {
		return nil, fmt.Errorf("read page contets: %w", err)
	}
	return &Page{pageNr: pageCount, decodedStream: decoded, lines: nil}, nil
}

func (r *DocumentReader) readObjectHeader() *Object {
	_id := new(bytes.Buffer)
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if IsWhiteSpace(char) {
			break
		}
		_id.WriteByte(char)
	}
	_gen := new(bytes.Buffer)
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if IsWhiteSpace(char) {
			break
		}
		_gen.WriteByte(char)
	}
	r.index += len(ObjectBegin)

	id, _ := strconv.Atoi(_id.String())
	gen, _ := strconv.Atoi(_gen.String())
	return New_Object_Indirect(id, gen)
}

func (r *DocumentReader) readObjectFromOffset(offset int) (*Object, error) {
	if offset >= len(r.buffer) {
		return nil, fmt.Errorf("reader.readObjectFromOffest: offset %d is out of range of buffer with size %d", offset, len(r.buffer))
	}
	//fmt.Printf("jumping to offset %d\n", offset)
	r.index = offset
	object := r.readObject()
	return object, nil
}

func (r *DocumentReader) readObject() *Object {
	object := r.readObjectHeader()

	for r.index < len(r.buffer) {
		token, err := r.readNextToken()
		if err != nil {
			fmt.Println("readObject:", err)
			return object
		}
		switch token {
		case ObjectEnd:
			return object
		case DictionaryBegin:
			object.AddIndirectObject(r.readDictionary())
		case ArrayBegin:
			object.AddIndirectObject(r.readArray())
		case NameDelimiter:
			object.AddIndirectObject(r.readName())
		case LiteralStringBegin:
			object.AddIndirectObject(r.readLiteralString())
		case Null:
			object.AddIndirectObject(New_Object_Null())
		case StreamBegin:
			object.AddStream(func(length int) *Object { return r.readStream(length) })
		}
	}

	return object
}

func (r *DocumentReader) readName() *Object {
	buffer := make([]byte, 0, 64)
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			fmt.Println("readName:", err)
			break
		}
		if IsWhiteSpace(char) || IsDelimiter(char) {
			r.unreadByte()
			break
		}
		buffer = append(buffer, char)
	}
	//fmt.Println(string(buffer))
	return New_Object_Name(Object_Name(buffer))
}

func (r *DocumentReader) readStream(length int) *Object {
	r.index += len(eol) // clear the EOL after 'stream'

	stream := make([]byte, length)
	copy(stream, r.buffer[r.index:r.index+length])
	r.index += length

	return New_Object_Stream(stream)
}

func (r *DocumentReader) readLiteralString() *Object {
	buffer := make([]byte, 0, 64)
	pOpenCount := 0
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			break
		}
		switch char {
		case '(':
			pOpenCount++
		case ')':
			pOpenCount--
		}
		if pOpenCount < 0 {
			break
		}
		buffer = append(buffer, char)
	}
	return New_Object_LiteralString(Object_LiteralString(buffer))
}

func (r *DocumentReader) readDictionary() *Object {
	dictionary := map[string]*Object{}
	for r.index < len(r.buffer) {
		token, err := r.readNextToken()
		if err != nil {
			fmt.Println("readDictionary:", err)
			return New_Object_Dictionary(dictionary)
		}
		switch token {
		case DictionaryEnd:
			return New_Object_Dictionary(dictionary)
		case NameDelimiter:
			key := r.readName()
			token, err := r.readNextToken()
			if err != nil {
				fmt.Println("readDictionary: value", err)
				return New_Object_Dictionary(dictionary)
			}
			switch token {
			case NameDelimiter:
				dictionary[key.String()] = r.readName()
			case ArrayBegin:
				dictionary[key.String()] = r.readArray()
			case DictionaryBegin:
				dictionary[key.String()] = r.readDictionary()
			case LiteralStringBegin:
				dictionary[key.String()] = r.readLiteralString()
			case BooleanTrue, BooleanFalse:
				dictionary[key.String()] = New_Object_Boolean(token == "true")
			case ObjectRef:
				dictionary[key.String()] = r.readObjectRef()
			case Null:
				dictionary[key.String()] = New_Object_Null()
			default: // numbers
				dictionary[key.String()] = r.readNumeric()
			}
		}
	}
	return New_Object_Dictionary(dictionary)
}

func (r *DocumentReader) readObjectRef() *Object {
	index := r.index - 2
	_gen := make([]byte, 0, 4)
	for index >= 0 {
		index--
		if !IsNumber(r.buffer[index]) {
			break
		}
		_gen = append(_gen, r.buffer[index])
	}
	_id := make([]byte, 0, 8)
	for index >= 0 {
		index--
		if !IsNumber(r.buffer[index]) {
			break
		}
		_id = append(_id, r.buffer[index])
	}
	slices.Reverse(_id)
	id, _ := strconv.Atoi(string(_id))
	slices.Reverse(_gen)
	gen, _ := strconv.Atoi(string(_gen))
	return New_Object_Ref(id, gen)
}

func (r *DocumentReader) readArray() *Object {
	array := make([]*Object, 0, 16)
	for r.index < len(r.buffer) {
		token, err := r.readNextToken()
		if err != nil {
			fmt.Println("readArray:", err)
			break
		}
		switch token {
		case ArrayEnd:
			return New_Object_Array(array)
		case NameDelimiter:
			array = append(array, r.readName())
		case LiteralStringBegin:
			array = append(array, r.readLiteralString())
		case HexadecimalStringBegin:
			array = append(array, r.readHexadecimalString())
		case ObjectRef:
			array = append(array, r.readObjectRef())
		case BooleanTrue, BooleanFalse:
			array = append(array, New_Object_Boolean(token == "true"))
		case Null:
			array = append(array, New_Object_Null())
		case DictionaryBegin:
			array = append(array, r.readDictionary())
		//case StreamBegin:
		//	array = append(array, dr.readStream())
		default: // numbers
			array = append(array, r.readNumeric())

		// These should NEVER happen!!!
		case ObjectEnd:
			panic(fmt.Errorf("tokenizer: unexpected 'endobj' at %d", r.index))
		case DictionaryEnd:
			panic(fmt.Errorf("tokenizer: unexpected '>>' at %d", r.index))
		case StreamEnd:
			panic(fmt.Errorf("tokenizer: unexpected 'endstream' at %d", r.index))
		}
	}
	return New_Object_Array(array)
}

func (r *DocumentReader) readHexadecimalString() *Object {
	buffer := make([]byte, 0, 128)
	for r.index < len(r.buffer) {
		token, err := r.readNextToken()
		if err != nil {
			fmt.Println("readHexadecimalString:", err)
			break
		}
		switch token {
		case HexadecimalStringEnd:
			return New_Object_HexadeciamalString(buffer)
		default:
			buffer = append(buffer, []byte(token)...)
		}
	}

	return New_Object_HexadeciamalString(buffer)
}

func (r *DocumentReader) readNumeric() *Object {
	buffer := make([]byte, 0, 16)
	r.unreadByte() // unread 1 since it has been read as a token already

	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			fmt.Println("readNumeric:", err)
			return New_Object_Null()
		}
		//fmt.Print(string(char))
		if IsNumerical(char) {
			r.unreadByte() // go back once to get char
			break
		}
	}
	for r.index < len(r.buffer) {
		char, err := r.readNextByte()
		if err != nil {
			fmt.Println("readNumeric:", err)
			return New_Object_Null()
		}
		//fmt.Print(string(char))

		if !IsNumerical(char) {
			r.unreadByte() // unread the non-numerical char
			break
		}
		buffer = append(buffer, char)
	}
	fistBuf, lastBuf, isReal := bytes.Cut(buffer, []byte("."))
	if isReal {
		// real
		fist, _ := strconv.ParseInt(string(fistBuf), 10, 32)
		last, _ := strconv.ParseInt(string(lastBuf), 10, 32)
		return New_Object_RealNumber(int(fist), int(last))
	}
	// integer
	integer, _ := strconv.ParseInt(string(buffer), 10, 32)
	return New_Object_IntegerNumber(int(integer))
}
