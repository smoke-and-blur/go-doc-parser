package entity

import "mime/multipart"

type Descriptor struct {
	Header      *multipart.FileHeader
	FitlerEarly bool
}

type Entry struct {
	ID
	Count int
}

type Group struct {
	Entries []Entry
	Total   int
}

type Page struct {
	Filename       string
	KnownGroups    []Group
	UnknownEntries []Entry
}

type Footnote struct {
	ID
	Comments []string
}

type Data struct {
	Records      []Record
	Pages        []Page
	Total        []Entry
	TotalUnknown []Entry
	Footnotes    []Footnote
	Summary      string
}

type ID struct {
	QualifiedName
	Hint string
}

type Record struct {
	// N
	// Problem
	ID
	EndHour uint64
	Comment string
}

type QualifiedName struct {
	Type string
	Name string
}
