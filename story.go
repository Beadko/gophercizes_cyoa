package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tmpl *template.Template

var defaultHandlerTmpl = `
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
            <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>
</html>`

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)

	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, h.s["intro"]); err != nil {
		fmt.Printf("Error executing story %s: %v", h.s["intro"], err)
	}

}

func NewHandler(s Story) http.Handler {
	return handler{s}
}
