package markdown

import (
	"bytes"
	"testing"
)

func TestHeadingParser_Parse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Node
	}{
		{
			name: "h1",
			in:   "# Title\nbody",
			want: &HeadingNode{Level: 1, Text: "Title"},
		},
		{
			name: "h2",
			in:   "## Title",
			want: &HeadingNode{Level: 2, Text: "Title"},
		},
		{
			name: "inline markdown is kept for glamour",
			in:   "# Using `kyma` **today**",
			want: &HeadingNode{Level: 1, Text: "Using `kyma` **today**"},
		},
		{
			name: "closing sequence is dropped",
			in:   "# Title ##",
			want: &HeadingNode{Level: 1, Text: "Title"},
		},
		{
			name: "hash glued to the text is not a closing sequence",
			in:   "# C#",
			want: &HeadingNode{Level: 1, Text: "C#"},
		},
		{
			name: "indented up to three spaces",
			in:   "   # Title",
			want: &HeadingNode{Level: 1, Text: "Title"},
		},
		{
			name: "indented code block",
			in:   "    # Title",
			want: nil,
		},
		{
			name: "hashtag",
			in:   "#hashtag",
			want: nil,
		},
		{
			name: "not at the start of the line",
			in:   "issue # 42",
			want: nil,
		},
		{
			name: "empty heading",
			in:   "# ",
			want: nil,
		},
		{
			name: "level without a font size",
			in:   "### Title",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewHeadingParser(1, 2)
			r := bytes.NewReader([]byte(tt.in))

			// Advance to the first '#', the trigger byte the parser expects
			// to have been consumed already.
			for {
				b, err := r.ReadByte()
				if err != nil {
					t.Fatalf("no '#' in %q", tt.in)
				}
				if b == '#' {
					break
				}
			}

			got := p.Parse(r, nil)
			if Dump(got) != Dump(tt.want) {
				t.Errorf("Parse() got:\n%s\nwant:\n%s", Dump(got), Dump(tt.want))
			}
		})
	}
}

func TestMarkdownParser_HeadingsWithFontSize(t *testing.T) {
	p := NewMarkdownParser()
	p.Register(Prioritized[Parser](NewHeadingParser(1), 1))

	got := Dump(p.Parse([]byte("# Title\n\nbody\n\n## Section\n\n    # code\n")))
	want := Dump(node(&MarkdownRootNode{},
		&HeadingNode{Level: 1, Text: "Title"},
		&GlamourNode{Text: "\nbody\n\n## Section\n\n    # code\n"},
	))
	if got != want {
		t.Errorf("Parse() got:\n%s\nwant:\n%s", got, want)
	}
}
