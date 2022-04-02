package main

import "github.com/gonutz/prototype/draw"

const (
	blockFieldWidth  = 10
	blockFieldHeight = 18
)

func newBlocks() *blocks {
	return &blocks{
		thisPiece: randomTetromino(),
		nextPiece: randomTetromino(),
	}
}

type blocksField [blockFieldWidth][blockFieldHeight]blockKind

type blocks struct {
	thisPiece tetromino
	nextPiece tetromino
	field     blocksField
}

type blockKind int

const (
	blockEmpty blockKind = iota
	blockL
	blockJ
	blockO
	blockI
	blockS
	blockZ
	blockT
	firstBlock = blockL
	lastBlock  = blockT
)

func randomTetromino() tetromino {
	return tetromino{
		x:    3,
		y:    0,
		kind: firstBlock + blockKind(randInt())%(lastBlock-firstBlock+1),
	}
}

type tetromino struct {
	x, y     int
	kind     blockKind
	rotation int
}

type xy struct {
	x, y int
}

func (t *tetromino) parts() [4]xy {
	x, y := t.x, t.y
	switch t.kind {
	case blockL:
		switch t.rotation {
		case 0:
			return [4]xy{{x, y + 1}, {x + 1, y + 1}, {x + 2, y + 1}, {x, y + 2}}
		case 1:
			return [4]xy{{x, y}, {x + 1, y}, {x + 1, y + 1}, {x + 1, y + 2}}
		case 2:
			return [4]xy{{x, y + 1}, {x + 1, y + 1}, {x + 2, y + 1}, {x + 2, y}}
		case 3:
			return [4]xy{{x + 1, y}, {x + 1, y + 1}, {x + 1, y + 2}, {x + 2, y + 2}}
		}
	case blockJ:
		switch t.rotation {
		case 0:
			return [4]xy{{x, y + 1}, {x + 1, y + 1}, {x + 2, y + 1}, {x + 2, y + 2}}
		case 1:
			return [4]xy{{x + 1, y}, {x + 1, y + 1}, {x + 1, y + 2}, {x, y + 2}}
		case 2:
			return [4]xy{{x, y}, {x, y + 1}, {x + 1, y + 1}, {x + 2, y + 1}}
		case 3:
			return [4]xy{{x + 1, y}, {x + 2, y}, {x + 1, y + 1}, {x + 1, y + 2}}
		}
	case blockO:
		return [4]xy{{x + 1, y + 1}, {x + 2, y + 1}, {x + 1, y + 2}, {x + 2, y + 2}}
	case blockI:
		switch t.rotation {
		case 0, 2:
			return [4]xy{{x, y + 1}, {x + 1, y + 1}, {x + 2, y + 1}, {x + 3, y + 1}}
		case 1, 3:
			return [4]xy{{x + 1, y - 1}, {x + 1, y}, {x + 1, y + 1}, {x + 1, y + 2}}
		}
	case blockS:
		switch t.rotation {
		case 0, 2:
			return [4]xy{{x + 1, y + 1}, {x + 2, y + 1}, {x, y + 2}, {x + 1, y + 2}}
		case 1, 3:
			return [4]xy{{x, y}, {x, y + 1}, {x + 1, y + 1}, {x + 1, y + 2}}
		}
	case blockZ:
		switch t.rotation {
		case 0, 2:
			return [4]xy{{x, y + 1}, {x + 1, y + 1}, {x + 1, y + 2}, {x + 2, y + 2}}
		case 1, 3:
			return [4]xy{{x, y + 1}, {x, y + 2}, {x + 1, y}, {x + 1, y + 1}}
		}
	case blockT:
		switch t.rotation {
		case 0:
			return [4]xy{{x, y + 1}, {x + 1, y + 1}, {x + 2, y + 1}, {x + 1, y + 2}}
		case 1:
			return [4]xy{{x + 1, y}, {x + 1, y + 1}, {x + 1, y + 2}, {x, y + 1}}
		case 2:
			return [4]xy{{x, y + 1}, {x + 1, y + 1}, {x + 2, y + 1}, {x + 1, y}}
		case 3:
			return [4]xy{{x + 1, y}, {x + 1, y + 1}, {x + 1, y + 2}, {x + 2, y + 1}}
		}
	}
	panic("unknown tetromino kind")
}

var blockColors = []draw.Color{
	blockL: draw.RGB(1, 0, 0),
	blockJ: draw.RGB(0, 1, 1),
	blockO: draw.RGB(1, 1, 0),
	blockI: draw.RGB(1, 1, 1),
	blockS: draw.RGB(0, 1, 0),
	blockZ: draw.RGB(1, 0, 1),
	blockT: draw.RGB(0, 0, 1),
}

func (b *blocks) update(window draw.Window) gameMode {
	if window.WasKeyPressed(draw.KeyEscape) {
		return newMenu()
	}

	if window.WasKeyPressed(draw.KeyUp) {
		delta := 1
		if window.IsKeyDown(draw.KeyLeftShift) || window.IsKeyDown(draw.KeyRightShift) {
			delta = 3
		}

		old := b.thisPiece.rotation
		b.thisPiece.rotation = (b.thisPiece.rotation + delta) % 4
		if collides(&b.field, &b.thisPiece) {
			b.thisPiece.rotation = old
		}
	}

	if window.WasKeyPressed(draw.KeyLeft) {
		b.thisPiece.x--
		if collides(&b.field, &b.thisPiece) {
			b.thisPiece.x++
		}
	}

	if window.WasKeyPressed(draw.KeyRight) {
		b.thisPiece.x++
		if collides(&b.field, &b.thisPiece) {
			b.thisPiece.x--
		}
	}

	resetPieceInGround := func() {
		b.thisPiece.y--
		b.field.place(&b.thisPiece)
		b.thisPiece = b.nextPiece
		b.nextPiece = randomTetromino()
		b.field.clearFullRows()
	}

	if window.WasKeyPressed(draw.KeyDown) {
		// Drop one down.
		b.thisPiece.y++
		if collides(&b.field, &b.thisPiece) {
			resetPieceInGround()
		}
	}

	if window.WasKeyPressed(draw.KeySpace) {
		// Drop all the way to the floor.
		for !collides(&b.field, &b.thisPiece) {
			b.thisPiece.y++
		}
		resetPieceInGround()
	}

	windowW, windowH := window.Size()
	tileWidth := windowW / blockFieldWidth
	tileHeight := windowH / blockFieldHeight
	tileSize := min(tileWidth, tileHeight)
	fieldW := tileSize * blockFieldWidth
	fieldH := tileSize * blockFieldHeight
	xOffset := (windowW - fieldW) / 2
	yOffset := (windowH - fieldH) / 2
	borderColor := draw.RGB(0, 0, 0)

	// Draw background.
	window.FillRect(
		xOffset,
		yOffset,
		blockFieldWidth*tileSize,
		blockFieldHeight*tileSize,
		draw.RGB(0.83, 0.77, 0.65),
	)

	drawTile := func(tileX, tileY int, color draw.Color) {
		window.FillRect(tileX, tileY, tileSize, tileSize, borderColor)
		b := tileSize / 15
		window.FillRect(tileX+b, tileY+b, tileSize-2*b, tileSize-2*b, color)
	}

	for y := 0; y < blockFieldHeight; y++ {
		for x := 0; x < blockFieldWidth; x++ {
			kind := b.field[x][y]
			if kind != blockEmpty {
				tileX := xOffset + x*tileSize
				tileY := yOffset + y*tileSize
				color := blockColors[kind]
				drawTile(tileX, tileY, color)
			}
		}
	}

	color := blockColors[b.thisPiece.kind]
	for _, part := range b.thisPiece.parts() {
		tileX := xOffset + part.x*tileSize
		tileY := yOffset + part.y*tileSize
		drawTile(tileX, tileY, color)
	}

	return b
}

func collides(f *blocksField, t *tetromino) bool {
	for _, part := range t.parts() {
		x, y := part.x, part.y

		// Check if we hit left/right/bottom walls, there is no wall at the top.
		if x < 0 || x >= blockFieldWidth || y >= blockFieldHeight {
			return true
		}

		// Check if the part is inside a solid block of the field.
		if 0 <= x && x < blockFieldWidth &&
			0 <= y && y < blockFieldHeight &&
			f[x][y] != blockEmpty {
			return true
		}
	}

	return false
}

func (f *blocksField) place(t *tetromino) {
	for _, part := range t.parts() {
		x, y := part.x, part.y
		if 0 <= x && x < blockFieldWidth &&
			0 <= y && y < blockFieldHeight {
			f[x][y] = t.kind
		}
	}
}

func (f *blocksField) clearFullRows() {
	for y := 0; y < blockFieldHeight; y++ {
		if f.rowFull(y) {
			f.dropRowsInto(y)
		}
	}
}

func (f *blocksField) rowFull(y int) bool {
	for x := 0; x < blockFieldWidth; x++ {
		if f[x][y] == blockEmpty {
			return false
		}
	}
	return true
}

func (f *blocksField) dropRowsInto(fullY int) {
	for y := fullY; y >= 1; y-- {
		for x := 0; x < blockFieldWidth; x++ {
			f[x][y] = f[x][y-1]
		}
	}
	for x := 0; x < blockFieldWidth; x++ {
		f[x][0] = blockEmpty
	}
}
