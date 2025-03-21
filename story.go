package cyoa

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func LoadStory(fileName string) (Story, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(f)

	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func PlayStory(story Story) error {
	sc := bufio.NewScanner(os.Stdin)

	for {
		playChapters(story, sc)
		fmt.Println("\nWould you like to play again?")
		fmt.Println("1. Yes, let's do this!")
		fmt.Println("2. No, I'm bored now")

		fmt.Print("\nEnter your choice (1-2): ")
		if !sc.Scan() {
			return fmt.Errorf("Error reading input.")
		}

		choice := sc.Text()
		if choice != "1" {
			fmt.Println("Thanks for playing!")
			return nil
		}
	}
}

func playChapters(story Story, sc *bufio.Scanner) error {
	start := "intro"

	for {
		ch, exists := story[start]
		if !exists {
			return fmt.Errorf("Error: Chapter '%s' is not found", start)
		}
		fmt.Printf("\n%s\n", strings.ToUpper(ch.Title))
		fmt.Println(strings.Repeat("-", len(ch.Title)))

		for _, paragraph := range ch.Paragraphs {
			fmt.Println(paragraph)
			fmt.Println()
		}

		if len(ch.Options) == 0 {
			fmt.Println("THE END")
			return nil
		}
		fmt.Println("What would you like to do?")
		for i, option := range ch.Options {
			fmt.Printf("%d. %s\n", i+1, option.Text)
		}

		for {
			fmt.Print("\nEnter your choice (1-" + strconv.Itoa(len(ch.Options)) + "): ")
			if !sc.Scan() {
				fmt.Println("Error reading input.")
				return nil
			}
			input := sc.Text()
			choice, err := strconv.Atoi(input)
			if err != nil || choice < 1 || choice > len(ch.Options) {
				fmt.Printf("Please enter a number between 1 and %d.\n", len(ch.Options))
				continue
			}
			start = ch.Options[choice-1].Chapter
			break
		}
	}
}

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

func PathFn(base string) func(r *http.Request) string {
	return func(r *http.Request) string {
		p := strings.TrimSpace(r.URL.Path)
		p = strings.TrimPrefix(p, base)
		p = strings.Trim(p, "/")

		if p == "" {
			p = "intro"
		}

		return p
	}
}

func NewHandler(s Story, t *template.Template, opts ...HandlerOption) http.Handler {
	h := handler{s, t, PathFn("/")}
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
