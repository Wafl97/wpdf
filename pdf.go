package wpdf

import (
	"fmt"
	"os"
	//. "github.com/Wafl97/wpdf/errors"
)

type PDF struct {
	//reader    *DocumentReader
	Header    Header
	Trailer   Trailer
	pages     []*Page
	pageCount int
}

func Open(filename string) (*PDF, error) {
	// Read file bytes
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	pdf := new(PDF)

	// Read header
	reader := NewDocumentReader(bytes)
	header, err := reader.ReadHeader()
	if err != nil {
		return nil, err
	}
	pdf.Header = *header

	// Read trailer
	trailer, err := reader.ReadTrailer()
	if err != nil {
		return nil, err
	}
	pdf.Trailer = *trailer

	// Parse document objects
	pdf.pageCount, err = reader.ReadPageCount()
	if err != nil {
		return nil, err
	}
	pdf.pages = make([]*Page, 0, pdf.pageCount)

	reader.BuildPageTree(&pdf.pages)

	return pdf, nil
}

func (pdf *PDF) PageCount() int {
	return pdf.pageCount
}

func (pdf *PDF) Page(pageNr int) (*Page, error) {
	// TODO: dont know if it should stay as 1 indexed
	pageNrIdx := pageNr - 1 // make it 1 based indexing for users
	if pageNrIdx < 0 || pageNrIdx > pdf.pageCount {
		return nil, fmt.Errorf("wpdf: cannot return page %d out of %d pages", pageNr, pdf.pageCount)
	}
	if pdf.pages[pageNrIdx] == nil {
		return nil, fmt.Errorf("wpdf: page %d chould not be found", pageNr)
	}
	if pdf.pages[pageNrIdx].decodedStream == nil {
		return nil, fmt.Errorf("wpdf: page %d has no contents", pageNr)
		// find page and read stream and decode it
	}
	return pdf.pages[pageNrIdx], nil
}

func (pdf *PDF) IteratePages(iter func(index int, p *Page)) error {
	for pageNrIdx := 1; pageNrIdx < len(pdf.pages); pageNrIdx++ {
		page, err := pdf.Page(pageNrIdx)
		if err != nil {
			return err
		}
		iter(pageNrIdx, page)
	}
	return nil
}
