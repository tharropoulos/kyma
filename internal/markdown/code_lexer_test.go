package markdown

import (
	"testing"

	"github.com/alecthomas/chroma/v2"
)

func TestCodeLexer_TypeScriptGenericsAreNotTags(t *testing.T) {
	code := `type TupleOf<L extends number, T, Acc extends T[] = []> =
  Acc["length"] extends L ? Acc : TupleOf<L, T, [...Acc, T]>;
`
	for _, language := range []string{"ts", "typescript", "TS"} {
		tokens, err := chroma.Tokenise(codeLexer(language), nil, code)
		if err != nil {
			t.Fatalf("%s: Tokenise() failed: %v", language, err)
		}
		for _, token := range tokens {
			switch token.Type {
			case chroma.Error, chroma.NameTag, chroma.NameAttribute:
				t.Errorf("%s: %q lexed as %s", language, token.Value, token.Type)
			}
		}
	}
}

func TestCodeLexer_TSXKeepsTags(t *testing.T) {
	tokens, err := chroma.Tokenise(codeLexer("tsx"), nil, `const el = <div className="x" />;`)
	if err != nil {
		t.Fatalf("Tokenise() failed: %v", err)
	}
	for _, token := range tokens {
		if token.Type == chroma.NameTag && token.Value == "div" {
			return
		}
	}
	t.Errorf("tsx did not lex <div> as a tag: %v", tokens)
}
