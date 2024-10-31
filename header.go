package wpdf

import "fmt"

type Header struct {
	version Version
}

func (h *Header) String() string {
	return fmt.Sprintf("Header:\n\tVersion: %s\n", h.version.String())
}

func (h *Header) GetVersion() Version {
	return h.version
}
