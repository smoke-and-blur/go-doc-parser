package processor

import (
	"archive/zip"
	"bytes"
	. "go-doc-parser/internal/entity"
	"go-doc-parser/internal/parser"
	"io"

	"github.com/fumiama/go-docx"
)

func Process(files []*zip.File) (out []Record) {
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
		// put filename into each record lol
		records := parser.ParseTable(file.Name, table)

		out = append(out, records...)
	}
	return
}
