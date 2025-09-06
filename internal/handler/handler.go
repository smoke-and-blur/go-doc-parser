package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go-doc-parser/internal/entity"
	"html/template"
	"io"
	"net/http"
)

func Handler(process func([]*zip.File) entity.Data) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// for k, v := range r.Header {
		// 	for _, v := range v {
		// 		fmt.Println(k, v)
		// 	}
		// }

		if r.Method != http.MethodPost {
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to read all:", err)
			return
		}

		reader, err := zip.NewReader(bytes.NewReader(body), r.ContentLength)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data := process(reader.File)

		tpl, err := template.ParseFiles("template_new.gohtml")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to parse the templates:", err)
			return
		}

		tpl.Execute(w, data)
	}
}
