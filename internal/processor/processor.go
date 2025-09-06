package processor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go-doc-parser/internal/entity"
	"go-doc-parser/internal/generator"
	"go-doc-parser/internal/parser"
	"io"

	"github.com/fumiama/go-docx"
)

func NewProcessor(dictionary [][]entity.QualifiedName) func(files []*zip.File) (out entity.Data) {
	cases := generator.Plural("випадку", "випадках", "випадках")

	return func(files []*zip.File) (out entity.Data) {

		overall := map[entity.ID]int{}
		overallComments := map[entity.ID][]string{}

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

			p := Processor{
				Known:     map[entity.QualifiedName]int{},
				Unknown:   map[entity.ID]int{},
				Commented: map[entity.ID][]string{},
			}

			for _, group := range dictionary {
				for _, name := range group {
					p.Known[name] = 0
				}
			}

			p.Process(records)

			for id, comments := range p.Commented {
				for _, comment := range comments {
					if len(comment) > 0 {
						overallComments[id] = append(overallComments[id], comment)
					}
				}
			}

			page := entity.Page{
				Filename: file.FileHeader.Name,
			}

			for _, group := range dictionary {
				total := 0

				entries := []entity.Entry{}

				for _, name := range group {
					count := p.Known[name]

					entries = append(
						entries,
						entity.Entry{
							entity.ID{
								name, "",
							}, count,
						},
					)

					total += count
				}

				overall[entity.ID{group[0], ""}] += total

				page.KnownGroups = append(
					page.KnownGroups,
					entity.Group{
						Entries: entries,
						Total:   total,
					},
				)
			}

			for k, count := range p.Unknown {
				page.UnknownEntries = append(page.UnknownEntries, entity.Entry{k, count})
				overall[k] += count
			}

			out.Pages = append(out.Pages, page)
		}

		summary := ""

		for i, group := range dictionary {
			name := group[0]
			id := entity.ID{name, ""}
			total := overall[id]

			delete(overall, id)

			comment := "ОПДК не виявлено"

			times := len(overallComments[id])

			if times > 0 {
				comment = fmt.Sprintf("в %s M затриманих", cases(times))
			}

			summary += fmt.Sprintf("%d. %s - польотів: %d, %s;\n", i+1, group[0].Name, total, comment)

			out.Total = append(out.Total, entity.Entry{entity.ID{name, ""}, total})
		}

		for id, total := range overall {
			t := id.Type
			if len(t) < 1 {
				t = "інше"
			}

			out.TotalUnknown = append(out.TotalUnknown, entity.Entry{id, total})
		}

		for key, comments := range overallComments {
			out.Footnotes = append(out.Footnotes, entity.Footnote{key, comments})
		}

		out.Summary = summary

		return
	}
}

type Processor struct {
	Known     map[entity.QualifiedName]int
	Unknown   map[entity.ID]int
	Commented map[entity.ID][]string
}

func (p Processor) Process(records []entity.Record) {
	for _, record := range records {
		// filter records happened after 18:00???
		// if record.EndHour < 18 {
		// 	continue
		// }

		count, ok := p.Known[record.QualifiedName]
		if !ok {
			count := p.Unknown[record.ID]
			p.Unknown[record.ID] = count + 1
			p.Commented[record.ID] = append(p.Commented[record.ID], record.Comment)
			continue
		}

		if len(record.Comment) > 0 {
			id := entity.ID{record.QualifiedName, ""}
			p.Commented[id] = append(p.Commented[id], record.Comment)
		}

		p.Known[record.QualifiedName] = count + 1
	}
}
