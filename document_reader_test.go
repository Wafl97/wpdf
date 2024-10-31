package wpdf_test

import (
	"os"
	"testing"

	"github.com/Wafl97/wpdf"
)

func TestOpenPdf(t *testing.T) {
	pdf, err := wpdf.Open(".\\test_pdfs\\Microservice_Development_Tool.pdf")
	//pdf, err := wpdf.Open(".\\resources\\lm-info.pdf")
	//pdf, err := wpdf.Open(".\\resources\\Score card_læsevejledning_06.pdf")

	if err != nil {
		t.Fatal(err)
	}

	printPage := func(i int, p *wpdf.Page) {
		t.Log("page", i, string(p.DecodedStream()))
	}
	if err := pdf.IteratePages(printPage); err != nil {
		t.Fatal(err)
	}
}

func TestReadHeader(t *testing.T) {
	bytes, err := os.ReadFile(".\\test_pdfs\\Microservice_Development_Tool.pdf")
	//rawDictionary, err := os.ReadFile(".\\resources\\lm-info.pdf")
	//rawDictionary, err := os.ReadFile(".\\resources\\Score card_læsevejledning_06.pdf")
	if err != nil {
		t.Fatal(err)
	}

	reader := wpdf.NewDocumentReader(bytes)

	header, err := reader.ReadHeader()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(header.String())
}

func TestReadTrailer(t *testing.T) {
	bytes, err := os.ReadFile(".\\test_pdfs\\Microservice_Development_Tool.pdf")
	//bytes, err := os.ReadFile(".\\resources\\lm-info.pdf")
	//bytes, err := os.ReadFile(".\\resources\\Score card_læsevejledning_06.pdf")
	if err != nil {
		t.Error(err)
	}

	reader := wpdf.NewDocumentReader(bytes)

	trailer, err := reader.ReadTrailer()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(trailer.String())

	//bytes, err := os.ReadFile(".\\test_pdfs\\Microservice_Development_Tool.pdf")
	bytes, err = os.ReadFile(".\\resources\\lm-info.pdf")
	//bytes, err := os.ReadFile(".\\resources\\Score card_læsevejledning_06.pdf")
	if err != nil {
		t.Error(err)
	}

	reader = wpdf.NewDocumentReader(bytes)

	trailer, err = reader.ReadTrailer()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(trailer.String())
	//t.Log(trailer.Ref(1))

	//bytes, err := os.ReadFile(".\\test_pdfs\\Microservice_Development_Tool.pdf")
	//bytes, err := os.ReadFile(".\\resources\\lm-info.pdf")
	bytes, err = os.ReadFile(".\\resources\\Score card_læsevejledning_06.pdf")
	if err != nil {
		t.Error(err)
	}

	reader = wpdf.NewDocumentReader(bytes)

	trailer, err = reader.ReadTrailer()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(trailer.String())
}
