package menu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
)

// MenuItem represents a single option in the menu.
type MenuItem struct {
	Label    string
	OnSelect func()
}

// Menu is a reusable vertical list menu with keyboard navigation.
type Menu struct {
	items       []MenuItem
	selected    int
	visible     bool
	fontSize    float64
	itemSpacing float64
	color       color.Color
	selectColor color.Color
	onNavigate  func()
	onSelect    func()
	onCancel    func()
}

// NewMenu creates a new menu with default settings.
func NewMenu() *Menu {
	return &Menu{
		fontSize:    24,
		itemSpacing: 4,
		color:       color.Gray{Y: 128},
		selectColor: color.White,
	}
}

// AddItem adds a menu item with the given label and callback.
func (m *Menu) AddItem(label string, callback func()) {
	m.items = append(m.items, MenuItem{Label: label, OnSelect: callback})
}

// SetVisible sets the visibility of the menu.
func (m *Menu) SetVisible(visible bool) {
	m.visible = visible
}

// Visible returns whether the menu is visible.
func (m *Menu) Visible() bool {
	return m.visible
}

// SetFontSize sets the font size for menu items.
func (m *Menu) SetFontSize(size float64) {
	m.fontSize = size
}

// SetItemSpacing sets the spacing between menu items.
func (m *Menu) SetItemSpacing(spacing float64) {
	m.itemSpacing = spacing
}

// UpdateItemLabel updates the label of a menu item at the given index.
func (m *Menu) UpdateItemLabel(index int, label string) {
	if index < 0 || index >= len(m.items) {
		return
	}
	m.items[index].Label = label
}

// SetColor sets the color for unselected items.
func (m *Menu) SetColor(c color.Color) {
	m.color = c
}

// SetSelectColor sets the color for the selected item.
func (m *Menu) SetSelectColor(c color.Color) {
	m.selectColor = c
}

// SetOnNavigate sets a callback for when navigation occurs (e.g., play sound).
func (m *Menu) SetOnNavigate(fn func()) {
	m.onNavigate = fn
}

// SetOnSelect sets a callback for when an item is selected (e.g., play sound).
func (m *Menu) SetOnSelect(fn func()) {
	m.onSelect = fn
}

// SetOnCancel sets a callback for when cancel is pressed.
func (m *Menu) SetOnCancel(fn func()) {
	m.onCancel = fn
}

// SelectedIndex returns the currently selected item index.
func (m *Menu) SelectedIndex() int {
	return m.selected
}

// Update handles input for menu navigation.
func (m *Menu) Update() {
	if !m.visible {
		return
	}

	// Navigate up
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.NavigateUp()
		return
	}

	// Navigate down
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.NavigateDown()
		return
	}

	// Select
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		m.Select()
		return
	}

	// Cancel
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		m.Cancel()
		return
	}
}

// NavigateUp moves selection up with wrap-around.
func (m *Menu) NavigateUp() {
	if len(m.items) == 0 {
		return
	}
	m.selected--
	if m.selected < 0 {
		m.selected = len(m.items) - 1
	}
	if m.onNavigate != nil {
		m.onNavigate()
	}
}

// NavigateDown moves selection down with wrap-around.
func (m *Menu) NavigateDown() {
	if len(m.items) == 0 {
		return
	}
	m.selected++
	if m.selected >= len(m.items) {
		m.selected = 0
	}
	if m.onNavigate != nil {
		m.onNavigate()
	}
}

// Select triggers the selected item's callback.
func (m *Menu) Select() {
	if len(m.items) == 0 || m.selected < 0 || m.selected >= len(m.items) {
		return
	}
	if m.onSelect != nil {
		m.onSelect()
	}
	if m.items[m.selected].OnSelect != nil {
		m.items[m.selected].OnSelect()
	}
}

// Cancel triggers the cancel callback.
func (m *Menu) Cancel() {
	if m.onCancel != nil {
		m.onCancel()
	}
}

// Draw renders the menu centered at the given position.
func (m *Menu) Draw(screen *ebiten.Image, font *font.FontText, centerX, centerY int) {
	if !m.visible || font == nil || len(m.items) == 0 {
		return
	}

	// Calculate total height to center vertically
	itemHeight := int(m.fontSize + m.itemSpacing)
	totalHeight := len(m.items)*itemHeight - int(m.itemSpacing)
	startY := centerY - totalHeight/2

	for i, item := range m.items {
		y := startY + i*itemHeight

		op := &text.DrawOptions{}
		if i == m.selected {
			op.ColorScale.ScaleWithColor(m.selectColor)
		} else {
			op.ColorScale.ScaleWithColor(m.color)
		}

		// Measure text width for precise centering
		face := font.NewFace(m.fontSize)
		textWidth, _ := text.Measure(item.Label, face, 0)
		op.GeoM.Translate(float64(centerX)-textWidth/2, float64(y))

		font.Draw(screen, item.Label, m.fontSize, op)
	}
}
