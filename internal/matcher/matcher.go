package matcher

type Position struct {
	Start int
	End   int
}

type Match struct {
	Line      int
	Positions []Position
	Data      []byte
}
