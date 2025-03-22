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

func LoadStory(fileName string) (Story, string, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, "", err
	}

	var story Story
	if err := json.Unmarshal(fileContent, &story); err != nil {
		return nil, "", err
	}

	start := FindStartingChapter(story)
	if start == "" {
		return nil, "", fmt.Errorf("error: story has no chapters")
	}

	return story, start, nil
}

func FindStartingChapter(story Story) string {
	if _, exists := story["intro"]; exists {
		return "intro"
	}

	linkedCh := make(map[string]bool)
	for _, ch := range story {
		for _, option := range ch.Options {
			linkedCh[option.Chapter] = true
		}
	}
	var startChs []string
	for chName := range story {
		if !linkedCh[chName] {
			startChs = append(startChs, chName)
		}
	}

	if len(startChs) == 1 {
		return startChs[0]
	}

	return ""
}

func PlayStory(story Story, start string) error {
	sc := bufio.NewScanner(os.Stdin)

	for {
		loadChapters(story, sc, start)
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

func loadChapters(story Story, sc *bufio.Scanner, start string) error {
	currentChapter := start

	for {
		ch, exists := story[currentChapter]
		if !exists {
			return fmt.Errorf("error: chapter '%s' not found", start)
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
			currentChapter = ch.Options[choice-1].Chapter
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

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func PathFn(base string, ch string) func(r *http.Request) string {
	return func(r *http.Request) string {
		p := strings.TrimSpace(r.URL.Path)
		p = strings.TrimPrefix(p, base)
		p = strings.Trim(p, "/")

		if p == "" {
			p = ch
		}

		return p
	}
}

func NewHandler(s Story, t *template.Template, ch string, opts ...HandlerOption) http.Handler {
	h := handler{s, t, PathFn("/", ch)}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := h.pathFn(r)

	if ch, ok := h.s[p]; ok {
		if err := h.t.Execute(w, ch); err != nil {
			log.Printf("Error executing story %s: %v", p, err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
