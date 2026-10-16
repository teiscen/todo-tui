package backend

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// HexCode is a "#RRGGBB" color. Turning one into a lipgloss style is the job
// of the ui package, so backend has no rendering dependencies.
type HexCode string

func (c HexCode) HexCodeToInt() (int, error) {
	str := strings.TrimPrefix(string(c), "#")

	val32, err := strconv.ParseInt(str, 16, 32)
	if err != nil {
		return 0, err
	}
	return int(val32), nil
}

func (c HexCode) ParseRGB() (r int, g int, b int, err error) {
	colorInt, err := c.HexCodeToInt()
	if err != nil {
		return 0, 0, 0, err
	}

	// Not neccessary but keeping it for parity with Green and Blue
	r = (colorInt & 0xFF0000) >> 16
	g = (colorInt & 0x00FF00) >> 8
	b = (colorInt & 0x0000FF) >> 0

	return r, g, b, nil
}

func CombineRGB(r, g, b int) (HexCode, error) {
	isValid := func(i int) bool {
		return i >= 0 && i <= 0xFF
	}
	if !isValid(r) || !isValid(g) || !isValid(b) {
		return HexCode("#FFFFFF"), fmt.Errorf("rgb value must be between 0 and 255")
	}

	combined := (r << 16) | (g << 8) | b

	// %06X:0-pad left with 0 not space, 6-how much to pad, X-convert int to Hex
	return HexCode(fmt.Sprintf("#%06X", combined)), nil
}

// Averages out the colors, but inteded to be used with a muted color
// such as #808080
func (c HexCode) MuteColorTarget(target HexCode, strenght float64) (HexCode, error) {
	r, g, b, err := c.ParseRGB()
	if err != nil {
		return HexCode("#FFFFFF"), err
	}

	tr, tg, tb, err := target.ParseRGB()
	if err != nil {
		return HexCode("#FFFFFF"), err
	}

	lerp := func(from, to int) int {
		return int(math.Round(float64(from) + float64(to-from)*strenght))
	}

	return CombineRGB(lerp(r, tr), lerp(g, tg), lerp(b, tb))
}
