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

		end, err := strconv.ParseUint(endTime[:2], 10, 64)
		if err != nil {
			// report problem
			continue
		}

		parts := []string{}

		for _, p := range row.TableCells[6].Paragraphs {
			parts = append(parts, strings.Join(strings.Fields(p.String()), " "))
		}

		item := parts[0]

		hint := ""

		if len(parts) > 1 {
			hint = strings.Join(parts[1:], " ")
		}

		p := NameParser{
			In: []rune(item),
		}

		q := p.ParseName()

		paragraphs := []string{}

		for _, p := range row.TableCells[8].Paragraphs {

			p := p.String()

			p = strings.Join(strings.Fields(p), " ")

			if p == "ОПДК не виявлено" {
				continue
			}

			if len(p) > 0 {
				paragraphs = append(paragraphs, p)
			}
		}

		// if len(paragraphs) < 1 {
		// 	continue
		// }

		record := entity.Record{
			ID: entity.ID{
				q,
				hint,
			},
			EndHour: end,
			Comment: strings.Join(paragraphs, "\n"),
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
