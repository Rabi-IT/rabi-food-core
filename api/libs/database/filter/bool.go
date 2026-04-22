package filter

// Bool is an optional boolean filter for queries.
// Zero value (BoolEmpty) means no filter is applied.
type Bool int8

const (
	BoolEmpty Bool = 0
	True      Bool = 1
	False     Bool = -1
)

func (f Bool) IsEmpty() bool { return f == 0 }
func (f Bool) Value() bool   { return f == 1 }
