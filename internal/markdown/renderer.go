package markdown

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/museslabs/kyma/internal/config"
	"github.com/museslabs/kyma/internal/img"
	"github.com/museslabs/kyma/internal/img/placeholder"
)

// glamourMargin is the horizontal margin glamour reserves around a document,
// on top of the word wrap width it is configured with.
const glamourMargin = 2

// proseWidth is how wide a slide's content is drawn, unless the slide is only a
// grid: what glamour's default 80 column wrap comes out as once its margin is
// taken. Keeping it narrower than the slide is what leaves the slide's layout
// room to center or right align it, as a block padded to the slide's full width
// cannot move.
const proseWidth = 80 - glamourMargin

type RendererOption func(*Renderer) error

type Renderer struct {
	tr      *glamour.TermRenderer
	wrapped map[int]*glamour.TermRenderer
	parser  *MarkdownParser
	options rendererOptions
	// dimmed draws the lines a code block leaves out of its highlighted range,
	// in the color the theme sets on code blocks.
	dimmed lipgloss.Style
}

type rendererOptions struct {
	imgBackend img.ImageBackend
	theme      string
	fontSize   map[int]int
}

func NewRenderer(theme string, options ...RendererOption) (*Renderer, error) {
	tr, err := glamour.NewTermRenderer(glamour.WithStylePath(theme))
	if err != nil {
		return nil, err
	}

	p := NewMarkdownParser()
	p.Register(Prioritized[Parser](NewImageParser(), 1))
	p.Register(Prioritized[Parser](NewCodeBlockParser(), 1))
	p.Register(Prioritized[Parser](NewGridParser(), 1))
	p.Register(Prioritized[Parser](NewGridRowParser(), 2))
	p.Register(Prioritized[Parser](NewGridColumnParser(), 3))

	// Looked up once, since a custom theme is read from disk and a slide is
	// re-rendered on every animation frame.
	dimmed := lipgloss.NewStyle()
	if color := config.CodeBlockColor(theme); color != "" {
		dimmed = dimmed.Foreground(lipgloss.Color(color))
	}

	r := &Renderer{
		tr:      tr,
		wrapped: map[int]*glamour.TermRenderer{},
		parser:  p,
		options: rendererOptions{
			imgBackend: img.Get("chafa"),
			theme:      theme,
		},
		dimmed: dimmed,
	}
	for _, o := range options {
		if err := o(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// termRenderer returns a glamour renderer whose output is exactly width cells
// wide. Glamour lays its document out inside a two cell margin that it adds on
// top of the wrap width, so the wrap has to be asked for two cells short.
// Renderers are cached because a slide is re-rendered on every animation frame.
func (r *Renderer) termRenderer(width int) (*glamour.TermRenderer, error) {
	if width <= glamourMargin {
		return r.tr, nil
	}

	if tr, ok := r.wrapped[width]; ok {
		return tr, nil
	}

	tr, err := glamour.NewTermRenderer(
		glamour.WithStylePath(r.options.theme),
		glamour.WithWordWrap(width-glamourMargin),
	)
	if err != nil {
		return nil, err
	}
	r.wrapped[width] = tr

	return tr, nil
}

func (r *Renderer) Render(in string, animating bool, width, height int) (string, error) {
	return r.RenderBytes([]byte(in), animating, width, height)
}

func (r *Renderer) RenderBytes(in []byte, animating bool, width, height int) (string, error) {
	var b strings.Builder

	// Clear kitty images
	if !animating {
		b.WriteString(placeholder.Clear())
	}

	if err := r.renderNode(r.parser.Parse(in), animating, width, height, &b); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (r *Renderer) renderNode(n Node, animating bool, width, height int, b *strings.Builder) error {
	if n == nil {
		return nil
	}

	switch n.Kind() {
	case NodeKindMarkdownRoot:
		// A slide that is only a grid, as every master layout is, lays out the
		// whole slide and fills it. It ends on the grid's last line, as the line
		// break after it, or the blank lines around its tags, would push a grid
		// as tall as the slide past its bottom.
		if onlyGrids(n.Children()) {
			var grids strings.Builder
			for _, c := range n.Children() {
				if c.Kind() == NodeKindGlamour {
					continue
				}
				if err := r.renderNode(c, animating, width, height, &grids); err != nil {
					return err
				}
			}
			b.WriteString(strings.TrimSuffix(grids.String(), "\n"))
			return nil
		}

		// Anything else placed directly on the slide is drawn as wide as one
		// block, so it lines up and the slide's layout aligns it as a whole, and
		// to the room glamour's text has once its margin is taken.
		block := min(width-2*glamourMargin, proseWidth)
		for _, c := range n.Children() {
			w := block
			switch c.Kind() {
			case NodeKindGlamour, NodeKindHeading:
				// Glamour leaves its margin out of what it draws.
				w = block + 2*glamourMargin
			case NodeKindImage:
				// An image keeps the size it asked for.
				w = width
			}
			if err := r.renderNode(c, animating, w, height, b); err != nil {
				return err
			}
		}
		return nil

	case NodeKindGlamour:
		n := n.(*GlamourNode)

		tr, err := r.termRenderer(width)
		if err != nil {
			return err
		}

		out, err := tr.Render(n.Text)
		if err != nil {
			return err
		}
		b.WriteString(out)

	case NodeKindHeading:
		if err := r.renderHeading(n.(*HeadingNode), animating, width, b); err != nil {
			return err
		}

	case NodeKindError:
		n := n.(*ErrorNode)
		fmt.Fprintf(b, "%s\n", lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render("⚠ "+n.Message))

	case NodeKindImage:
		n := n.(*ImageNode)

		// Keep an image from spilling out of the column it was placed in.
		imgWidth, imgHeight := n.Width, n.Height
		if width > 0 && imgWidth > width {
			if imgWidth > 0 {
				imgHeight = imgHeight * width / imgWidth
			}
			imgWidth = width
		}

		// An image drawn as text, in symbols or in kitty's placeholders, goes
		// wherever the text around it does, so transitions draw it as it is.
		// Images kept to be drawn elsewhere always go in placeholders, at
		// their full resolution, whatever this terminal can show.
		if symbols := !placeholder.Keeping() && r.options.imgBackend.SymbolsOnly(); symbols || placeholder.Enabled() {
			out, err := r.options.imgBackend.Render(n.Path, imgWidth, imgHeight, symbols)
			if err != nil {
				fmt.Fprintf(b, "[Error rendering image: %s]", n.Label)
				break
			}
			b.WriteString(out)
			break
		}

		limg, err := r.options.imgBackend.Render(n.Path, imgWidth, imgHeight, true)
		if err != nil {
			fmt.Fprintf(b, "[Error rendering image: %s]", n.Label)
			break
		}

		himg, err := r.options.imgBackend.Render(n.Path, imgWidth, imgHeight, false)
		if err != nil {
			fmt.Fprintf(b, "[Error rendering image: %s]", n.Label)
			break
		}

		if !animating {
			b.WriteString(ansi.SaveCursor)
			b.WriteString(himg)
			b.WriteString(ansi.RestoreCursor)
			// Under the image the text only holds its place: drawing the
			// symbols there too shows them through every transparent pixel.
			b.WriteString(blank(limg))
		} else {
			b.WriteString(limg)
		}

	case NodeKindCodeBlock:
		n := n.(*CodeBlockNode)

		lines := strings.Split(n.Code, "\n")

		var renderedContent string
		if n.Language != "" {
			lexer := lexers.Get(n.Language)
			if lexer == nil {
				lexer = lexers.Fallback
			}
			lexer = chroma.Coalesce(lexer)
			style := config.GetChromaStyle(r.options.theme)

			renderedContent = r.renderHighlightedCode(n.Code, lines, n, lexer, style)
		} else {
			renderedContent = r.renderPlainCode(lines, n)
		}

		codeWidth := 78
		if width > 0 {
			codeWidth = width
		}
		b.WriteString(lipgloss.NewStyle().Width(codeWidth).Render(renderedContent))
		b.WriteString("\n")

	case NodeKindGrid, NodeKindGridRow, NodeKindGridColumn:
		out, err := r.renderNodeBlock(n, animating, sizing{
			Width:  gridWidth(width),
			Budget: remainingHeight(height, b),
		})
		if err != nil {
			return err
		}
		b.WriteString(out)
		b.WriteString("\n")
		return nil

	default:
		return fmt.Errorf("invalid node kind: %d", n.Kind())
	}

	for _, c := range n.Children() {
		if err := r.renderNode(c, animating, width, height, b); err != nil {
			return err
		}
	}

	return nil
}

// onlyGrids reports whether the nodes placed directly on a slide are grids and
// nothing else, apart from the blank lines glamour holds around their tags.
func onlyGrids(nodes []Node) bool {
	for _, n := range nodes {
		switch n.Kind() {
		case NodeKindGrid, NodeKindGridRow, NodeKindGridColumn:
		case NodeKindGlamour:
			if strings.TrimSpace(n.(*GlamourNode).Text) != "" {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// blank returns s with every line turned into blank cells as wide as it was
// drawn. The cells are no-break spaces, which show nothing but, unlike spaces,
// still count as ink when a column lines its content up.
func blank(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat("\u00a0", ansi.StringWidth(line))
	}
	return strings.Join(lines, "\n")
}

func WithImageBackend(backend string) RendererOption {
	return func(r *Renderer) error {
		r.options.imgBackend = img.Get(backend)
		return nil
	}
}

// WithFontSize draws headings at the scale set for their level, keyed by level.
// It is only meant for terminals that implement kitty's text sizing protocol,
// see [textsize.Supported]. Levels left at scale 1 stay with glamour.
func WithFontSize(sizes map[int]int) RendererOption {
	return func(r *Renderer) error {
		var levels []int
		for level, size := range sizes {
			if size > 1 {
				levels = append(levels, level)
			}
		}
		if len(levels) == 0 {
			return nil
		}

		r.options.fontSize = sizes
		r.parser.Register(Prioritized[Parser](NewHeadingParser(levels...), 1))
		return nil
	}
}

func (r *Renderer) formatLineNumber(lineNum, width int) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingRight(1)

	lineNumStr := strconv.Itoa(lineNum)
	paddedLineNum := fmt.Sprintf("%*s", width-1, lineNumStr)
	return style.Render(paddedLineNum)
}

func (r *Renderer) getLineNumberWidth(startLine, totalLines int) int {
	if totalLines == 0 {
		return 0
	}
	maxLineNum := startLine + totalLines - 1
	return len(strconv.Itoa(maxLineNum)) + 2
}

func (r *Renderer) renderPlainCode(lines []string, info *CodeBlockNode) string {
	var result strings.Builder
	lineNumberWidth := 0

	if info.ShowLineNumbers {
		lineNumberWidth = r.getLineNumberWidth(info.StartLine, len(lines))
	}

	for i, line := range lines {
		displayLineNum := info.StartLine + i

		if info.ShowLineNumbers {
			result.WriteString(r.formatLineNumber(displayLineNum, lineNumberWidth))
		}

		result.WriteString(line)

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func (r *Renderer) renderHighlightedCode(
	content string,
	lines []string,
	info *CodeBlockNode,
	lexer chroma.Lexer,
	style *chroma.Style,
) string {
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		return r.renderPlainCode(lines, info)
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return r.renderPlainCode(lines, info)
	}

	var formattedBuf strings.Builder
	if err := formatter.Format(&formattedBuf, style, iterator); err != nil {
		return r.renderPlainCode(lines, info)
	}

	formattedLines := strings.Split(formattedBuf.String(), "\n")

	// Ensure we have the same number of lines
	for len(formattedLines) < len(lines) {
		formattedLines = append(formattedLines, "")
	}
	if len(formattedLines) > len(lines) {
		formattedLines = formattedLines[:len(lines)]
	}

	var result strings.Builder
	lineNumberWidth := 0

	if info.ShowLineNumbers {
		lineNumberWidth = r.getLineNumberWidth(info.StartLine, len(lines))
	}

	for i, line := range lines {
		displayLineNum := info.StartLine + i
		relativeLineNum := i + 1

		if info.ShowLineNumbers {
			result.WriteString(r.formatLineNumber(displayLineNum, lineNumberWidth))
		}

		if r.shouldHighlightLine(relativeLineNum, info.Ranges) {
			formattedLine := ""
			if i < len(formattedLines) {
				formattedLine = strings.TrimRight(formattedLines[i], " \t\n\r")
			}
			if formattedLine == "" {
				formattedLine = line
			}
			result.WriteString(formattedLine)
		} else {
			result.WriteString(r.dimmed.Render(line))
		}

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func (r *Renderer) shouldHighlightLine(lineNum int, ranges []CodeBlockLineRange) bool {
	if len(ranges) == 0 {
		return true // highlight all lines if no ranges specified
	}

	for _, r := range ranges {
		if lineNum >= r.Start && lineNum <= r.End {
			return true
		}
	}
	return false
}
