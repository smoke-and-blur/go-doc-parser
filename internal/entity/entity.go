package entity

type ID struct {
	ShortID
	Hint string
}

type ShortID struct {
	Type string
	Name string
}

type Event struct {
	Start   uint64 // hour?
	End     uint64 // hour?
	Comment string
}

type Record struct {
	Filename string
	ID
	Event
}
