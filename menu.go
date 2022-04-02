package main

import "github.com/gonutz/prototype/draw"

const (
	itemStartGame = iota
	itemToggleFullscreen
	itemExitGame
)

func newMenu() *menu {
	return &menu{
		items: []string{
			"Start Game",
			"Toggle Fullscreen (F11)",
			"Quit (Escape)",
		},
	}
}

type menu struct {
	selected int
	items    []string
}

func (m *menu) update(window draw.Window) gameMode {
	if window.WasKeyPressed(draw.KeyEscape) {
		window.Close()
		return m
	}

	if window.WasKeyPressed(draw.KeyEnter) {
		switch m.selected {
		case itemStartGame:
			return newBlocks()
		case itemToggleFullscreen:
			toggleFullscreen(window)
		case itemExitGame:
			window.Close()
			return m
		}
	}

	if window.WasKeyPressed(draw.KeyDown) {
		m.selected = (m.selected + 1) % len(m.items)
	}
	if window.WasKeyPressed(draw.KeyUp) {
		m.selected = (m.selected + len(m.items) - 1) % len(m.items)
	}
	if window.WasKeyPressed(draw.KeyHome) {
		m.selected = 0
	}
	if window.WasKeyPressed(draw.KeyEnd) {
		m.selected = len(m.items) - 1
	}

	windowW, windowH := window.Size()
	smallTextScale := float32(windowH) / 400
	largeTextScale := float32(windowH) / 300
	_, smallTextH := window.GetScaledTextSize("H", smallTextScale)
	_, largeTextH := window.GetScaledTextSize("H", largeTextScale)
	textMargin := smallTextH / 2
	totalTextH := (len(m.items)-1)*smallTextH + largeTextH + (len(m.items)-1)*textMargin
	y := (windowH - totalTextH) / 2
	for i, item := range m.items {
		scale := smallTextScale
		color := draw.RGB(0.8, 0.8, 0.9)
		if i == m.selected {
			item = ">   " + item + "   <"
			scale = largeTextScale
			color = draw.RGB(1, 0.8, 0.8)
		}
		textW, textH := window.GetScaledTextSize(item, scale)
		x := (windowW - textW) / 2
		window.DrawScaledText(item, x, y, scale, color)

		y += textH + textMargin
	}

	return m
}
