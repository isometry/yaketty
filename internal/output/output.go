package output

import "fmt"

type OutputStyle interface {
	Render(name, words string)
}

type Text struct{}

func (t Text) Render(name, words string) {
	fmt.Printf("\033[1m%s\033[0m: %s\n\n", name, words)
}

type JSON struct{}

func (j JSON) Render(name, words string) {
	fmt.Printf(`{"name": %q, "words": %q}`+"\n", name, words)
}
