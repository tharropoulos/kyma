package markdown

import (
	"strings"
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// codeLexer is the lexer a code block in language is highlighted with.
func codeLexer(language string) chroma.Lexer {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	if lexer.Config().Name == "TypeScript" && !strings.EqualFold(language, "tsx") {
		return typeScriptLexer()
	}
	return chroma.Coalesce(lexer)
}

// typeScriptLexer is chroma's TypeScript lexer without the JSX rules it shares
// with TSX. Those read every generic, as in TupleOf<L extends number>, as a
// tag, and mark what a tag cannot hold, such as the T[] = [] that follows, as
// errors.
var typeScriptLexer = sync.OnceValue(func() chroma.Lexer {
	base := lexers.Get("typescript").(*chroma.RegexLexer)
	rules, err := base.Rules()
	if err != nil {
		return chroma.Coalesce(base)
	}
	return chroma.Coalesce(chroma.MustNewLexer(base.Config(), func() chroma.Rules {
		// An empty state leaves nothing for root to include.
		return rules.Merge(chroma.Rules{"jsx": nil})
	}))
})
