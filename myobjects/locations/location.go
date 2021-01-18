package locations

import "fmt"

//Location ...
type Location struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

//New ...
func New(row int, col int) *Location {
	return &Location{
		Row: row,
		Col: col,
	}
}

func (l *Location) String() string {
	return fmt.Sprintf("(%d,%d)", l.Row, l.Col)
}
