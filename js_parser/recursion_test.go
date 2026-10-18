package js_parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ije/esbuild-internal/config"
	"github.com/ije/esbuild-internal/logger"
	"github.com/ije/esbuild-internal/test"
)

func TestRecursionLimit(t *testing.T) {
	for _, tc := range []struct {
		name, before, open, value, close, after string
	}{
		{"parentheses", "export default ", "(", "0", ")", ";"},
		{"unary", "export default ", "!", "x", "", ";"},
		{"assignment", "", "x=", "0", "", ";"},
		{"exponent", "export default ", "x**", "x", "", ";"},
		{"conditional", "export default ", "x?", "0", ":0", ";"},
		{"array", "export default ", "[", "0", "]", ";"},
		{"object", "export default ", "({x:", "0", "})", ";"},
		{"call", "export default ", "f(", "0", ")", ";"},
		{"call chain", "export default x", "", "", "()", ";"},
		{"member chain", "export default x", "", "", ".x", ";"},
		{"template", "export default ", "`a${", "0", "}`", ";"},
		{"if", "", "if(x)", "x();", "", ""},
		{"block", "", "{", "x();", "}", ""},
		{"function", "", "function f(){", "x();", "}", ""},
		{"binding", "let ", "[", "x", "]", "=value;"},
		{"jsx", "export default ", "<>", "x", "</>", ";"},
		{"type", "type X=", "(", "number", ")", ";"},
		{"type argument", "let x:", "Array<", "number", ">", ";"},
		{"type assertion", "export default ", "<X>", "x", "", ";"},
		{"arrow binding", "(", "[", "x", "]", ")=>0;"},
		{"arrow return type", "(x):", "Array<", "number", ">", "=>x;"},
		{"namespace", "namespace a", ".a", "", "", "{}"},
		{"class modifiers", "class X {", "readonly ", "x", "", ": number}"},
	} {
		for _, depth := range []int{16, 2 * maxRecursionDepth} {
			t.Run(fmt.Sprintf("%s/%d", tc.name, depth), func(t *testing.T) {
				options := config.Options{TS: config.TSOptions{Parse: true}}
				options.JSX.Parse = tc.name == "jsx"
				code := tc.before + strings.Repeat(tc.open, depth) + tc.value + strings.Repeat(tc.close, depth) + tc.after
				log := logger.NewDeferLog(logger.DeferLogNoVerboseOrDebug, nil)
				_, ok := Parse(log, test.SourceForTest(code), OptionsFromConfig(&options))
				msgs := log.Done()
				if depth == 16 {
					if !ok || len(msgs) != 0 {
						t.Fatalf("ordinary nesting failed: %v", msgs)
					}
				} else if ok || len(msgs) != 1 || msgs[0].Kind != logger.Error || msgs[0].Data.Text != "Source nesting is too deep" {
					t.Fatalf("expected nesting error, got ok=%v, messages=%v", ok, msgs)
				}
			})
		}
	}
}

func TestRecursionLimitIgnoresText(t *testing.T) {
	text := strings.Repeat("(", 2*maxRecursionDepth)
	for _, code := range []string{
		"export default '" + text + "';",
		"/*" + text + "*/ export default 0;",
		"export default /" + strings.Repeat(`\(`, 2*maxRecursionDepth) + "/;",
		strings.Repeat("f();", 2*maxRecursionDepth),
	} {
		log := logger.NewDeferLog(logger.DeferLogNoVerboseOrDebug, nil)
		_, ok := Parse(log, test.SourceForTest(code), OptionsFromConfig(&config.Options{}))
		if msgs := log.Done(); !ok || len(msgs) != 0 {
			t.Fatalf("valid source failed: %v", msgs)
		}
	}
}
