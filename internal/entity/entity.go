package entity

type EventGroup struct {
	ID
	Events []Event
}

type CommentGroup struct {
	ID
	Comments []string
}

type Page struct {
	Filename            string
	SelectedSupergroups [][]EventGroup
	OtherGroups         []EventGroup
}

type Data struct {
	Pages              []Page
	AggregatedSelected []EventGroup
	AggregatedOther    []EventGroup
	AggregatedComments []CommentGroup
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
	Start uint64 // hour?
	End   uint64 // hour?
}

type Record struct {
	ID
	Event
	Comment string
}
