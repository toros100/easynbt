package pointer

import "github.com/toros100/easynbt/nbt"

//go:generate go run ../../. -types=List

// yes, this is stupid
type List struct {
	Head int8
	Tail nbt.Option[*List]
}
