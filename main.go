package main

import (
	"go-doc-parser/internal/handler"
	"go-doc-parser/internal/processor"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	process := processor.Process

	http.ListenAndServe(
		"0.0.0.0:"+port,
		http.HandlerFunc(
			handler.Handler(process),
		),
	)
}
