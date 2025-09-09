package entity

type Group struct {
	ID
	Events []Event
}

type Page struct {
	Filename            string
	SelectedSupergroups [][]Group
	OtherGroups         []Group
}

type Data struct {
	Pages              []Page
	AggregatedSelected []Group
	AggregatedOther    []Group
	AggregatedComments []Group
	Summary            string
}

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
	ID
	Event
}
