package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	cyoa "github.com/Beadko/gophercizes_cyoa"
)

var (
	port         = flag.Int("port", 8080, "The port to start the CYOA web application on")
	fileName     = flag.String("story", "gopher.json", "the JSON file with the CYOA story")
	templateFile = flag.String("template", "gopher.html", "the HTML template for CYOA story")
	cli          = flag.Bool("cli", false, "Run in CLI mode instead of web server")
	basePath     = "/story"
)

func main() {
	flag.Parse()
	fmt.Printf("Using the story %s.\n", *fileName)

	story, err := cyoa.LoadStory(*fileName)
	if err != nil {
		fmt.Printf("Could not parse the struct from %s: %v\n", *fileName, err)
		return
	}
	if *cli {
		fmt.Printf("Starting CLI version of Create Your Own Adventure")
		if err := cyoa.PlayStory(story); err != nil {
			fmt.Printf("Could not play a story '%s': %v", story, err)
		}

	} else {
		tmpl, err := template.ParseFiles(*templateFile)
		if err != nil {
			fmt.Printf("Could not parse the template file %s: %s.\n", *templateFile, err)
		}
		h := cyoa.NewHandler(story, tmpl, cyoa.WithPathFunc(cyoa.PathFn(basePath)))
		fmt.Printf("Starting the server on port: %d\n", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
	}
}
