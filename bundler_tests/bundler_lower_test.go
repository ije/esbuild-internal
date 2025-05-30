package bundler_tests

// This file contains tests for "lowering" syntax, which means converting it to
// older JavaScript. For example, "a ** b" becomes a call to "Math.pow(a, b)"
// when lowered. Which syntax is lowered is determined by the language target.

import (
	"testing"

	"github.com/ije/esbuild-internal/compat"
	"github.com/ije/esbuild-internal/config"
)

var lower_suite = suite{
	name: "lower",
}

func TestLowerOptionalCatchNameCollisionNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				try {}
				catch { var e, e2 }
				var e3
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2018),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerObjectSpreadNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.jsx": `
				let tests = [
					{...a, ...b},
					{a, b, ...c},
					{...a, b, c},
					{a, ...b, c},
					{a, b, ...c, ...d, e, f, ...g, ...h, i, j},
				]
				let jsx = [
					<div {...a} {...b}/>,
					<div a b {...c}/>,
					<div {...a} b c/>,
					<div a {...b} c/>,
					<div a b {...c} {...d} e f {...g} {...h} i j/>,
				]
			`,
		},
		entryPaths: []string{"/entry.jsx"},
		options: config.Options{
			UnsupportedJSFeatures: es(2017),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerExponentiationOperatorNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				let tests = {
					// Exponentiation operator
					0: a ** b ** c,
					1: (a ** b) ** c,

					// Exponentiation assignment operator
					2: a **= b,
					3: a.b **= c,
					4: a[b] **= c,
					5: a().b **= c,
					6: a()[b] **= c,
					7: a[b()] **= c,
					8: a()[b()] **= c,

					// These all should not need capturing (no object identity)
					9: a[0] **= b,
					10: a[false] **= b,
					11: a[null] **= b,
					12: a[void 0] **= b,
					13: a[123n] **= b,
					14: a[this] **= b,

					// These should need capturing (have object identitiy)
					15: a[/x/] **= b,
					16: a[{}] **= b,
					17: a[[]] **= b,
					18: a[() => {}] **= b,
					19: a[function() {}] **= b,
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2015),
			AbsOutputFile:         "/out.js",
		},
		expectedScanLog: `entry.js: WARNING: Big integer literals are not available in the configured target environment and may crash at run-time
`,
	})
}

func TestLowerPrivateFieldAssignments2015NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#x
					unary() {
						this.#x++
						this.#x--
						++this.#x
						--this.#x
					}
					binary() {
						this.#x = 1
						this.#x += 1
						this.#x -= 1
						this.#x *= 1
						this.#x /= 1
						this.#x %= 1
						this.#x **= 1
						this.#x <<= 1
						this.#x >>= 1
						this.#x >>>= 1
						this.#x &= 1
						this.#x |= 1
						this.#x ^= 1
						this.#x &&= 1
						this.#x ||= 1
						this.#x ??= 1
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2015),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateFieldAssignments2019NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#x
					unary() {
						this.#x++
						this.#x--
						++this.#x
						--this.#x
					}
					binary() {
						this.#x = 1
						this.#x += 1
						this.#x -= 1
						this.#x *= 1
						this.#x /= 1
						this.#x %= 1
						this.#x **= 1
						this.#x <<= 1
						this.#x >>= 1
						this.#x >>>= 1
						this.#x &= 1
						this.#x |= 1
						this.#x ^= 1
						this.#x &&= 1
						this.#x ||= 1
						this.#x ??= 1
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2019),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateFieldAssignments2020NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#x
					unary() {
						this.#x++
						this.#x--
						++this.#x
						--this.#x
					}
					binary() {
						this.#x = 1
						this.#x += 1
						this.#x -= 1
						this.#x *= 1
						this.#x /= 1
						this.#x %= 1
						this.#x **= 1
						this.#x <<= 1
						this.#x >>= 1
						this.#x >>>= 1
						this.#x &= 1
						this.#x |= 1
						this.#x ^= 1
						this.#x &&= 1
						this.#x ||= 1
						this.#x ??= 1
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateFieldAssignmentsNextNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#x
					unary() {
						this.#x++
						this.#x--
						++this.#x
						--this.#x
					}
					binary() {
						this.#x = 1
						this.#x += 1
						this.#x -= 1
						this.#x *= 1
						this.#x /= 1
						this.#x %= 1
						this.#x **= 1
						this.#x <<= 1
						this.#x >>= 1
						this.#x >>>= 1
						this.#x &= 1
						this.#x |= 1
						this.#x ^= 1
						this.#x &&= 1
						this.#x ||= 1
						this.#x ??= 1
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			AbsOutputFile: "/out.js",
		},
	})
}

func TestLowerPrivateFieldOptionalChain2019NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#x
					foo() {
						this?.#x.y
						this?.y.#x
						this.#x?.y
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2019),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateFieldOptionalChain2020NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#x
					foo() {
						this?.#x.y
						this?.y.#x
						this.#x?.y
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateFieldOptionalChainNextNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#x
					foo() {
						this?.#x.y
						this?.y.#x
						this.#x?.y
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			AbsOutputFile: "/out.js",
		},
	})
}

func TestTSLowerPrivateFieldOptionalChain2015NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				class Foo {
					#x
					foo() {
						this?.#x.y
						this?.y.#x
						this.#x?.y
					}
				}
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			UnsupportedJSFeatures: es(2015),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestTSLowerPrivateStaticMembers2015NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				class Foo {
					static #x
					static get #y() {}
					static set #y(x) {}
					static #z() {}
					foo() {
						Foo.#x += 1
						Foo.#y += 1
						Foo.#z()
					}
				}
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			UnsupportedJSFeatures: es(2015),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestTSLowerPrivateFieldAndMethodAvoidNameCollision2015(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				export class WeakMap {
					#x
				}
				export class WeakSet {
					#y() {}
				}
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2015),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateGetterSetter2015(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					get #foo() { return this.foo }
					set #bar(val) { this.bar = val }
					get #prop() { return this.prop }
					set #prop(val) { this.prop = val }
					foo(fn) {
						fn().#foo
						fn().#bar = 1
						fn().#prop
						fn().#prop = 2
					}
					unary(fn) {
						fn().#prop++;
						fn().#prop--;
						++fn().#prop;
						--fn().#prop;
					}
					binary(fn) {
						fn().#prop = 1;
						fn().#prop += 1;
						fn().#prop -= 1;
						fn().#prop *= 1;
						fn().#prop /= 1;
						fn().#prop %= 1;
						fn().#prop **= 1;
						fn().#prop <<= 1;
						fn().#prop >>= 1;
						fn().#prop >>>= 1;
						fn().#prop &= 1;
						fn().#prop |= 1;
						fn().#prop ^= 1;
						fn().#prop &&= 1;
						fn().#prop ||= 1;
						fn().#prop ??= 1;
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2015),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateGetterSetter2019(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					get #foo() { return this.foo }
					set #bar(val) { this.bar = val }
					get #prop() { return this.prop }
					set #prop(val) { this.prop = val }
					foo(fn) {
						fn().#foo
						fn().#bar = 1
						fn().#prop
						fn().#prop = 2
					}
					unary(fn) {
						fn().#prop++;
						fn().#prop--;
						++fn().#prop;
						--fn().#prop;
					}
					binary(fn) {
						fn().#prop = 1;
						fn().#prop += 1;
						fn().#prop -= 1;
						fn().#prop *= 1;
						fn().#prop /= 1;
						fn().#prop %= 1;
						fn().#prop **= 1;
						fn().#prop <<= 1;
						fn().#prop >>= 1;
						fn().#prop >>>= 1;
						fn().#prop &= 1;
						fn().#prop |= 1;
						fn().#prop ^= 1;
						fn().#prop &&= 1;
						fn().#prop ||= 1;
						fn().#prop ??= 1;
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2019),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateGetterSetter2020(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					get #foo() { return this.foo }
					set #bar(val) { this.bar = val }
					get #prop() { return this.prop }
					set #prop(val) { this.prop = val }
					foo(fn) {
						fn().#foo
						fn().#bar = 1
						fn().#prop
						fn().#prop = 2
					}
					unary(fn) {
						fn().#prop++;
						fn().#prop--;
						++fn().#prop;
						--fn().#prop;
					}
					binary(fn) {
						fn().#prop = 1;
						fn().#prop += 1;
						fn().#prop -= 1;
						fn().#prop *= 1;
						fn().#prop /= 1;
						fn().#prop %= 1;
						fn().#prop **= 1;
						fn().#prop <<= 1;
						fn().#prop >>= 1;
						fn().#prop >>>= 1;
						fn().#prop &= 1;
						fn().#prop |= 1;
						fn().#prop ^= 1;
						fn().#prop &&= 1;
						fn().#prop ||= 1;
						fn().#prop ??= 1;
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateGetterSetterNext(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					get #foo() { return this.foo }
					set #bar(val) { this.bar = val }
					get #prop() { return this.prop }
					set #prop(val) { this.prop = val }
					foo(fn) {
						fn().#foo
						fn().#bar = 1
						fn().#prop
						fn().#prop = 2
					}
					unary(fn) {
						fn().#prop++;
						fn().#prop--;
						++fn().#prop;
						--fn().#prop;
					}
					binary(fn) {
						fn().#prop = 1;
						fn().#prop += 1;
						fn().#prop -= 1;
						fn().#prop *= 1;
						fn().#prop /= 1;
						fn().#prop %= 1;
						fn().#prop **= 1;
						fn().#prop <<= 1;
						fn().#prop >>= 1;
						fn().#prop >>>= 1;
						fn().#prop &= 1;
						fn().#prop |= 1;
						fn().#prop ^= 1;
						fn().#prop &&= 1;
						fn().#prop ||= 1;
						fn().#prop ??= 1;
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestLowerPrivateMethod2019(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					#field
					#method() {}
					baseline() {
						a().foo
						b().foo(x)
						c()?.foo(x)
						d().foo?.(x)
						e()?.foo?.(x)
					}
					privateField() {
						a().#field
						b().#field(x)
						c()?.#field(x)
						d().#field?.(x)
						e()?.#field?.(x)
						f()?.foo.#field(x).bar()
					}
					privateMethod() {
						a().#method
						b().#method(x)
						c()?.#method(x)
						d().#method?.(x)
						e()?.#method?.(x)
						f()?.foo.#method(x).bar()
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2019),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateMethod2020(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					#field
					#method() {}
					baseline() {
						a().foo
						b().foo(x)
						c()?.foo(x)
						d().foo?.(x)
						e()?.foo?.(x)
					}
					privateField() {
						a().#field
						b().#field(x)
						c()?.#field(x)
						d().#field?.(x)
						e()?.#field?.(x)
						f()?.foo.#field(x).bar()
					}
					privateMethod() {
						a().#method
						b().#method(x)
						c()?.#method(x)
						d().#method?.(x)
						e()?.#method?.(x)
						f()?.foo.#method(x).bar()
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateMethodNext(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					#field
					#method() {}
					baseline() {
						a().foo
						b().foo(x)
						c()?.foo(x)
						d().foo?.(x)
						e()?.foo?.(x)
					}
					privateField() {
						a().#field
						b().#field(x)
						c()?.#field(x)
						d().#field?.(x)
						e()?.#field?.(x)
						f()?.foo.#field(x).bar()
					}
					privateMethod() {
						a().#method
						b().#method(x)
						c()?.#method(x)
						d().#method?.(x)
						e()?.#method?.(x)
						f()?.foo.#method(x).bar()
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestLowerPrivateClassExpr2020NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export let Foo = class {
					#field
					#method() {}
					static #staticField
					static #staticMethod() {}
					foo() {
						this.#field = this.#method()
						Foo.#staticField = Foo.#staticMethod()
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateMethodWithModifiers2020(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo {
					*#g() {}
					async #a() {}
					async *#ag() {}

					static *#sg() {}
					static async #sa() {}
					static async *#sag() {}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerAsync2016NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				async function foo(bar) {
					await bar
					return [this, arguments]
				}
				class Foo {async foo() {}}
				new (class Bar extends class { } {
					constructor() {
						let x = 1;
						(async () => {
							console.log("before super", x);  // (1) Sync phase
							await 1;
							console.log("after super", x);   // (2) Async phase
						})();
						super();
						x = 2;
					}
				})();
				export default [
					foo,
					Foo,
					async function() {},
					async () => {},
					{async foo() {}},
					class {async foo() {}},
					function() {
						return async (bar) => {
							await bar
							return [this, arguments]
						}
					},
				]
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerAsync2017NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				async function foo(bar) {
					await bar
					return arguments
				}
				class Foo {async foo() {}}
				export default [
					foo,
					Foo,
					async function() {},
					async () => {},
					{async foo() {}},
					class {async foo() {}},
					function() {
						return async (bar) => {
							await bar
							return [this, arguments]
						}
					},
				]
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2017),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerAsyncThis2016CommonJS(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				exports.foo = async () => this
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerAsyncThis2016ES6(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export {bar} from "./other"
				export let foo = async () => this
			`,
			"/other.js": `
				export let bar = async () => {}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
		debugLogs: true,
		expectedScanLog: `entry.js: DEBUG: Top-level "this" will be replaced with undefined since this file is an ECMAScript module
entry.js: NOTE: This file is considered to be an ECMAScript module because of the "export" keyword here:
`,
	})
}

func TestLowerAsyncES5(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import './fn-stmt'
				import './fn-expr'
				import './arrow-1'
				import './arrow-2'
				import './export-def-1'
				import './export-def-2'
				import './obj-method'
			`,
			"/fn-stmt.js":      `async function foo() {}`,
			"/fn-expr.js":      `(async function() {})`,
			"/arrow-1.js":      `(async () => {})`,
			"/arrow-2.js":      `(async x => {})`,
			"/export-def-1.js": `export default async function foo() {}`,
			"/export-def-2.js": `export default async function() {}`,
			"/obj-method.js":   `({async foo() {}})`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(5),
			AbsOutputFile:         "/out.js",
		},
		expectedScanLog: `arrow-1.js: ERROR: Transforming async functions to the configured target environment is not supported yet
arrow-2.js: ERROR: Transforming async functions to the configured target environment is not supported yet
export-def-1.js: ERROR: Transforming async functions to the configured target environment is not supported yet
export-def-2.js: ERROR: Transforming async functions to the configured target environment is not supported yet
fn-expr.js: ERROR: Transforming async functions to the configured target environment is not supported yet
fn-stmt.js: ERROR: Transforming async functions to the configured target environment is not supported yet
obj-method.js: ERROR: Transforming async functions to the configured target environment is not supported yet
`,
	})
}

func TestLowerAsyncSuperES2017NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Derived extends Base {
					async test(key) {
						return [
							await super.foo,
							await super[key],
							await ([super.foo] = [0]),
							await ([super[key]] = [0]),

							await (super.foo = 1),
							await (super[key] = 1),
							await (super.foo += 2),
							await (super[key] += 2),

							await ++super.foo,
							await ++super[key],
							await super.foo++,
							await super[key]++,

							await super.foo.name,
							await super[key].name,
							await super.foo?.name,
							await super[key]?.name,

							await super.foo(1, 2),
							await super[key](1, 2),
							await super.foo?.(1, 2),
							await super[key]?.(1, 2),

							await (() => super.foo)(),
							await (() => super[key])(),
							await (() => super.foo())(),
							await (() => super[key]())(),

							await super.foo` + "``" + `,
							await super[key]` + "``" + `,
						]
					}
				}

				// This covers a bug that caused a compiler crash
				let fn = async () => class extends Base {
					a = super.a
					b = () => super.b
					c() { return super.c }
					d() { return () => super.d }
				}

				// This covers a bug that generated bad code
				class Derived2 extends Base {
					async a() { return class { [super.foo] = 123 } }
					b = async () => class { [super.foo] = 123 }
				}

				// This covers putting the generated temporary variable inside the loop
				for (let i = 0; i < 3; i++) {
					objs.push({
						__proto__: {
							foo() { return i },
						},
						async bar() { return super.foo() },
					})
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2017),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerAsyncSuperES2016NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Derived extends Base {
					async test(key) {
						return [
							await super.foo,
							await super[key],
							await ([super.foo] = [0]),
							await ([super[key]] = [0]),

							await (super.foo = 1),
							await (super[key] = 1),
							await (super.foo += 2),
							await (super[key] += 2),

							await ++super.foo,
							await ++super[key],
							await super.foo++,
							await super[key]++,

							await super.foo.name,
							await super[key].name,
							await super.foo?.name,
							await super[key]?.name,

							await super.foo(1, 2),
							await super[key](1, 2),
							await super.foo?.(1, 2),
							await super[key]?.(1, 2),

							await (() => super.foo)(),
							await (() => super[key])(),
							await (() => super.foo())(),
							await (() => super[key]())(),

							await super.foo` + "``" + `,
							await super[key]` + "``" + `,
						]
					}
				}

				// This covers a bug that caused a compiler crash
				let fn = async () => class extends Base {
					a = super.a
					b = () => super.b
					c() { return super.c }
					d() { return () => super.d }
				}

				// This covers a bug that generated bad code
				class Derived2 extends Base {
					async a() { return class { [super.foo] = 123 } }
					b = async () => class { [super.foo] = 123 }
				}

				// This covers putting the generated temporary variable inside the loop
				for (let i = 0; i < 3; i++) {
					objs.push({
						__proto__: {
							foo() { return i },
						},
						async bar() { return super.foo() },
					})
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerStaticAsyncSuperES2021NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Derived extends Base {
					static test = async (key) => {
						return [
							await super.foo,
							await super[key],
							await ([super.foo] = [0]),
							await ([super[key]] = [0]),

							await (super.foo = 1),
							await (super[key] = 1),
							await (super.foo += 2),
							await (super[key] += 2),

							await ++super.foo,
							await ++super[key],
							await super.foo++,
							await super[key]++,

							await super.foo.name,
							await super[key].name,
							await super.foo?.name,
							await super[key]?.name,

							await super.foo(1, 2),
							await super[key](1, 2),
							await super.foo?.(1, 2),
							await super[key]?.(1, 2),

							await (() => super.foo)(),
							await (() => super[key])(),
							await (() => super.foo())(),
							await (() => super[key]())(),

							await super.foo` + "``" + `,
							await super[key]` + "``" + `,
						]
					}
				}

				// This covers a bug that caused a compiler crash
				let fn = async () => class extends Base {
					static a = super.a
					static b = () => super.b
					static c() { return super.c }
					static d() { return () => super.d }
				}

				// This covers a bug that generated bad code
				class Derived2 extends Base {
					static async a() { return class { [super.foo] = 123 } }
					static b = async () => class { [super.foo] = 123 }
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2021),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerStaticAsyncSuperES2016NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Derived extends Base {
					static test = async (key) => {
						return [
							await super.foo,
							await super[key],
							await ([super.foo] = [0]),
							await ([super[key]] = [0]),

							await (super.foo = 1),
							await (super[key] = 1),
							await (super.foo += 2),
							await (super[key] += 2),

							await ++super.foo,
							await ++super[key],
							await super.foo++,
							await super[key]++,

							await super.foo.name,
							await super[key].name,
							await super.foo?.name,
							await super[key]?.name,

							await super.foo(1, 2),
							await super[key](1, 2),
							await super.foo?.(1, 2),
							await super[key]?.(1, 2),

							await (() => super.foo)(),
							await (() => super[key])(),
							await (() => super.foo())(),
							await (() => super[key]())(),

							await super.foo` + "``" + `,
							await super[key]` + "``" + `,
						]
					}
				}

				// This covers a bug that caused a compiler crash
				let fn = async () => class extends Base {
					static a = super.a
					static b = () => super.b
					static c() { return super.c }
					static d() { return () => super.d }
				}

				// This covers a bug that generated bad code
				class Derived2 extends Base {
					static async a() { return class { [super.foo] = 123 } }
					static b = async () => class { [super.foo] = 123 }
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerStaticSuperES2021NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Derived extends Base {
					static test = key => {
						return [
							super.foo,
							super[key],
							([super.foo] = [0]),
							([super[key]] = [0]),

							(super.foo = 1),
							(super[key] = 1),
							(super.foo += 2),
							(super[key] += 2),

							++super.foo,
							++super[key],
							super.foo++,
							super[key]++,

							super.foo.name,
							super[key].name,
							super.foo?.name,
							super[key]?.name,

							super.foo(1, 2),
							super[key](1, 2),
							super.foo?.(1, 2),
							super[key]?.(1, 2),

							(() => super.foo)(),
							(() => super[key])(),
							(() => super.foo())(),
							(() => super[key]())(),

							super.foo` + "``" + `,
							super[key]` + "``" + `,
						]
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2021),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerStaticSuperES2016NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Derived extends Base {
					static test = key => {
						return [
							super.foo,
							super[key],
							([super.foo] = [0]),
							([super[key]] = [0]),

							(super.foo = 1),
							(super[key] = 1),
							(super.foo += 2),
							(super[key] += 2),

							++super.foo,
							++super[key],
							super.foo++,
							super[key]++,

							super.foo.name,
							super[key].name,
							super.foo?.name,
							super[key]?.name,

							super.foo(1, 2),
							super[key](1, 2),
							super.foo?.(1, 2),
							super[key]?.(1, 2),

							(() => super.foo)(),
							(() => super[key])(),
							(() => super.foo())(),
							(() => super[key]())(),

							super.foo` + "``" + `,
							super[key]` + "``" + `,
						]
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerAsyncArrowSuperES2016(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export { default as foo1 } from "./foo1"
				export { default as foo2 } from "./foo2"
				export { default as foo3 } from "./foo3"
				export { default as foo4 } from "./foo4"
				export { default as bar1 } from "./bar1"
				export { default as bar2 } from "./bar2"
				export { default as bar3 } from "./bar3"
				export { default as bar4 } from "./bar4"
				export { default as baz1 } from "./baz1"
				export { default as baz2 } from "./baz2"
				import "./outer"
			`,
			"/foo1.js": `export default class extends x { foo1() { return async () => super.foo('foo1') } }`,
			"/foo2.js": `export default class extends x { foo2() { return async () => () => super.foo('foo2') } }`,
			"/foo3.js": `export default class extends x { foo3() { return () => async () => super.foo('foo3') } }`,
			"/foo4.js": `export default class extends x { foo4() { return async () => async () => super.foo('foo4') } }`,
			"/bar1.js": `export default class extends x { bar1 = async () => super.foo('bar1') }`,
			"/bar2.js": `export default class extends x { bar2 = async () => () => super.foo('bar2') }`,
			"/bar3.js": `export default class extends x { bar3 = () => async () => super.foo('bar3') }`,
			"/bar4.js": `export default class extends x { bar4 = async () => async () => super.foo('bar4') }`,
			"/baz1.js": `export default class extends x { async baz1() { return () => super.foo('baz1') } }`,
			"/baz2.js": `export default class extends x { async baz2() { return () => () => super.foo('baz2') } }`,
			"/outer.js": `
				// Helper functions for "super" shouldn't be inserted into this outer function
				export default (async function () {
					class y extends z {
						foo = async () => super.foo()
					}
					await new y().foo()()
				})()
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerAsyncArrowSuperSetterES2016(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export { default as foo1 } from "./foo1"
				export { default as foo2 } from "./foo2"
				export { default as foo3 } from "./foo3"
				export { default as foo4 } from "./foo4"
				export { default as bar1 } from "./bar1"
				export { default as bar2 } from "./bar2"
				export { default as bar3 } from "./bar3"
				export { default as bar4 } from "./bar4"
				export { default as baz1 } from "./baz1"
				export { default as baz2 } from "./baz2"
				import "./outer"
			`,
			"/foo1.js": `export default class extends x { foo1() { return async () => super.foo = 'foo1' } }`,
			"/foo2.js": `export default class extends x { foo2() { return async () => () => super.foo = 'foo2' } }`,
			"/foo3.js": `export default class extends x { foo3() { return () => async () => super.foo = 'foo3' } }`,
			"/foo4.js": `export default class extends x { foo4() { return async () => async () => super.foo = 'foo4' } }`,
			"/bar1.js": `export default class extends x { bar1 = async () => super.foo = 'bar1' }`,
			"/bar2.js": `export default class extends x { bar2 = async () => () => super.foo = 'bar2' }`,
			"/bar3.js": `export default class extends x { bar3 = () => async () => super.foo = 'bar3' }`,
			"/bar4.js": `export default class extends x { bar4 = async () => async () => super.foo = 'bar4' }`,
			"/baz1.js": `export default class extends x { async baz1() { return () => super.foo = 'baz1' } }`,
			"/baz2.js": `export default class extends x { async baz2() { return () => () => super.foo = 'baz2' } }`,
			"/outer.js": `
				// Helper functions for "super" shouldn't be inserted into this outer function
				export default (async function () {
					class y extends z {
						foo = async () => super.foo = 'foo'
					}
					await new y().foo()()
				})()
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerStaticAsyncArrowSuperES2016(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export { default as foo1 } from "./foo1"
				export { default as foo2 } from "./foo2"
				export { default as foo3 } from "./foo3"
				export { default as foo4 } from "./foo4"
				export { default as bar1 } from "./bar1"
				export { default as bar2 } from "./bar2"
				export { default as bar3 } from "./bar3"
				export { default as bar4 } from "./bar4"
				export { default as baz1 } from "./baz1"
				export { default as baz2 } from "./baz2"
				import "./outer"
			`,
			"/foo1.js": `export default class extends x { static foo1() { return async () => super.foo('foo1') } }`,
			"/foo2.js": `export default class extends x { static foo2() { return async () => () => super.foo('foo2') } }`,
			"/foo3.js": `export default class extends x { static foo3() { return () => async () => super.foo('foo3') } }`,
			"/foo4.js": `export default class extends x { static foo4() { return async () => async () => super.foo('foo4') } }`,
			"/bar1.js": `export default class extends x { static bar1 = async () => super.foo('bar1') }`,
			"/bar2.js": `export default class extends x { static bar2 = async () => () => super.foo('bar2') }`,
			"/bar3.js": `export default class extends x { static bar3 = () => async () => super.foo('bar3') }`,
			"/bar4.js": `export default class extends x { static bar4 = async () => async () => super.foo('bar4') }`,
			"/baz1.js": `export default class extends x { static async baz1() { return () => super.foo('baz1') } }`,
			"/baz2.js": `export default class extends x { static async baz2() { return () => () => super.foo('baz2') } }`,
			"/outer.js": `
				// Helper functions for "super" shouldn't be inserted into this outer function
				export default (async function () {
					class y extends z {
						static foo = async () => super.foo()
					}
					await y.foo()()
				})()
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerStaticAsyncArrowSuperSetterES2016(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export { default as foo1 } from "./foo1"
				export { default as foo2 } from "./foo2"
				export { default as foo3 } from "./foo3"
				export { default as foo4 } from "./foo4"
				export { default as bar1 } from "./bar1"
				export { default as bar2 } from "./bar2"
				export { default as bar3 } from "./bar3"
				export { default as bar4 } from "./bar4"
				export { default as baz1 } from "./baz1"
				export { default as baz2 } from "./baz2"
				import "./outer"
			`,
			"/foo1.js": `export default class extends x { static foo1() { return async () => super.foo = 'foo1' } }`,
			"/foo2.js": `export default class extends x { static foo2() { return async () => () => super.foo = 'foo2' } }`,
			"/foo3.js": `export default class extends x { static foo3() { return () => async () => super.foo = 'foo3' } }`,
			"/foo4.js": `export default class extends x { static foo4() { return async () => async () => super.foo = 'foo4' } }`,
			"/bar1.js": `export default class extends x { static bar1 = async () => super.foo = 'bar1' }`,
			"/bar2.js": `export default class extends x { static bar2 = async () => () => super.foo = 'bar2' }`,
			"/bar3.js": `export default class extends x { static bar3 = () => async () => super.foo = 'bar3' }`,
			"/bar4.js": `export default class extends x { static bar4 = async () => async () => super.foo = 'bar4' }`,
			"/baz1.js": `export default class extends x { static async baz1() { return () => super.foo = 'baz1' } }`,
			"/baz2.js": `export default class extends x { static async baz2() { return () => () => super.foo = 'baz2' } }`,
			"/outer.js": `
				// Helper functions for "super" shouldn't be inserted into this outer function
				export default (async function () {
					class y extends z {
						static foo = async () => super.foo = 'foo'
					}
					await y.foo()()
				})()
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2016),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateSuperES2022(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export { default as foo1 } from "./foo1"
				export { default as foo2 } from "./foo2"
				export { default as foo3 } from "./foo3"
				export { default as foo4 } from "./foo4"
				export { default as foo5 } from "./foo5"
				export { default as foo6 } from "./foo6"
				export { default as foo7 } from "./foo7"
				export { default as foo8 } from "./foo8"
			`,
			"/foo1.js": `export default class extends x { #foo() { super.foo() } }`,
			"/foo2.js": `export default class extends x { #foo() { super.foo++ } }`,
			"/foo3.js": `export default class extends x { static #foo() { super.foo() } }`,
			"/foo4.js": `export default class extends x { static #foo() { super.foo++ } }`,
			"/foo5.js": `export default class extends x { #foo = () => { super.foo() } }`,
			"/foo6.js": `export default class extends x { #foo = () => { super.foo++ } }`,
			"/foo7.js": `export default class extends x { static #foo = () => { super.foo() } }`,
			"/foo8.js": `export default class extends x { static #foo = () => { super.foo++ } }`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2022),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerPrivateSuperES2021(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export { default as foo1 } from "./foo1"
				export { default as foo2 } from "./foo2"
				export { default as foo3 } from "./foo3"
				export { default as foo4 } from "./foo4"
				export { default as foo5 } from "./foo5"
				export { default as foo6 } from "./foo6"
				export { default as foo7 } from "./foo7"
				export { default as foo8 } from "./foo8"
			`,
			"/foo1.js": `export default class extends x { #foo() { super.foo() } }`,
			"/foo2.js": `export default class extends x { #foo() { super.foo++ } }`,
			"/foo3.js": `export default class extends x { static #foo() { super.foo() } }`,
			"/foo4.js": `export default class extends x { static #foo() { super.foo++ } }`,
			"/foo5.js": `export default class extends x { #foo = () => { super.foo() } }`,
			"/foo6.js": `export default class extends x { #foo = () => { super.foo++ } }`,
			"/foo7.js": `export default class extends x { static #foo = () => { super.foo() } }`,
			"/foo8.js": `export default class extends x { static #foo = () => { super.foo++ } }`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2021),
			AbsOutputFile:         "/out.js",
		},
	})
}

// https://github.com/evanw/esbuild/issues/2158
func TestLowerPrivateSuperStaticBundleIssue2158(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class Foo extends Object {
					static FOO;
					constructor() {
						super();
					}
					#foo;
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestLowerClassField2020NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#foo = 123
					#bar
					foo = 123
					bar
					static #s_foo = 123
					static #s_bar
					static s_foo = 123
					static s_bar
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerClassFieldNextNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#foo = 123
					#bar
					foo = 123
					bar
					static #s_foo = 123
					static #s_bar
					static s_foo = 123
					static s_bar
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			AbsOutputFile: "/out.js",
		},
	})
}

func TestTSLowerClassField2020NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				class Foo {
					#foo = 123
					#bar
					foo = 123
					bar
					static #s_foo = 123
					static #s_bar
					static s_foo = 123
					static s_bar
				}
			`,
			"/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": false
				}
			}`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestTSLowerClassPrivateFieldNextNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				class Foo {
					#foo = 123
					#bar
					foo = 123
					bar
					static #s_foo = 123
					static #s_bar
					static s_foo = 123
					static s_bar
				}
			`,
			"/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": false
				}
			}`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			AbsOutputFile: "/out.js",
		},
	})
}

func TestLowerClassFieldStrictTsconfigJson2020(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import loose from './loose'
				import strict from './strict'
				console.log(loose, strict)
			`,
			"/loose/index.js": `
				export default class {
					foo
				}
			`,
			"/loose/tsconfig.json": `
				{
					"compilerOptions": {
						"useDefineForClassFields": false
					}
				}
			`,
			"/strict/index.js": `
				export default class {
					foo
				}
			`,
			"/strict/tsconfig.json": `
				{
					"compilerOptions": {
						"useDefineForClassFields": true
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestTSLowerClassFieldStrictTsconfigJson2020(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import loose from './loose'
				import strict from './strict'
				console.log(loose, strict)
			`,
			"/loose/index.ts": `
				export default class {
					foo
				}
			`,
			"/loose/tsconfig.json": `
				{
					"compilerOptions": {
						"useDefineForClassFields": false
					}
				}
			`,
			"/strict/index.ts": `
				export default class {
					foo
				}
			`,
			"/strict/tsconfig.json": `
				{
					"compilerOptions": {
						"useDefineForClassFields": true
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2020),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestTSLowerObjectRest2017NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				const { ...local_const } = {};
				let { ...local_let } = {};
				var { ...local_var } = {};
				let arrow_fn = ({ ...x }) => { };
				let fn_expr = function ({ ...x } = default_value) {};
				let class_expr = class { method(x, ...[y, { ...z }]) {} };

				function fn_stmt({ a = b(), ...x }, { c = d(), ...y }) {}
				class class_stmt { method({ ...x }) {} }
				namespace ns { export let { ...x } = {} }
				try { } catch ({ ...catch_clause }) {}

				for (const { ...for_in_const } in { abc }) {}
				for (let { ...for_in_let } in { abc }) {}
				for (var { ...for_in_var } in { abc }) ;
				for (const { ...for_of_const } of [{}]) ;
				for (let { ...for_of_let } of [{}]) x()
				for (var { ...for_of_var } of [{}]) x()
				for (const { ...for_const } = {}; x; x = null) {}
				for (let { ...for_let } = {}; x; x = null) {}
				for (var { ...for_var } = {}; x; x = null) {}
				for ({ ...x } in { abc }) {}
				for ({ ...x } of [{}]) {}
				for ({ ...x } = {}; x; x = null) {}

				({ ...assign } = {});
				({ obj_method({ ...x }) {} });

				// Check for used return values
				({ ...x } = x);
				for ({ ...x } = x; 0; ) ;
				console.log({ ...x } = x);
				console.log({ x, ...xx } = { x });
				console.log({ x: { ...xx } } = { x });
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			UnsupportedJSFeatures: es(2017),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestTSLowerObjectRest2018NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				const { ...local_const } = {};
				let { ...local_let } = {};
				var { ...local_var } = {};
				let arrow_fn = ({ ...x }) => { };
				let fn_expr = function ({ ...x } = default_value) {};
				let class_expr = class { method(x, ...[y, { ...z }]) {} };

				function fn_stmt({ a = b(), ...x }, { c = d(), ...y }) {}
				class class_stmt { method({ ...x }) {} }
				namespace ns { export let { ...x } = {} }
				try { } catch ({ ...catch_clause }) {}

				for (const { ...for_in_const } in { abc }) {}
				for (let { ...for_in_let } in { abc }) {}
				for (var { ...for_in_var } in { abc }) ;
				for (const { ...for_of_const } of [{}]) ;
				for (let { ...for_of_let } of [{}]) x()
				for (var { ...for_of_var } of [{}]) x()
				for (const { ...for_const } = {}; x; x = null) {}
				for (let { ...for_let } = {}; x; x = null) {}
				for (var { ...for_var } = {}; x; x = null) {}
				for ({ ...x } in { abc }) {}
				for ({ ...x } of [{}]) {}
				for ({ ...x } = {}; x; x = null) {}

				({ ...assign } = {});
				({ obj_method({ ...x }) {} });

				// Check for used return values
				({ ...x } = x);
				for ({ ...x } = x; 0; ) ;
				console.log({ ...x } = x);
				console.log({ x, ...xx } = { x });
				console.log({ x: { ...xx } } = { x });
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			UnsupportedJSFeatures: es(2018),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestClassSuperThisIssue242NoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				export class A {}

				export class B extends A {
					#e: string
					constructor(c: { d: any }) {
						super()
						this.#e = c.d ?? 'test'
					}
					f() {
						return this.#e
					}
				}
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			UnsupportedJSFeatures: es(2019),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerExportStarAsNameCollisionNoBundle(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export * as ns from 'path'
				let ns = 123
				export {ns as sn}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			UnsupportedJSFeatures: es(2019),
			AbsOutputFile:         "/out.js",
		},
	})
}

func TestLowerExportStarAsNameCollision(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import * as test from './nested'
				console.log(test.foo, test.oof)
				export * as ns from 'path1'
				let ns = 123
				export {ns as sn}
			`,
			"/nested.js": `
				export * as foo from 'path2'
				let foo = 123
				export {foo as oof}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			UnsupportedJSFeatures: es(2019),
			AbsOutputFile:         "/out.js",
			ExternalSettings: config.ExternalSettings{
				PreResolve: config.ExternalMatchers{Exact: map[string]bool{
					"path1": true,
					"path2": true,
				}},
			},
		},
	})
}

func TestLowerStrictModeSyntax(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import './for-in'
			`,
			"/for-in.js": `
				if (test)
					for (var a = b in {}) ;
				for (var x = y in {}) ;
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			OutputFormat:  config.FormatESModule,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestLowerForbidStrictModeSyntax(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import './with'
				import './delete-1'
				import './delete-2'
				import './delete-3'
			`,
			"/with.js": `
				with (x) y
			`,
			"/delete-1.js": `
				delete x
			`,
			"/delete-2.js": `
				delete (y)
			`,
			"/delete-3.js": `
				delete (1 ? z : z)
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			OutputFormat:  config.FormatESModule,
			AbsOutputFile: "/out.js",
		},
		expectedScanLog: `delete-1.js: ERROR: Delete of a bare identifier cannot be used with the "esm" output format due to strict mode
delete-2.js: ERROR: Delete of a bare identifier cannot be used with the "esm" output format due to strict mode
with.js: ERROR: With statements cannot be used with the "esm" output format due to strict mode
`,
	})
}

func TestLowerPrivateClassFieldOrder(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#foo = 123 // This must be set before "bar" is initialized
					bar = this.#foo
				}
				console.log(new Foo().bar === 123)
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateField,
		},
	})
}

func TestLowerPrivateClassMethodOrder(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					bar = this.#foo()
					#foo() { return 123 } // This must be set before "bar" is initialized
				}
				console.log(new Foo().bar === 123)
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateMethod,
		},
	})
}

func TestLowerPrivateClassAccessorOrder(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					bar = this.#foo
					get #foo() { return 123 } // This must be set before "bar" is initialized
				}
				console.log(new Foo().bar === 123)
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateAccessor,
		},
	})
}

func TestLowerPrivateClassStaticFieldOrder(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					static #foo = 123 // This must be set before "bar" is initialized
					static bar = Foo.#foo
				}
				console.log(Foo.bar === 123)

				class FooThis {
					static #foo = 123 // This must be set before "bar" is initialized
					static bar = this.#foo
				}
				console.log(FooThis.bar === 123)
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateStaticField,
		},
	})
}

func TestLowerPrivateClassStaticMethodOrder(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					static bar = Foo.#foo()
					static #foo() { return 123 } // This must be set before "bar" is initialized
				}
				console.log(Foo.bar === 123)

				class FooThis {
					static bar = this.#foo()
					static #foo() { return 123 } // This must be set before "bar" is initialized
				}
				console.log(FooThis.bar === 123)
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateStaticMethod,
		},
	})
}

func TestLowerPrivateClassStaticAccessorOrder(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					static bar = Foo.#foo
					static get #foo() { return 123 } // This must be set before "bar" is initialized
				}
				console.log(Foo.bar === 123)

				class FooThis {
					static bar = this.#foo
					static get #foo() { return 123 } // This must be set before "bar" is initialized
				}
				console.log(FooThis.bar === 123)
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateStaticAccessor,
		},
	})
}

func TestLowerPrivateClassBrandCheckUnsupported(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#foo
					#bar
					baz() {
						return [
							this.#foo,
							this.#bar,
							#foo in this,
						]
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateBrandCheck,
		},
	})
}

func TestLowerPrivateClassBrandCheckSupported(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class Foo {
					#foo
					#bar
					baz() {
						return [
							this.#foo,
							this.#bar,
							#foo in this,
						]
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModePassThrough,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestLowerTemplateObject(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				x = () => [
					tag` + "`x`" + `,
					tag` + "`\\xFF`" + `,
					tag` + "`\\x`" + `,
					tag` + "`\\u`" + `,
				]
				y = () => [
					tag` + "`x${y}z`" + `,
					tag` + "`\\xFF${y}z`" + `,
					tag` + "`x${y}\\z`" + `,
					tag` + "`x${y}\\u`" + `,
				]
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.TemplateLiteral,
		},
	})
}

// See https://github.com/evanw/esbuild/issues/1424 for more information
func TestLowerPrivateClassFieldStaticIssue1424(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class T {
					#a() { return 'a'; }
					#b() { return 'b'; }
					static c;
					d() { console.log(this.#a()); }
				}
				new T().d();
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.ClassPrivateMethod,
		},
	})
}

// See https://github.com/evanw/esbuild/issues/1493 for more information
func TestLowerNullishCoalescingAssignmentIssue1493(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export class A {
					#a;
					f() {
						this.#a ??= 1;
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: compat.LogicalAssignment,
		},
	})
}

func TestStaticClassBlockESNext(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class A {
					static {}
					static {
						this.thisField++
						A.classField++
						super.superField = super.superField + 1
						super.superField++
					}
				}
				let B = class {
					static {}
					static {
						this.thisField++
						super.superField = super.superField + 1
						super.superField++
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestStaticClassBlockES2021(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				class A {
					static {}
					static {
						this.thisField++
						A.classField++
						super.superField = super.superField + 1
						super.superField++
					}
				}
				let B = class {
					static {}
					static {
						this.thisField++
						super.superField = super.superField + 1
						super.superField++
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: es(2021),
		},
	})
}

func TestLowerRegExpNameCollision(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export function foo(RegExp) {
					return new RegExp(/./d, 'd')
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: es(2021),
		},
	})
}

func TestLowerForAwait2017(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export default [
					async () => { for await (x of y) z(x) },
					async () => { for await (x.y of y) z(x) },
					async () => { for await (let x of y) z(x) },
					async () => { for await (const x of y) z(x) },
					async () => { label: for await (const x of y) break label },
					async () => { label: for await (const x of y) continue label },
				]
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: es(2017),
		},
	})
}

func TestLowerForAwait2015(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				export default [
					async () => { for await (x of y) z(x) },
					async () => { for await (x.y of y) z(x) },
					async () => { for await (let x of y) z(x) },
					async () => { for await (const x of y) z(x) },
					async () => { label: for await (const x of y) break label },
					async () => { label: for await (const x of y) continue label },
				]
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputFile:         "/out.js",
			UnsupportedJSFeatures: es(2015),
		},
	})
}

func TestLowerNestedFunctionDirectEval(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/1.js": "if (foo) { function x() {} }",
			"/2.js": "if (foo) { function x() {} eval('') }",
			"/3.js": "if (foo) { function x() {} if (bar) { eval('') } }",
			"/4.js": "if (foo) { eval(''); function x() {} }",
			"/5.js": "'use strict'; if (foo) { function x() {} }",
			"/6.js": "'use strict'; if (foo) { function x() {} eval('') }",
			"/7.js": "'use strict'; if (foo) { function x() {} if (bar) { eval('') } }",
			"/8.js": "'use strict'; if (foo) { eval(''); function x() {} }",
		},
		entryPaths: []string{
			"/1.js",
			"/2.js",
			"/3.js",
			"/4.js",
			"/5.js",
			"/6.js",
			"/7.js",
			"/8.js",
		},
		options: config.Options{
			Mode:         config.ModePassThrough,
			AbsOutputDir: "/out",
		},
	})
}

func TestJavaScriptDecoratorsESNext(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				@x.y()
				@(new y.x)
				export default class Foo {
					@x @y mUndef
					@x @y mDef = 1
					@x @y method() { return new Foo }
					@x @y static sUndef
					@x @y static sDef = new Foo
					@x @y static sMethod() { return new Foo }
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModePassThrough,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestJavaScriptAutoAccessorESNext(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/js-define.js": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
			`,
			"/ts-define/ts-define.ts": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
				class Normal { accessor a = b; c = d }
				class Private { accessor #a = b; c = d }
				class StaticNormal { static accessor a = b; static c = d }
				class StaticPrivate { static accessor #a = b; static c = d }
			`,
			"/ts-define/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": true,
				},
			}`,
			"/ts-assign/ts-assign.ts": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
				class Normal { accessor a = b; c = d }
				class Private { accessor #a = b; c = d }
				class StaticNormal { static accessor a = b; static c = d }
				class StaticPrivate { static accessor #a = b; static c = d }
			`,
			"/ts-assign/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": false,
				},
			}`,
		},
		entryPaths: []string{
			"/js-define.js",
			"/ts-define/ts-define.ts",
			"/ts-assign/ts-assign.ts",
		},
		options: config.Options{
			Mode:         config.ModePassThrough,
			AbsOutputDir: "/out",
		},
	})
}

func TestJavaScriptAutoAccessorES2022(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/js-define.js": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
			`,
			"/ts-define/ts-define.ts": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
				class Normal { accessor a = b; c = d }
				class Private { accessor #a = b; c = d }
				class StaticNormal { static accessor a = b; static c = d }
				class StaticPrivate { static accessor #a = b; static c = d }
			`,
			"/ts-define/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": true,
				},
			}`,
			"/ts-assign/ts-assign.ts": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
				class Normal { accessor a = b; c = d }
				class Private { accessor #a = b; c = d }
				class StaticNormal { static accessor a = b; static c = d }
				class StaticPrivate { static accessor #a = b; static c = d }
			`,
			"/ts-assign/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": false,
				},
			}`,
		},
		entryPaths: []string{
			"/js-define.js",
			"/ts-define/ts-define.ts",
			"/ts-assign/ts-assign.ts",
		},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: es(2022),
		},
	})
}

func TestJavaScriptAutoAccessorES2021(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/js-define.js": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
			`,
			"/ts-define/ts-define.ts": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
				class Normal { accessor a = b; c = d }
				class Private { accessor #a = b; c = d }
				class StaticNormal { static accessor a = b; static c = d }
				class StaticPrivate { static accessor #a = b; static c = d }
			`,
			"/ts-define/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": true,
				},
			}`,
			"/ts-assign/ts-assign.ts": `
				class Foo {
					accessor one = 1
					accessor #two = 2
					accessor [three()] = 3

					static accessor four = 4
					static accessor #five = 5
					static accessor [six()] = 6
				}
				class Normal { accessor a = b; c = d }
				class Private { accessor #a = b; c = d }
				class StaticNormal { static accessor a = b; static c = d }
				class StaticPrivate { static accessor #a = b; static c = d }
			`,
			"/ts-assign/tsconfig.json": `{
				"compilerOptions": {
					"useDefineForClassFields": false,
				},
			}`,
		},
		entryPaths: []string{
			"/js-define.js",
			"/ts-define/ts-define.ts",
			"/ts-assign/ts-assign.ts",
		},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: es(2021),
		},
	})
}

func TestLowerUsing(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				using a = b
				await using c = d
				if (nested) {
					using x = 1
					await using y = 2
				}

				function foo() {
					using a = b
					if (nested) {
						using x = 1
					}
				}

				async function bar() {
					using a = b
					await using c = d
					if (nested) {
						using x = 1
						await using y = 2
					}
				}
			`,
			"/loops.js": `
				for (using a of b) c(() => a)
				for (await using d of e) f(() => d)
				for await (using g of h) i(() => g)
				for await (await using j of k) l(() => j)

				if (nested) {
					for (using a of b) c(() => a)
					for (await using d of e) f(() => d)
					for await (using g of h) i(() => g)
					for await (await using j of k) l(() => j)
				}

				function foo() {
					for (using a of b) c(() => a)
				}

				async function bar() {
					for (using a of b) c(() => a)
					for (await using d of e) f(() => d)
					for await (using g of h) i(() => g)
					for await (await using j of k) l(() => j)
				}
			`,
			"/switch.js": `
				using x = y
				switch (foo) {
					case 0: using c = d
					default: using e = f
				}
				switch (foo) {
					case 0: await using c = d
					default: using e = f
				}

				async function foo() {
					using x = y
					switch (foo) {
						case 0: using c = d
						default: using e = f
					}
					switch (foo) {
						case 0: await using c = d
						default: using e = f
					}
				}
			`,
		},
		entryPaths: []string{
			"/entry.js",
			"/loops.js",
			"/switch.js",
		},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.Using,
		},
	})
}

func TestLowerUsingUnsupportedAsync(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				function foo() {
					using a = b
					if (nested) {
						using x = 1
					}
				}

				async function bar() {
					using a = b
					await using c = d
					if (nested) {
						using x = 1
						await using y = 2
					}
				}
			`,
			"/loops.js": `
				for (using a of b) c(() => a)

				if (nested) {
					for (using a of b) c(() => a)
				}

				function foo() {
					for (using a of b) c(() => a)
				}

				async function bar() {
					for (using a of b) c(() => a)
					for (await using d of e) f(() => d)
				}
			`,
			"/switch.js": `
				using x = y
				switch (foo) {
					case 0: using c = d
					default: using e = f
				}

				async function foo() {
					using x = y
					switch (foo) {
						case 0: using c = d
						default: using e = f
					}
					switch (foo) {
						case 0: await using c = d
						default: using e = f
					}
				}
			`,
		},
		entryPaths: []string{
			"/entry.js",
			"/loops.js",
			"/switch.js",
		},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.AsyncAwait | compat.TopLevelAwait,
		},
	})
}

func TestLowerUsingUnsupportedUsingAndAsync(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				function foo() {
					using a = b
					if (nested) {
						using x = 1
					}
				}

				async function bar() {
					using a = b
					await using c = d
					if (nested) {
						using x = 1
						await using y = 2
					}
				}
			`,
			"/loops.js": `
				for (using a of b) c(() => a)

				if (nested) {
					for (using a of b) c(() => a)
				}

				function foo() {
					for (using a of b) c(() => a)
				}

				async function bar() {
					for (using a of b) c(() => a)
					for (await using d of e) f(() => d)
				}
			`,
			"/switch.js": `
				using x = y
				switch (foo) {
					case 0: using c = d
					default: using e = f
				}

				async function foo() {
					using x = y
					switch (foo) {
						case 0: using c = d
						default: using e = f
					}
					switch (foo) {
						case 0: await using c = d
						default: using e = f
					}
				}
			`,
		},
		entryPaths: []string{
			"/entry.js",
			"/loops.js",
			"/switch.js",
		},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.Using | compat.AsyncAwait | compat.TopLevelAwait,
		},
	})
}

func TestLowerUsingHoisting(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/hoist-use-strict.js": `
				"use strict"
				using a = b
				function foo() {
					"use strict"
					using a = b
				}
			`,
			"/hoist-directive.js": `
				"use wtf"
				using a = b
				function foo() {
					"use wtf"
					using a = b
				}
			`,
			"/hoist-import.js": `
				using a = b
				import "./foo"
				using c = d
			`,
			"/hoist-export-star.js": `
				using a = b
				export * from './foo'
				using c = d
			`,
			"/hoist-export-from.js": `
				using a = b
				export {x, y} from './foo'
				using c = d
			`,
			"/hoist-export-clause.js": `
				using a = b
				export {a, c as 'c!'}
				using c = d
			`,
			"/hoist-export-local-direct.js": `
				using a = b
				export var ac1 = [a, c], { x: [x1] } = foo
				export let a1 = a, { y: [y1] } = foo
				export const c1 = c, { z: [z1] } = foo
				var ac2 = [a, c], { x: [x2] } = foo
				let a2 = a, { y: [y2] } = foo
				const c2 = c, { z: [z2] } = foo
				using c = d
			`,
			"/hoist-export-local-indirect.js": `
				using a = b
				var ac1 = [a, c], { x: [x1] } = foo
				let a1 = a, { y: [y1] } = foo
				const c1 = c, { z: [z1] } = foo
				var ac2 = [a, c], { x: [x2] } = foo
				let a2 = a, { y: [y2] } = foo
				const c2 = c, { z: [z2] } = foo
				using c = d
				export {x1, y1, z1}
			`,
			"/hoist-export-class-direct.js": `
				using a = b
				export class Foo1 { ac = [a, c] }
				export class Bar1 { ac = [a, c, Bar1] }
				class Foo2 { ac = [a, c] }
				class Bar2 { ac = [a, c, Bar2] }
				using c = d
			`,
			"/hoist-export-class-indirect.js": `
				using a = b
				class Foo1 { ac = [a, c] }
				class Bar1 { ac = [a, c, Bar1] }
				class Foo2 { ac = [a, c] }
				class Bar2 { ac = [a, c, Bar2] }
				using c = d
				export {Foo1, Bar1}
			`,
			"/hoist-export-function-direct.js": `
				using a = b
				export function foo1() { return [a, c] }
				export function bar1() { return [a, c, bar1] }
				function foo2() { return [a, c] }
				function bar2() { return [a, c, bar2] }
				using c = d
			`,
			"/hoist-export-function-indirect.js": `
				using a = b
				function foo1() { return [a, c] }
				function bar1() { return [a, c, bar1] }
				function foo2() { return [a, c] }
				function bar2() { return [a, c, bar2] }
				using c = d
				export {foo1, bar1}
			`,
			"/hoist-export-default-class-name-unused.js": `
				using a = b
				export default class Foo {
					ac = [a, c]
				}
				using c = d
			`,
			"/hoist-export-default-class-name-used.js": `
				using a = b
				export default class Foo {
					ac = [a, c, Foo]
				}
				using c = d
			`,
			"/hoist-export-default-class-anonymous.js": `
				using a = b
				export default class {
					ac = [a, c]
				}
				using c = d
			`,
			"/hoist-export-default-function-name-unused.js": `
				using a = b
				export default function foo() {
					return [a, c]
				}
				using c = d
			`,
			"/hoist-export-default-function-name-used.js": `
				using a = b
				export default function foo() {
					return [a, c, foo]
				}
				using c = d
			`,
			"/hoist-export-default-function-anonymous.js": `
				using a = b
				export default function() {
					return [a, c]
				}
				using c = d
			`,
			"/hoist-export-default-expr.js": `
				using a = b
				export default [a, c]
				using c = d
			`,
		},
		entryPaths: []string{
			"/hoist-use-strict.js",
			"/hoist-directive.js",
			"/hoist-import.js",
			"/hoist-export-star.js",
			"/hoist-export-from.js",
			"/hoist-export-clause.js",
			"/hoist-export-local-direct.js",
			"/hoist-export-local-indirect.js",
			"/hoist-export-class-direct.js",
			"/hoist-export-class-indirect.js",
			"/hoist-export-function-direct.js",
			"/hoist-export-function-indirect.js",
			"/hoist-export-default-class-name-unused.js",
			"/hoist-export-default-class-name-used.js",
			"/hoist-export-default-class-anonymous.js",
			"/hoist-export-default-function-name-unused.js",
			"/hoist-export-default-function-name-used.js",
			"/hoist-export-default-function-anonymous.js",
			"/hoist-export-default-expr.js",
		},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.Using,
		},
	})
}

func TestLowerUsingInsideTSNamespace(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				namespace ns {
					export let a = b
					using c = d
					export let e = f
				}
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.Using,
		},
	})
}

func TestLowerAsyncGenerator(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				async function* foo() {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				}
				foo = async function* () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				}
				foo = { async *bar () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				} }
				class Foo { async *bar () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				} }
				Foo = class { async *bar () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				} }
				async function bar() {
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				}
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.AsyncGenerator,
		},
	})
}

func TestLowerAsyncGeneratorNoAwait(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.ts": `
				async function* foo() {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				}
				foo = async function* () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				}
				foo = { async *bar () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				} }
				class Foo { async *bar () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				} }
				Foo = class { async *bar () {
					yield
					yield x
					yield *x
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				} }
				async function bar() {
					await using x = await y
					for await (let x of y) {}
					for await (await using x of y) {}
				}
			`,
		},
		entryPaths: []string{"/entry.ts"},
		options: config.Options{
			Mode:                  config.ModePassThrough,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.AsyncGenerator | compat.AsyncAwait,
		},
	})
}

func TestJavaScriptDecoratorsBundleIssue3768(t *testing.T) {
	lower_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/base-instance-method.js":   `class Foo { @dec foo() { return Foo } }`,
			"/base-instance-field.js":    `class Foo { @dec foo = Foo }`,
			"/base-instance-accessor.js": `class Foo { @dec accessor foo = Foo }`,

			"/base-static-method.js":   `class Foo { @dec static foo() { return Foo } }`,
			"/base-static-field.js":    `class Foo { @dec static foo = Foo }`,
			"/base-static-accessor.js": `class Foo { @dec static accessor foo = Foo }`,

			"/derived-instance-method.js":   `class Foo extends Bar { @dec foo() { return Foo } }`,
			"/derived-instance-field.js":    `class Foo extends Bar { @dec foo = Foo }`,
			"/derived-instance-accessor.js": `class Foo extends Bar { @dec accessor foo = Foo }`,

			"/derived-static-method.js":   `class Foo extends Bar { @dec static foo() { return Foo } }`,
			"/derived-static-field.js":    `class Foo extends Bar { @dec static foo = Foo }`,
			"/derived-static-accessor.js": `class Foo extends Bar { @dec static accessor foo = Foo }`,
		},
		entryPaths: []string{"/*"},
		options: config.Options{
			Mode:                  config.ModeBundle,
			AbsOutputDir:          "/out",
			UnsupportedJSFeatures: compat.Decorators,
		},
	})
}
