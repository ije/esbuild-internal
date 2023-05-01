# Changelog

## 0.17.18

* Fix non-default JSON import error with `export {} from` ([#3070](https://github.com/evanw/esbuild/issues/3070))

    This release fixes a bug where esbuild incorrectly identified statements of the form `export { default as x } from "y" assert { type: "json" }` as a non-default import. The bug did not affect code of the form `import { default as x } from ...` (only code that used the `export` keyword).

* Fix a crash with an invalid subpath import ([#3067](https://github.com/evanw/esbuild/issues/3067))

    Previously esbuild could crash when attempting to generate a friendly error message for an invalid [subpath import](https://nodejs.org/api/packages.html#subpath-imports) (i.e. an import starting with `#`). This happened because esbuild originally only supported the `exports` field and the code for that error message was not updated when esbuild later added support for the `imports` field. This crash has been fixed.

## 0.17.17

* Fix CSS nesting transform for top-level `&` ([#3052](https://github.com/evanw/esbuild/issues/3052))

    Previously esbuild could crash with a stack overflow when lowering CSS nesting rules with a top-level `&`, such as in the code below. This happened because esbuild's CSS nesting transform didn't handle top-level `&`, causing esbuild to inline the top-level selector into itself. This release handles top-level `&` by replacing it with [the `:scope` pseudo-class](https://drafts.csswg.org/selectors-4/#the-scope-pseudo):

    ```css
    /* Original code */
    &,
    a {
      .b {
        color: red;
      }
    }

    /* New output (with --target=chrome90) */
    :is(:scope, a) .b {
      color: red;
    }
    ```

* Support `exports` in `package.json` for `extends` in `tsconfig.json` ([#3058](https://github.com/evanw/esbuild/issues/3058))

    TypeScript 5.0 added the ability to use `extends` in `tsconfig.json` to reference a path in a package whose `package.json` file contains an `exports` map that points to the correct location. This doesn't automatically work in esbuild because `tsconfig.json` affects esbuild's path resolution, so esbuild's normal path resolution logic doesn't apply.

    This release adds support for doing this by adding some additional code that attempts to resolve the `extends` path using the `exports` field. The behavior should be similar enough to esbuild's main path resolution logic to work as expected.

    Note that esbuild always treats this `extends` import as a `require()` import since that's what TypeScript appears to do. Specifically the `require` condition will be active and the `import` condition will be inactive.

* Fix watch mode with `NODE_PATH` ([#3062](https://github.com/evanw/esbuild/issues/3062))

    Node has a rarely-used feature where you can extend the set of directories that node searches for packages using the `NODE_PATH` environment variable. While esbuild supports this too, previously a bug prevented esbuild's watch mode from picking up changes to imported files that were contained directly in a `NODE_PATH` directory. You're supposed to use `NODE_PATH` for packages, but some people abuse this feature by putting files in that directory instead (e.g. `node_modules/some-file.js` instead of `node_modules/some-pkg/some-file.js`). The watch mode bug happens when you do this because esbuild first tries to read `some-file.js` as a directory and then as a file. Watch mode was incorrectly waiting for `some-file.js` to become a valid directory. This release fixes this edge case bug by changing watch mode to watch `some-file.js` as a file when this happens.

## 0.17.16

* Fix CSS nesting transform for triple-nested rules that start with a combinator ([#3046](https://github.com/evanw/esbuild/issues/3046))

    This release fixes a bug with esbuild where triple-nested CSS rules that start with a combinator were not transformed correctly for older browsers. Here's an example of such a case before and after this bug fix:

    ```css
    /* Original input */
    .a {
      color: red;
      > .b {
        color: green;
        > .c {
          color: blue;
        }
      }
    }

    /* Old output (with --target=chrome90) */
    .a {
      color: red;
    }
    .a > .b {
      color: green;
    }
    .a .b > .c {
      color: blue;
    }

    /* New output (with --target=chrome90) */
    .a {
      color: red;
    }
    .a > .b {
      color: green;
    }
    .a > .b > .c {
      color: blue;
    }
    ```

* Support `--inject` with a file loaded using the `copy` loader ([#3041](https://github.com/evanw/esbuild/issues/3041))

    This release now allows you to use `--inject` with a file that is loaded using the `copy` loader. The `copy` loader copies the imported file to the output directory verbatim and rewrites the path in the `import` statement to point to the copied output file. When used with `--inject`, this means the injected file will be copied to the output directory as-is and a bare `import` statement for that file will be inserted in any non-copy output files that esbuild generates.

    Note that since esbuild doesn't parse the contents of copied files, esbuild will not expose any of the export names as usable imports when you do this (in the way that esbuild's `--inject` feature is typically used). However, any side-effects that the injected file has will still occur.

## 0.17.15

* Allow keywords as type parameter names in mapped types ([#3033](https://github.com/evanw/esbuild/issues/3033))

    TypeScript allows type keywords to be used as parameter names in mapped types. Previously esbuild incorrectly treated this as an error. Code that does this is now supported:

    ```ts
    type Foo = 'a' | 'b' | 'c'
    type A = { [keyof in Foo]: number }
    type B = { [infer in Foo]: number }
    type C = { [readonly in Foo]: number }
    ```

* Add annotations for re-exported modules in node ([#2486](https://github.com/evanw/esbuild/issues/2486), [#3029](https://github.com/evanw/esbuild/issues/3029))

    Node lets you import named imports from a CommonJS module using ESM import syntax. However, the allowed names aren't derived from the properties of the CommonJS module. Instead they are derived from an arbitrary syntax-only analysis of the CommonJS module's JavaScript AST.

    To accommodate node doing this, esbuild's ESM-to-CommonJS conversion adds a special non-executable "annotation" for node that describes the exports that node should expose in this scenario. It takes the form `0 && (module.exports = { ... })` and comes at the end of the file (`0 && expr` means `expr` is never evaluated).

    Previously esbuild didn't do this for modules re-exported using the `export * from` syntax. Annotations for these re-exports will now be added starting with this release:

    ```js
    // Original input
    export { foo } from './foo'
    export * from './bar'

    // Old output (with --format=cjs --platform=node)
    ...
    0 && (module.exports = {
      foo
    });

    // New output (with --format=cjs --platform=node)
    ...
    0 && (module.exports = {
      foo,
      ...require("./bar")
    });
    ```

    Note that you need to specify both `--format=cjs` and `--platform=node` to get these node-specific annotations.

* Avoid printing an unnecessary space in between a number and a `.` ([#3026](https://github.com/evanw/esbuild/pull/3026))

    JavaScript typically requires a space in between a number token and a `.` token to avoid the `.` being interpreted as a decimal point instead of a member expression. However, this space is not required if the number token itself contains a decimal point, an exponent, or uses a base other than 10. This release of esbuild now avoids printing the unnecessary space in these cases:

    ```js
    // Original input
    foo(1000 .x, 0 .x, 0.1 .x, 0.0001 .x, 0xFFFF_0000_FFFF_0000 .x)

    // Old output (with --minify)
    foo(1e3 .x,0 .x,.1 .x,1e-4 .x,0xffff0000ffff0000 .x);

    // New output (with --minify)
    foo(1e3.x,0 .x,.1.x,1e-4.x,0xffff0000ffff0000.x);
    ```

* Fix server-sent events with live reload when writing to the file system root ([#3027](https://github.com/evanw/esbuild/issues/3027))

    This release fixes a bug where esbuild previously failed to emit server-sent events for live reload when `outdir` was the file system root, such as `/`. This happened because `/` is the only path on Unix that cannot have a trailing slash trimmed from it, which was fixed by improved path handling.

## 0.17.14

* Allow the TypeScript 5.0 `const` modifier in object type declarations ([#3021](https://github.com/evanw/esbuild/issues/3021))

    The new TypeScript 5.0 `const` modifier was added to esbuild in version 0.17.5, and works with classes, functions, and arrow expressions. However, support for it wasn't added to object type declarations (e.g. interfaces) due to an oversight. This release adds support for these cases, so the following TypeScript 5.0 code can now be built with esbuild:

    ```ts
    interface Foo { <const T>(): T }
    type Bar = { new <const T>(): T }
    ```

* Implement preliminary lowering for CSS nesting ([#1945](https://github.com/evanw/esbuild/issues/1945))

    Chrome has [implemented the new CSS nesting specification](https://developer.chrome.com/articles/css-nesting/) in version 112, which is currently in beta but will become stable very soon. So CSS nesting is now a part of the web platform!

    This release of esbuild can now transform nested CSS syntax into non-nested CSS syntax for older browsers. The transformation relies on the `:is()` pseudo-class in many cases, so the transformation is only guaranteed to work when targeting browsers that support `:is()` (e.g. Chrome 88+). You'll need to set esbuild's [`target`](https://esbuild.github.io/api/#target) to the browsers you intend to support to tell esbuild to do this transformation. You will get a warning if you use CSS nesting syntax with a `target` which includes older browsers that don't support `:is()`.

    The lowering transformation looks like this:

    ```css
    /* Original input */
    a.btn {
      color: #333;
      &:hover { color: #444 }
      &:active { color: #555 }
    }

    /* New output (with --target=chrome88) */
    a.btn {
      color: #333;
    }
    a.btn:hover {
      color: #444;
    }
    a.btn:active {
      color: #555;
    }
    ```

    More complex cases may generate the `:is()` pseudo-class:

    ```css
    /* Original input */
    div, p {
      .warning, .error {
        padding: 20px;
      }
    }

    /* New output (with --target=chrome88) */
    :is(div, p) :is(.warning, .error) {
      padding: 20px;
    }
    ```

    In addition, esbuild now has a special warning message for nested style rules that start with an identifier. This isn't allowed in CSS because the syntax would be ambiguous with the existing declaration syntax. The new warning message looks like this:

    ```
    ▲ [WARNING] A nested style rule cannot start with "p" because it looks like the start of a declaration [css-syntax-error]

        <stdin>:1:7:
          1 │ main { p { margin: auto } }
            │        ^
            ╵        :is(p)

      To start a nested style rule with an identifier, you need to wrap the identifier in ":is(...)" to
      prevent the rule from being parsed as a declaration.
    ```

    Keep in mind that the transformation in this release is a preliminary implementation. CSS has many features that interact in complex ways, and there may be some edge cases that don't work correctly yet.

* Minification now removes unnecessary `&` CSS nesting selectors

    This release introduces the following CSS minification optimizations:

    ```css
    /* Original input */
    a {
      font-weight: bold;
      & {
        color: blue;
      }
      & :hover {
        text-decoration: underline;
      }
    }

    /* Old output (with --minify) */
    a{font-weight:700;&{color:#00f}& :hover{text-decoration:underline}}

    /* New output (with --minify) */
    a{font-weight:700;:hover{text-decoration:underline}color:#00f}
    ```

* Minification now removes duplicates from CSS selector lists

    This release introduces the following CSS minification optimization:

    ```css
    /* Original input */
    div, div { color: red }

    /* Old output (with --minify) */
    div,div{color:red}

    /* New output (with --minify) */
    div{color:red}
    ```

## 0.17.13

* Work around an issue with `NODE_PATH` and Go's WebAssembly internals ([#3001](https://github.com/evanw/esbuild/issues/3001))

    Go's WebAssembly implementation returns `EINVAL` instead of `ENOTDIR` when using the `readdir` syscall on a file. This messes up esbuild's implementation of node's module resolution algorithm since encountering `ENOTDIR` causes esbuild to continue its search (since it's a normal condition) while other encountering other errors causes esbuild to fail with an I/O error (since it's an unexpected condition). You can encounter this issue in practice if you use node's legacy `NODE_PATH` feature to tell esbuild to resolve node modules in a custom directory that was not installed by npm. This release works around this problem by converting `EINVAL` into `ENOTDIR` for the `readdir` syscall.

* Fix a minification bug with CSS `@layer` rules that have parsing errors ([#3016](https://github.com/evanw/esbuild/issues/3016))

    CSS at-rules [require either a `{}` block or a semicolon at the end](https://www.w3.org/TR/css-syntax-3/#consume-at-rule). Omitting both of these causes esbuild to treat the rule as an unknown at-rule. Previous releases of esbuild had a bug that incorrectly removed unknown at-rules without any children during minification if the at-rule token matched an at-rule that esbuild can handle. Specifically [cssnano](https://cssnano.co/) can generate `@layer` rules with parsing errors, and empty `@layer` rules cannot be removed because they have side effects (`@layer` didn't exist when esbuild's CSS support was added, so esbuild wasn't written to handle this). This release changes esbuild to no longer discard `@layer` rules with parsing errors when minifying (the rule `@layer c` has a parsing error):

    ```css
    /* Original input */
    @layer a {
      @layer b {
        @layer c
      }
    }

    /* Old output (with --minify) */
    @layer a.b;

    /* New output (with --minify) */
    @layer a.b.c;
    ```

* Unterminated strings in CSS are no longer an error

    The CSS specification provides [rules for handling parsing errors](https://www.w3.org/TR/CSS22/syndata.html#parsing-errors). One of those rules is that user agents must close strings upon reaching the end of a line (i.e., before an unescaped line feed, carriage return or form feed character), but then drop the construct (declaration or rule) in which the string was found. For example:

    ```css
    p {
      color: green;
      font-family: 'Courier New Times
      color: red;
      color: green;
    }
    ```

    ...would be treated the same as:

    ```css
    p { color: green; color: green; }
    ```

    ...because the second declaration (from `font-family` to the semicolon after `color: red`) is invalid and is dropped.

    Previously using this CSS with esbuild failed to build due to a syntax error, even though the code can be interpreted by a browser. With this release, the code now produces a warning instead of an error, and esbuild prints the invalid CSS such that it stays invalid in the output:

    ```css
    /* esbuild's new non-minified output: */
    p {
      color: green;
      font-family: 'Courier New Times
      color: red;
      color: green;
    }
    ```

    ```css
    /* esbuild's new minified output: */
    p{font-family:'Courier New Times
    color: red;color:green}
    ```

## 0.17.12

* Fix a crash when parsing inline TypeScript decorators ([#2991](https://github.com/evanw/esbuild/issues/2991))

    Previously esbuild's TypeScript parser crashed when parsing TypeScript decorators if the definition of the decorator was inlined into the decorator itself:

    ```ts
    @(function sealed(constructor: Function) {
      Object.seal(constructor);
      Object.seal(constructor.prototype);
    })
    class Foo {}
    ```

    This crash was not noticed earlier because this edge case did not have test coverage. The crash is fixed in this release.

## 0.17.11

* Fix the `alias` feature to always prefer the longest match ([#2963](https://github.com/evanw/esbuild/issues/2963))

    It's possible to configure conflicting aliases such as `--alias:a=b` and `--alias:a/c=d`, which is ambiguous for the import path `a/c/x` (since it could map to either `b/c/x` or `d/x`). Previously esbuild would pick the first matching `alias`, which would non-deterministically pick between one of the possible matches. This release fixes esbuild to always deterministically pick the longest possible match.

* Minify calls to some global primitive constructors ([#2962](https://github.com/evanw/esbuild/issues/2962))

    With this release, esbuild's minifier now replaces calls to `Boolean`/`Number`/`String`/`BigInt` with equivalent shorter code when relevant:

    ```js
    // Original code
    console.log(
      Boolean(a ? (b | c) !== 0 : (c & d) !== 0),
      Number(e ? '1' : '2'),
      String(e ? '1' : '2'),
      BigInt(e ? 1n : 2n),
    )

    // Old output (with --minify)
    console.log(Boolean(a?(b|c)!==0:(c&d)!==0),Number(e?"1":"2"),String(e?"1":"2"),BigInt(e?1n:2n));

    // New output (with --minify)
    console.log(!!(a?b|c:c&d),+(e?"1":"2"),e?"1":"2",e?1n:2n);
    ```

* Adjust some feature compatibility tables for node ([#2940](https://github.com/evanw/esbuild/issues/2940))

    This release makes the following adjustments to esbuild's internal feature compatibility tables for node, which tell esbuild which versions of node are known to support all aspects of that feature:

    * `class-private-brand-checks`: node v16.9+ => node v16.4+ (a decrease)
    * `hashbang`: node v12.0+ => node v12.5+ (an increase)
    * `optional-chain`: node v16.9+ => node v16.1+ (a decrease)
    * `template-literal`: node v4+ => node v10+ (an increase)

    Each of these adjustments was identified by comparing against data from the `node-compat-table` package and was manually verified using old node executables downloaded from https://nodejs.org/download/release/.

## 0.17.10

* Update esbuild's handling of CSS nesting to match the latest specification changes ([#1945](https://github.com/evanw/esbuild/issues/1945))

    The syntax for the upcoming CSS nesting feature has [recently changed](https://webkit.org/blog/13813/try-css-nesting-today-in-safari-technology-preview/). The `@nest` prefix that was previously required in some cases is now gone, and nested rules no longer have to start with `&` (as long as they don't start with an identifier or function token).

    This release updates esbuild's pass-through handling of CSS nesting syntax to match the latest specification changes. So you can now use esbuild to bundle CSS containing nested rules and try them out in a browser that supports CSS nesting (which includes nightly builds of both Chrome and Safari).

    However, I'm not implementing lowering of nested CSS to non-nested CSS for older browsers yet. While the syntax has been decided, the semantics are still in flux. In particular, there is still some debate about changing the fundamental way that CSS nesting works. For example, you might think that the following CSS is equivalent to a `.outer .inner button { ... }` rule:

    ```css
    .inner button {
      .outer & {
        color: red;
      }
    }
    ```

    But instead it's actually equivalent to a `.outer :is(.inner button) { ... }` rule which unintuitively also matches the following DOM structure:

    ```html
    <div class="inner">
      <div class="outer">
        <button></button>
      </div>
    </div>
    ```

    The `:is()` behavior is preferred by browser implementers because it's more memory-efficient, but the straightforward translation into a `.outer .inner button { ... }` rule is preferred by developers used to the existing CSS preprocessing ecosystem (e.g. SASS). It seems premature to commit esbuild to specific semantics for this syntax at this time given the ongoing debate.

* Fix cross-file CSS rule deduplication involving `url()` tokens ([#2936](https://github.com/evanw/esbuild/issues/2936))

    Previously cross-file CSS rule deduplication didn't handle `url()` tokens correctly. These tokens contain references to import paths which may be internal (i.e. in the bundle) or external (i.e. not in the bundle). When comparing two `url()` tokens for equality, the underlying import paths should be compared instead of their references. This release of esbuild fixes `url()` token comparisons. One side effect is that `@font-face` rules should now be deduplicated correctly across files:

    ```css
    /* Original code */
    @import "data:text/css, \
      @import 'http://example.com/style.css'; \
      @font-face { src: url(http://example.com/font.ttf) }";
    @import "data:text/css, \
      @font-face { src: url(http://example.com/font.ttf) }";

    /* Old output (with --bundle --minify) */
    @import"http://example.com/style.css";@font-face{src:url(http://example.com/font.ttf)}@font-face{src:url(http://example.com/font.ttf)}

    /* New output (with --bundle --minify) */
    @import"http://example.com/style.css";@font-face{src:url(http://example.com/font.ttf)}
    ```

## 0.17.9

* Parse rest bindings in TypeScript types ([#2937](https://github.com/evanw/esbuild/issues/2937))

    Previously esbuild was unable to parse the following valid TypeScript code:

    ```ts
    let tuple: (...[e1, e2, ...es]: any) => any
    ```

    This release includes support for parsing code like this.

* Fix TypeScript code translation for certain computed `declare` class fields ([#2914](https://github.com/evanw/esbuild/issues/2914))

    In TypeScript, the key of a computed `declare` class field should only be preserved if there are no decorators for that field. Previously esbuild always preserved the key, but esbuild will now remove the key to match the output of the TypeScript compiler:

    ```ts
    // Original code
    declare function dec(a: any, b: any): any
    declare const removeMe: unique symbol
    declare const keepMe: unique symbol
    class X {
        declare [removeMe]: any
        @dec declare [keepMe]: any
    }

    // Old output
    var _a;
    class X {
    }
    removeMe, _a = keepMe;
    __decorateClass([
      dec
    ], X.prototype, _a, 2);

    // New output
    var _a;
    class X {
    }
    _a = keepMe;
    __decorateClass([
      dec
    ], X.prototype, _a, 2);
    ```

* Fix a crash with path resolution error generation ([#2913](https://github.com/evanw/esbuild/issues/2913))

    In certain situations, a module containing an invalid import path could previously cause esbuild to crash when it attempts to generate a more helpful error message. This crash has been fixed.

## 0.17.8

* Fix a minification bug with non-ASCII identifiers ([#2910](https://github.com/evanw/esbuild/issues/2910))

    This release fixes a bug with esbuild where non-ASCII identifiers followed by a keyword were incorrectly not separated by a space. This bug affected both the `in` and `instanceof` keywords. Here's an example of the fix:

    ```js
    // Original code
    π in a

    // Old output (with --minify --charset=utf8)
    πin a;

    // New output (with --minify --charset=utf8)
    π in a;
    ```

* Fix a regression with esbuild's WebAssembly API in version 0.17.6 ([#2911](https://github.com/evanw/esbuild/issues/2911))

    Version 0.17.6 of esbuild updated the Go toolchain to version 1.20.0. This had the unfortunate side effect of increasing the amount of stack space that esbuild uses (presumably due to some changes to Go's WebAssembly implementation) which could cause esbuild's WebAssembly-based API to crash with a stack overflow in cases where it previously didn't crash. One such case is the package `grapheme-splitter` which contains code that looks like this:

    ```js
    if (
      (0x0300 <= code && code <= 0x036F) ||
      (0x0483 <= code && code <= 0x0487) ||
      (0x0488 <= code && code <= 0x0489) ||
      (0x0591 <= code && code <= 0x05BD) ||
      // ... many hundreds of lines later ...
    ) {
      return;
    }
    ```

    This edge case involves a chain of binary operators that results in an AST over 400 nodes deep. Normally this wouldn't be a problem because Go has growable call stacks, so the call stack would just grow to be as large as needed. However, WebAssembly byte code deliberately doesn't expose the ability to manipulate the stack pointer, so Go's WebAssembly translation is forced to use the fixed-size WebAssembly call stack. So esbuild's WebAssembly implementation is vulnerable to stack overflow in cases like these.

    It's not unreasonable for this to cause a stack overflow, and for esbuild's answer to this problem to be "don't write code like this." That's how many other AST-manipulation tools handle this problem. However, it's possible to implement AST traversal using iteration instead of recursion to work around limited call stack space. This version of esbuild implements this code transformation for esbuild's JavaScript parser and printer, so esbuild's WebAssembly implementation is now able to process the `grapheme-splitter` package (at least when compiled with Go 1.20.0 and run with node's WebAssembly implementation).

## 0.17.7

* Change esbuild's parsing of TypeScript instantiation expressions to match TypeScript 4.8+ ([#2907](https://github.com/evanw/esbuild/issues/2907))

    This release updates esbuild's implementation of instantiation expression erasure to match [microsoft/TypeScript#49353](https://github.com/microsoft/TypeScript/pull/49353). The new rules are as follows (copied from TypeScript's PR description):

    > When a potential type argument list is followed by
    >
    > * a line break,
    > * an `(` token,
    > * a template literal string, or
    > * any token except `<` or `>` that isn't the start of an expression,
    >
    > we consider that construct to be a type argument list. Otherwise we consider the construct to be a `<` relational expression followed by a `>` relational expression.

* Ignore `sideEffects: false` for imported CSS files ([#1370](https://github.com/evanw/esbuild/issues/1370), [#1458](https://github.com/evanw/esbuild/pull/1458), [#2905](https://github.com/evanw/esbuild/issues/2905))

    This release ignores the `sideEffects` annotation in `package.json` for CSS files that are imported into JS files using esbuild's `css` loader. This means that these CSS files are no longer be tree-shaken.

    Importing CSS into JS causes esbuild to automatically create a CSS entry point next to the JS entry point containing the bundled CSS. Previously packages that specified some form of `"sideEffects": false` could potentially cause esbuild to consider one or more of the JS files on the import path to the CSS file to be side-effect free, which would result in esbuild removing that CSS file from the bundle. This was problematic because the removal of that CSS is outwardly observable, since all CSS is global, so it was incorrect for previous versions of esbuild to tree-shake CSS files imported into JS files.

* Add constant folding for certain additional equality cases ([#2394](https://github.com/evanw/esbuild/issues/2394), [#2895](https://github.com/evanw/esbuild/issues/2895))

    This release adds constant folding for expressions similar to the following:

    ```js
    // Original input
    console.log(
      null === 'foo',
      null === undefined,
      null == undefined,
      false === 0,
      false == 0,
      1 === true,
      1 == true,
    )

    // Old output
    console.log(
      null === "foo",
      null === void 0,
      null == void 0,
      false === 0,
      false == 0,
      1 === true,
      1 == true
    );

    // New output
    console.log(
      false,
      false,
      true,
      false,
      true,
      false,
      true
    );
    ```

## 0.17.6

* Fix a CSS parser crash on invalid CSS ([#2892](https://github.com/evanw/esbuild/issues/2892))

    Previously the following invalid CSS caused esbuild's parser to crash:

    ```css
    @media screen
    ```

    The crash was caused by trying to construct a helpful error message assuming that there was an opening `{` token, which is not the case here. This release fixes the crash.

* Inline TypeScript enums that are referenced before their declaration

    Previously esbuild inlined enums within a TypeScript file from top to bottom, which meant that references to TypeScript enum members were only inlined within the same file if they came after the enum declaration. With this release, esbuild will now inline enums even when they are referenced before they are declared:

    ```ts
    // Original input
    export const foo = () => Foo.FOO
    const enum Foo { FOO = 0 }

    // Old output (with --tree-shaking=true)
    export const foo = () => Foo.FOO;
    var Foo = /* @__PURE__ */ ((Foo2) => {
      Foo2[Foo2["FOO"] = 0] = "FOO";
      return Foo2;
    })(Foo || {});

    // New output (with --tree-shaking=true)
    export const foo = () => 0 /* FOO */;
    ```

    This makes esbuild's TypeScript output smaller and faster when processing code that does this. I noticed this issue when I ran the TypeScript compiler's source code through esbuild's bundler. Now that the TypeScript compiler is going to be bundled with esbuild in the upcoming TypeScript 5.0 release, improvements like this will also improve the TypeScript compiler itself!

* Fix esbuild installation on Arch Linux ([#2785](https://github.com/evanw/esbuild/issues/2785), [#2812](https://github.com/evanw/esbuild/issues/2812), [#2865](https://github.com/evanw/esbuild/issues/2865))

    Someone made an unofficial `esbuild` package for Linux that adds the `ESBUILD_BINARY_PATH=/usr/bin/esbuild` environment variable to the user's default environment. This breaks all npm installations of esbuild for users with this unofficial Linux package installed, which has affected many people. Most (all?) people who encounter this problem haven't even installed this unofficial package themselves; instead it was installed for them as a dependency of another Linux package. The problematic change to add the `ESBUILD_BINARY_PATH` environment variable was reverted in the latest version of this unofficial package. However, old versions of this unofficial package are still there and will be around forever. With this release, `ESBUILD_BINARY_PATH` is now ignored by esbuild's install script when it's set to the value `/usr/bin/esbuild`. This should unbreak using npm to install `esbuild` in these problematic Linux environments.

    Note: The `ESBUILD_BINARY_PATH` variable is an undocumented way to override the location of esbuild's binary when esbuild's npm package is installed, which is necessary to substitute your own locally-built esbuild binary when debugging esbuild's npm package. It's only meant for very custom situations and should absolutely not be forced on others by default, especially without their knowledge. I may remove the code in esbuild's installer that reads `ESBUILD_BINARY_PATH` in the future to prevent these kinds of issues. It will unfortunately make debugging esbuild harder. If `ESBUILD_BINARY_PATH` is ever removed, it will be done in a "breaking change" release.

## 0.17.5

* Parse `const` type parameters from TypeScript 5.0

    The TypeScript 5.0 beta announcement adds [`const` type parameters](https://devblogs.microsoft.com/typescript/announcing-typescript-5-0-beta/#const-type-parameters) to the language. You can now add the `const` modifier on a type parameter of a function, method, or class like this:

    ```ts
    type HasNames = { names: readonly string[] };
    const getNamesExactly = <const T extends HasNames>(arg: T): T["names"] => arg.names;
    const names = getNamesExactly({ names: ["Alice", "Bob", "Eve"] });
    ```

    The type of `names` in the above example is `readonly ["Alice", "Bob", "Eve"]`. Marking the type parameter as `const` behaves as if you had written `as const` at every use instead. The above code is equivalent to the following TypeScript, which was the only option before TypeScript 5.0:

    ```ts
    type HasNames = { names: readonly string[] };
    const getNamesExactly = <T extends HasNames>(arg: T): T["names"] => arg.names;
    const names = getNamesExactly({ names: ["Alice", "Bob", "Eve"] } as const);
    ```

    You can read [the announcement](https://devblogs.microsoft.com/typescript/announcing-typescript-5-0-beta/#const-type-parameters) for more information.

* Make parsing generic `async` arrow functions more strict in `.tsx` files

    Previously esbuild's TypeScript parser incorrectly accepted the following code as valid:

    ```tsx
    let fn = async <T> () => {};
    ```

    The official TypeScript parser rejects this code because it thinks it's the identifier `async` followed by a JSX element starting with `<T>`. So with this release, esbuild will now reject this syntax in `.tsx` files too. You'll now have to add a comma after the type parameter to get generic arrow functions like this to parse in `.tsx` files:

    ```tsx
    let fn = async <T,> () => {};
    ```

* Allow the `in` and `out` type parameter modifiers on class expressions

    TypeScript 4.7 added the `in` and `out` modifiers on the type parameters of classes, interfaces, and type aliases. However, while TypeScript supported them on both class expressions and class statements, previously esbuild only supported them on class statements due to an oversight. This release now allows these modifiers on class expressions too:

    ```ts
    declare let Foo: any;
    Foo = class <in T> { };
    Foo = class <out T> { };
    ```

* Update `enum` constant folding for TypeScript 5.0

    TypeScript 5.0 contains an [updated definition of what it considers a constant expression](https://github.com/microsoft/TypeScript/pull/50528):

    > An expression is considered a *constant expression* if it is
    >
    > * a number or string literal,
    > * a unary `+`, `-`, or `~` applied to a numeric constant expression,
    > * a binary `+`, `-`, `*`, `/`, `%`, `**`, `<<`, `>>`, `>>>`, `|`, `&`, `^` applied to two numeric constant expressions,
    > * a binary `+` applied to two constant expressions whereof at least one is a string,
    > * a template expression where each substitution expression is a constant expression,
    > * a parenthesized constant expression,
    > * a dotted name (e.g. `x.y.z`) that references a `const` variable with a constant expression initializer and no type annotation,
    > * a dotted name that references an enum member with an enum literal type, or
    > * a dotted name indexed by a string literal (e.g. `x.y["z"]`) that references an enum member with an enum literal type.

    This impacts esbuild's implementation of TypeScript's `const enum` feature. With this release, esbuild will now attempt to follow these new rules. For example, you can now initialize an `enum` member with a template literal expression that contains a numeric constant:

    ```ts
    // Original input
    const enum Example {
      COUNT = 100,
      ERROR = `Expected ${COUNT} items`,
    }
    console.log(
      Example.COUNT,
      Example.ERROR,
    )

    // Old output (with --tree-shaking=true)
    var Example = /* @__PURE__ */ ((Example2) => {
      Example2[Example2["COUNT"] = 100] = "COUNT";
      Example2[Example2["ERROR"] = `Expected ${100 /* COUNT */} items`] = "ERROR";
      return Example2;
    })(Example || {});
    console.log(
      100 /* COUNT */,
      Example.ERROR
    );

    // New output (with --tree-shaking=true)
    console.log(
      100 /* COUNT */,
      "Expected 100 items" /* ERROR */
    );
    ```

    These rules are not followed exactly due to esbuild's limitations. The rule about dotted references to `const` variables is not followed both because esbuild's enum processing is done in an isolated module setting and because doing so would potentially require esbuild to use a type system, which it doesn't have. For example:

    ```ts
    // The TypeScript compiler inlines this but esbuild doesn't:
    declare const x = 'foo'
    const enum Foo { X = x }
    console.log(Foo.X)
    ```

    Also, the rule that requires converting numbers to a string currently only followed for 32-bit signed integers and non-finite numbers. This is done to avoid accidentally introducing a bug if esbuild's number-to-string operation doesn't exactly match the behavior of a real JavaScript VM. Currently esbuild's number-to-string constant folding is conservative for safety.

* Forbid definite assignment assertion operators on class methods

    In TypeScript, class methods can use the `?` optional property operator but not the `!` definite assignment assertion operator (while class fields can use both):

    ```ts
    class Foo {
      // These are valid TypeScript
      a?
      b!
      x?() {}

      // This is invalid TypeScript
      y!() {}
    }
    ```

    Previously esbuild incorrectly allowed the definite assignment assertion operator with class methods. This will no longer be allowed starting with this release.

## 0.17.4

* Implement HTTP `HEAD` requests in serve mode ([#2851](https://github.com/evanw/esbuild/issues/2851))

    Previously esbuild's serve mode only responded to HTTP `GET` requests. With this release, esbuild's serve mode will also respond to HTTP `HEAD` requests, which are just like HTTP `GET` requests except that the body of the response is omitted.

* Permit top-level await in dead code branches ([#2853](https://github.com/evanw/esbuild/issues/2853))

    Adding top-level await to a file has a few consequences with esbuild:

    1. It causes esbuild to assume that the input module format is ESM, since top-level await is only syntactically valid in ESM. That prevents you from using `module` and `exports` for exports and also enables strict mode, which disables certain syntax and changes how function hoisting works (among other things).
    2. This will cause esbuild to fail the build if either top-level await isn't supported by your language target (e.g. it's not supported in ES2021) or if top-level await isn't supported by the chosen output format (e.g. it's not supported with CommonJS).
    3. Doing this will prevent you from using `require()` on this file or on any file that imports this file (even indirectly), since the `require()` function doesn't return a promise and so can't represent top-level await.

    This release relaxes these rules slightly: rules 2 and 3 will now no longer apply when esbuild has identified the code branch as dead code, such as when it's behind an `if (false)` check. This should make it possible to use esbuild to convert code into different output formats that only uses top-level await conditionally. This release does not relax rule 1. Top-level await will still cause esbuild to unconditionally consider the input module format to be ESM, even when the top-level `await` is in a dead code branch. This is necessary because whether the input format is ESM or not affects the whole file, not just the dead code branch.

* Fix entry points where the entire file name is the extension ([#2861](https://github.com/evanw/esbuild/issues/2861))

    Previously if you passed esbuild an entry point where the file extension is the entire file name, esbuild would use the parent directory name to derive the name of the output file. For example, if you passed esbuild a file `./src/.ts` then the output name would be `src.js`. This bug happened because esbuild first strips the file extension to get `./src/` and then joins the path with the working directory to get the absolute path (e.g. `join("/working/dir", "./src/")` gives `/working/dir/src`). However, the join operation also canonicalizes the path which strips the trailing `/`. Later esbuild uses the "base name" operation to extract the name of the output file. Since there is no trailing `/`, esbuild returns `"src"` as the base name instead of `""`, which causes esbuild to incorrectly include the directory name in the output file name. This release fixes this bug by deferring the stripping of the file extension until after all path manipulations have been completed. So now the file `./src/.ts` will generate an output file named `.js`.

* Support replacing property access expressions with inject

    At a high level, this change means the `inject` feature can now replace all of the same kinds of names as the `define` feature. So `inject` is basically now a more powerful version of `define`, instead of previously only being able to do some of the things that `define` could do.

    Soem background is necessary to understand this change if you aren't already familiar with the `inject` feature. The `inject` feature lets you replace references to global variable with a shim. It works like this:

    1. Put the shim in its own file
    2. Export the shim as the name of the global variable you intend to replace
    3. Pass the file to esbuild using the `inject` feature

    For example, if you inject the following file using `--inject:./injected.js`:

    ```js
    // injected.js
    let processShim = { cwd: () => '/' }
    export { processShim as process }
    ```

    Then esbuild will replace all references to `process` with the `processShim` variable, which will cause `process.cwd()` to return `'/'`. This feature is sort of abusing the ESM export alias syntax to specify the mapping of global variables to shims. But esbuild works this way because using this syntax for that purpose is convenient and terse.

    However, if you wanted to replace a property access expression, the process was more complicated and not as nice. You would have to:

    1. Put the shim in its own file
    2. Export the shim as some random name
    3. Pass the file to esbuild using the `inject` feature
    4. Use esbuild's `define` feature to map the property access expression to the random name you made in step 2

    For example, if you inject the following file using `--inject:./injected2.js --define:process.cwd=someRandomName`:

    ```js
    // injected2.js
    let cwdShim = () => '/'
    export { cwdShim as someRandomName }
    ```

    Then esbuild will replace all references to `process.cwd` with the `cwdShim` variable, which will also cause `process.cwd()` to return `'/'` (but which this time will not mess with other references to `process`, which might be desirable).

    With this release, using the inject feature to replace a property access expression is now as simple as using it to replace an identifier. You can now use JavaScript's ["arbitrary module namespace identifier names"](https://github.com/tc39/ecma262/pull/2154) feature to specify the property access expression directly using a string literal. For example, if you inject the following file using `--inject:./injected3.js`:

    ```js
    // injected3.js
    let cwdShim = () => '/'
    export { cwdShim as 'process.cwd' }
    ```

    Then esbuild will now replace all references to `process.cwd` with the `cwdShim` variable, which will also cause `process.cwd()` to return `'/'` (but which will also not mess with other references to `process`).

    In addition to inserting a shim for a global variable that doesn't exist, another use case is replacing references to static methods on global objects with cached versions to both minify them better and to make access to them potentially faster. For example:

    ```js
    // Injected file
    let cachedMin = Math.min
    let cachedMax = Math.max
    export {
      cachedMin as 'Math.min',
      cachedMax as 'Math.max',
    }

    // Original input
    function clampRGB(r, g, b) {
      return {
        r: Math.max(0, Math.min(1, r)),
        g: Math.max(0, Math.min(1, g)),
        b: Math.max(0, Math.min(1, b)),
      }
    }

    // Old output (with --minify)
    function clampRGB(a,t,m){return{r:Math.max(0,Math.min(1,a)),g:Math.max(0,Math.min(1,t)),b:Math.max(0,Math.min(1,m))}}

    // New output (with --minify)
    var a=Math.min,t=Math.max;function clampRGB(h,M,m){return{r:t(0,a(1,h)),g:t(0,a(1,M)),b:t(0,a(1,m))}}
    ```

## 0.17.3

* Fix incorrect CSS minification for certain rules ([#2838](https://github.com/evanw/esbuild/issues/2838))

    Certain rules such as `@media` could previously be minified incorrectly. Due to a typo in the duplicate rule checker, two known `@`-rules that share the same hash code were incorrectly considered to be equal. This problem was made worse by the rule hashing code considering two unknown declarations (such as CSS variables) to have the same hash code, which also isn't optimal from a performance perspective. Both of these issues have been fixed:

    ```css
    /* Original input */
    @media (prefers-color-scheme: dark) { body { --VAR-1: #000; } }
    @media (prefers-color-scheme: dark) { body { --VAR-2: #000; } }

    /* Old output (with --minify) */
    @media (prefers-color-scheme: dark){body{--VAR-2: #000}}

    /* New output (with --minify) */
    @media (prefers-color-scheme: dark){body{--VAR-1: #000}}@media (prefers-color-scheme: dark){body{--VAR-2: #000}}
    ```

## 0.17.2

* Add `onDispose` to the plugin API ([#2140](https://github.com/evanw/esbuild/issues/2140), [#2205](https://github.com/evanw/esbuild/issues/2205))

    If your plugin wants to perform some cleanup after it's no longer going to be used, you can now use the `onDispose` API to register a callback for cleanup-related tasks. For example, if a plugin starts a long-running child process then it may want to terminate that process when the plugin is discarded. Previously there was no way to do this. Here's an example:

    ```js
    let examplePlugin = {
      name: 'example',
      setup(build) {
        build.onDispose(() => {
          console.log('This plugin is no longer used')
        })
      },
    }
    ```

    These `onDispose` callbacks will be called after every `build()` call regardless of whether the build failed or not as well as after the first `dispose()` call on a given build context.

## 0.17.1

* Make it possible to cancel a build ([#2725](https://github.com/evanw/esbuild/issues/2725))

    The context object introduced in version 0.17.0 has a new `cancel()` method. You can use it to cancel a long-running build so that you can start a new one without needing to wait for the previous one to finish. When this happens, the previous build should always have at least one error and have no output files (i.e. it will be a failed build).

    Using it might look something like this:

    * JS:

        ```js
        let ctx = await esbuild.context({
          // ...
        })

        let rebuildWithTimeLimit = timeLimit => {
          let timeout = setTimeout(() => ctx.cancel(), timeLimit)
          return ctx.rebuild().finally(() => clearTimeout(timeout))
        }

        let build = await rebuildWithTimeLimit(500)
        ```

    * Go:

        ```go
        ctx, err := api.Context(api.BuildOptions{
          // ...
        })
        if err != nil {
          return
        }

        rebuildWithTimeLimit := func(timeLimit time.Duration) api.BuildResult {
          t := time.NewTimer(timeLimit)
          go func() {
            <-t.C
            ctx.Cancel()
          }()
          result := ctx.Rebuild()
          t.Stop()
          return result
        }

        build := rebuildWithTimeLimit(500 * time.Millisecond)
        ```

    This API is a quick implementation and isn't maximally efficient, so the build may continue to do some work for a little bit before stopping. For example, I have added stop points between each top-level phase of the bundler and in the main module graph traversal loop, but I haven't added fine-grained stop points within the internals of the linker. How quickly esbuild stops can be improved in future releases. This means you'll want to wait for `cancel()` and/or the previous `rebuild()` to finish (i.e. await the returned promise in JavaScript) before starting a new build, otherwise `rebuild()` will give you the just-canceled build that still hasn't ended yet. Note that `onEnd` callbacks will still be run regardless of whether or not the build was canceled.

* Fix server-sent events without `servedir` ([#2827](https://github.com/evanw/esbuild/issues/2827))

    The server-sent events for live reload were incorrectly using `servedir` to calculate the path to modified output files. This means events couldn't be sent when `servedir` wasn't specified. This release uses the internal output directory (which is always present) instead of `servedir` (which might be omitted), so live reload should now work when `servedir` is not specified.

* Custom entry point output paths now work with the `copy` loader ([#2828](https://github.com/evanw/esbuild/issues/2828))

    Entry points can optionally provide custom output paths to change the path of the generated output file. For example, `esbuild foo=abc.js bar=xyz.js --outdir=out` generates the files `out/foo.js` and `out/bar.js`. However, this previously didn't work when using the `copy` loader due to an oversight. This bug has been fixed. For example, you can now do `esbuild foo=abc.html bar=xyz.html --outdir=out --loader:.html=copy` to generate the files `out/foo.html` and `out/bar.html`.

* The JS API can now take an array of objects ([#2828](https://github.com/evanw/esbuild/issues/2828))

    Previously it was not possible to specify two entry points with the same custom output path using the JS API, although it was possible to do this with the Go API and the CLI. This will not cause a collision if both entry points use different extensions (e.g. if one uses `.js` and the other uses `.css`). You can now pass the JS API an array of objects to work around this API limitation:

    ```js
    // The previous API didn't let you specify duplicate output paths
    let result = await esbuild.build({
      entryPoints: {
        // This object literal contains a duplicate key, so one is ignored
        'dist': 'foo.js',
        'dist': 'bar.css',
      },
    })

    // You can now specify duplicate output paths as an array of objects
    let result = await esbuild.build({
      entryPoints: [
        { in: 'foo.js', out: 'dist' },
        { in: 'bar.css', out: 'dist' },
      ],
    })
    ```

## 0.17.0

**This release deliberately contains backwards-incompatible changes.** To avoid automatically picking up releases like this, you should either be pinning the exact version of `esbuild` in your `package.json` file (recommended) or be using a version range syntax that only accepts patch upgrades such as `^0.16.0` or `~0.16.0`. See npm's documentation about [semver](https://docs.npmjs.com/cli/v6/using-npm/semver/) for more information.

At a high level, the breaking changes in this release fix some long-standing issues with the design of esbuild's incremental, watch, and serve APIs. This release also introduces some exciting new features such as live reloading. In detail:

* Move everything related to incremental builds to a new `context` API ([#1037](https://github.com/evanw/esbuild/issues/1037), [#1606](https://github.com/evanw/esbuild/issues/1606), [#2280](https://github.com/evanw/esbuild/issues/2280), [#2418](https://github.com/evanw/esbuild/issues/2418))

    This change removes the `incremental` and `watch` options as well as the `serve()` method, and introduces a new `context()` method. The context method takes the same arguments as the `build()` method but only validates its arguments and does not do an initial build. Instead, builds can be triggered using the `rebuild()`, `watch()`, and `serve()` methods on the returned context object. The new context API looks like this:

    ```js
    // Create a context for incremental builds
    const context = await esbuild.context({
      entryPoints: ['app.ts'],
      bundle: true,
    })

    // Manually do an incremental build
    const result = await context.rebuild()

    // Enable watch mode
    await context.watch()

    // Enable serve mode
    await context.serve()

    // Dispose of the context
    context.dispose()
    ```

    The switch to the context API solves a major issue with the previous API which is that if the initial build fails, a promise is thrown in JavaScript which prevents you from accessing the returned result object. That prevented you from setting up long-running operations such as watch mode when the initial build contained errors. It also makes tearing down incremental builds simpler as there is now a single way to do it instead of three separate ways.

    In addition, this release also makes some subtle changes to how incremental builds work. Previously every call to `rebuild()` started a new build. If you weren't careful, then builds could actually overlap. This doesn't cause any problems with esbuild itself, but could potentially cause problems with plugins (esbuild doesn't even give you a way to identify which overlapping build a given plugin callback is running on). Overlapping builds also arguably aren't useful, or at least aren't useful enough to justify the confusion and complexity that they bring. With this release, there is now only ever a single active build per context. Calling `rebuild()` before the previous rebuild has finished now "merges" with the existing rebuild instead of starting a new build.

* Allow using `watch` and `serve` together ([#805](https://github.com/evanw/esbuild/issues/805), [#1650](https://github.com/evanw/esbuild/issues/1650), [#2576](https://github.com/evanw/esbuild/issues/2576))

    Previously it was not possible to use watch mode and serve mode together. The rationale was that watch mode is one way of automatically rebuilding your project and serve mode is another (since serve mode automatically rebuilds on every request). However, people want to combine these two features to make "live reloading" where the browser automatically reloads the page when files are changed on the file system.

    This release now allows you to use these two features together. You can only call the `watch()` and `serve()` APIs once each per context, but if you call them together on the same context then esbuild will automatically rebuild both when files on the file system are changed *and* when the server serves a request.

* Support "live reloading" through server-sent events ([#802](https://github.com/evanw/esbuild/issues/802))

    [Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) are a simple way to pass one-directional messages asynchronously from the server to the client. Serve mode now provides a `/esbuild` endpoint with an `change` event that triggers every time esbuild's output changes. So you can now implement simple "live reloading" (i.e. reloading the page when a file is edited and saved) like this:

    ```js
    new EventSource('/esbuild').addEventListener('change', () => location.reload())
    ```

    The event payload is a JSON object with the following shape:

    ```ts
    interface ChangeEvent {
      added: string[]
      removed: string[]
      updated: string[]
    }
    ```

    This JSON should also enable more complex live reloading scenarios. For example, the following code hot-swaps changed CSS `<link>` tags in place without reloading the page (but still reloads when there are other types of changes):

    ```js
    new EventSource('/esbuild').addEventListener('change', e => {
      const { added, removed, updated } = JSON.parse(e.data)
      if (!added.length && !removed.length && updated.length === 1) {
        for (const link of document.getElementsByTagName("link")) {
          const url = new URL(link.href)
          if (url.host === location.host && url.pathname === updated[0]) {
            const next = link.cloneNode()
            next.href = updated[0] + '?' + Math.random().toString(36).slice(2)
            next.onload = () => link.remove()
            link.parentNode.insertBefore(next, link.nextSibling)
            return
          }
        }
      }
      location.reload()
    })
    ```

    Implementing live reloading like this has a few known caveats:

    * These events only trigger when esbuild's output changes. They do not trigger when files unrelated to the build being watched are changed. If your HTML file references other files that esbuild doesn't know about and those files are changed, you can either manually reload the page or you can implement your own live reloading infrastructure instead of using esbuild's built-in behavior.

    * The `EventSource` API is supposed to automatically reconnect for you. However, there's a bug in Firefox that breaks this if the server is ever temporarily unreachable: https://bugzilla.mozilla.org/show_bug.cgi?id=1809332. Workarounds are to use any other browser, to manually reload the page if this happens, or to write more complicated code that manually closes and re-creates the `EventSource` object if there is a connection error. I'm hopeful that this bug will be fixed.

    * Browser vendors have decided to not implement HTTP/2 without TLS. This means that each `/esbuild` event source will take up one of your precious 6 simultaneous per-domain HTTP/1.1 connections. So if you open more than six HTTP tabs that use this live-reloading technique, you will be unable to use live reloading in some of those tabs (and other things will likely also break). The workaround is to enable HTTPS, which is now possible to do in esbuild itself (see below).

* Add built-in support for HTTPS ([#2169](https://github.com/evanw/esbuild/issues/2169))

    You can now tell esbuild's built-in development server to use HTTPS instead of HTTP. This is sometimes necessary because browser vendors have started making modern web features unavailable to HTTP websites. Previously you had to put a proxy in front of esbuild to enable HTTPS since esbuild's development server only supported HTTP. But with this release, you can now enable HTTPS with esbuild without an additional proxy.

    To enable HTTPS with esbuild:

    1. Generate a self-signed certificate. There are many ways to do this. Here's one way, assuming you have `openssl` installed:

        ```
        openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 9999 -nodes -subj /CN=127.0.0.1
        ```

    2. Add `--keyfile=key.pem` and `--certfile=cert.pem` to your esbuild development server command
    3. Click past the scary warning in your browser when you load your page

    If you have more complex needs than this, you can still put a proxy in front of esbuild and use that for HTTPS instead. Note that if you see the message "Client sent an HTTP request to an HTTPS server" when you load your page, then you are using the incorrect protocol. Replace `http://` with `https://` in your browser's URL bar.

    Keep in mind that esbuild's HTTPS support has nothing to do with security. The only reason esbuild now supports HTTPS is because browsers have made it impossible to do local development with certain modern web features without jumping through these extra hoops. *Please do not use esbuild's development server for anything that needs to be secure.* It's only intended for local development and no considerations have been made for production environments whatsoever.

* Better support copying `index.html` into the output directory ([#621](https://github.com/evanw/esbuild/issues/621), [#1771](https://github.com/evanw/esbuild/issues/1771))

    Right now esbuild only supports JavaScript and CSS as first-class content types. Previously this meant that if you were building a website with a HTML file, a JavaScript file, and a CSS file, you could use esbuild to build the JavaScript file and the CSS file into the output directory but not to copy the HTML file into the output directory. You needed a separate `cp` command for that.

    Or so I thought. It turns out that the `copy` loader added in version 0.14.44 of esbuild is sufficient to have esbuild copy the HTML file into the output directory as well. You can add something like `index.html --loader:.html=copy` and esbuild will copy `index.html` into the output directory for you. The benefits of this are a) you don't need a separate `cp` command and b) the `index.html` file will automatically be re-copied when esbuild is in watch mode and the contents of `index.html` are edited. This also goes for other non-HTML file types that you might want to copy.

    This pretty much already worked. The one thing that didn't work was that esbuild's built-in development server previously only supported implicitly loading `index.html` (e.g. loading `/about/index.html` when you visit `/about/`) when `index.html` existed on the file system. Previously esbuild didn't support implicitly loading `index.html` if it was a build result. That bug has been fixed with this release so it should now be practical to use the `copy` loader to do this.

* Fix `onEnd` not being called in serve mode ([#1384](https://github.com/evanw/esbuild/issues/1384))

    Previous releases had a bug where plugin `onEnd` callbacks weren't called when using the top-level `serve()` API. This API no longer exists and the internals have been reimplemented such that `onEnd` callbacks should now always be called at the end of every build.

* Incremental builds now write out build results differently ([#2104](https://github.com/evanw/esbuild/issues/2104))

    Previously build results were always written out after every build. However, this could cause the output directory to fill up with files from old builds if code splitting was enabled, since the file names for code splitting chunks contain content hashes and old files were not deleted.

    With this release, incremental builds in esbuild will now delete old output files from previous builds that are no longer relevant. Subsequent incremental builds will also no longer overwrite output files whose contents haven't changed since the previous incremental build.

* The `onRebuild` watch mode callback was removed ([#980](https://github.com/evanw/esbuild/issues/980), [#2499](https://github.com/evanw/esbuild/issues/2499))

    Previously watch mode accepted an `onRebuild` callback which was called whenever watch mode rebuilt something. This was not great in practice because if you are running code after a build, you likely want that code to run after every build, not just after the second and subsequent builds. This release removes option to provide an `onRebuild` callback. You can create a plugin with an `onEnd` callback instead. The `onEnd` plugin API already exists, and is a way to run some code after every build.

* You can now return errors from `onEnd` ([#2625](https://github.com/evanw/esbuild/issues/2625))

    It's now possible to add additional build errors and/or warnings to the current build from within your `onEnd` callback by returning them in an array. This is identical to how the `onStart` callback already works. The evaluation of `onEnd` callbacks have been moved around a bit internally to make this possible.

    Note that the build will only fail (i.e. reject the promise) if the additional errors are returned from `onEnd`. Adding additional errors to the result object that's passed to `onEnd` won't affect esbuild's behavior at all.

* Print URLs and ports from the Go and JS APIs ([#2393](https://github.com/evanw/esbuild/issues/2393))

    Previously esbuild's CLI printed out something like this when serve mode is active:

    ```
     > Local:   http://127.0.0.1:8000/
     > Network: http://192.168.0.1:8000/
    ```

    The CLI still does this, but now the JS and Go serve mode APIs will do this too. This only happens when the log level is set to `verbose`, `debug`, or `info` but not when it's set to `warning`, `error`, or `silent`.

### Upgrade guide for existing code:

* Rebuild (a.k.a. incremental build):

    Before:

    ```js
    const result = await esbuild.build({ ...buildOptions, incremental: true });
    builds.push(result);
    for (let i = 0; i < 4; i++) builds.push(await result.rebuild());
    await result.rebuild.dispose(); // To free resources
    ```

    After:

    ```js
    const ctx = await esbuild.context(buildOptions);
    for (let i = 0; i < 5; i++) builds.push(await ctx.rebuild());
    await ctx.dispose(); // To free resources
    ```

    Previously the first build was done differently than subsequent builds. Now both the first build and subsequent builds are done using the same API.

* Serve:

    Before:

    ```js
    const serveResult = await esbuild.serve(serveOptions, buildOptions);
    ...
    serveResult.stop(); await serveResult.wait; // To free resources
    ```

    After:

    ```js
    const ctx = await esbuild.context(buildOptions);
    const serveResult = await ctx.serve(serveOptions);
    ...
    await ctx.dispose(); // To free resources
    ```

* Watch:

    Before:

    ```js
    const result = await esbuild.build({ ...buildOptions, watch: true });
    ...
    result.stop(); // To free resources
    ```

    After:

    ```js
    const ctx = await esbuild.context(buildOptions);
    await ctx.watch();
    ...
    await ctx.dispose(); // To free resources
    ```

* Watch with `onRebuild`:

    Before:

    ```js
    const onRebuild = (error, result) => {
      if (error) console.log('subsequent build:', error);
      else console.log('subsequent build:', result);
    };
    try {
      const result = await esbuild.build({ ...buildOptions, watch: { onRebuild } });
      console.log('first build:', result);
      ...
      result.stop(); // To free resources
    } catch (error) {
      console.log('first build:', error);
    }
    ```

    After:

    ```js
    const plugins = [{
      name: 'my-plugin',
      setup(build) {
        let count = 0;
        build.onEnd(result => {
          if (count++ === 0) console.log('first build:', result);
          else console.log('subsequent build:', result);
        });
      },
    }];
    const ctx = await esbuild.context({ ...buildOptions, plugins });
    await ctx.watch();
    ...
    await ctx.dispose(); // To free resources
    ```

    The `onRebuild` function has now been removed. The replacement is to make a plugin with an `onEnd` callback.

    Previously `onRebuild` did not fire for the first build (only for subsequent builds). This was usually problematic, so using `onEnd` instead of `onRebuild` is likely less error-prone. But if you need to emulate the old behavior of `onRebuild` that ignores the first build, then you'll need to manually count and ignore the first build in your plugin (as demonstrated above).

Notice how all of these API calls are now done off the new context object. You should now be able to use all three kinds of incremental builds (`rebuild`, `serve`, and `watch`) together on the same context object. Also notice how calling `dispose` on the context is now the common way to discard the context and free resources in all of these situations.

## 0.16.17

* Fix additional comment-related regressions ([#2814](https://github.com/evanw/esbuild/issues/2814))

    This release fixes more edge cases where the new comment preservation behavior that was added in 0.16.14 could introduce syntax errors. Specifically:

    ```js
    x = () => (/* comment */ {})
    for ((/* comment */ let).x of y) ;
    function *f() { yield (/* comment */class {}) }
    ```

    These cases caused esbuild to generate code with a syntax error in version 0.16.14 or above. These bugs have now been fixed.

## 0.16.16

* Fix a regression caused by comment preservation ([#2805](https://github.com/evanw/esbuild/issues/2805))

    The new comment preservation behavior that was added in 0.16.14 introduced a regression where comments in certain locations could cause esbuild to omit certain necessary parentheses in the output. The outermost parentheses were incorrectly removed for the following syntax forms, which then introduced syntax errors:

    ```js
    (/* comment */ { x: 0 }).x;
    (/* comment */ function () { })();
    (/* comment */ class { }).prototype;
    ```

    This regression has been fixed.

## 0.16.15

* Add `format` to input files in the JSON metafile data

    When `--metafile` is enabled, input files may now have an additional `format` field that indicates the export format used by this file. When present, the value will either be `cjs` for CommonJS-style exports or `esm` for ESM-style exports. This can be useful in bundle analysis.

    For example, esbuild's new [Bundle Size Analyzer](https://esbuild.github.io/analyze/) now uses this information to visualize whether ESM or CommonJS was used for each directory and file of source code (click on the CJS/ESM bar at the top).

    This information is helpful when trying to reduce the size of your bundle. Using the ESM variant of a dependency instead of the CommonJS variant always results in a faster and smaller bundle because it omits CommonJS wrappers, and also may result in better tree-shaking as it allows esbuild to perform tree-shaking at the statement level instead of the module level.

* Fix a bundling edge case with dynamic import ([#2793](https://github.com/evanw/esbuild/issues/2793))

    This release fixes a bug where esbuild's bundler could produce incorrect output. The problematic edge case involves the entry point importing itself using a dynamic `import()` expression in an imported file, like this:

    ```js
    // src/a.js
    export const A = 42;

    // src/b.js
    export const B = async () => (await import(".")).A

    // src/index.js
    export * from "./a"
    export * from "./b"
    ```

* Remove new type syntax from type declarations in the `esbuild` package ([#2798](https://github.com/evanw/esbuild/issues/2798))

    Previously you needed to use TypeScript 4.3 or newer when using the `esbuild` package from TypeScript code due to the use of a getter in an interface in `node_modules/esbuild/lib/main.d.ts`. This release removes this newer syntax to allow people with versions of TypeScript as far back as TypeScript 3.5 to use this latest version of the `esbuild` package. Here is change that was made to esbuild's type declarations:

    ```diff
     export interface OutputFile {
       /** "text" as bytes */
       contents: Uint8Array;
       /** "contents" as text (changes automatically with "contents") */
    -  get text(): string;
    +  readonly text: string;
     }
    ```

## 0.16.14

* Preserve some comments in expressions ([#2721](https://github.com/evanw/esbuild/issues/2721))

    Various tools give semantic meaning to comments embedded inside of expressions. For example, Webpack and Vite have special "magic comments" that can be used to affect code splitting behavior:

    ```js
    import(/* webpackChunkName: "foo" */ '../foo');
    import(/* @vite-ignore */ dynamicVar);
    new Worker(/* webpackChunkName: "bar" */ new URL("../bar.ts", import.meta.url));
    new Worker(new URL('./path', import.meta.url), /* @vite-ignore */ dynamicOptions);
    ```

    Since esbuild can be used as a preprocessor for these tools (e.g. to strip TypeScript types), it can be problematic if esbuild doesn't do additional work to try to retain these comments. Previously esbuild special-cased Webpack comments in these specific locations in the AST. But Vite would now like to use similar comments, and likely other tools as well.

    So with this release, esbuild now will attempt to preserve some comments inside of expressions in more situations than before. This behavior is mainly intended to preserve these special "magic comments" that are meant for other tools to consume, although esbuild will no longer only preserve Webpack-specific comments so it should now be tool-agnostic. There is no guarantee that all such comments will be preserved (especially when `--minify-syntax` is enabled). So this change does *not* mean that esbuild is now usable as a code formatter. In particular comment preservation is more likely to happen with leading comments than with trailing comments. You should put comments that you want to be preserved *before* the relevant expression instead of after it. Also note that this change does not retain any more statement-level comments than before (i.e. comments not embedded inside of expressions). Comment preservation is not enabled when `--minify-whitespace` is enabled (which is automatically enabled when you use `--minify`).

## 0.16.13

* Publish a new bundle visualization tool

    While esbuild provides bundle metadata via the `--metafile` flag, previously esbuild left analysis of it completely up to third-party tools (well, outside of the rudimentary `--analyze` flag). However, the esbuild website now has a built-in bundle visualization tool:

    * https://esbuild.github.io/analyze/

    You can pass `--metafile` to esbuild to output bundle metadata, then upload that JSON file to this tool to visualize your bundle. This is helpful for answering questions such as:

    * Which packages are included in my bundle?
    * How did a specific file get included?
    * How small did a specific file compress to?
    * Was a specific file tree-shaken or not?

    I'm publishing this tool because I think esbuild should provide *some* answer to "how do I visualize my bundle" without requiring people to reach for third-party tools. At the moment the tool offers two types of visualizations: a radial "sunburst chart" and a linear "flame chart". They serve slightly different but overlapping use cases (e.g. the sunburst chart is more keyboard-accessible while the flame chart is easier with the mouse). This tool may continue to evolve over time.

* Fix `--metafile` and `--mangle-cache` with `--watch` ([#1357](https://github.com/evanw/esbuild/issues/1357))

    The CLI calls the Go API and then also writes out the metafile and/or mangle cache JSON files if those features are enabled. This extra step is necessary because these files are returned by the Go API as in-memory strings. However, this extra step accidentally didn't happen for all builds after the initial build when watch mode was enabled. This behavior used to work but it was broken in version 0.14.18 by the introduction of the mangle cache feature. This release fixes the combination of these features, so the metafile and mangle cache features should now work with watch mode. This behavior was only broken for the CLI, not for the JS or Go APIs.

* Add an `original` field to the metafile

    The metadata file JSON now has an additional field: each import in an input file now contains the pre-resolved path in the `original` field in addition to the post-resolved path in the `path` field. This means it's now possible to run certain additional analysis over your bundle. For example, you should be able to use this to detect when the same package subpath is represented multiple times in the bundle, either because multiple versions of a package were bundled or because a package is experiencing the [dual-package hazard](https://nodejs.org/api/packages.html#dual-package-hazard).

## 2022

All esbuild versions published in the year 2022 (versions 0.14.11 through 0.16.12) can be found in [CHANGELOG-2022.md](./CHANGELOG-2022.md).

## 2021

All esbuild versions published in the year 2021 (versions 0.8.29 through 0.14.10) can be found in [CHANGELOG-2021.md](./CHANGELOG-2021.md).

## 2020

All esbuild versions published in the year 2020 (versions 0.3.0 through 0.8.28) can be found in [CHANGELOG-2020.md](./CHANGELOG-2020.md).
