//go:build ignore

package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", ":8080", "Port to serve on")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("output")))
	log.Printf("Serving output/ on http://localhost%s", *port)
	log.Printf("  - Domain docs: http://localhost%s/docs/index.html", *port)
	log.Printf("  - Diff reports: http://localhost%s/tmp/", *port)
	log.Printf("Press Ctrl+C to stop.\n")
	log.Fatal(http.ListenAndServe(*port, nil))
}
