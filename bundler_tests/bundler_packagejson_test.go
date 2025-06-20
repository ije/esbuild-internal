package bundler_tests

import (
	"testing"

	"github.com/ije/esbuild-internal/config"
)

var packagejson_suite = suite{
	name: "packagejson",
}

func TestPackageJsonMain(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./custom-main.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/custom-main.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBadMain(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./does-not-exist.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonSyntaxErrorComment(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					// Single-line comment
					"a": 1
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/node_modules/demo-pkg/package.json: ERROR: JSON does not support comments
`,
	})
}

func TestPackageJsonSyntaxErrorTrailingComma(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"a": 1,
					"b": 2,
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/node_modules/demo-pkg/package.json: ERROR: JSON does not support trailing commas
`,
	})
}

func TestPackageJsonModule(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./main.esm.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.esm.js": `
				export default function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserString(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"browser": "./browser"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/browser.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapRelativeToRelative(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main",
					"browser": {
						"./main.js": "./main-browser",
						"./lib/util.js": "./lib/util-browser"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				const util = require('./lib/util')
				module.exports = function() {
					return ['main', util]
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main-browser.js": `
				const util = require('./lib/util')
				module.exports = function() {
					return ['main-browser', util]
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/lib/util.js": `
				module.exports = 'util'
			`,
			"/Users/user/project/node_modules/demo-pkg/lib/util-browser.js": `
				module.exports = 'util-browser'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapRelativeToModule(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main",
					"browser": {
						"./util.js": "util-browser"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				const util = require('./util')
				module.exports = function() {
					return ['main', util]
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/util.js": `
				module.exports = 'util'
			`,
			"/Users/user/project/node_modules/util-browser/index.js": `
				module.exports = 'util-browser'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapRelativeDisabled(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main",
					"browser": {
						"./util-node.js": false
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				const util = require('./util-node')
				module.exports = function(obj) {
					return util.inspect(obj)
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/util-node.js": `
				module.exports = require('util')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapModuleToRelative(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"browser": {
						"node-pkg": "./node-pkg-browser"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/node-pkg-browser.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				const fn = require('node-pkg')
				module.exports = function() {
					return fn()
				}
			`,
			"/Users/user/project/node_modules/node-pkg/index.js": `
				module.exports = function() {
					return 234
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapModuleToModule(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"browser": {
						"node-pkg": "node-pkg-browser"
					}
				}
			`,
			"/Users/user/project/node_modules/node-pkg-browser/index.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				const fn = require('node-pkg')
				module.exports = function() {
					return fn()
				}
			`,
			"/Users/user/project/node_modules/node-pkg/index.js": `
				module.exports = function() {
					return 234
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapModuleDisabled(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"browser": {
						"node-pkg": false
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				const fn = require('node-pkg')
				module.exports = function() {
					return fn()
				}
			`,
			"/Users/user/project/node_modules/node-pkg/index.js": `
				module.exports = function() {
					return 234
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapNativeModuleDisabled(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"browser": {
						"fs": false
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				const fs = require('fs')
				module.exports = function() {
					return fs.readFile()
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserMapAvoidMissing(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'component-classes'
			`,
			"/Users/user/project/node_modules/component-classes/package.json": `
				{
					"browser": {
						"indexof": "component-indexof"
					}
				}
			`,
			"/Users/user/project/node_modules/component-classes/index.js": `
				try {
					var index = require('indexof');
				} catch (err) {
					var index = require('component-indexof');
				}
			`,
			"/Users/user/project/node_modules/component-indexof/index.js": `
				module.exports = function() {
					return 234
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserOverModuleBrowser(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./main.esm.js",
					"browser": "./main.browser.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.esm.js": `
				export default function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.browser.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformBrowser,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserOverMainNode(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./main.esm.js",
					"browser": "./main.browser.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.esm.js": `
				export default function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.browser.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNode,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserWithModuleBrowser(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./main.esm.js",
					"browser": {
						"./main.js": "./main.browser.js",
						"./main.esm.js": "./main.browser.esm.js"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.esm.js": `
				export default function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.browser.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.browser.esm.js": `
				export default function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformBrowser,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserWithMainNode(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./main.esm.js",
					"browser": {
						"./main.js": "./main.browser.js",
						"./main.esm.js": "./main.browser.esm.js"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.esm.js": `
				export default function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.browser.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.browser.esm.js": `
				export default function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNode,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserNodeModulesNoExt(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import {browser as a} from 'demo-pkg/no-ext'
				import {node as b} from 'demo-pkg/no-ext.js'
				import {browser as c} from 'demo-pkg/ext'
				import {browser as d} from 'demo-pkg/ext.js'
				console.log(a)
				console.log(b)
				console.log(c)
				console.log(d)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"browser": {
						"./no-ext": "./no-ext-browser.js",
						"./ext.js": "./ext-browser.js"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/no-ext.js": `
				export let node = 'node'
			`,
			"/Users/user/project/node_modules/demo-pkg/no-ext-browser.js": `
				export let browser = 'browser'
			`,
			"/Users/user/project/node_modules/demo-pkg/ext.js": `
				export let node = 'node'
			`,
			"/Users/user/project/node_modules/demo-pkg/ext-browser.js": `
				export let browser = 'browser'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserNodeModulesIndexNoExt(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import {browser as a} from 'demo-pkg/no-ext'
				import {node as b} from 'demo-pkg/no-ext/index.js'
				import {browser as c} from 'demo-pkg/ext'
				import {browser as d} from 'demo-pkg/ext/index.js'
				console.log(a)
				console.log(b)
				console.log(c)
				console.log(d)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"browser": {
						"./no-ext": "./no-ext-browser/index.js",
						"./ext/index.js": "./ext-browser/index.js"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/no-ext/index.js": `
				export let node = 'node'
			`,
			"/Users/user/project/node_modules/demo-pkg/no-ext-browser/index.js": `
				export let browser = 'browser'
			`,
			"/Users/user/project/node_modules/demo-pkg/ext/index.js": `
				export let node = 'node'
			`,
			"/Users/user/project/node_modules/demo-pkg/ext-browser/index.js": `
				export let browser = 'browser'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserNoExt(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import {browser as a} from './demo-pkg/no-ext'
				import {node as b} from './demo-pkg/no-ext.js'
				import {browser as c} from './demo-pkg/ext'
				import {browser as d} from './demo-pkg/ext.js'
				console.log(a)
				console.log(b)
				console.log(c)
				console.log(d)
			`,
			"/Users/user/project/src/demo-pkg/package.json": `
				{
					"browser": {
						"./no-ext": "./no-ext-browser.js",
						"./ext.js": "./ext-browser.js"
					}
				}
			`,
			"/Users/user/project/src/demo-pkg/no-ext.js": `
				export let node = 'node'
			`,
			"/Users/user/project/src/demo-pkg/no-ext-browser.js": `
				export let browser = 'browser'
			`,
			"/Users/user/project/src/demo-pkg/ext.js": `
				export let node = 'node'
			`,
			"/Users/user/project/src/demo-pkg/ext-browser.js": `
				export let browser = 'browser'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserIndexNoExt(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import {browser as a} from './demo-pkg/no-ext'
				import {node as b} from './demo-pkg/no-ext/index.js'
				import {browser as c} from './demo-pkg/ext'
				import {browser as d} from './demo-pkg/ext/index.js'
				console.log(a)
				console.log(b)
				console.log(c)
				console.log(d)
			`,
			"/Users/user/project/src/demo-pkg/package.json": `
				{
					"browser": {
						"./no-ext": "./no-ext-browser/index.js",
						"./ext/index.js": "./ext-browser/index.js"
					}
				}
			`,
			"/Users/user/project/src/demo-pkg/no-ext/index.js": `
				export let node = 'node'
			`,
			"/Users/user/project/src/demo-pkg/no-ext-browser/index.js": `
				export let browser = 'browser'
			`,
			"/Users/user/project/src/demo-pkg/ext/index.js": `
				export let node = 'node'
			`,
			"/Users/user/project/src/demo-pkg/ext-browser/index.js": `
				export let browser = 'browser'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

// See https://github.com/evanw/esbuild/issues/2002
func TestPackageJsonBrowserIssue2002A(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `require('pkg/sub')`,
			"/Users/user/project/src/node_modules/pkg/package.json": `{
				"browser": {
					"./sub": "./sub/foo.js"
				}
			}`,
			"/Users/user/project/src/node_modules/pkg/sub/foo.js":   `require('sub')`,
			"/Users/user/project/src/node_modules/sub/package.json": `{ "main": "./bar" }`,
			"/Users/user/project/src/node_modules/sub/bar.js":       `works()`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonBrowserIssue2002B(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `require('pkg/sub')`,
			"/Users/user/project/src/node_modules/pkg/package.json": `{
				"browser": {
					"./sub": "./sub/foo.js",
					"./sub/sub": "./sub/bar.js"
				}
			}`,
			"/Users/user/project/src/node_modules/pkg/sub/foo.js": `require('sub')`,
			"/Users/user/project/src/node_modules/pkg/sub/bar.js": `works()`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

// See https://github.com/evanw/esbuild/issues/2239
func TestPackageJsonBrowserIssue2002C(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `require('pkg/sub')`,
			"/Users/user/project/src/node_modules/pkg/package.json": `{
				"browser": {
					"./sub": "./sub/foo.js",
					"./sub/sub.js": "./sub/bar.js"
				}
			}`,
			"/Users/user/project/src/node_modules/pkg/sub/foo.js": `require('sub')`,
			"/Users/user/project/src/node_modules/sub/index.js":   `works()`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardImportOnly(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./module.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = 'main'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardRequireOnly(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				console.log(require('demo-pkg'))
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./module.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = 'main'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardImportAndRequireSameFile(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import value from 'demo-pkg'
				console.log(value, require('demo-pkg'))
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./module.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = 'main'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardImportAndRequireSeparateFiles(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import './test-main'
				import './test-module'
			`,
			"/Users/user/project/src/test-main.js": `
				console.log(require('demo-pkg'))
			`,
			"/Users/user/project/src/test-module.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./module.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = 'main'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardImportAndRequireForceModuleBeforeMain(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import './test-main'
				import './test-module'
			`,
			"/Users/user/project/src/test-main.js": `
				console.log(require('demo-pkg'))
			`,
			"/Users/user/project/src/test-module.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./module.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = 'main'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			MainFields:    []string{"module", "main"},
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardImportAndRequireImplicitMain(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import './test-index'
				import './test-module'
			`,
			"/Users/user/project/src/test-index.js": `
				console.log(require('demo-pkg'))
			`,
			"/Users/user/project/src/test-module.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"module": "./module.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				module.exports = 'index'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardImportAndRequireImplicitMainForceModuleBeforeMain(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import './test-index'
				import './test-module'
			`,
			"/Users/user/project/src/test-index.js": `
				console.log(require('demo-pkg'))
			`,
			"/Users/user/project/src/test-module.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"module": "./module.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/index.js": `
				module.exports = 'index'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			MainFields:    []string{"module", "main"},
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonDualPackageHazardImportAndRequireBrowser(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import './test-main'
				import './test-module'
			`,
			"/Users/user/project/src/test-main.js": `
				console.log(require('demo-pkg'))
			`,
			"/Users/user/project/src/test-module.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./module.js",
					"browser": {
						"./main.js": "./main.browser.js",
						"./module.js": "./module.browser.js"
					}
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = 'main'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.js": `
				export default 'module'
			`,
			"/Users/user/project/node_modules/demo-pkg/main.browser.js": `
				module.exports = 'browser main'
			`,
			"/Users/user/project/node_modules/demo-pkg/module.browser.js": `
				export default 'browser module'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonMainFieldsA(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"a": "./a.js",
					"b": "./b.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/a.js": `
				module.exports = 'a'
			`,
			"/Users/user/project/node_modules/demo-pkg/b.js": `
				export default 'b'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			MainFields:    []string{"a", "b"},
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonMainFieldsB(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import value from 'demo-pkg'
				console.log(value)
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"a": "./a.js",
					"b": "./b.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/a.js": `
				module.exports = 'a'
			`,
			"/Users/user/project/node_modules/demo-pkg/b.js": `
				export default 'b'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			MainFields:    []string{"b", "a"},
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonNeutralNoDefaultMainFields(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"main": "./main.js",
					"module": "./main.esm.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = function() {
					return 123
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.esm.js": `
				export default function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNeutral,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "demo-pkg"
Users/user/project/node_modules/demo-pkg/package.json: NOTE: The "main" field here was ignored. Main fields must be configured explicitly when using the "neutral" platform.
NOTE: You can mark the path "demo-pkg" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonNeutralExplicitMainFields(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import fn from 'demo-pkg'
				console.log(fn())
			`,
			"/Users/user/project/node_modules/demo-pkg/package.json": `
				{
					"hello": "./main.js",
					"module": "./main.esm.js"
				}
			`,
			"/Users/user/project/node_modules/demo-pkg/main.js": `
				module.exports = function() {
					return 123
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNeutral,
			MainFields:    []string{"hello"},
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonExportsErrorInvalidModuleSpecifier(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
				import 'pkg2'
				import 'pkg3'
				import 'pkg4'
				import 'pkg5'
				import 'pkg6'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{ "exports": { ".": "./%%" } }
			`,
			"/Users/user/project/node_modules/pkg2/package.json": `
				{ "exports": { ".": "./%2f" } }
			`,
			"/Users/user/project/node_modules/pkg3/package.json": `
				{ "exports": { ".": "./%2F" } }
			`,
			"/Users/user/project/node_modules/pkg4/package.json": `
				{ "exports": { ".": "./%5c" } }
			`,
			"/Users/user/project/node_modules/pkg5/package.json": `
				{ "exports": { ".": "./%5C" } }
			`,
			"/Users/user/project/node_modules/pkg6/package.json": `
				{ "exports": { ".": "./%31.js" } }
			`,
			"/Users/user/project/node_modules/pkg6/1.js": `
				console.log(1)
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The module specifier "./%%" is invalid:
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg2"
Users/user/project/node_modules/pkg2/package.json: NOTE: The module specifier "./%2f" is invalid:
NOTE: You can mark the path "pkg2" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg3"
Users/user/project/node_modules/pkg3/package.json: NOTE: The module specifier "./%2F" is invalid:
NOTE: You can mark the path "pkg3" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg4"
Users/user/project/node_modules/pkg4/package.json: NOTE: The module specifier "./%5c" is invalid:
NOTE: You can mark the path "pkg4" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg5"
Users/user/project/node_modules/pkg5/package.json: NOTE: The module specifier "./%5C" is invalid:
NOTE: You can mark the path "pkg5" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsErrorInvalidPackageConfiguration(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
				import 'pkg2/foo'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{ "exports": { ".": false } }
			`,
			"/Users/user/project/node_modules/pkg2/package.json": `
				{ "exports": { "./foo": false } }
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/node_modules/pkg1/package.json: WARNING: This value must be a string, an object, an array, or null
Users/user/project/node_modules/pkg2/package.json: WARNING: This value must be a string, an object, an array, or null
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The package configuration has an invalid value here:
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg2/foo"
Users/user/project/node_modules/pkg2/package.json: NOTE: The package configuration has an invalid value here:
NOTE: You can mark the path "pkg2/foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsErrorInvalidPackageTarget(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
				import 'pkg2'
				import 'pkg3'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{ "exports": { ".": "invalid" } }
			`,
			"/Users/user/project/node_modules/pkg2/package.json": `
				{ "exports": { ".": "./../pkg3" } }
			`,
			"/Users/user/project/node_modules/pkg3/package.json": `
				{ "exports": { ".": "./node_modules/pkg" } }
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The package target "invalid" is invalid because it doesn't start with "./":
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg2"
Users/user/project/node_modules/pkg2/package.json: NOTE: The package target "./../pkg3" is invalid because it contains invalid segment "..":
NOTE: You can mark the path "pkg2" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg3"
Users/user/project/node_modules/pkg3/package.json: NOTE: The package target "./node_modules/pkg" is invalid because it contains invalid segment "node_modules":
NOTE: You can mark the path "pkg3" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsErrorPackagePathNotExported(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1/foo'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{ "exports": { ".": {} } }
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1/foo"
Users/user/project/node_modules/pkg1/package.json: NOTE: The path "./foo" is not exported by package "pkg1":
NOTE: You can mark the path "pkg1/foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsErrorModuleNotFound(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{ "exports": { ".": "./foo.js" } }
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The module "./foo.js" was not found on the file system:
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsErrorUnsupportedDirectoryImport(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
				import 'pkg2'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{ "exports": { ".": "./foo/" } }
			`,
			"/Users/user/project/node_modules/pkg2/package.json": `
				{ "exports": { ".": "./foo" } }
			`,
			"/Users/user/project/node_modules/pkg2/foo/bar.js": `
				console.log(bar)
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The module "./foo" was not found on the file system:
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg2"
Users/user/project/node_modules/pkg2/package.json: NOTE: Importing the directory "./foo" is forbidden by this package:
Users/user/project/node_modules/pkg2/package.json: NOTE: The presence of "exports" here makes importing a directory forbidden:
NOTE: You can mark the path "pkg2" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonImportsErrorUnsupportedDirectoryImport(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import '#foo1/bar'
				import '#foo2/bar'
			`,
			"/Users/user/project/package.json": `
				{
					"imports": {
						"#foo1/*": "./foo1/*",
						"#foo2/bar": "./foo2/bar"
					}
				}
			`,
			"/Users/user/project/foo1/bar/index.js": `
				console.log(bar)
			`,
			"/Users/user/project/foo2/bar/index.js": `
				console.log(bar)
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "#foo1/bar"
Users/user/project/package.json: NOTE: Importing the directory "./foo1/bar" is forbidden by this package:
Users/user/project/package.json: NOTE: The presence of "imports" here makes importing a directory forbidden:
Users/user/project/src/entry.js: NOTE: Import from "/index.js" to get the file "Users/user/project/foo1/bar/index.js":
NOTE: You can mark the path "#foo1/bar" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "#foo2/bar"
Users/user/project/package.json: NOTE: Importing the directory "./foo2/bar" is forbidden by this package:
Users/user/project/package.json: NOTE: The presence of "imports" here makes importing a directory forbidden:
NOTE: You can mark the path "#foo2/bar" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsRequireOverImport(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				require('pkg')
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"import": "./import.js",
						"require": "./require.js",
						"default": "./default.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/import.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/require.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonExportsImportOverRequire(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg'
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"require": "./require.js",
						"import": "./import.js",
						"default": "./default.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/require.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/import.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonExportsDefaultOverImportAndRequire(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg'
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"default": "./default.js",
						"import": "./import.js",
						"require": "./require.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/require.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/import.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/default.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonExportsEntryPointImportOverRequire(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/node_modules/pkg/package.json": `
				{
					"exports": {
						"import": "./import.js",
						"require": "./require.js"
					},
					"module": "./module.js",
					"main": "./main.js"
				}
			`,
			"/node_modules/pkg/import.js": `
				console.log('SUCCESS')
			`,
			"/node_modules/pkg/require.js": `
				console.log('FAILURE')
			`,
			"/node_modules/pkg/module.js": `
				console.log('FAILURE')
			`,
			"/node_modules/pkg/main.js": `
				console.log('FAILURE')
			`,
		},
		entryPaths: []string{"pkg"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestPackageJsonExportsEntryPointRequireOnly(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/node_modules/pkg/package.json": `
				{
					"exports": {
						"require": "./require.js"
					},
					"module": "./module.js",
					"main": "./main.js"
				}
			`,
			"/node_modules/pkg/require.js": `
				console.log('FAILURE')
			`,
			"/node_modules/pkg/module.js": `
				console.log('FAILURE')
			`,
			"/node_modules/pkg/main.js": `
				console.log('FAILURE')
			`,
		},
		entryPaths: []string{"pkg"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
		expectedScanLog: `ERROR: Could not resolve "pkg"
node_modules/pkg/package.json: NOTE: The path "." is not currently exported by package "pkg":
node_modules/pkg/package.json: NOTE: None of the conditions in the package definition ("require") match any of the currently active conditions ("browser", "default", "import"):
`,
	})
}

func TestPackageJsonExportsEntryPointModuleOverMain(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/node_modules/pkg/package.json": `
				{
					"module": "./module.js",
					"main": "./main.js"
				}
			`,
			"/node_modules/pkg/module.js": `
				console.log('SUCCESS')
			`,
			"/node_modules/pkg/main.js": `
				console.log('FAILURE')
			`,
		},
		entryPaths: []string{"pkg"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestPackageJsonExportsEntryPointMainOnly(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/node_modules/pkg/package.json": `
				{
					"main": "./main.js"
				}
			`,
			"/node_modules/pkg/main.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"pkg"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

func TestPackageJsonExportsBrowser(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg'
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"node": "./node.js",
						"browser": "./browser.js",
						"default": "./default.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/node.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/browser.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			Platform:      config.PlatformBrowser,
		},
	})
}

func TestPackageJsonExportsNode(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg'
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"browser": "./browser.js",
						"node": "./node.js",
						"default": "./default.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/browser.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/node.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			Platform:      config.PlatformNode,
		},
	})
}

func TestPackageJsonExportsNeutral(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg'
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"node": "./node.js",
						"browser": "./browser.js",
						"default": "./default.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/node.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/browser.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg/default.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			Platform:      config.PlatformNeutral,
		},
	})
}

func TestPackageJsonExportsOrderIndependent(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1/foo/bar.js'
				import 'pkg2/foo/bar.js'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						"./": "./1/",
						"./foo/": "./2/"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/1/foo/bar.js": `
				console.log('FAILURE')
			`,
			"/Users/user/project/node_modules/pkg1/2/bar.js": `
				console.log('SUCCESS')
			`,
			"/Users/user/project/node_modules/pkg2/package.json": `
				{
					"exports": {
						"./foo/": "./1/",
						"./": "./2/"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg2/1/bar.js": `
				console.log('SUCCESS')
			`,
			"/Users/user/project/node_modules/pkg2/2/foo/bar.js": `
				console.log('FAILURE')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonExportsWildcard(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1/foo'
				import 'pkg1/foo2'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						"./foo*": "./file*.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/file.js": `
				console.log('SUCCESS')
			`,
			"/Users/user/project/node_modules/pkg1/file2.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonExportsErrorMissingTrailingSlash(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1/foo/bar'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{ "exports": { "./foo/": "./test" } }
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1/foo/bar"
Users/user/project/node_modules/pkg1/package.json: NOTE: The module specifier "./test" is invalid because it doesn't end in "/":
NOTE: You can mark the path "pkg1/foo/bar" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsCustomConditions(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						"custom1": "./custom1.js",
						"custom2": "./custom2.js",
						"default": "./default.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/custom2.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			Conditions:    []string{"custom2"},
		},
	})
}

func TestPackageJsonExportsNotExactMissingExtension(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1/foo/bar'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						"./foo/": "./dir/"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/dir/bar.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonExportsNotExactMissingExtensionPattern(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1/foo/bar'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						"./foo/*": "./dir/*"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/dir/bar.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1/foo/bar"
Users/user/project/node_modules/pkg1/package.json: NOTE: The module "./dir/bar" was not found on the file system:
Users/user/project/src/entry.js: NOTE: Import from "pkg1/foo/bar.js" to get the file "Users/user/project/node_modules/pkg1/dir/bar.js":
NOTE: You can mark the path "pkg1/foo/bar" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsExactMissingExtension(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1/foo/bar'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						"./foo/bar": "./dir/bar"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/dir/bar.js": `
				console.log('SUCCESS')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1/foo/bar"
Users/user/project/node_modules/pkg1/package.json: NOTE: The module "./dir/bar" was not found on the file system:
NOTE: You can mark the path "pkg1/foo/bar" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsNoConditionsMatch(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
				import 'pkg1/foo.js'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						".": {
							"what": "./foo.js"
						},
						"./foo.js": {
							"what": "./foo.js"
						}
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/foo.js": `
				console.log('FAILURE')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The path "." is not currently exported by package "pkg1":
Users/user/project/node_modules/pkg1/package.json: NOTE: None of the conditions in the package definition ("what") match any of the currently active conditions ("browser", "default", "import"):
Users/user/project/node_modules/pkg1/package.json: NOTE: Consider enabling the "what" condition if this package expects it to be enabled. You can use 'Conditions: []string{"what"}' to do that:
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1/foo.js"
Users/user/project/node_modules/pkg1/package.json: NOTE: The path "./foo.js" is not currently exported by package "pkg1":
Users/user/project/node_modules/pkg1/package.json: NOTE: None of the conditions in the package definition ("what") match any of the currently active conditions ("browser", "default", "import"):
Users/user/project/node_modules/pkg1/package.json: NOTE: Consider enabling the "what" condition if this package expects it to be enabled. You can use 'Conditions: []string{"what"}' to do that:
NOTE: You can mark the path "pkg1/foo.js" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsMustUseRequire(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg1'
				import 'pkg1/foo.js'
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						".": {
							"require": "./foo.js"
						},
						"./foo.js": {
							"require": "./foo.js"
						}
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/foo.js": `
				console.log('FAILURE')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The path "." is not currently exported by package "pkg1":
Users/user/project/node_modules/pkg1/package.json: NOTE: None of the conditions in the package definition ("require") match any of the currently active conditions ("browser", "default", "import"):
Users/user/project/src/entry.js: NOTE: Consider using a "require()" call to import this file, which will work because the "require" condition is supported by this package:
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1/foo.js"
Users/user/project/node_modules/pkg1/package.json: NOTE: The path "./foo.js" is not currently exported by package "pkg1":
Users/user/project/node_modules/pkg1/package.json: NOTE: None of the conditions in the package definition ("require") match any of the currently active conditions ("browser", "default", "import"):
Users/user/project/src/entry.js: NOTE: Consider using a "require()" call to import this file, which will work because the "require" condition is supported by this package:
NOTE: You can mark the path "pkg1/foo.js" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonExportsMustUseImport(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				require('pkg1')
				require('pkg1/foo.js')
			`,
			"/Users/user/project/node_modules/pkg1/package.json": `
				{
					"exports": {
						".": {
							"import": "./foo.js"
						},
						"./foo.js": {
							"import": "./foo.js"
						}
					}
				}
			`,
			"/Users/user/project/node_modules/pkg1/foo.js": `
				console.log('FAILURE')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1"
Users/user/project/node_modules/pkg1/package.json: NOTE: The path "." is not currently exported by package "pkg1":
Users/user/project/node_modules/pkg1/package.json: NOTE: None of the conditions in the package definition ("import") match any of the currently active conditions ("browser", "default", "require"):
Users/user/project/src/entry.js: NOTE: Consider using an "import" statement to import this file, which will work because the "import" condition is supported by this package:
NOTE: You can mark the path "pkg1" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle. You can also surround this "require" call with a try/catch block to handle this failure at run-time instead of bundle-time.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg1/foo.js"
Users/user/project/node_modules/pkg1/package.json: NOTE: The path "./foo.js" is not currently exported by package "pkg1":
Users/user/project/node_modules/pkg1/package.json: NOTE: None of the conditions in the package definition ("import") match any of the currently active conditions ("browser", "default", "require"):
Users/user/project/src/entry.js: NOTE: Consider using an "import" statement to import this file, which will work because the "import" condition is supported by this package:
NOTE: You can mark the path "pkg1/foo.js" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle. You can also surround this "require" call with a try/catch block to handle this failure at run-time instead of bundle-time.
`,
	})
}

func TestPackageJsonExportsReverseLookup(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				require('pkg/path/to/real/file')
				require('pkg/path/to/other/file')
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"./lib/te*": {
							"default": "./path/to/re*.js"
						},
						"./extra/": {
							"default": "./path/to/"
						}
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/path/to/real/file.js":  ``,
			"/Users/user/project/node_modules/pkg/path/to/other/file.js": ``,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg/path/to/real/file"
Users/user/project/node_modules/pkg/package.json: NOTE: The path "./path/to/real/file" is not exported by package "pkg":
Users/user/project/node_modules/pkg/package.json: NOTE: The file "./path/to/real/file.js" is exported at path "./lib/teal/file":
Users/user/project/src/entry.js: NOTE: Import from "pkg/lib/teal/file" to get the file "Users/user/project/node_modules/pkg/path/to/real/file.js":
NOTE: You can mark the path "pkg/path/to/real/file" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle. You can also surround this "require" call with a try/catch block to handle this failure at run-time instead of bundle-time.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg/path/to/other/file"
Users/user/project/node_modules/pkg/package.json: NOTE: The path "./path/to/other/file" is not exported by package "pkg":
Users/user/project/node_modules/pkg/package.json: NOTE: The file "./path/to/other/file.js" is exported at path "./extra/other/file.js":
Users/user/project/src/entry.js: NOTE: Import from "pkg/extra/other/file.js" to get the file "Users/user/project/node_modules/pkg/path/to/other/file.js":
NOTE: You can mark the path "pkg/path/to/other/file" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle. You can also surround this "require" call with a try/catch block to handle this failure at run-time instead of bundle-time.
`,
	})
}

func TestPackageJsonExportsPatternTrailers(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'pkg/path/foo.js/bar.js'
				import 'pkg2/features/abc'
				import 'pkg2/features/xyz.js'
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"./path/*/bar.js": "./dir/baz-*"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/dir/baz-foo.js": `
				console.log('works')
			`,
			"/Users/user/project/node_modules/pkg2/package.json": `
				{
					"exports": {
						"./features/*": "./public/*.js",
						"./features/*.js": "./public/*.js"
					}
				}
			`,
			"/Users/user/project/node_modules/pkg2/public/abc.js": `
				console.log('abc')
			`,
			"/Users/user/project/node_modules/pkg2/public/xyz.js": `
				console.log('xyz')
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

// Node's package.json format for "exports" allows for arrays to be used as map
// values, like in the example below. Webpack's implementation interprets this
// as a way to specify several alternative directories to search for packages.
// See: https://webpack.js.org/guides/package-exports/#alternatives. However,
// this doesn't follow Node's specification for how "exports" should work:
// https://nodejs.org/api/esm.html#resolver-algorithm. Also no one else
// implements it this way (e.g. both Node and Rollup don't do this).
//
// This test case can only be built by Webpack. Implementations that follow the
// specification (including esbuild) will fail to build this test case. This
// test case only exists to document that esbuild doesn't follow Webpack's
// behavior here.
func TestPackageJsonExportsAlternatives(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import redApple from 'pkg/apples/red.js'
				import greenApple from 'pkg/apples/green.js'
				import redBook from 'pkg/books/red'
				import greenBook from 'pkg/books/green'
				console.log({redApple, greenApple, redBook, greenBook})
			`,
			"/Users/user/project/node_modules/pkg/package.json": `
				{
					"exports": {
						"./apples/": ["./good-apples/", "./bad-apples/"],
						"./books/*": ["./good-books/*-book.js", "./bad-books/*-book.js"]
					}
				}
			`,
			"/Users/user/project/node_modules/pkg/good-apples/green.js": `
				export default '🍏'
			`,
			"/Users/user/project/node_modules/pkg/bad-apples/red.js": `
				export default '🍎'
			`,
			"/Users/user/project/node_modules/pkg/good-books/green-book.js": `
				export default '📗'
			`,
			"/Users/user/project/node_modules/pkg/bad-books/red-book.js": `
				export default '📕'
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "pkg/apples/red.js"
Users/user/project/node_modules/pkg/package.json: NOTE: The module "./good-apples/red.js" was not found on the file system:
NOTE: You can mark the path "pkg/apples/red.js" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/entry.js: ERROR: Could not resolve "pkg/books/red"
Users/user/project/node_modules/pkg/package.json: NOTE: The module "./good-books/red-book.js" was not found on the file system:
NOTE: You can mark the path "pkg/books/red" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonImports(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/foo/entry.js": `
				import '#top-level'
				import '#nested/path.js'
				import '#star/c.js'
				import '#slash/d.js'
			`,
			"/Users/user/project/src/package.json": `
				{
					"imports": {
						"#top-level": "./a.js",
						"#nested/path.js": "./b.js",
						"#star/*": "./some-star/*",
						"#slash/": "./some-slash/"
					}
				}
			`,
			"/Users/user/project/src/a.js":            `console.log('a.js')`,
			"/Users/user/project/src/b.js":            `console.log('b.js')`,
			"/Users/user/project/src/some-star/c.js":  `console.log('c.js')`,
			"/Users/user/project/src/some-slash/d.js": `console.log('d.js')`,
		},
		entryPaths: []string{"/Users/user/project/src/foo/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonImportsRemapToOtherPackage(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import '#top-level'
				import '#nested/path.js'
				import '#star/c.js'
				import '#slash/d.js'
			`,
			"/Users/user/project/src/package.json": `
				{
					"imports": {
						"#top-level": "pkg/a.js",
						"#nested/path.js": "pkg/b.js",
						"#star/*": "pkg/some-star/*",
						"#slash/": "pkg/some-slash/"
					}
				}
			`,
			"/Users/user/project/src/node_modules/pkg/a.js":            `console.log('a.js')`,
			"/Users/user/project/src/node_modules/pkg/b.js":            `console.log('b.js')`,
			"/Users/user/project/src/node_modules/pkg/some-star/c.js":  `console.log('c.js')`,
			"/Users/user/project/src/node_modules/pkg/some-slash/d.js": `console.log('d.js')`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
	})
}

func TestPackageJsonImportsErrorMissingRemappedPackage(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import '#foo'
			`,
			"/Users/user/project/src/package.json": `
				{
					"imports": {
						"#foo": "bar"
					}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "#foo"
Users/user/project/src/package.json: NOTE: The remapped path "bar" could not be resolved:
NOTE: You can mark the path "#foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonImportsInvalidPackageConfiguration(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import '#foo'
			`,
			"/Users/user/project/src/package.json": `
				{
					"imports": "#foo"
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "#foo"
Users/user/project/src/package.json: NOTE: The package configuration has an invalid value here:
NOTE: You can mark the path "#foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
Users/user/project/src/package.json: WARNING: The value for "imports" must be an object
`,
	})
}

func TestPackageJsonImportsErrorEqualsHash(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import '#'
			`,
			"/Users/user/project/src/package.json": `
				{
					"imports": {}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "#"
Users/user/project/src/package.json: NOTE: This "imports" map was ignored because the module specifier "#" is invalid:
NOTE: You can mark the path "#" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonImportsErrorStartsWithHashSlash(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import '#/foo'
			`,
			"/Users/user/project/src/package.json": `
				{
					"imports": {}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "#/foo"
Users/user/project/src/package.json: NOTE: This "imports" map was ignored because the module specifier "#/foo" is invalid:
NOTE: You can mark the path "#/foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonMainFieldsErrorMessageDefault(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'foo'
			`,
			"/Users/user/project/node_modules/foo/package.json": `
				{
					"main": "./foo"
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "foo"
NOTE: You can mark the path "foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonMainFieldsErrorMessageNotIncluded(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'foo'
			`,
			"/Users/user/project/node_modules/foo/package.json": `
				{
					"main": "./foo"
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{"some", "fields"},
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "foo"
Users/user/project/node_modules/foo/package.json: NOTE: The "main" field here was ignored because the list of main fields to use is currently set to ["some", "fields"].
NOTE: You can mark the path "foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonMainFieldsErrorMessageEmpty(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/entry.js": `
				import 'foo'
			`,
			"/Users/user/project/node_modules/foo/package.json": `
				{
					"main": "./foo"
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
		expectedScanLog: `Users/user/project/src/entry.js: ERROR: Could not resolve "foo"
Users/user/project/node_modules/foo/package.json: NOTE: The "main" field here was ignored because the list of main fields to use is currently set to [].
NOTE: You can mark the path "foo" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestPackageJsonTypeShouldBeTypes(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/index.js": ``,
			"/Users/user/project/package.json": `
				{
					"main": "./src/index.js",
					"type": "./src/index.d.ts"
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/index.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
		expectedScanLog: `Users/user/project/package.json: WARNING: "./src/index.d.ts" is not a valid value for the "type" field
Users/user/project/package.json: NOTE: TypeScript type declarations use the "types" field, not the "type" field:
`,
	})
}

func TestPackageJsonImportSelfUsingRequire(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/index.js": `
				module.exports = 'index'
				console.log(
					require("xyz"),
					require("xyz/bar"),
				)
			`,
			"/Users/user/project/src/foo-import.js": `
				export default 'foo'
			`,
			"/Users/user/project/src/foo-require.js": `
				module.exports = 'foo'
			`,
			"/Users/user/project/package.json": `
				{
					"name": "xyz",
					"exports": {
						".": "./src/index.js",
						"./bar": {
							"import": "./src/foo-import.js",
							"require": "./src/foo-require.js"
						}
					}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/index.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
	})
}

func TestPackageJsonImportSelfUsingImport(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/index.js": `
				import xyz from "xyz"
				import foo from "xyz/bar"
				export default 'index'
				console.log(xyz, foo)
			`,
			"/Users/user/project/src/foo-import.js": `
				export default 'foo'
			`,
			"/Users/user/project/src/foo-require.js": `
				module.exports = 'foo'
			`,
			"/Users/user/project/package.json": `
				{
					"name": "xyz",
					"exports": {
						".": "./src/index.js",
						"./bar": {
							"import": "./src/foo-import.js",
							"require": "./src/foo-require.js"
						}
					}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/index.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
	})
}

func TestPackageJsonImportSelfUsingRequireScoped(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/index.js": `
				module.exports = 'index'
				console.log(
					require("@some-scope/xyz"),
					require("@some-scope/xyz/bar"),
				)
			`,
			"/Users/user/project/src/foo-import.js": `
				export default 'foo'
			`,
			"/Users/user/project/src/foo-require.js": `
				module.exports = 'foo'
			`,
			"/Users/user/project/package.json": `
				{
					"name": "@some-scope/xyz",
					"exports": {
						".": "./src/index.js",
						"./bar": {
							"import": "./src/foo-import.js",
							"require": "./src/foo-require.js"
						}
					}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/index.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
	})
}

func TestPackageJsonImportSelfUsingImportScoped(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/index.js": `
				import xyz from "@some-scope/xyz"
				import foo from "@some-scope/xyz/bar"
				export default 'index'
				console.log(xyz, foo)
			`,
			"/Users/user/project/src/foo-import.js": `
				export default 'foo'
			`,
			"/Users/user/project/src/foo-require.js": `
				module.exports = 'foo'
			`,
			"/Users/user/project/package.json": `
				{
					"name": "@some-scope/xyz",
					"exports": {
						".": "./src/index.js",
						"./bar": {
							"import": "./src/foo-import.js",
							"require": "./src/foo-require.js"
						}
					}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/index.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
	})
}

func TestPackageJsonImportSelfUsingRequireFailure(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/index.js": `
				require("xyz/src/foo.js")
			`,
			"/Users/user/project/src/foo.js": `
				module.exports = 'foo'
			`,
			"/Users/user/project/package.json": `
				{
					"name": "xyz",
					"exports": {
						".": "./src/index.js",
						"./bar": "./src/foo.js"
					}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/index.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
		expectedScanLog: `Users/user/project/src/index.js: ERROR: Could not resolve "xyz/src/foo.js"
Users/user/project/package.json: NOTE: The path "./src/foo.js" is not exported by package "xyz":
Users/user/project/package.json: NOTE: The file "./src/foo.js" is exported at path "./bar":
Users/user/project/src/index.js: NOTE: Import from "xyz/bar" to get the file "Users/user/project/src/foo.js":
NOTE: You can mark the path "xyz/src/foo.js" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle. You can also surround this "require" call with a try/catch block to handle this failure at run-time instead of bundle-time.
`,
	})
}

func TestPackageJsonImportSelfUsingImportFailure(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/Users/user/project/src/index.js": `
				import "xyz/src/foo.js"
			`,
			"/Users/user/project/src/foo.js": `
				export default 'foo'
			`,
			"/Users/user/project/package.json": `
				{
					"name": "xyz",
					"exports": {
						".": "./src/index.js",
						"./bar": "./src/foo.js"
					}
				}
			`,
		},
		entryPaths: []string{"/Users/user/project/src/index.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/Users/user/project/out.js",
			MainFields:    []string{},
		},
		expectedScanLog: `Users/user/project/src/index.js: ERROR: Could not resolve "xyz/src/foo.js"
Users/user/project/package.json: NOTE: The path "./src/foo.js" is not exported by package "xyz":
Users/user/project/package.json: NOTE: The file "./src/foo.js" is exported at path "./bar":
Users/user/project/src/index.js: NOTE: Import from "xyz/bar" to get the file "Users/user/project/src/foo.js":
NOTE: You can mark the path "xyz/src/foo.js" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

func TestCommonJSVariableInESMTypeModule(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js":     `module.exports = null`,
			"/package.json": `{ "type": "module" }`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
		expectedScanLog: `entry.js: WARNING: The CommonJS "module" variable is treated as a global variable in an ECMAScript module and may not work as expected
package.json: NOTE: This file is considered to be an ECMAScript module because the enclosing "package.json" file sets the type of this file to "module":
NOTE: Node's package format requires that CommonJS files in a "type": "module" package use the ".cjs" file extension.
`,
	})
}

func TestPackageJsonNodePathsIssue2752(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/src/entry.js": `
				import "pkg1"
				import "pkg2"
				import "@scope/pkg3/baz"
				import "@scope/pkg4"
			`,
			"/usr/lib/pkg/pkg1/package.json":          `{ "main": "./foo.js" }`,
			"/usr/lib/pkg/pkg1/foo.js":                `console.log('pkg1')`,
			"/lib/pkg/pkg2/package.json":              `{ "exports": { ".": "./bar.js" } }`,
			"/lib/pkg/pkg2/bar.js":                    `console.log('pkg2')`,
			"/var/lib/pkg/@scope/pkg3/package.json":   `{ "browser": { "./baz.js": "./baz-browser.js" } }`,
			"/var/lib/pkg/@scope/pkg3/baz-browser.js": `console.log('pkg3')`,
			"/tmp/pkg/@scope/pkg4/package.json":       `{ "exports": { ".": { "import": "./bat.js" } } }`,
			"/tmp/pkg/@scope/pkg4/bat.js":             `console.log('pkg4')`,
		},
		entryPaths: []string{"/src/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
			AbsNodePaths: []string{
				"/usr/lib/pkg",
				"/lib/pkg",
				"/var/lib/pkg",
				"/tmp/pkg",
			},
		},
	})
}

// See: https://github.com/evanw/esbuild/issues/3377
func TestPackageJsonReversePackageExportsIssue3377(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/lib/msw-config.ts": `
				import { setupWorker, type SetupWorker } from 'msw/browser'
				setupWorker();
			`,
			"/node_modules/msw/package.json": `{
				"exports": {
					"./browser": {
						"node": null,
						"require": "./lib/browser/index.js",
						"import": "./lib/browser/index.mjs",
						"default": "./lib/browser/index.js"
					}
				}
			}`,
			"/node_modules/msw/browser/package.json": `{
				"main": "../lib/browser/index.js",
				"module": "../lib/browser/index.mjs"
			}`,
			"/node_modules/msw/lib/browser/index.js":  `TEST FAILURE`,
			"/node_modules/msw/lib/browser/index.mjs": `TEST FAILURE`,
		},
		entryPaths: []string{"/lib/msw-config.ts"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
			Platform:      config.PlatformNode,
		},
		expectedScanLog: `lib/msw-config.ts: ERROR: Could not resolve "msw/browser"
node_modules/msw/package.json: NOTE: The path "./browser" cannot be imported from package "msw" because it was explicitly disabled by the package author here:
NOTE: You can mark the path "msw/browser" as external to exclude it from the bundle, which will remove this error and leave the unresolved path in the bundle.
`,
	})
}

// See: https://github.com/evanw/esbuild/issues/3367
func TestPackageJsonDisabledTypeModuleIssue3367(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import foo from 'foo'
				foo()
			`,
			"/package.json": `
				{
					"browser": {
						"foo": false
					}
				}
			`,
			"/node_modules/foo/package.json": `
				{
					"type": "module"
				}
			`,
			"/node_modules/foo/index.js": `
				export default function() {}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

// See: https://github.com/evanw/esbuild/issues/3485
func TestPackageJsonSubpathImportNodeBuiltinIssue3485(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import fs from '#fs'
				import http from '#http'
				fs.readFileSync()
				http.createServer()
			`,
			"/package.json": `
				{
					"imports": {
						"#fs": {
							"node": "fs",
							"default": "./empty.js"
						},
						"#http": {
							"node": "node:http",
							"default": "./empty.js"
						}
					}
				}
			`,
			"/empty.js": ``,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNode,
			AbsOutputFile: "/out.js",
		},
	})
}

// See: https://github.com/evanw/esbuild/issues/3867
func TestPackageJsonBadExportsImportAndRequireWarningIssue3867(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import "foo"
			`,
			"/node_modules/foo/package.json": `
				{
					"exports": {
						".": {
							"import": "./dist/node/index.js",
							"require": "./dist/node/index.cjs",
							"node": {
								"import": "./dist/node/index.js",
								"require": "./dist/node/index.cjs"
							},
							"browser": {
								"import": "./dist/browser/index.js",
								"require": "./dist/browser/index.cjs"
							},
							"worker": {
								"import": "./dist/browser/index.js",
								"require": "./dist/browser/index.cjs"
							}
						}
					}
				}
			`,
			"/node_modules/foo/dist/node/index.js":     ``,
			"/node_modules/foo/dist/node/index.cjs":    ``,
			"/node_modules/foo/dist/browser/index.js":  ``,
			"/node_modules/foo/dist/browser/index.cjs": ``,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNode,
			AbsOutputFile: "/out.js",
		},
		debugLogs: true,
		expectedScanLog: `node_modules/foo/package.json: DEBUG: The conditions "node" and "browser" and "worker" here will never be used as they come after both "import" and "require"
node_modules/foo/package.json: NOTE: The "import" condition comes earlier and will be used for all "import" statements:
node_modules/foo/package.json: NOTE: The "require" condition comes earlier and will be used for all "require" calls:
`,
	})
}

// See: https://github.com/evanw/esbuild/issues/3867
func TestPackageJsonBadExportsDefaultWarningIssue3867(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import "foo"
			`,
			"/node_modules/foo/package.json": `
				{
					"exports": {
						".": {
							"default": "./dist/node/index.js",
							"node": {
								"import": "./dist/node/index.js",
								"require": "./dist/node/index.cjs"
							},
							"browser": {
								"import": "./dist/browser/index.js",
								"require": "./dist/browser/index.cjs"
							},
							"worker": {
								"import": "./dist/browser/index.js",
								"require": "./dist/browser/index.cjs"
							}
						}
					}
				}
			`,
			"/node_modules/foo/dist/node/index.js":     ``,
			"/node_modules/foo/dist/node/index.cjs":    ``,
			"/node_modules/foo/dist/browser/index.js":  ``,
			"/node_modules/foo/dist/browser/index.cjs": ``,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNode,
			AbsOutputFile: "/out.js",
		},
		debugLogs: true,
		expectedScanLog: `node_modules/foo/package.json: DEBUG: The conditions "node" and "browser" and "worker" here will never be used as they come after "default"
node_modules/foo/package.json: NOTE: The "default" condition comes earlier and will always be chosen:
`,
	})
}

// See: https://github.com/evanw/esbuild/issues/3887
func TestPackageJsonExportsDefaultWarningIssue3887(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import "foo"
			`,
			"/node_modules/foo/dist/index.js": `
				success()
			`,
			"/node_modules/foo/package.json": `
				{
					"exports": {
						".": {
							"node": "./dist/index.js",
							"require": "./dist/index.js",
							"import": "./dist/index.esm.js",
							"default": "./dist/index.esm.js"
						}
					}
				}
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			Platform:      config.PlatformNode,
			AbsOutputFile: "/out.js",
		},
		debugLogs: true,
	})
}

// https://github.com/evanw/esbuild/issues/4144
func TestConfusingNameCollisionsIssue4144(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import { it } from 'mydependency'
				console.log(it())
			`,
			"/node_modules/mydependency/package.json": `
				{
					"main": "./package/index.js"
				}
			`,
			"/node_modules/mydependency/package/index.js": `
				export { it } from './utils'
				export let works = true
			`,
			"/node_modules/mydependency/package/utils/index.js": `
				export { it } from './utils'
			`,
			"/node_modules/mydependency/package/utils/utils.js": `
				// This should resolve to "../index.js" not "../../package.json"
				import { works } from '..'
				export function it() { return works }
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}

// https://github.com/evanw/esbuild/issues/4187
func TestPackageJsonBrowserMatchingTrailingSlashIssue4187(t *testing.T) {
	packagejson_suite.expectBundled(t, bundled{
		files: map[string]string{
			"/entry.js": `
				import axios from "axios"
			`,
			"/node_modules/axios/package.json": `
				{
					"browser": {
						"./node/index.js": "./browser/index.js"
					}
				}
			`,
			"/node_modules/axios/index.js": `
				module.exports = require('./node/');
			`,
			"/node_modules/axios/node/index.js": `
				module.exports = { get: () => new Promise('Node') }
			`,
			"/node_modules/axios/browser/index.js": `
				module.exports = { get: () => new Promise('Browser') }
			`,
		},
		entryPaths: []string{"/entry.js"},
		options: config.Options{
			Mode:          config.ModeBundle,
			AbsOutputFile: "/out.js",
		},
	})
}
