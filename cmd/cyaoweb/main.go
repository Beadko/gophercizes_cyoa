package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	cyoa "github.com/Beadko/gophercizes_cyoa"
)

func main() {
	port := flag.Int("port", 3000, "The port to start the CYOA web application on")
	fileName := flag.String("file", "gopher.json", "the JSON file with the CYOA stroy")
	flag.Parse()
	fmt.Printf("Using the story %s.\n", *fileName)

	f, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("Could not open file %s: %v/n", *fileName, err)
		return
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		fmt.Printf("Could not parse the struct from %s: %v\n", *fileName, err)
		return
	}
	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
