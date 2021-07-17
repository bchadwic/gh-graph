package main

import (
	lg "github.com/charmbracelet/lipgloss"
)

const (
	Catagories     = 5
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

func GetColors(stats *Stats) *ColorPalette {
	cp := &ColorPalette{
		Limits: make([]int, Catagories),
	}

	max := DefaultBestDay
	if max > stats.BestDay {
		max = stats.BestDay
	}

	// TODO : Find a better way to make catagories
	curr := max
	for i := 0; i < Catagories; i++ {
		cp.Limits[i] = curr
		curr /= Catagories
	}

	if lg.HasDarkBackground() {
		for i := 0; i < len(cp.Limits); i++ {
			cp.Colors = append(cp.Colors, Color{})
		}
	} else {
		for i := 0; i < len(cp.Limits); i++ {
			cp.Colors = append(cp.Colors, Color{})
		}
	}

	return nil
}
