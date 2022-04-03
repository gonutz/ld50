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

var (
	globalFullscreen = false
	globalMenu       = newMenu()
	globalControls   = newControlsMenu()
	globalHighscores = newHighscoresMenu()
)

func toggleFullscreen(window draw.Window) {
	globalFullscreen = !globalFullscreen
}

func main() {
	rand.Seed(time.Now().UnixNano())

	type settings struct {
		Fullscreen   bool
		LeftKey      int
		RightKey     int
		DropOneKey   int
		DropAllKey   int
		RotateCWKey  int
		RotateCCWKey int
		Highscores   []int
	}
	settingsPath := filepath.Join(os.Getenv("APPDATA"), "ld50.settings")
	if data, err := os.ReadFile(settingsPath); err == nil {
		var s settings
		if err := json.Unmarshal(data, &s); err == nil {
			globalFullscreen = s.Fullscreen
			setKey(&globalControls.left, s.LeftKey)
			setKey(&globalControls.right, s.RightKey)
			setKey(&globalControls.dropOne, s.DropOneKey)
			setKey(&globalControls.dropAll, s.DropAllKey)
			setKey(&globalControls.rotateCW, s.RotateCWKey)
			setKey(&globalControls.rotateCCW, s.RotateCCWKey)
			globalHighscores.scores = s.Highscores
		}
	}
	defer func() {
		s := settings{
			Fullscreen:   globalFullscreen,
			LeftKey:      int(globalControls.left),
			RightKey:     int(globalControls.right),
			DropOneKey:   int(globalControls.dropOne),
			DropAllKey:   int(globalControls.dropAll),
			RotateCWKey:  int(globalControls.rotateCW),
			RotateCCWKey: int(globalControls.rotateCCW),
			Highscores:   globalHighscores.scores,
		}
		if data, err := json.Marshal(&s); err == nil {
			os.WriteFile(settingsPath, data, 0666)
		}
	}()

	var mode gameMode = globalMenu

	draw.RunWindow("Schmetris - Delay the Inevitable", 1000, 800, func(window draw.Window) {
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

func setKey(key *draw.Key, to int) {
	k := draw.Key(to)
	for _, valid := range allKeys {
		if k == valid {
			*key = k
		}
	}
}
