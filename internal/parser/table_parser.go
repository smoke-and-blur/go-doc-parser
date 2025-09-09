package parser

import (
	"go-doc-parser/internal/entity"
	"strconv"
	"strings"

	"github.com/fumiama/go-docx"
)

func ParseTable(tag string, table *docx.Table) (out []entity.Record) {
	for _, row := range table.TableRows[1:] {
		endTime := row.TableCells[5].Paragraphs[0].String()
		if len(endTime) < 2 {
			continue
		}

		// TODO: start time, minutes?

		end, err := strconv.ParseUint(endTime[:2], 10, 64)
		if err != nil {
			// report problem
			continue
		}

		parts := []string{}

		for _, paragraph := range row.TableCells[6].Paragraphs {
			parts = append(parts, strings.Join(strings.Fields(paragraph.String()), " "))
		}

		item := parts[0]

		hint := ""

		if len(parts) > 1 {
			hint = strings.Join(parts[1:], " ")
		}

		parser := NameParser{
			In: []rune(item),
		}

		shortID := parser.ParseName()

		paragraphs := []string{}

		for _, p := range row.TableCells[8].Paragraphs {

			paragraph := p.String()

			paragraph = strings.Join(strings.Fields(paragraph), " ")

			// ???
			if paragraph == "ОПДК не виявлено" {
				continue
			}

			if len(paragraph) > 0 {
				paragraphs = append(paragraphs, paragraph)
			}
		}

		record := entity.Record{
			ID: entity.ID{
				shortID,
				hint,
			},
			Event: entity.Event{
				Start:   0,
				End:     end,
				Comment: strings.Join(paragraphs, "\n"),
			},
		}

		out = append(out, record)
	}

	return
}

func FindFirstTable(doc *docx.Docx) *docx.Table {
	for _, it := range doc.Document.Body.Items {
		switch it := it.(type) {
		case *docx.Table:
			return it
		}
	}

	return nil
}
