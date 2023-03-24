package types

import "github.com/gogf/gf/v2/util/gconv"

// NullUint64 represents an uint64 that may be null.
// NullUint64 implements the Scanner interface, so
// it can be used as a scan destination, similar to NullString.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool // Valid is true if Uint64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullUint64) Scan(value any) error {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	n.Uint64 = gconv.Uint64(value)
	return nil
}
