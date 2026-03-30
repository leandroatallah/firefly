package scene

import (
	"image"
	"log"
	"math"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/tilemap"
)

type FlipType int

const (
	FlipTypeSmooth FlipType = iota
	FlipTypeInstant
)

type ScreenFlipper struct {
	cam     *camera.Controller
	player  body.Movable
	tilemap *tilemap.Tilemap

	screenWidth  float64
	screenHeight float64

	// Configuration
	FlipStrategy       func(dx, dy int) FlipType
	PlayerPushDistance float64 // pixels

	// State
	isFlipping    bool
	flipSourceX   float64
	flipSourceY   float64
	flipTargetX   float64
	flipTargetY   float64
	playerSourceX float64
	playerTargetX float64
	playerSourceY float64
	playerTargetY float64
	flipProgress  float64

	// Hooks
	OnFlipStart  func()
	OnFlipFinish func()

	// Room Management
	rooms       []image.Rectangle
	currentRoom *image.Rectangle

	context *app.AppContext
}

func NewScreenFlipper(cam *camera.Controller, player body.Movable, tm *tilemap.Tilemap, context *app.AppContext) *ScreenFlipper {
	cfg := config.Get()
	return &ScreenFlipper{
		cam:                cam,
		player:             player,
		tilemap:            tm,
		context:            context,
		screenWidth:        float64(cfg.ScreenWidth),
		screenHeight:       float64(cfg.ScreenHeight),
		PlayerPushDistance: 0.0,
	}
}

func (sf *ScreenFlipper) IsFlipping() bool {
	return sf.isFlipping
}

func (sf *ScreenFlipper) Update() {
	if sf.isFlipping {
		sf.updateFlip()
		return
	}

	sf.checkTrigger()
}

func (sf *ScreenFlipper) SnapToCurrentRoom() {
	if sf.player == nil || sf.cam == nil || sf.tilemap == nil {
		return
	}

	sf.ensureRooms()
	if len(sf.rooms) == 0 {
		return
	}

	if sf.currentRoom == nil {
		sf.updateCurrentRoom()
	}
	if sf.currentRoom == nil {
		return
	}

	sf.cam.SetBounds(sf.currentRoom)

	px, py := sf.player.GetPositionMin()
	w, h := sf.player.GetShape().Width(), sf.player.GetShape().Height()
	centerX := float64(px) + float64(w)/2
	centerY := float64(py) + float64(h)/2

	halfW := sf.screenWidth / 2
	halfH := sf.screenHeight / 2

	minCamX := float64(sf.currentRoom.Min.X) + halfW
	maxCamX := float64(sf.currentRoom.Max.X) - halfW
	minCamY := float64(sf.currentRoom.Min.Y) + halfH
	maxCamY := float64(sf.currentRoom.Max.Y) - halfH

	if centerX < minCamX {
		centerX = minCamX
	}
	if centerX > maxCamX {
		centerX = maxCamX
	}
	if centerY < minCamY {
		centerY = minCamY
	}
	if centerY > maxCamY {
		centerY = maxCamY
	}

	sf.cam.SetCenter(centerX, centerY)
}

func (sf *ScreenFlipper) checkTrigger() {
	if sf.player == nil || sf.cam == nil || sf.tilemap == nil {
		return
	}

	sf.ensureRooms()

	// Initialize current room if needed
	if sf.currentRoom == nil {
		sf.updateCurrentRoom()
		// Snap camera to room
		if sf.currentRoom != nil {
			sf.cam.SetBounds(sf.currentRoom)
		}
		return
	}

	// Check bounds of CURRENT ROOM, not camera view
	left := float64(sf.currentRoom.Min.X)
	right := float64(sf.currentRoom.Max.X)
	top := float64(sf.currentRoom.Min.Y)
	bottom := float64(sf.currentRoom.Max.Y)

	// Player Position
	px, py := sf.player.GetPositionMin()
	w, h := sf.player.GetShape().Width(), sf.player.GetShape().Height()
	pCenterX := float64(px) + float64(w)/2
	pCenterY := float64(py) + float64(h)/2

	// Thresholds (small buffer to avoid floating point issues)
	buffer := 2.0

	if pCenterX < left-buffer {
		sf.triggerFlip(-1, 0)
	} else if pCenterX > right+buffer {
		sf.triggerFlip(1, 0)
	} else if pCenterY < top-buffer {
		sf.triggerFlip(0, -1)
	} else if pCenterY > bottom+buffer {
		sf.triggerFlip(0, 1)
	}
}

func (sf *ScreenFlipper) updateCurrentRoom() {
	px, py := sf.player.GetPositionMin()
	w, h := sf.player.GetShape().Width(), sf.player.GetShape().Height()
	center := image.Point{X: px + w/2, Y: py + h/2}

	for i := range sf.rooms {
		if center.In(sf.rooms[i]) {
			sf.currentRoom = &sf.rooms[i]
			return
		}
	}
}

func (sf *ScreenFlipper) ensureRooms() {
	if len(sf.rooms) > 0 {
		return
	}

	// 1. Try to load from "Camera" layer
	layer, found := sf.tilemap.FindLayerByName("Camera")
	if found {
		for _, obj := range layer.Objects {
			r := image.Rect(int(obj.X), int(obj.Y), int(obj.X+obj.Width), int(obj.Y+obj.Height))
			sf.rooms = append(sf.rooms, r)
		}
	}

	// 2. Fallback: Generate Grid
	if len(sf.rooms) == 0 {
		mw := sf.tilemap.Width * sf.tilemap.Tilewidth
		mh := sf.tilemap.Height * sf.tilemap.Tileheight
		cols := int(math.Ceil(float64(mw) / sf.screenWidth))
		rows := int(math.Ceil(float64(mh) / sf.screenHeight))

		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				r := image.Rect(
					x*int(sf.screenWidth),
					y*int(sf.screenHeight),
					(x+1)*int(sf.screenWidth),
					(y+1)*int(sf.screenHeight),
				)
				sf.rooms = append(sf.rooms, r)
			}
		}
	}
}

func (sf *ScreenFlipper) triggerFlip(dx, dy int) {
	// Calculate Player Target (pushed into next room)
	px, py := sf.player.GetPositionMin()
	playerSourceX := float64(px)
	playerSourceY := float64(py)
	playerTargetX := playerSourceX
	playerTargetY := playerSourceY

	w, h := sf.player.GetShape().Width(), sf.player.GetShape().Height()

	// Check for collisions and adjust push distance if needed
	pushDist := sf.PlayerPushDistance
	if sf.context != nil && sf.context.Space != nil && (dx != 0 || dy != 0) {
		// Try to find a valid position, decreasing push distance if blocked
		// We use a step of 4 pixels to check
		for pushDist >= 0 {
			tx := playerSourceX
			ty := playerSourceY
			if dx != 0 {
				tx += float64(dx) * pushDist
			}
			if dy != 0 {
				ty += float64(dy) * pushDist
			}

			// Check collision at this potential target
			rect := image.Rect(int(tx), int(ty), int(tx)+w, int(ty)+h)
			cols := sf.context.Space.Query(rect)
			blocked := false
			for _, c := range cols {
				// Ignore self and non-obstructive bodies
				if c.ID() != sf.player.ID() && c.IsObstructive() {
					blocked = true
					break
				}
			}

			if !blocked {
				// Found valid position
				playerTargetX = tx
				playerTargetY = ty
				break
			}

			pushDist -= 4 // Back off
		}
	} else {
		// No collision check available or needed
		if dx != 0 {
			playerTargetX += float64(dx) * sf.PlayerPushDistance
		}
		if dy != 0 {
			playerTargetY += float64(dy) * sf.PlayerPushDistance
		}
	}

	// Find Next Room based on Player Target
	targetCenter := image.Point{
		X: int(playerTargetX) + w/2,
		Y: int(playerTargetY) + h/2,
	}

	var nextRoom *image.Rectangle
	for i := range sf.rooms {
		if targetCenter.In(sf.rooms[i]) {
			nextRoom = &sf.rooms[i]
			break
		}
	}

	if nextRoom == nil {
		log.Printf("Cannot flip: No room found for target position %v", targetCenter)
		return
	}

	// Calculate Camera Target Position for the New Room
	// We want to center the camera on the player, but clamped to the New Room.
	// Temporarily calculate clamped position
	halfW := sf.screenWidth / 2
	halfH := sf.screenHeight / 2

	minCamX := float64(nextRoom.Min.X) + halfW
	maxCamX := float64(nextRoom.Max.X) - halfW
	minCamY := float64(nextRoom.Min.Y) + halfH
	maxCamY := float64(nextRoom.Max.Y) - halfH

	targetCamX := float64(targetCenter.X)
	targetCamY := float64(targetCenter.Y)

	if targetCamX < minCamX {
		targetCamX = minCamX
	}
	if targetCamX > maxCamX {
		targetCamX = maxCamX
	}
	if targetCamY < minCamY {
		targetCamY = minCamY
	}
	if targetCamY > maxCamY {
		targetCamY = maxCamY
	}

	// Current Camera Position
	camX, camY := sf.cam.Kamera().Center()

	// Determine Strategy
	flipType := FlipTypeSmooth
	if sf.FlipStrategy != nil {
		flipType = sf.FlipStrategy(dx, dy)
	}

	// Update Current Room immediately for Instant, or after for Smooth?
	// Actually, during flip, we are "between" rooms.
	// But we should unclamp the camera now.
	sf.cam.SetBounds(nil)

	if flipType == FlipTypeInstant {
		if sf.OnFlipStart != nil {
			sf.OnFlipStart()
		}
		sf.currentRoom = nextRoom
		sf.cam.SetBounds(nextRoom) // Re-clamp immediately
		sf.cam.SetCenter(targetCamX, targetCamY)
		sf.player.SetPosition(int(playerTargetX), int(playerTargetY))
		if sf.OnFlipFinish != nil {
			sf.OnFlipFinish()
		}
		return
	}

	// Start Flip
	sf.isFlipping = true
	sf.flipSourceX = camX
	sf.flipSourceY = camY
	sf.flipTargetX = targetCamX
	sf.flipTargetY = targetCamY
	sf.flipProgress = 0

	// Player Transition
	sf.playerSourceX = playerSourceX
	sf.playerSourceY = playerSourceY
	sf.playerTargetX = playerTargetX
	sf.playerTargetY = playerTargetY

	// Store next room to set it when finished
	sf.currentRoom = nextRoom

	if sf.OnFlipStart != nil {
		sf.OnFlipStart()
	}
}

func (sf *ScreenFlipper) updateFlip() {
	cfg := config.Get()
	speed := cfg.ScreenFlipSpeed
	if speed <= 0 {
		speed = 1.0 / 60.0 // Default fallback
	}

	sf.flipProgress += speed

	if sf.flipProgress >= 1.0 {
		sf.finishFlip()
		return
	}

	// Linear interpolation
	t := sf.flipProgress
	newCamX := sf.flipSourceX + (sf.flipTargetX-sf.flipSourceX)*t
	newCamY := sf.flipSourceY + (sf.flipTargetY-sf.flipSourceY)*t
	sf.cam.SetCenter(newCamX, newCamY)

	newPX := sf.playerSourceX + (sf.playerTargetX-sf.playerSourceX)*t
	newPY := sf.playerSourceY + (sf.playerTargetY-sf.playerSourceY)*t
	sf.player.SetPosition(int(newPX), int(newPY))
}

func (sf *ScreenFlipper) finishFlip() {
	sf.flipProgress = 1.0
	sf.isFlipping = false

	if sf.OnFlipFinish != nil {
		sf.OnFlipFinish()
	}

	// Force exact position
	sf.cam.SetCenter(sf.flipTargetX, sf.flipTargetY)
	sf.player.SetPosition(int(sf.playerTargetX), int(sf.playerTargetY))

	// Re-enable bounds for the new room
	if sf.currentRoom != nil {
		sf.cam.SetBounds(sf.currentRoom)
	}
}
