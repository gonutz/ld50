package main

import (
	"fmt"
	"sort"

	"github.com/gonutz/prototype/draw"
)

func newHighscoresMenu() *highscoresMenu {
	return &highscoresMenu{}
}

type highscoresMenu struct {
	scores []int
}

func (m *highscoresMenu) addScore(score int) {
	if score == 0 {
		return
	}

	m.scores = append(m.scores, score)
	sort.Slice(m.scores, func(i, j int) bool {
		return m.scores[i] > m.scores[j]
	})
	const maxScores = 5
	if len(m.scores) > maxScores {
		m.scores = m.scores[:maxScores]
	}
}

func (m *highscoresMenu) update(window draw.Window) gameMode {
	if window.WasKeyPressed(draw.KeyEscape) ||
		window.WasKeyPressed(draw.KeyEnter) || window.WasKeyPressed(draw.KeyNumEnter) {
		return globalMenu
	}

	var text string
	for i, score := range m.scores {
		text += fmt.Sprintf("%d.   %5d\n", i+1, score)
	}
	if len(m.scores) == 0 {
		text += " No scores yet.\n"
	}

	windowW, windowH := window.Size()
	smallTextScale := float32(windowH) / 400
	largeTextScale := float32(windowH) / 300
	textW, textH := window.GetScaledTextSize(text, smallTextScale)
	x := (windowW - textW) / 2
	y := (windowH - textH) / 2
	window.DrawScaledText(text, x, y, smallTextScale, draw.RGB(0.8, 0.8, 0.9))

	text = ">   Back to menu   <"
	textW, _ = window.GetScaledTextSize(text, largeTextScale)
	window.DrawScaledText(
		text,
		(windowW-textW)/2,
		y+textH,
		largeTextScale,
		draw.RGB(1, 0.8, 0.8),
	)

	return m
}
