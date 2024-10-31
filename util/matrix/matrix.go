package matrix

import (
	"fmt"
	"strings"
)

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float64 | float32
}

type Matrix[T Number] struct {
	m      int
	n      int
	values []T
}

func (m *Matrix[T]) String() string {
	sb := new(strings.Builder)
	for x := 0; x < m.m; x++ {
		for y := 0; y < m.n; y++ {
			value, _ := m.At(x, y)
			sb.WriteString(fmt.Sprintf("%v ", value))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *Matrix[T]) Rows() int {
	return m.n
}

func (m *Matrix[T]) Row(r int) ([]T, error) {
	row := make([]T, m.n)
	for x := 0; x < m.n; x++ {
		value, err := m.At(x, r)
		if err != nil {
			return nil, err
		}
		row[x] = value
	}
	return row, nil
}

func (m *Matrix[T]) calcIndex(x, y int) int {
	return x + y*m.m
}

func (m *Matrix[T]) At(x, y int) (T, error) {
	actualIndex := m.calcIndex(x, y)
	if actualIndex > len(m.values) {
		return 0, fmt.Errorf("matrix.SetAt: the indecies %d, %d are out of range of matrix with dimensions %d x %d", x, y, m.m, m.n)
	}
	return m.values[actualIndex], nil
}

func (m *Matrix[T]) SetAt(x, y int, value T) error {
	actualIndex := m.calcIndex(x, y)
	if actualIndex > len(m.values) {
		return fmt.Errorf("matrix.SetAt: the indecies %d, %d are out of range of matrix with dimensions %d x %d", x, y, m.m, m.n)
	}
	m.values[actualIndex] = value
	return nil
}

func (m *Matrix[T]) AddAt(x, y int, value T) error {
	actualIndex := m.calcIndex(x, y)
	if actualIndex > len(m.values) {
		return fmt.Errorf("matrix.SetAt: the indecies %d, %d are out of range of matrix with dimensions %d x %d", x, y, m.m, m.n)
	}
	m.values[actualIndex] += value
	return nil
}

func New[T Number](m, n int) *Matrix[T] {
	return &Matrix[T]{
		m:      m,
		n:      n,
		values: make([]T, n*m),
	}
}

func From[T Number](m, n int, values []T) (*Matrix[T], error) {
	if len(values) != m*n {
		return nil, fmt.Errorf("matix.From: size of values does not match n * m")
	}
	matrix := New[T](m, n)
	matrix.values = values
	return matrix, nil
}

func Multiply[T Number](a, b *Matrix[T]) (*Matrix[T], error) {
	if a.n != b.m {
		return nil, fmt.Errorf("matrix.Multiply: dimensions of a and b are not valid for multiplying")
	}
	m := a.m
	p := b.n
	n := b.m

	c := New[T](m, p)

	for x := 0; x < m; x++ {
		for y := 0; y < p; y++ {
			var prod T
			for i := 0; i < n; i++ {
				valueA, err := a.At(x, i)
				if err != nil {
					return nil, err
				}
				valueB, err := b.At(i, y)
				if err != nil {
					return nil, err
				}
				prod += valueA * valueB
			}
			if err := c.SetAt(x, y, prod); err != nil {
				return nil, fmt.Errorf("matrix.Multiply: %w", err)
			}
		}
	}
	return c, nil
}
