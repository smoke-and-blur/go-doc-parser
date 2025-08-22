package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/fumiama/go-docx"
)

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

type Processor struct {
	Known     map[QualifiedName]int
	Unknown   map[ID]int
	Commented map[ID]string
}

type QualifiedName struct {
	Type string
	Name string
}

type Parser struct {
	In       []rune
	Position int
}

// narrow down options character by character
func (p *Parser) matchChar(options [][]rune, i int) [][]rune {
	if len(options) == 1 && len(options[0]) == i { // reached the end of the word?
		p.Position += i
		return options
	}

	out := [][]rune{}

	for _, option := range options {
		if i >= len(option) {
			continue
		}

		if p.In[p.Position+i] == option[i] {
			out = append(out, option)
		}
	}

	if len(out) < 1 {
		return nil
	}

	i++

	return p.matchChar(out, i)
}

func (p *Parser) MatchChar(options ...string) string {
	runeOptions := [][]rune{}

	for _, option := range options {
		runeOptions = append(runeOptions, []rune(option))
	}

	out := p.matchChar(runeOptions, 0)

	if len(out) > 0 {
		return string(out[0])
	}

	return ""
}

func (p *Parser) SkipSpace() {
	if unicode.IsSpace(rune(p.In[p.Position])) {
		p.Position++
		p.SkipSpace()
	}
}

func (p *Parser) trimQuotes() (out []rune) {
	r := p.In[p.Position]

	// fmt.Printf("%c %t\n", r, isQuotation)

	if unicode.Is(unicode.Quotation_Mark, r) {
		p.Position++
	}

	out = p.In[p.Position:len(p.In)]

	r = out[len(out)-1]

	if unicode.Is(unicode.Quotation_Mark, r) {
		out = out[:len(out)-1]
	}

	p.Position = len(p.In)

	return out
}

func (p *Parser) ParseName() QualifiedName {
	t := p.MatchChar("віпс", "впс")

	// fmt.Printf("%s %d %c\n", t, p.Position, p.In[p.Position])

	p.SkipSpace()
	// fmt.Printf("%d %c\n", p.Position, p.In[p.Position])

	// n := p.filterSymbols([]rune{})

	name := p.trimQuotes()

	// fmt.Printf("%s %d %c\n", string(n), p.Position, p.In[p.Position-1])

	return QualifiedName{t, string(name)}
}

func (p Processor) Process(records []Record, filterEarly bool) {
	for _, record := range records {
		// filter records happened after 18:00???
		if filterEarly && record.EndHour < 18 {
			continue
		}

		count, ok := p.Known[record.QualifiedName]
		if !ok {
			count := p.Unknown[record.ID]
			p.Unknown[record.ID] = count + 1
			p.Commented[record.ID] = record.Comment
			continue
		}

		if len(record.Comment) > 0 {
			id := ID{record.QualifiedName, ""}
			p.Commented[id] = record.Comment
		}

		p.Known[record.QualifiedName] = count + 1
	}
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

func ParseTable(table *docx.Table) (out []Record) {
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

		item := row.TableCells[6].Paragraphs[0].String()

		hint := ""

		if len(row.TableCells[6].Paragraphs) > 1 {
			hint = row.TableCells[6].Paragraphs[1].String()
		}

		p := Parser{
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

		record := Record{
			ID: ID{
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

func main() {
	// dictionary sets up the right order for the names
	dictionary := [][]QualifiedName{}

	data := os.Getenv("DATA")

	err := json.Unmarshal([]byte(data), &dictionary)
	if err != nil {
		fmt.Println("failed to unmarshal the data:", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	process := NewProcessor(dictionary)

	http.ListenAndServe(
		"0.0.0.0:"+port,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				r.ParseMultipartForm(64_000_000)

				if r.MultipartForm == nil {
					return
				}

				descriptors := []Descriptor{}

				// TODO: consider making it in parallel and sorting back again?

				for name, headers := range r.MultipartForm.File {
					filterEarly := name == "filtered_files"

					for _, header := range headers {
						descriptors = append(descriptors, Descriptor{header, filterEarly})
					}
				}

				process(w, descriptors)
			},
		),
	)
}

type Descriptor struct {
	Header      *multipart.FileHeader
	FitlerEarly bool
}

func NewProcessor(dictionary [][]QualifiedName) func(w io.Writer, descriptors []Descriptor) {
	return func(w io.Writer, descriptors []Descriptor) {

		overall := map[ID]int{}
		overallComments := map[ID][]string{}

		for _, descriptor := range descriptors {
			file, err := descriptor.Header.Open()
			if err != nil {
				return //
			}
			defer file.Close()

			doc, err := docx.Parse(file, descriptor.Header.Size)
			if err != nil {
				panic(err)
				// TODO: ???
			}

			table := FindFirstTable(doc)

			records := ParseTable(table)

			p := Processor{
				Known:     map[QualifiedName]int{},
				Unknown:   map[ID]int{},
				Commented: map[ID]string{},
			}

			for _, group := range dictionary {
				for _, name := range group {
					p.Known[name] = 0
				}
			}

			p.Process(records, descriptor.FitlerEarly)

			fmt.Fprintf(w, "%s\n18:00-00:00: %t\n\n", descriptor.Header.Filename, descriptor.FitlerEarly)

			for _, group := range dictionary {
				total := 0

				for _, name := range group {
					count := p.Known[name]
					fmt.Fprintf(w, "%s | %s | %d\n", name.Type, name.Name, count)
					total += count

					comment := p.Commented[ID{name, ""}]
					if len(comment) > 0 {
						id := ID{group[0], ""}
						overallComments[id] = append(overallComments[id], comment)
					}
				}

				overall[ID{group[0], ""}] += total

				fmt.Fprintf(w, "%d\n", total)
				fmt.Fprintln(w)

			}

			for k, count := range p.Unknown {
				fmt.Fprintf(w, "%s | %s | %s | %d\n", "інше", k.Name, k.Hint, count)
				overall[k] += count
			}

			if len(p.Unknown) > 0 {
				fmt.Fprintf(w, "\n")
			}

		}

		fmt.Fprintf(w, "...\n\nTOTAL\n\n")

		completeOutput := ""

		for i, group := range dictionary {
			name := group[0]
			id := ID{name, ""}
			total := overall[id]

			delete(overall, id)

			comment := "ОПДК не виявлено"

			times := len(overallComments[id])

			if times > 0 {
				comment = fmt.Sprintf("в %d випадках M затриманих", times)
			}

			completeOutput += fmt.Sprintf("%d. %s - польотів: %d, %s;\n", i+1, group[0].Name, total, comment)

			fmt.Fprintf(w, "%s | %s | %d\n", name.Type, name.Name, total)
		}

		fmt.Fprintln(w)

		for id, total := range overall {
			fmt.Fprintf(w, "%s | %s | %s | %d\n", "інше", id.Name, id.Hint, total)
		}

		if len(overall) > 0 {
			fmt.Fprintln(w)
		}

		for key, comments := range overallComments {
			fmt.Fprintf(w, "%s | %s | %s\n", key.Type, key.Name, key.Hint)
			for _, comment := range comments {
				fmt.Fprintf(w, "%s\n", comment)
			}
			fmt.Fprintln(w)
		}

		fmt.Fprintf(w, "%s", completeOutput)
	}
}
