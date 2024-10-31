package matrix_test

import (
	"github.com/Wafl97/wpdf/util/matrix"
	"testing"
)

func TestMatrixMultiply(t *testing.T) {
	a, err := matrix.From(3, 3, []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < a.Rows(); y++ {
		row, err := a.Row(y)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(row)
	}
	t.Log("by")
	b, err := matrix.From(3, 2, []uint8{5, 4, 3, 2, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < b.Rows(); y++ {
		row, err := b.Row(y)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(row)
	}
	t.Log("=")
	c, err := matrix.Multiply(a, b)
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < c.Rows(); y++ {
		row, err := c.Row(y)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(row)
	}
}
