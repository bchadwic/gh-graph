package color

import (
	"fmt"

	"github.com/bchadwic/gh-graph/pkg/stats"
	lg "github.com/charmbracelet/lipgloss"
)

const (
	Catagories     = 5
	GroupFormRate  = 2
	DefaultBestDay = 100
)

type ColorPalette struct {
	Colors []Color
	Limits []int
}

type Color struct {
	R uint8
	G uint8
	B uint8
}

func (cp *ColorPalette) Initialize(stats *stats.Stats) *ColorPalette {
	cp.Limits = make([]int, Catagories)

	max := DefaultBestDay
	if max > stats.BestDay {
		max = stats.BestDay
	}

	curr := max
	for i := Catagories - 1; i >= 0; i-- {
		cp.Limits[i] = curr
		curr = (curr / GroupFormRate) + 1
	}

	if lg.HasDarkBackground() {
		for i := 0; i < Catagories; i++ {
			cp.Colors = append(cp.Colors, Color{
				R: 30,
				G: uint8(i+1) * 50,
				B: 30,
			})
		}
	} else {
		for i := Catagories - 1; i >= 0; i-- {
			cp.Colors = append(cp.Colors, Color{
				R: 30,
				G: uint8(i+1) * 50,
				B: 30,
			})
		}
	}

	return cp
}

func (cp *ColorPalette) GetColor(colorIndex int) string {
	for i, e := range cp.Limits {
		if colorIndex < e {
			return cp.Colors[i].Hex()
		}
	}
	return cp.Colors[Catagories-1].Hex()
}

func (c *Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
