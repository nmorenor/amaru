package engine

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Spritesheet is a class that stores a set of tiles from a file, used by tilemaps and animations
type Spritesheet struct {
	texture       *ebiten.Image         // The original texture
	width, height float64               // The dimensions of the total texture
	cells         []SpriteRegion        // The dimensions of each sprite
	cache         map[int]*ebiten.Image // The cell cache cells
}

type Point struct {
	X, Y int
}

// SpriteRegion holds the position data for each sprite on the sheet
type SpriteRegion struct {
	Position      Point
	Width, Height int
}

// NewAsymmetricSpritesheetFromTexture creates a new AsymmetricSpriteSheet from a
// TextureResource. The data provided is the location and size of the sprites
func NewAsymmetricSpritesheetFromTexture(tr *ebiten.Image, spriteRegions []SpriteRegion) *Spritesheet {
	sheet := &Spritesheet{
		texture: tr,
		width:   float64(tr.Bounds().Dx()),
		height:  float64(tr.Bounds().Dy()),
		cells:   spriteRegions,
		cache:   make(map[int]*ebiten.Image),
	}
	return sheet
}

// NewSpritesheetFromTexture creates a new spritesheet from a texture resource.
func NewSpritesheetFromTexture(tr *ebiten.Image, cellWidth, cellHeight int) *Spritesheet {
	spriteRegions := generateSymmetricSpriteRegions(float64(tr.Bounds().Dx()), float64(tr.Bounds().Dy()), cellWidth, cellHeight, 0, 0)
	return NewAsymmetricSpritesheetFromTexture(tr, spriteRegions)
}

// NewSpritesheetWithBorderFromTexture creates a new spritesheet from a texture resource.
// This sheet has sprites of a uniform width and height, but also have borders around
// each sprite to prevent bleeding over
func NewSpritesheetWithBorderFromTexture(tr *ebiten.Image, cellWidth, cellHeight, borderWidth, borderHeight int) *Spritesheet {
	spriteRegions := generateSymmetricSpriteRegions(float64(tr.Bounds().Dx()), float64(tr.Bounds().Dy()), cellWidth, cellHeight, borderWidth, borderHeight)
	return NewAsymmetricSpritesheetFromTexture(tr, spriteRegions)
}

// Cell gets the region at the index i, updates and pulls from cache if need be
func (s *Spritesheet) Cell(index int) *ebiten.Image {
	if r, ok := s.cache[index]; ok {
		return r
	}

	cell := s.cells[index]
	width := cell.Width
	height := cell.Height
	x := (index % len(s.cells)) * width
	y := (index / len(s.cells)) * height

	rect := image.Rect(x, y, x+width, y+height)
	img := s.texture.SubImage(rect).(*ebiten.Image)
	s.cache[index] = img

	return s.cache[index]
}

// Drawable returns the drawable for a given index
func (s *Spritesheet) Drawable(index int) *ebiten.Image {
	return s.Cell(index)
}

// Drawables returns all the drawables on the sheet
func (s *Spritesheet) Drawables() []*ebiten.Image {
	drawables := make([]*ebiten.Image, s.CellCount())

	for i := 0; i < s.CellCount(); i++ {
		drawables[i] = s.Drawable(i)
	}

	return drawables
}

// CellCount returns the number of cells on the sheet
func (s *Spritesheet) CellCount() int {
	return len(s.cells)
}

// Cells returns all the cells on the sheet
func (s *Spritesheet) Cells() []*ebiten.Image {
	cellsNo := s.CellCount()
	cells := make([]*ebiten.Image, cellsNo)
	for i := 0; i < cellsNo; i++ {
		cells[i] = s.Cell(i)
	}

	return cells
}

// Width is the amount of tiles on the x-axis of the spritesheet
// only if the sprite sheet is symmetric with no border.
func (s Spritesheet) Width() float64 {
	return s.width / float64(s.Cell(0).Bounds().Dx())
}

// Height is the amount of tiles on the y-axis of the spritesheet
// only if the sprite sheet is symmetric with no border.
func (s Spritesheet) Height() float64 {
	return s.height / float64(s.Cell(0).Bounds().Dy())
}

func generateSymmetricSpriteRegions(totalWidth, totalHeight float64, cellWidth, cellHeight, borderWidth, borderHeight int) []SpriteRegion {
	var spriteRegions []SpriteRegion

	for y := 0; y <= int(math.Floor(totalHeight-1)); y += cellHeight + borderHeight {
		for x := 0; x <= int(math.Floor(totalWidth-1)); x += cellWidth + borderWidth {
			spriteRegion := SpriteRegion{
				Position: Point{X: x, Y: y},
				Width:    cellWidth,
				Height:   cellHeight,
			}
			spriteRegions = append(spriteRegions, spriteRegion)
		}
	}

	return spriteRegions
}
