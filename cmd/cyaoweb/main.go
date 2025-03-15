package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	cyao "github.com/Beadko/gophercizes_cyoa"
)

func main() {
	fileName := flag.String("file", "gopher.json", "the JSON file with the CYOA stroy")
	flag.Parse()
	fmt.Printf("Using the story %s.\n", *fileName)

	f, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("Could not open file %s: %v/n", *fileName, err)
		return
	}

	d := json.NewDecoder(f)

	var story cyao.Story
	if err := d.Decode(&story); err != nil {
		fmt.Printf("Could not parse the struct from %s: %v\n", *fileName, err)
		return
	}

	fmt.Printf("%+v\n", story)
}
