package main

import "github.com/gonutz/prototype/draw"

const (
	blockFieldWidth      = 10
	blockFieldHeight     = 18
	minDropDelay         = 3
	linesPerAcceleration = 10
)

func newBlocks() *blocks {
	initialDropDelay := 60
	return &blocks{
		thisPiece:  randomTetromino(),
		nextPiece:  randomTetromino(),
		dropDelay:  initialDropDelay,
		nextDropIn: initialDropDelay,
		// We increase the drop speed every linesPerAcceleration lines.
		nextDropAccelerationLines: linesPerAcceleration,
	}
}

type blocksField [blockFieldWidth][blockFieldHeight]blockKind

type blocks struct {
	thisPiece  tetromino
	nextPiece  tetromino
	field      blocksField
	left       bool
	leftTimer  int
	right      bool
	rightTimer int
	down       bool
	downTimer  int
	dropDelay  int
	nextDropIn int
	lines      int
	score      int
	// nextDropAccelerationLines remembers after how many lines total the next
	// drop speed increase happens.
	nextDropAccelerationLines int
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
		return globalMenu
	}

	rotateBy := func(delta int) {
		old := b.thisPiece.rotation
		b.thisPiece.rotation = (b.thisPiece.rotation + delta) % 4
		if collides(&b.field, &b.thisPiece) {
			b.thisPiece.rotation = old
		}
	}
	if window.WasKeyPressed(globalControls.rotateCW) {
		rotateBy(1)
	}
	if window.WasKeyPressed(globalControls.rotateCCW) {
		rotateBy(3)
	}

	moveX := func(dx int) {
		b.thisPiece.x += dx
		if collides(&b.field, &b.thisPiece) {
			b.thisPiece.x -= dx
		}
	}

	const (
		leftRightInitialDelay = 13
		leftRightRepeatDelay  = 4
	)

	b.leftTimer--
	leftDown := window.IsKeyDown(globalControls.left)
	if leftDown && !b.left || b.left && b.leftTimer <= 0 {
		moveX(-1)
		b.leftTimer = leftRightRepeatDelay
		if !b.left {
			b.leftTimer = leftRightInitialDelay
		}
	}
	b.left = leftDown

	b.rightTimer--
	rightDown := window.IsKeyDown(globalControls.right)
	if rightDown && !b.right || b.right && b.rightTimer <= 0 {
		moveX(1)
		b.rightTimer = leftRightRepeatDelay
		if !b.right {
			b.rightTimer = leftRightInitialDelay
		}
	}
	b.right = rightDown

	resetPieceInGround := func() {
		b.thisPiece.y--
		b.field.place(&b.thisPiece)
		b.thisPiece = b.nextPiece
		b.nextPiece = randomTetromino()
		b.clearFullRows()
		if b.lines >= b.nextDropAccelerationLines {
			b.dropDelay = b.dropDelay * 3 / 4
			b.nextDropAccelerationLines += 10
		}
		if b.dropDelay < minDropDelay {
			b.dropDelay = minDropDelay
		}
		b.down = false
	}

	wasDropped := false

	b.downTimer--
	downDown := window.IsKeyDown(globalControls.dropOne)
	if downDown && !b.down || b.down && b.downTimer <= 0 {
		// Drop one down.
		b.thisPiece.y++
		if collides(&b.field, &b.thisPiece) {
			resetPieceInGround()
		}
		wasDropped = true

		b.downTimer = 3
		if !b.down {
			b.downTimer = 20
		}
	}
	b.down = downDown

	if window.WasKeyPressed(globalControls.dropAll) {
		// Drop all the way to the floor.
		for !collides(&b.field, &b.thisPiece) {
			b.thisPiece.y++
		}
		resetPieceInGround()
		wasDropped = true
	}

	b.nextDropIn--
	if !wasDropped && b.nextDropIn <= 0 {
		// Drop one down.
		b.thisPiece.y++
		if collides(&b.field, &b.thisPiece) {
			resetPieceInGround()
		}
		wasDropped = true
	}

	if wasDropped {
		b.nextDropIn = b.dropDelay
	}

	windowW, windowH := window.Size()
	// We want to have 6 blocks left and 6 blocks right of the field. We want to
	// display the preview block on the right and the main game should be
	// centered. Thus 6 for the preview block on the right, 6 on the left to
	// have it be symetrical.
	tileWidth := windowW / (blockFieldWidth + 12)
	tileHeight := windowH / blockFieldHeight
	tileSize := min(tileWidth, tileHeight)
	fieldW := tileSize * blockFieldWidth
	fieldH := tileSize * blockFieldHeight
	xOffset := (windowW - fieldW) / 2
	yOffset := (windowH - fieldH) / 2
	borderColor := draw.RGB(0, 0, 0)
	fieldBackground := draw.RGB(0.83, 0.77, 0.65)

	// Draw background.
	window.FillRect(
		xOffset,
		yOffset,
		blockFieldWidth*tileSize,
		blockFieldHeight*tileSize,
		fieldBackground,
	)

	drawColoredTile := func(tileX, tileY int, kind blockKind, transparent bool) {
		color := blockColors[kind]
		if transparent {
			color.A = 0.4
		}
		if !transparent {
			window.FillRect(tileX, tileY, tileSize, tileSize, borderColor)
		}
		b := tileSize / 15
		window.FillRect(tileX+b, tileY+b, tileSize-2*b, tileSize-2*b, color)
	}
	drawTile := func(tileX, tileY int, kind blockKind) {
		drawColoredTile(tileX, tileY, kind, false)
	}
	drawTransparentTile := func(tileX, tileY int, kind blockKind) {
		drawColoredTile(tileX, tileY, kind, true)
	}

	// Draw the field, the pieces that are already solid.
	for y := 0; y < blockFieldHeight; y++ {
		for x := 0; x < blockFieldWidth; x++ {
			kind := b.field[x][y]
			if kind != blockEmpty {
				tileX := xOffset + x*tileSize
				tileY := yOffset + y*tileSize
				drawTile(tileX, tileY, kind)
			}
		}
	}

	// Draw currently active tetromino's drop shadow.
	dropped := b.thisPiece
	for !collides(&b.field, &dropped) {
		dropped.y++
	}
	dropped.y-- // Un-collide.
	for _, part := range dropped.parts() {
		tileX := xOffset + part.x*tileSize
		tileY := yOffset + part.y*tileSize
		drawTransparentTile(tileX, tileY, b.thisPiece.kind)
	}

	// Draw the currently active tetromino separately on top.
	for _, part := range b.thisPiece.parts() {
		tileX := xOffset + part.x*tileSize
		tileY := yOffset + part.y*tileSize
		drawTile(tileX, tileY, b.thisPiece.kind)
	}

	// Draw the preview tetromino.
	minX, maxX, minY, maxY := 999, -999, 999, -999
	for _, part := range b.nextPiece.parts() {
		x, y := part.x, part.y
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	pieceW := (maxX - minX + 1) * tileSize
	pieceH := (maxY - minY + 1) * tileSize
	fieldRight := windowW/2 + (blockFieldWidth*tileSize)/2
	previewX := fieldRight + tileSize/2
	previewY := yOffset + (blockFieldHeight-5)*tileSize
	previewSize := 5 * tileSize
	xOffset = previewX + (previewSize-pieceW)/2
	yOffset = previewY + (previewSize-pieceH)/2
	window.FillRect(previewX, previewY, previewSize, previewSize, fieldBackground)
	for _, part := range b.nextPiece.parts() {
		x := (part.x - minX) * tileSize
		y := (part.y - minY) * tileSize
		drawTile(xOffset+x, yOffset+y, b.nextPiece.kind)
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

func (b *blocks) clearFullRows() {
	rows := 0
	for y := 0; y < blockFieldHeight; y++ {
		if b.field.rowFull(y) {
			rows++
			b.field.dropRowsInto(y)
		}
	}
	scores := []int{0, 1, 3, 6, 10}
	b.score += scores[rows]
	b.lines += rows
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
