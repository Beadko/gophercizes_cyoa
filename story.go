package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

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

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	p := strings.TrimSpace(r.URL.Path)
	if p == "" || p == "/" {
		p = "/intro"
	}
	return p[1:]
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tmpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := h.pathFn(r)

	if ch, ok := h.s[p]; ok {
		if err := h.t.Execute(w, ch); err != nil {
			log.Printf("Error executing story %s: %v", h.s["intro"], err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
