package output

import "fmt"

type OutputStyle interface {
	Render(name, words string)
}

type Text struct{}

func (t Text) Render(name, words string) {
	fmt.Printf("\033[1m%s\033[0m: %s\n\n", name, words)
}
