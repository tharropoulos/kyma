package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
)

func TestRenderer_BasicRender(t *testing.T) {
	r, err := NewRenderer("dark")
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	got, gotErr := r.Render(
		`
Glamour
=======

A casual introduction. 你好世界!


## Let’s talk about artichokes

The _artichoke_ is mentioned as a garden plant in the 8th century BC by Homer
**and** Hesiod. The naturally occurring variant of the artichoke, the cardoon,
which is native to the Mediterranean area, also has records of use as a food
among the ancient Greeks and Romans. Pliny the Elder mentioned growing of
_carduus_ in Carthage and Cordoba.

> He holds him with a skinny hand,
> ‘There was a ship,’ quoth he.
> ‘Hold off! unhand me, grey-beard loon!’
> An artichoke, dropt he.

--Samuel Taylor Coleridge, [The Rime of the Ancient Mariner][rime]

[rime]: https://poetryfoundation.org/poems/43997/
		`, false, 0, 0)
	if gotErr != nil {
		t.Errorf("Render() failed: %v", gotErr)
	}

	golden.RequireEqual(t, got)
}

func TestRenderer_RenderCodeBlockWithHighlightedLines(t *testing.T) {
	r, err := NewRenderer("dark")
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	codeBlock := `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
	`

	got, gotErr := r.Render(fmt.Sprintf("# Slide\n```go{3,6}\n%s\n```", codeBlock), false, 0, 0)
	if gotErr != nil {
		t.Errorf("Render() failed: %v", gotErr)
	}

	golden.RequireEqual(t, got)
}

// TestRenderer_DimsLinesOutsideTheHighlight checks that a line a code block
// leaves out of its highlighted range is drawn in the theme's code block color,
// and a highlighted line is not.
func TestRenderer_DimsLinesOutsideTheHighlight(t *testing.T) {
	profile := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(profile) })

	theme := filepath.Join(t.TempDir(), "theme.json")
	if err := os.WriteFile(theme, []byte(`{"code_block": {"color": "#808080"}}`), 0o644); err != nil {
		t.Fatalf("could not write theme: %v", err)
	}

	r, err := NewRenderer(theme)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	got, err := r.Render("```go{2}\ndimmed()\nhighlighted()\n```", false, 0, 0)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	gray := "\x1b[38;2;128;128;128m"
	if !strings.Contains(got, gray+"dimmed()") {
		t.Errorf("line outside the range is not gray:\n%q", got)
	}
	if strings.Contains(got, gray+"highlighted()") {
		t.Errorf("highlighted line is gray:\n%q", got)
	}
}

func TestRenderer_RenderCodeBlockWithLineNumbers(t *testing.T) {
	r, err := NewRenderer("dark")
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	codeBlock := `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}`

	got, gotErr := r.Render(
		fmt.Sprintf("# Slide\n```go{3,6} --numbered\n%s\n```\n", codeBlock),
		false, 0, 0,
	)
	if gotErr != nil {
		t.Errorf("Render() failed: %v", gotErr)
	}

	golden.RequireEqual(t, got)
}

func TestRenderer_RenderCodeBlockWithStartFromLine(t *testing.T) {
	r, err := NewRenderer("dark")
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	codeBlock := `func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.head; e != nil; e = e.next {
			if !yield(e.val) {
				return
			}
		}
	}
}`

	got, gotErr := r.Render(
		fmt.Sprintf("# Slide\n```go{2-8} --numbered --start-at-line 10\n%s\n```\n", codeBlock),
		false, 0, 0,
	)
	if gotErr != nil {
		t.Errorf("Render() failed: %v", gotErr)
	}

	golden.RequireEqual(t, got)
}

func TestRenderer_GridLayout(t *testing.T) {
	r, err := NewRenderer("dark")
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	got, gotErr := r.Render(
		`
Heading
=======

[grid]
[column]
## Grid column 1 data
[/column]

[column]
## Grid column 2

> Lorem ipsum dolor sit amet, consectetur adipiscing elit.
[/column]

[column]
## Grid column 3

- item 1
- item 2
- item 3
[/column]
[/grid]

[Epilogue](https://kyma.ink)
		`, false, 0, 0)
	if gotErr != nil {
		t.Errorf("Render() failed: %v", gotErr)
	}

	golden.RequireEqual(t, got)
}
