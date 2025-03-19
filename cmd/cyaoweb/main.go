package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

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
	tmpl := template.Must(template.New("").Parse(storyTmpl))
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tmpl), cyoa.WithPathFunc(pathFn))
	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	p := strings.TrimSpace(r.URL.Path)
	if p == "/story" || p == "/story/" {
		p = "/story/intro"
	}
	return p[len("/story/"):]
}

var storyTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}
            <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
	<style>
	        body {
            font-family: 'Georgia', serif;
            background-color: #121212;
            color: #f5f5f5;
            text-align: center;
            margin: 0;
            padding: 20px;
        }
        h1 {
            font-size: 2.5em;
            margin-bottom: 20px;
            text-shadow: 2px 2px 5px rgba(255, 255, 255, 0.2);
        }
        p {
            font-size: 1.2em;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto 20px;
            padding: 0 15px;
        }
        ul {
            list-style: none;
            padding: 0;
        }
        li {
            margin: 15px 0;
        }
        a {
            display: inline-block;
            background: #0077cc;
            color: white;
            padding: 12px 20px;
            border-radius: 5px;
            text-decoration: none;
            font-weight: bold;
            transition: background 0.3s ease, transform 0.2s;
        }
        a:hover {
            background: #0055aa;
            transform: scale(1.05);
        }
	</style>
</body>
</html>`
