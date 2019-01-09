package kernel

type Id string

// CommandType type
type CommandType int

const (
	UpdateCard CommandType = CommandType(iota)
	RemoveCard
	ExcludeChild
	InsertBefore
	InsertAfter
	AppendChild
	LayoutBoard
)
