package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/gonutz/prototype/draw"
)

func randInt() int {
	return rand.Int()
}

type gameMode interface {
	update(window draw.Window) gameMode
}

var globalFullscreen = false

func toggleFullscreen(window draw.Window) {
	globalFullscreen = !globalFullscreen
}

func main() {
	rand.Seed(time.Now().UnixNano())

	type settings struct {
		Fullscreen bool
	}
	settingsPath := filepath.Join(os.Getenv("APPDATA"), "ld50.settings")
	if data, err := os.ReadFile(settingsPath); err == nil {
		var s settings
		if err := json.Unmarshal(data, &s); err == nil {
			globalFullscreen = s.Fullscreen
		}
	}
	defer func() {
		s := settings{
			Fullscreen: globalFullscreen,
		}
		if data, err := json.Marshal(&s); err == nil {
			os.WriteFile(settingsPath, data, 0666)
		}
	}()

	var mode gameMode
	startGame := func() {
		//mode = newBlocks()
		mode = newMenu()
	}
	startGame()

	draw.RunWindow("Delay the Inevitable", 1000, 800, func(window draw.Window) {
		// F11 is not game state specific. It always toggles full screen.
		// Instead of every state needing to handle it, we handle it globally
		// here.
		if window.WasKeyPressed(draw.KeyF11) {
			toggleFullscreen(window)
		}

		window.SetFullscreen(globalFullscreen)
		window.ShowCursor(!globalFullscreen)

		mode = mode.update(window)
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
