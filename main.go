package main

import (
	"encoding/json"
	"fmt"
	"go-doc-parser/internal/entity"
	"go-doc-parser/internal/handler"
	"go-doc-parser/internal/processor"
	"net/http"
	"os"
)

func main() {
	// dictionary sets up the right order for the names
	dictionary := [][]entity.ID{}

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

	process := processor.NewProcessor(dictionary)

	http.ListenAndServe(
		"0.0.0.0:"+port,
		http.HandlerFunc(
			handler.Handler(process),
		),
	)
}
