package wpdf_errors

import (
	"errors"
)

var (
	ErrNotYetImplemented = errors.New("not yet implemented")
	ErrInvalidHeader     = errors.New("invalid pdf header")
	ErrInvlidObject      = errors.New("invlid pdf object")
)
