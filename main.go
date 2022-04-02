package main

import "github.com/gonutz/prototype/draw"

func main() {
	fullscreen := false
	draw.RunWindow("Delay the Inevitable", 1000, 800, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
			return
		}

		if window.WasKeyPressed(draw.KeyF11) {
			fullscreen = !fullscreen
		}
		window.SetFullscreen(fullscreen)

		window.DrawText("ESC to quit\nF11 to toggle fullscreen", 0, 0, draw.White)
	})
}
