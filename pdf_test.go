package wpdf_test

import (
	"testing"

	"github.com/Wafl97/wpdf"
)

func TestPdf(t *testing.T) {
	a, err := wpdf.Open(".\\resources\\lm-info.pdf")
	if err != nil {
		t.Fatal("PDF a =", err)
	}
	t.Logf("a.Version.String(): %v\n", a.Header.GetVersion().String())
	t.Log(a.Header.GetVersion().String(), "is at most", wpdf.PDFv1_2.String(), a.Header.GetVersion().IsAtMost(wpdf.PDFv1_2))

	b, err := wpdf.Open(".\\resources\\Score card_læsevejledning_06.pdf")
	if err != nil {
		t.Log("PDF b =", err)
	}
	t.Logf("b.Version.String(): %v\n", b.Header.GetVersion().String())
	t.Log(b.Header.GetVersion().String(), "is at least", wpdf.PDFv1_2.String(), b.Header.GetVersion().IsAtLeast(wpdf.PDFv1_2))

	t.Log(a.Header.GetVersion().Is(wpdf.PDFv1_4))
	t.Log(b.Header.GetVersion().Is(wpdf.PDFv1_7))
}
