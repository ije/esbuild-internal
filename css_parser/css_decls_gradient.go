package css_parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ije/esbuild-internal/compat"
	"github.com/ije/esbuild-internal/css_ast"
	"github.com/ije/esbuild-internal/css_lexer"
	"github.com/ije/esbuild-internal/logger"
)

type gradientKind uint8

const (
	linearGradient gradientKind = iota
	radialGradient
	conicGradient
)

type parsedGradient struct {
	leadingTokens []css_ast.Token
	colorStops    []colorStop
	kind          gradientKind
	repeating     bool
}

type colorStop struct {
	positions []css_ast.Token
	color     css_ast.Token
	midpoint  css_ast.Token // Absent if "midpoint.Kind == css_lexer.T(0)"
}

func parseGradient(token css_ast.Token) (gradient parsedGradient, success bool) {
	if token.Kind != css_lexer.TFunction {
		return
	}

	switch strings.ToLower(token.Text) {
	case "linear-gradient":
		gradient.kind = linearGradient

	case "radial-gradient":
		gradient.kind = radialGradient

	case "conic-gradient":
		gradient.kind = conicGradient

	case "repeating-linear-gradient":
		gradient.kind = linearGradient
		gradient.repeating = true

	case "repeating-radial-gradient":
		gradient.kind = radialGradient
		gradient.repeating = true

	case "repeating-conic-gradient":
		gradient.kind = conicGradient
		gradient.repeating = true

	default:
		return
	}

	// Bail if any token is a "var()" since it may introduce commas
	tokens := *token.Children
	for _, t := range tokens {
		if t.Kind == css_lexer.TFunction && strings.EqualFold(t.Text, "var") {
			return
		}
	}

	// Try to strip the initial tokens
	if len(tokens) > 0 && !looksLikeColor(tokens[0]) {
		i := 0
		for i < len(tokens) && tokens[i].Kind != css_lexer.TComma {
			i++
		}
		gradient.leadingTokens = tokens[:i]
		if i < len(tokens) {
			tokens = tokens[i+1:]
		} else {
			tokens = nil
		}
	}

	// Try to parse the color stops
	for len(tokens) > 0 {
		// Parse the color
		color := tokens[0]
		if !looksLikeColor(color) {
			return
		}
		tokens = tokens[1:]

		// Parse up to two positions
		var positions []css_ast.Token
		for len(positions) < 2 && len(tokens) > 0 {
			position := tokens[0]
			if position.Kind.IsNumeric() || (position.Kind == css_lexer.TFunction && strings.EqualFold(position.Text, "calc")) {
				positions = append(positions, position)
			} else {
				break
			}
			tokens = tokens[1:]
		}

		// Parse the comma
		var midpoint css_ast.Token
		if len(tokens) > 0 {
			if tokens[0].Kind != css_lexer.TComma {
				return
			}
			tokens = tokens[1:]
			if len(tokens) == 0 {
				return
			}

			// Parse the midpoint, if any
			if len(tokens) > 0 && tokens[0].Kind.IsNumeric() {
				midpoint = tokens[0]
				tokens = tokens[1:]

				// Followed by a mandatory comma
				if len(tokens) == 0 || tokens[0].Kind != css_lexer.TComma {
					return
				}
				tokens = tokens[1:]
			}
		}

		// Add the color stop
		gradient.colorStops = append(gradient.colorStops, colorStop{
			color:     color,
			positions: positions,
			midpoint:  midpoint,
		})
	}

	success = true
	return
}

func (p *parser) generateGradient(token css_ast.Token, gradient parsedGradient) css_ast.Token {
	var children []css_ast.Token
	commaToken := p.commaToken(token.Loc)

	children = append(children, gradient.leadingTokens...)
	for _, stop := range gradient.colorStops {
		if len(children) > 0 {
			children = append(children, commaToken)
		}
		if len(stop.positions) == 0 && stop.midpoint.Kind == css_lexer.T(0) {
			stop.color.Whitespace &= ^css_ast.WhitespaceAfter
		}
		children = append(children, stop.color)
		children = append(children, stop.positions...)
		if stop.midpoint.Kind != css_lexer.T(0) {
			children = append(children, commaToken, stop.midpoint)
		}
	}

	token.Children = &children
	return token
}

func (p *parser) lowerAndMinifyGradient(token css_ast.Token, wouldClipColor *bool) css_ast.Token {
	gradient, ok := parseGradient(token)
	if !ok {
		return token
	}

	lowerMidpoints := p.options.unsupportedCSSFeatures.Has(compat.GradientMidpoints)
	lowerColorSpaces := p.options.unsupportedCSSFeatures.Has(compat.ColorFunctions)
	lowerInterpolation := p.options.unsupportedCSSFeatures.Has(compat.GradientInterpolation)

	// Assume that if the browser doesn't support color spaces in gradients, then
	// it doesn't correctly interpolate non-sRGB colors even when a color space
	// is not specified. This is the case for Firefox 120, for example, which has
	// support for the "color()" syntax but not for color spaces in gradients.
	// There is no entry in our feature support matrix for this edge case so we
	// make this assumption instead.
	//
	// Note that this edge case means we have to _replace_ the original gradient
	// with the expanded one instead of inserting a fallback before it. Otherwise
	// Firefox 120 would use the original gradient instead of the fallback because
	// it supports the syntax, but just renders it incorrectly.
	if lowerInterpolation {
		lowerColorSpaces = true
	}

	// Potentially expand the gradient to handle unsupported features
	didExpand := false
	if lowerMidpoints || lowerColorSpaces || lowerInterpolation {
		if colorStops, ok := tryToParseColorStops(gradient); ok {
			hasColorSpace := false
			hasMidpoint := false
			for _, stop := range colorStops {
				if stop.hasColorSpace {
					hasColorSpace = true
				}
				if stop.midpoint != nil {
					hasMidpoint = true
				}
			}
			remaining, colorSpace, hueMethod, hasInterpolation := removeColorInterpolation(gradient.leadingTokens)
			if (hasInterpolation && lowerInterpolation) || (hasColorSpace && lowerColorSpaces) || (hasMidpoint && lowerMidpoints) {
				if hasInterpolation {
					tryToExpandGradient(token.Loc, &gradient, colorStops, remaining, colorSpace, hueMethod)
				} else {
					if hasColorSpace {
						colorSpace = colorSpace_oklab
					} else {
						colorSpace = colorSpace_srgb
					}
					tryToExpandGradient(token.Loc, &gradient, colorStops, gradient.leadingTokens, colorSpace, shorterHue)
				}
				didExpand = true
			}
		}
	}

	// Lower all colors in the gradient stop
	for i, stop := range gradient.colorStops {
		gradient.colorStops[i].color = p.lowerAndMinifyColor(stop.color, wouldClipColor)
	}

	if p.options.unsupportedCSSFeatures.Has(compat.GradientDoublePosition) {
		// Replace double positions with duplicated single positions
		for _, stop := range gradient.colorStops {
			if len(stop.positions) > 1 {
				gradient.colorStops = switchToSinglePositions(gradient.colorStops)
				break
			}
		}
	} else if p.options.minifySyntax {
		// Replace duplicated single positions with double positions
		for i, stop := range gradient.colorStops {
			if i > 0 && len(stop.positions) == 1 {
				if prev := gradient.colorStops[i-1]; len(prev.positions) == 1 && prev.midpoint.Kind == css_lexer.T(0) &&
					css_ast.TokensEqual([]css_ast.Token{prev.color}, []css_ast.Token{stop.color}, nil) {
					gradient.colorStops = switchToDoublePositions(gradient.colorStops)
					break
				}
			}
		}
	}

	if p.options.minifySyntax || didExpand {
		gradient.colorStops = removeImpliedPositions(gradient.kind, gradient.colorStops)
	}

	return p.generateGradient(token, gradient)
}

func removeImpliedPositions(kind gradientKind, colorStops []colorStop) []colorStop {
	if len(colorStops) == 0 {
		return colorStops
	}

	positions := make([]valueWithUnit, len(colorStops))
	for i, stop := range colorStops {
		if len(stop.positions) == 1 {
			if pos, ok := tryToParseValue(stop.positions[0], kind); ok {
				positions[i] = pos
				continue
			}
		}
		positions[i].value = math.NaN()
	}

	start := 0
	for start < len(colorStops) {
		if startPos := positions[start]; !math.IsNaN(startPos.value) {
			end := start + 1
		run:
			for colorStops[end-1].midpoint.Kind == css_lexer.T(0) && end < len(colorStops) {
				endPos := positions[end]
				if math.IsNaN(endPos.value) || endPos.unit != startPos.unit {
					break
				}

				// Check that all values in this run are implied. Interpolation is done
				// using the start and end positions instead of the first and second
				// positions because it's more accurate.
				for i := start + 1; i < end; i++ {
					t := float64(i-start) / float64(end-start)
					impliedValue := startPos.value + (endPos.value-startPos.value)*t
					if math.Abs(positions[i].value-impliedValue) > 0.01 {
						break run
					}
				}
				end++
			}

			// Clear out all implied values
			if end-start > 1 {
				for i := start + 1; i+1 < end; i++ {
					colorStops[i].positions = nil
				}
				start = end - 1
				continue
			}
		}
		start++
	}

	if first := colorStops[0].positions; len(first) == 1 &&
		((first[0].Kind == css_lexer.TPercentage && first[0].PercentageValue() == "0") ||
			(first[0].Kind == css_lexer.TDimension && first[0].DimensionValue() == "0")) {
		colorStops[0].positions = nil
	}

	if last := colorStops[len(colorStops)-1].positions; len(last) == 1 &&
		last[0].Kind == css_lexer.TPercentage && last[0].PercentageValue() == "100" {
		colorStops[len(colorStops)-1].positions = nil
	}

	return colorStops
}

func switchToSinglePositions(double []colorStop) (single []colorStop) {
	for _, stop := range double {
		for i := range stop.positions {
			stop.positions[i].Whitespace = css_ast.WhitespaceBefore
		}
		for len(stop.positions) > 1 {
			clone := stop
			clone.positions = stop.positions[:1]
			clone.midpoint = css_ast.Token{}
			single = append(single, clone)
			stop.positions = stop.positions[1:]
		}
		single = append(single, stop)
	}
	return
}

func switchToDoublePositions(single []colorStop) (double []colorStop) {
	for i := 0; i < len(single); i++ {
		stop := single[i]
		if i+1 < len(single) && len(stop.positions) == 1 && stop.midpoint.Kind == css_lexer.T(0) {
			if next := single[i+1]; len(next.positions) == 1 &&
				css_ast.TokensEqual([]css_ast.Token{stop.color}, []css_ast.Token{next.color}, nil) {
				double = append(double, colorStop{
					color:     stop.color,
					positions: []css_ast.Token{stop.positions[0], next.positions[0]},
					midpoint:  next.midpoint,
				})
				i++
				continue
			}
		}
		double = append(double, stop)
	}
	return
}

func removeColorInterpolation(tokens []css_ast.Token) ([]css_ast.Token, colorSpace, hueMethod, bool) {
	for i := 0; i+1 < len(tokens); i++ {
		if in := tokens[i]; in.Kind == css_lexer.TIdent && strings.EqualFold(in.Text, "in") {
			if space := tokens[i+1]; space.Kind == css_lexer.TIdent {
				var colorSpace colorSpace
				hueMethod := shorterHue
				start := i
				end := i + 2

				// Parse the color space
				switch strings.ToLower(space.Text) {
				case "a98-rgb":
					colorSpace = colorSpace_a98_rgb
				case "display-p3":
					colorSpace = colorSpace_display_p3
				case "hsl":
					colorSpace = colorSpace_hsl
				case "hwb":
					colorSpace = colorSpace_hwb
				case "lab":
					colorSpace = colorSpace_lab
				case "lch":
					colorSpace = colorSpace_lch
				case "oklab":
					colorSpace = colorSpace_oklab
				case "oklch":
					colorSpace = colorSpace_oklch
				case "prophoto-rgb":
					colorSpace = colorSpace_prophoto_rgb
				case "rec2020":
					colorSpace = colorSpace_rec2020
				case "srgb":
					colorSpace = colorSpace_srgb
				case "srgb-linear":
					colorSpace = colorSpace_srgb_linear
				case "xyz":
					colorSpace = colorSpace_xyz
				case "xyz-d50":
					colorSpace = colorSpace_xyz_d50
				case "xyz-d65":
					colorSpace = colorSpace_xyz_d65
				default:
					return nil, 0, 0, false
				}

				// Parse the optional hue mode for polar color spaces
				if colorSpace.isPolar() && i+3 < len(tokens) {
					if hue := tokens[i+3]; hue.Kind == css_lexer.TIdent && strings.EqualFold(hue.Text, "hue") {
						if method := tokens[i+2]; method.Kind == css_lexer.TIdent {
							switch strings.ToLower(method.Text) {
							case "shorter":
								hueMethod = shorterHue
							case "longer":
								hueMethod = longerHue
							case "increasing":
								hueMethod = increasingHue
							case "decreasing":
								hueMethod = decreasingHue
							default:
								return nil, 0, 0, false
							}
							end = i + 4
						}
					}
				}

				// Remove all parsed tokens
				remaining := append(append([]css_ast.Token{}, tokens[:start]...), tokens[end:]...)
				if n := len(remaining); n > 0 {
					remaining[0].Whitespace &= ^css_ast.WhitespaceBefore
					remaining[n-1].Whitespace &= ^css_ast.WhitespaceAfter
				}
				return remaining, colorSpace, hueMethod, true
			}
		}
	}

	return nil, 0, 0, false
}

type valueWithUnit struct {
	unit  string
	value float64
}

type parsedColorStop struct {
	// Position information (may be a sum of two different units)
	positionTerms []valueWithUnit

	// Color midpoint (a.k.a. transition hint) information
	midpoint *valueWithUnit

	// Non-premultiplied color information in XYZ space
	x, y, z, alpha float64

	// Non-premultiplied color information in sRGB space
	r, g, b float64

	// Premultiplied color information in the interpolation color space
	v0, v1, v2 float64

	// True if the original color has a color space
	hasColorSpace bool
}

func tryToParseColorStops(gradient parsedGradient) ([]parsedColorStop, bool) {
	var colorStops []parsedColorStop

	for _, stop := range gradient.colorStops {
		color, ok := parseColor(stop.color)
		if !ok {
			return nil, false
		}
		var r, g, b float64
		if !color.hasColorSpace {
			r = float64(hexR(color.hex)) / 255
			g = float64(hexG(color.hex)) / 255
			b = float64(hexB(color.hex)) / 255
			color.x, color.y, color.z = lin_srgb_to_xyz(lin_srgb(r, g, b))
		} else {
			r, g, b = gam_srgb(xyz_to_lin_srgb(color.x, color.y, color.z))
		}
		parsedStop := parsedColorStop{
			x:             color.x,
			y:             color.y,
			z:             color.z,
			r:             r,
			g:             g,
			b:             b,
			alpha:         float64(hexA(color.hex)) / 255,
			hasColorSpace: color.hasColorSpace,
		}

		for i, position := range stop.positions {
			if position, ok := tryToParseValue(position, gradient.kind); ok {
				parsedStop.positionTerms = []valueWithUnit{position}
			} else {
				return nil, false
			}

			// Expand double positions
			if i+1 < len(stop.positions) {
				colorStops = append(colorStops, parsedStop)
			}
		}

		if stop.midpoint.Kind != css_lexer.T(0) {
			if midpoint, ok := tryToParseValue(stop.midpoint, gradient.kind); ok {
				parsedStop.midpoint = &midpoint
			} else {
				return nil, false
			}
		}

		colorStops = append(colorStops, parsedStop)
	}

	// Automatically fill in missing positions
	if len(colorStops) > 0 {
		type stopInfo struct {
			fromPos   valueWithUnit
			toPos     valueWithUnit
			fromCount int32
			toCount   int32
		}

		// Fill in missing positions for the endpoints first
		if first := &colorStops[0]; len(first.positionTerms) == 0 {
			first.positionTerms = []valueWithUnit{{value: 0, unit: "%"}}
		}
		if last := &colorStops[len(colorStops)-1]; len(last.positionTerms) == 0 {
			last.positionTerms = []valueWithUnit{{value: 100, unit: "%"}}
		}

		// Set all positions to be greater than the position before them
		for i, stop := range colorStops {
			var prevPos valueWithUnit
			for j := i - 1; j >= 0; j-- {
				prev := colorStops[j]
				if prev.midpoint != nil {
					prevPos = *prev.midpoint
					break
				}
				if len(prev.positionTerms) == 1 {
					prevPos = prev.positionTerms[0]
					break
				}
			}
			if len(stop.positionTerms) == 1 {
				if prevPos.unit == stop.positionTerms[0].unit {
					stop.positionTerms[0].value = math.Max(prevPos.value, stop.positionTerms[0].value)
				}
				prevPos = stop.positionTerms[0]
			}
			if stop.midpoint != nil && prevPos.unit == stop.midpoint.unit {
				stop.midpoint.value = math.Max(prevPos.value, stop.midpoint.value)
			}
		}

		// Scan over all other stops with missing positions
		infos := make([]stopInfo, len(colorStops))
		for i, stop := range colorStops {
			if len(stop.positionTerms) == 1 {
				continue
			}
			info := &infos[i]

			// Scan backward
			for from := i - 1; from >= 0; from-- {
				fromStop := colorStops[from]
				info.fromCount++
				if fromStop.midpoint != nil {
					info.fromPos = *fromStop.midpoint
					break
				}
				if len(fromStop.positionTerms) == 1 {
					info.fromPos = fromStop.positionTerms[0]
					break
				}
			}

			// Scan forward
			for to := i; to < len(colorStops); to++ {
				info.toCount++
				if toStop := colorStops[to]; toStop.midpoint != nil {
					info.toPos = *toStop.midpoint
					break
				}
				if to+1 < len(colorStops) {
					if toStop := colorStops[to+1]; len(toStop.positionTerms) == 1 {
						info.toPos = toStop.positionTerms[0]
						break
					}
				}
			}
		}

		// Then fill in all other missing positions
		for i, stop := range colorStops {
			if len(stop.positionTerms) != 1 {
				info := infos[i]
				t := float64(info.fromCount) / float64(info.fromCount+info.toCount)
				if info.fromPos.unit == info.toPos.unit {
					colorStops[i].positionTerms = []valueWithUnit{{
						value: info.fromPos.value + (info.toPos.value-info.fromPos.value)*t,
						unit:  info.fromPos.unit,
					}}
				} else {
					colorStops[i].positionTerms = []valueWithUnit{{
						value: info.fromPos.value * (1 - t),
						unit:  info.fromPos.unit,
					}, {
						value: info.toPos.value * t,
						unit:  info.toPos.unit,
					}}
				}
			}
		}

		// Midpoints are only supported if they use the same units as their neighbors
		for i, stop := range colorStops {
			if stop.midpoint != nil {
				next := colorStops[i+1]
				if len(stop.positionTerms) != 1 || stop.midpoint.unit != stop.positionTerms[0].unit ||
					len(next.positionTerms) != 1 || stop.midpoint.unit != next.positionTerms[0].unit {
					return nil, false
				}
			}
		}
	}

	return colorStops, true
}

func tryToParseValue(token css_ast.Token, kind gradientKind) (result valueWithUnit, success bool) {
	if kind == conicGradient {
		// <angle-percentage>
		switch token.Kind {
		case css_lexer.TDimension:
			degrees, ok := degreesForAngle(token)
			if !ok {
				return
			}
			result.value = degrees * (100.0 / 360)
			result.unit = "%"

		case css_lexer.TPercentage:
			percent, err := strconv.ParseFloat(token.PercentageValue(), 64)
			if err != nil {
				return
			}
			result.value = percent
			result.unit = "%"

		default:
			return
		}
	} else {
		// <length-percentage>
		switch token.Kind {
		case css_lexer.TNumber:
			zero, err := strconv.ParseFloat(token.Text, 64)
			if err != nil || zero != 0 {
				return
			}
			result.value = 0
			result.unit = "%"

		case css_lexer.TDimension:
			dimensionValue, err := strconv.ParseFloat(token.DimensionValue(), 64)
			if err != nil {
				return
			}
			result.value = dimensionValue
			result.unit = token.DimensionUnit()

		case css_lexer.TPercentage:
			percentageValue, err := strconv.ParseFloat(token.PercentageValue(), 64)
			if err != nil {
				return
			}
			result.value = percentageValue
			result.unit = "%"

		default:
			return
		}
	}

	success = true
	return
}

func tryToExpandGradient(
	loc logger.Loc,
	gradient *parsedGradient,
	colorStops []parsedColorStop,
	remaining []css_ast.Token,
	colorSpace colorSpace,
	hueMethod hueMethod,
) bool {
	// Convert color stops into the interpolation color space
	for i := range colorStops {
		stop := &colorStops[i]
		v0, v1, v2 := xyz_to_colorSpace(stop.x, stop.y, stop.z, colorSpace)
		stop.v0, stop.v1, stop.v2 = premultiply(v0, v1, v2, stop.alpha, colorSpace)
	}

	// Duplicate the endpoints if they should wrap around to themselves
	if hueMethod == longerHue && colorSpace.isPolar() && len(colorStops) > 0 {
		if first := colorStops[0]; len(first.positionTerms) == 1 {
			if first.positionTerms[0].value < 0 {
				colorStops[0].positionTerms[0].value = 0
			} else if first.positionTerms[0].value > 0 {
				first.midpoint = nil
				first.positionTerms = []valueWithUnit{{value: 0, unit: first.positionTerms[0].unit}}
				colorStops = append([]parsedColorStop{first}, colorStops...)
			}
		}
		if last := colorStops[len(colorStops)-1]; len(last.positionTerms) == 1 {
			if last.positionTerms[0].unit != "%" || last.positionTerms[0].value < 100 {
				last.positionTerms = []valueWithUnit{{value: 100, unit: "%"}}
				colorStops = append(colorStops, last)
			}
		}
	}

	var newColorStops []colorStop
	var generateColorStops func(
		int, parsedColorStop, parsedColorStop,
		float64, float64, float64, float64, float64, float64, float64, float64,
		float64, float64, float64, float64, float64, float64, float64, float64,
	)

	generateColorStops = func(
		depth int,
		from parsedColorStop, to parsedColorStop,
		prevX, prevY, prevZ, prevR, prevG, prevB, prevA, prevT float64,
		nextX, nextY, nextZ, nextR, nextG, nextB, nextA, nextT float64,
	) {
		if depth > 4 {
			return
		}

		t := (prevT + nextT) / 2
		positionT := t

		// Handle midpoints (which we have already checked uses the same units)
		if from.midpoint != nil {
			fromPos := from.positionTerms[0].value
			toPos := to.positionTerms[0].value
			stopPos := fromPos + (toPos-fromPos)*t
			H := (from.midpoint.value - fromPos) / (toPos - fromPos)
			P := (stopPos - fromPos) / (toPos - fromPos)
			if H <= 0 {
				positionT = 1
			} else if H >= 1 {
				positionT = 0
			} else {
				positionT = math.Pow(P, -1/math.Log2(H))
			}
		}

		v0, v1, v2 := interpolateColors(from.v0, from.v1, from.v2, to.v0, to.v1, to.v2, colorSpace, hueMethod, positionT)
		a := from.alpha + (to.alpha-from.alpha)*positionT
		v0, v1, v2 = unpremultiply(v0, v1, v2, a, colorSpace)
		x, y, z := colorSpace_to_xyz(v0, v1, v2, colorSpace)

		// Stop when the color is similar enough to the sRGB midpoint
		const epsilon = 4.0 / 255
		r, g, b := gam_srgb(xyz_to_lin_srgb(x, y, z))
		dr := r*a - (prevR*prevA+nextR*nextA)/2
		dg := g*a - (prevG*prevA+nextG*nextA)/2
		db := b*a - (prevB*prevA+nextB*nextA)/2
		if d := dr*dr + dg*dg + db*db; d < epsilon*epsilon {
			return
		}

		// Recursive split before this stop
		generateColorStops(depth+1, from, to,
			prevX, prevY, prevZ, prevR, prevG, prevB, prevA, prevT,
			x, y, z, r, g, b, a, t)

		// Generate this stop
		color := makeColorToken(loc, x, y, z, a)
		positionTerms := interpolatePositions(from.positionTerms, to.positionTerms, t)
		position := makePositionToken(loc, positionTerms)
		position.Whitespace = css_ast.WhitespaceBefore
		newColorStops = append(newColorStops, colorStop{
			color:     color,
			positions: []css_ast.Token{position},
		})

		// Recursive split after this stop
		generateColorStops(depth+1, from, to,
			x, y, z, r, g, b, a, t,
			nextX, nextY, nextZ, nextR, nextG, nextB, nextA, nextT)
	}

	for i, stop := range colorStops {
		color := makeColorToken(loc, stop.x, stop.y, stop.z, stop.alpha)
		position := makePositionToken(loc, stop.positionTerms)
		position.Whitespace = css_ast.WhitespaceBefore
		newColorStops = append(newColorStops, colorStop{
			color:     color,
			positions: []css_ast.Token{position},
		})

		// Generate new color stops in between as needed
		if i+1 < len(colorStops) {
			next := colorStops[i+1]
			generateColorStops(0, stop, next,
				stop.x, stop.y, stop.z, stop.r, stop.g, stop.b, stop.alpha, 0,
				next.x, next.y, next.z, next.r, next.g, next.b, next.alpha, 1)
		}
	}

	gradient.leadingTokens = remaining
	gradient.colorStops = newColorStops
	return true
}

func formatFloat(value float64, decimals int) string {
	return strings.TrimSuffix(strings.TrimRight(strconv.FormatFloat(value, 'f', decimals, 64), "0"), ".")
}

func makeDimensionOrPercentToken(loc logger.Loc, value float64, unit string) (token css_ast.Token) {
	token.Loc = loc
	token.Text = formatFloat(value, 2)
	if unit == "%" {
		token.Kind = css_lexer.TPercentage
	} else {
		token.Kind = css_lexer.TDimension
		token.UnitOffset = uint16(len(token.Text))
	}
	token.Text += unit
	return
}

func makePositionToken(loc logger.Loc, positionTerms []valueWithUnit) css_ast.Token {
	if len(positionTerms) == 1 {
		return makeDimensionOrPercentToken(loc, positionTerms[0].value, positionTerms[0].unit)
	}

	children := make([]css_ast.Token, 0, 1+2*len(positionTerms))
	for i, term := range positionTerms {
		if i > 0 {
			children = append(children, css_ast.Token{
				Loc:        loc,
				Kind:       css_lexer.TDelimPlus,
				Text:       "+",
				Whitespace: css_ast.WhitespaceBefore | css_ast.WhitespaceAfter,
			})
		}
		children = append(children, makeDimensionOrPercentToken(loc, term.value, term.unit))
	}

	return css_ast.Token{
		Loc:      loc,
		Kind:     css_lexer.TFunction,
		Text:     "calc",
		Children: &children,
	}
}

func makeColorToken(loc logger.Loc, x float64, y float64, z float64, a float64) (color css_ast.Token) {
	color.Loc = loc
	alpha := uint32(math.Round(a * 255))
	if hex, ok := tryToConvertToHexWithoutClipping(x, y, z, alpha); ok {
		color.Kind = css_lexer.THash
		if alpha == 255 {
			color.Text = fmt.Sprintf("%06x", hex>>8)
		} else {
			color.Text = fmt.Sprintf("%08x", hex)
		}
	} else {
		children := []css_ast.Token{
			{
				Loc:        loc,
				Kind:       css_lexer.TIdent,
				Text:       "xyz",
				Whitespace: css_ast.WhitespaceAfter,
			},
			{
				Loc:        loc,
				Kind:       css_lexer.TNumber,
				Text:       formatFloat(x, 3),
				Whitespace: css_ast.WhitespaceBefore | css_ast.WhitespaceAfter,
			},
			{
				Loc:        loc,
				Kind:       css_lexer.TNumber,
				Text:       formatFloat(y, 3),
				Whitespace: css_ast.WhitespaceBefore | css_ast.WhitespaceAfter,
			},
			{
				Loc:        loc,
				Kind:       css_lexer.TNumber,
				Text:       formatFloat(z, 3),
				Whitespace: css_ast.WhitespaceBefore,
			},
		}
		if a < 1 {
			children = append(children,
				css_ast.Token{
					Loc:        loc,
					Kind:       css_lexer.TDelimSlash,
					Text:       "/",
					Whitespace: css_ast.WhitespaceBefore | css_ast.WhitespaceAfter,
				},
				css_ast.Token{
					Loc:        loc,
					Kind:       css_lexer.TNumber,
					Text:       formatFloat(a, 3),
					Whitespace: css_ast.WhitespaceBefore,
				},
			)
		}
		color.Kind = css_lexer.TFunction
		color.Text = "color"
		color.Children = &children
	}
	return
}

func interpolateHues(a, b, t float64, hueMethod hueMethod) float64 {
	a /= 360
	b /= 360
	a -= math.Floor(a)
	b -= math.Floor(b)

	switch hueMethod {
	case shorterHue:
		delta := b - a
		if delta > 0.5 {
			a++
		}
		if delta < -0.5 {
			b++
		}

	case longerHue:
		delta := b - a
		if delta > 0 && delta < 0.5 {
			a++
		}
		if delta > -0.5 && delta <= 0 {
			b++
		}

	case increasingHue:
		if b < a {
			b++
		}

	case decreasingHue:
		if a < b {
			a++
		}
	}

	return (a + (b-a)*t) * 360
}

func interpolateColors(
	a0, a1, a2 float64, b0, b1, b2 float64,
	colorSpace colorSpace, hueMethod hueMethod, t float64,
) (v0 float64, v1 float64, v2 float64) {
	v1 = a1 + (b1-a1)*t

	switch colorSpace {
	case colorSpace_hsl, colorSpace_hwb:
		v2 = a2 + (b2-a2)*t
		v0 = interpolateHues(a0, b0, t, hueMethod)

	case colorSpace_lch, colorSpace_oklch:
		v0 = a0 + (b0-a0)*t
		v2 = interpolateHues(a2, b2, t, hueMethod)

	default:
		v0 = a0 + (b0-a0)*t
		v2 = a2 + (b2-a2)*t
	}

	return v0, v1, v2
}

func interpolatePositions(a []valueWithUnit, b []valueWithUnit, t float64) (result []valueWithUnit) {
	findUnit := func(unit string) int {
		for i, x := range result {
			if x.unit == unit {
				return i
			}
		}
		result = append(result, valueWithUnit{unit: unit})
		return len(result) - 1
	}

	// "result += a * (1 - t)"
	for _, term := range a {
		result[findUnit(term.unit)].value += term.value * (1 - t)
	}

	// "result += b * t"
	for _, term := range b {
		result[findUnit(term.unit)].value += term.value * t
	}

	// Remove an extra zero value for neatness. We don't remove all
	// of them because it may be important to retain a single zero.
	if len(result) > 1 {
		for i, term := range result {
			if term.value == 0 {
				copy(result[i:], result[i+1:])
				result = result[:len(result)-1]
				break
			}
		}
	}

	return
}

func premultiply(v0, v1, v2, alpha float64, colorSpace colorSpace) (float64, float64, float64) {
	if alpha < 1 {
		switch colorSpace {
		case colorSpace_hsl, colorSpace_hwb:
			v2 *= alpha
		case colorSpace_lch, colorSpace_oklch:
			v0 *= alpha
		default:
			v0 *= alpha
			v2 *= alpha
		}
		v1 *= alpha
	}
	return v0, v1, v2
}

func unpremultiply(v0, v1, v2, alpha float64, colorSpace colorSpace) (float64, float64, float64) {
	if alpha > 0 && alpha < 1 {
		switch colorSpace {
		case colorSpace_hsl, colorSpace_hwb:
			v2 /= alpha
		case colorSpace_lch, colorSpace_oklch:
			v0 /= alpha
		default:
			v0 /= alpha
			v2 /= alpha
		}
		v1 /= alpha
	}
	return v0, v1, v2
}
