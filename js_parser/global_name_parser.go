package js_parser

import (
	"github.com/ije/esbuild-internal/helpers"
	"github.com/ije/esbuild-internal/js_lexer"
	"github.com/ije/esbuild-internal/logger"
)

func ParseGlobalName(log logger.Log, source logger.Source) (result []string, ok bool) {
	ok = true
	defer func() {
		r := recover()
		if _, isLexerPanic := r.(js_lexer.LexerPanic); isLexerPanic {
			ok = false
		} else if r != nil {
			panic(r)
		}
	}()

	lexer := js_lexer.NewLexerGlobalName(log, source)

	// Start off with an identifier
	result = append(result, lexer.Identifier.String)
	lexer.Expect(js_lexer.TIdentifier)

	// Follow with dot or index expressions
	for lexer.Token != js_lexer.TEndOfFile {
		switch lexer.Token {
		case js_lexer.TDot:
			lexer.Next()
			if !lexer.IsIdentifierOrKeyword() {
				lexer.Expect(js_lexer.TIdentifier)
			}
			result = append(result, lexer.Identifier.String)
			lexer.Next()

		case js_lexer.TOpenBracket:
			lexer.Next()
			result = append(result, helpers.UTF16ToString(lexer.StringLiteral()))
			lexer.Expect(js_lexer.TStringLiteral)
			lexer.Expect(js_lexer.TCloseBracket)

		default:
			lexer.Expect(js_lexer.TDot)
		}
	}

	return
}
