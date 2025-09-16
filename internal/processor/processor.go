package processor

import (
	"archive/zip"
	"bytes"
	. "go-doc-parser/internal/entity"
	"go-doc-parser/internal/parser"
	"io"

	"github.com/fumiama/go-docx"
)

func NewProcessor(dictionary [][]ID) func(files []*zip.File) (out Data) {
	return func(files []*zip.File) (out Data) {

		aggregateComments := map[ID][]Event{}

		for _, file := range files {
			opened, _ := file.Open()

			reader, _ := io.ReadAll(opened)

			doc, err := docx.Parse(bytes.NewReader(reader), int64(file.FileHeader.UncompressedSize64))
			if err != nil {
				panic(err)
				// TODO: ???
			}

			// inject
			table := parser.FindFirstTable(doc)

			// inject
			records := parser.ParseTable(file.Name, table)

			p := Collector{
				EventsBySelectedIDs: map[ShortID][]Event{},
				EventsByOtherIDs:    map[ID][]Event{},
				CommentsByIDs:       map[ID][]Event{},
			}

			for _, group := range dictionary {
				for _, id := range group {
					p.EventsBySelectedIDs[id.ShortID] = []Event{} // empty array to make sure the key exists
				}
			}

			p.Collect(records)

			for id, comments := range p.CommentsByIDs {
				aggregateComments[id] = append(aggregateComments[id], comments...)
			}

			selectedSupergroups := [][]Group{}

			for _, groupIDs := range dictionary {
				groups := []Group{}
				for _, groupID := range groupIDs {
					events := p.EventsBySelectedIDs[groupID.ShortID]
					groups = append(groups, Group{ID: groupID, Events: events})
				}

				selectedSupergroups = append(selectedSupergroups, groups)
			}

			otherGroups := []Group{}

			for id, group := range p.EventsByOtherIDs {
				otherGroups = append(otherGroups, Group{id, group})
			}

			page := Page{
				Filename:            file.FileHeader.Name,
				SelectedSupergroups: selectedSupergroups,
				OtherGroups:         otherGroups,
			}

			out.Pages = append(out.Pages, page)
		}

		summary := ""

		for id, comments := range aggregateComments {
			out.AggregatedComments = append(out.AggregatedComments, Group{id, comments})
		}

		out.Summary = summary

		return
	}
}

// collect all events selected and other plus comments and group them by id
type Collector struct {
	EventsBySelectedIDs map[ShortID][]Event // need to fill in empty items for all selected ids before using
	EventsByOtherIDs    map[ID][]Event
	CommentsByIDs       map[ID][]Event
}

func (p Collector) Collect(records []Record) {
	for _, record := range records {
		if len(record.Comment) > 0 {
			p.CommentsByIDs[record.ID] = append(p.CommentsByIDs[record.ID], record.Event)
		}

		_, ok := p.EventsBySelectedIDs[record.ShortID] // ignore hint within ID for selected events
		if !ok {
			// not found, so register it as other event
			p.EventsByOtherIDs[record.ID] = append(p.EventsByOtherIDs[record.ID], record.Event) // use the full ID for this
			continue
		}

		p.EventsBySelectedIDs[record.ShortID] = append(p.EventsBySelectedIDs[record.ShortID], record.Event)
	}
}
