package wpdf

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/Wafl97/wpdf/errors"
)

// Version is used to represent the version of a given PDF.
// A version consist of both major and minor.
type Version struct {
	Major int
	Minor int
}

var (
	PDFv1_0 = Version{Major: 1, Minor: 0} //%PDF-1.0
	PDFv1_1 = Version{Major: 1, Minor: 1} //%PDF-1.1
	PDFv1_2 = Version{Major: 1, Minor: 2} //%PDF-1.2
	PDFv1_3 = Version{Major: 1, Minor: 3} //%PDF-1.3
	PDFv1_4 = Version{Major: 1, Minor: 4} //%PDF-1.4
	PDFv1_5 = Version{Major: 1, Minor: 5} //%PDF-1.5
	PDFv1_6 = Version{Major: 1, Minor: 6} //%PDF-1.6
	PDFv1_7 = Version{Major: 1, Minor: 7} //%PDF-1.7
)

func VersionFromString(str string) (*Version, error) {
	vStr, isValid := strings.CutPrefix(str, "%PDF-")
	if !isValid {
		return nil, fmt.Errorf("%w: not a pdf", ErrInvalidHeader)
	}
	vArr := strings.Split(vStr, ".")
	if len(vArr) != 2 {
		return nil, fmt.Errorf("%w: invalid version", ErrInvalidHeader)
	}
	major, err := strconv.Atoi(vArr[0])
	if err != nil {
		return nil, fmt.Errorf("%w: missing major version: %w", ErrInvalidHeader, err)
	}
	minor, err := strconv.Atoi(vArr[1])
	if err != nil {
		return nil, fmt.Errorf("%w: missing minor version: %w", ErrInvalidHeader, err)
	}
	return &Version{Major: major, Minor: minor}, nil
}

// String returns the string representation of v.
//
// # Examples:
//
//	PDFv1_2.String() : "PDF-1.2"
func (v Version) String() string {
	return fmt.Sprintf("PDF-%d.%d", v.Major, v.Minor)
}

// Is checks if v == other.
//
// # examples:
//
//	v1.3 == v1.2 : false
//	v1.3 == v1.3 : true
//	v1.3 == v1.4 : false
func (v Version) Is(other Version) bool {
	return v.Major == other.Major && v.Minor == other.Minor
}

// IsAtMost checks if v <= other.
//
// # examples:
//
//	v1.3 <= v1.2 : false
//	v1.3 <= v1.3 : true
//	v1.3 <= v1.4 : true
func (v Version) IsAtMost(other Version) bool {
	return v.Major <= other.Major && v.Minor <= other.Minor
}

// IsAtLeast checks if v >= other.
//
// # examples:
//
//	v1.3 >= v1.2 : true
//	v1.3 >= v1.3 : true
//	v1.3 >= v1.4 : false
func (v Version) IsAtLeast(other Version) bool {
	return v.Major >= other.Major && v.Minor >= other.Minor
}
