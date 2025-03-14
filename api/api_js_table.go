// This file was automatically generated by "js_table.ts"

package api

import "github.com/ije/esbuild-internal/compat"

type EngineName uint8

const (
	EngineChrome EngineName = iota
	EngineDeno
	EngineEdge
	EngineFirefox
	EngineHermes
	EngineIE
	EngineIOS
	EngineNode
	EngineOpera
	EngineRhino
	EngineSafari
)

func convertEngineName(engine EngineName) compat.Engine {
	switch engine {
	case EngineChrome:
		return compat.Chrome
	case EngineDeno:
		return compat.Deno
	case EngineEdge:
		return compat.Edge
	case EngineFirefox:
		return compat.Firefox
	case EngineHermes:
		return compat.Hermes
	case EngineIE:
		return compat.IE
	case EngineIOS:
		return compat.IOS
	case EngineNode:
		return compat.Node
	case EngineOpera:
		return compat.Opera
	case EngineRhino:
		return compat.Rhino
	case EngineSafari:
		return compat.Safari
	default:
		panic("Invalid engine name")
	}
}
