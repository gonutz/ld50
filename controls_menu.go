package main

import (
	"strings"

	"github.com/gonutz/prototype/draw"
)

func newControlsMenu() *controlsMenu {
	return &controlsMenu{
		left:      draw.KeyLeft,
		right:     draw.KeyRight,
		dropOne:   draw.KeyDown,
		dropAll:   draw.KeyS,
		rotateCW:  draw.KeyD,
		rotateCCW: draw.KeyA,
		editing:   -1,
	}
}

type controlsMenu struct {
	left      draw.Key
	right     draw.Key
	dropOne   draw.Key
	dropAll   draw.Key
	rotateCW  draw.Key
	rotateCCW draw.Key
	selected  int
	editing   int
}

func (m *controlsMenu) update(window draw.Window) gameMode {
	if window.WasKeyPressed(draw.KeyEscape) {
		if m.editing != -1 {
			m.editing = -1
		} else {
			return globalMenu
		}
	}

	// Be carful not to create another leftpad desaster.
	leftPad := func(s string) string {
		for len(s) < 12 {
			s = " " + s
		}
		return s
	}
	rightPad := func(s string) string {
		for len(s) < 8 {
			s = s + " "
		}
		return s
	}
	items := []string{
		leftPad(" Move Left: ") + rightPad(m.left.String()),
		leftPad("Move Right: ") + rightPad(m.right.String()),
		leftPad(" Move Down: ") + rightPad(m.dropOne.String()),
		leftPad("Drop Block: ") + rightPad(m.dropAll.String()),
		leftPad(" Rotate CW: ") + rightPad(m.rotateCW.String()),
		leftPad("Rotate CCW: ") + rightPad(m.rotateCCW.String()),
		"  Done  ",
	}
	if m.editing != -1 {
		s := items[m.editing]
		i := strings.Index(s, ": ")
		s = s[:i+2] + rightPad("...")
		items[m.editing] = s
	}

	if m.editing != -1 {
		for _, key := range allKeys {
			if window.WasKeyPressed(key) {
				keys := []*draw.Key{
					&m.left,
					&m.right,
					&m.dropOne,
					&m.dropAll,
					&m.rotateCW,
					&m.rotateCCW,
				}
				*keys[m.editing] = key
				m.editing = -1
				break
			}
		}
	} else {
		if window.WasKeyPressed(draw.KeyEnter) || window.WasKeyPressed(draw.KeyNumEnter) {
			if m.selected == len(items)-1 {
				return globalMenu
			}
			m.editing = m.selected
		}

		if window.WasKeyPressed(draw.KeyDown) {
			m.selected = (m.selected + 1) % len(items)
		}
		if window.WasKeyPressed(draw.KeyUp) {
			m.selected = (m.selected + len(items) - 1) % len(items)
		}
		if window.WasKeyPressed(draw.KeyHome) {
			m.selected = 0
		}
		if window.WasKeyPressed(draw.KeyEnd) {
			m.selected = len(items) - 1
		}
	}

	windowW, windowH := window.Size()
	smallTextScale := float32(windowH) / 400
	largeTextScale := float32(windowH) / 300
	_, smallTextH := window.GetScaledTextSize("H", smallTextScale)
	_, largeTextH := window.GetScaledTextSize("H", largeTextScale)
	textMargin := smallTextH / 2
	totalTextH := (len(items)-1)*smallTextH + largeTextH + (len(items)-1)*textMargin
	y := (windowH - totalTextH) / 2
	for i, item := range items {
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

var allKeys = []draw.Key{
	draw.KeyA,
	draw.KeyB,
	draw.KeyC,
	draw.KeyD,
	draw.KeyE,
	draw.KeyF,
	draw.KeyG,
	draw.KeyH,
	draw.KeyI,
	draw.KeyJ,
	draw.KeyK,
	draw.KeyL,
	draw.KeyM,
	draw.KeyN,
	draw.KeyO,
	draw.KeyP,
	draw.KeyQ,
	draw.KeyR,
	draw.KeyS,
	draw.KeyT,
	draw.KeyU,
	draw.KeyV,
	draw.KeyW,
	draw.KeyX,
	draw.KeyY,
	draw.KeyZ,
	draw.Key0,
	draw.Key1,
	draw.Key2,
	draw.Key3,
	draw.Key4,
	draw.Key5,
	draw.Key6,
	draw.Key7,
	draw.Key8,
	draw.Key9,
	draw.KeyNum0,
	draw.KeyNum1,
	draw.KeyNum2,
	draw.KeyNum3,
	draw.KeyNum4,
	draw.KeyNum5,
	draw.KeyNum6,
	draw.KeyNum7,
	draw.KeyNum8,
	draw.KeyNum9,
	draw.KeyF1,
	draw.KeyF2,
	draw.KeyF3,
	draw.KeyF4,
	draw.KeyF5,
	draw.KeyF6,
	draw.KeyF7,
	draw.KeyF8,
	draw.KeyF9,
	draw.KeyF10,
	draw.KeyF11,
	draw.KeyF12,
	draw.KeyF13,
	draw.KeyF14,
	draw.KeyF15,
	draw.KeyF16,
	draw.KeyF17,
	draw.KeyF18,
	draw.KeyF19,
	draw.KeyF20,
	draw.KeyF21,
	draw.KeyF22,
	draw.KeyF23,
	draw.KeyF24,
	draw.KeyEnter,
	draw.KeyNumEnter,
	draw.KeyLeftControl,
	draw.KeyRightControl,
	draw.KeyLeftShift,
	draw.KeyRightShift,
	draw.KeyLeftAlt,
	draw.KeyRightAlt,
	draw.KeyLeft,
	draw.KeyRight,
	draw.KeyUp,
	draw.KeyDown,
	draw.KeySpace,
	draw.KeyBackspace,
	draw.KeyTab,
	draw.KeyHome,
	draw.KeyEnd,
	draw.KeyPageDown,
	draw.KeyPageUp,
	draw.KeyDelete,
	draw.KeyInsert,
	draw.KeyNumAdd,
	draw.KeyNumSubtract,
	draw.KeyNumMultiply,
	draw.KeyNumDivide,
}
