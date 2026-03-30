package sequences

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
)

// CameraZoomCommand sets the camera zoom level and always rewinds back
type CameraZoomCommand struct {
	Zoom            float64
	Duration        int    // Duration of zoom in frames (0 for instant)
	Delay           int    // Frames to wait at peak zoom before rewinding
	OutDuration     int    // Duration of zoom out in frames (0 to use same as Duration)
	TargetID        string // Optional: body/actor ID to center camera on
	currentZoom     float64
	targetZoom      float64
	timer           int
	camera          *camera.Controller
	startX          float64
	startY          float64
	phase           int // 0=zooming in, 1=waiting, 2=zooming out
	rewindDur       int
	targetActor     actors.ActorEntity
	origFollowing   bool       // Store original following state
	origFollowTarget body.Body // Store original follow target
	sameTarget      bool // True if target is same as current follow target
}

func (c *CameraZoomCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	currentScene := ctx.SceneManager.CurrentScene()

	// Try to get camera from TilemapScene or scenes that embed it
	var cam *camera.Controller

	if tilemapScene, ok := currentScene.(*scene.TilemapScene); ok {
		cam = tilemapScene.Camera()
	} else if phasesScene, ok := currentScene.(interface{ Camera() *camera.Controller }); ok {
		// For scenes that embed TilemapScene like PhasesScene
		cam = phasesScene.Camera()
	}

	c.camera = cam
	if c.camera != nil {
		c.currentZoom = c.camera.Kamera().ZoomFactor

		// Store original position BEFORE any movement
		c.startX, c.startY = c.camera.Kamera().Center()

		// Store original following state and target
		c.origFollowing = c.camera.IsFollowing()
		c.origFollowTarget = c.camera.FollowTarget()

		c.targetZoom = c.Zoom
		if c.targetZoom <= 0 {
			c.targetZoom = 1.0 // Default zoom
		}
		c.timer = 0
		c.phase = 0 // Start with zoom in phase

		// Set zoom out duration (use same as Duration if not specified)
		if c.OutDuration <= 0 {
			c.rewindDur = c.Duration
		} else {
			c.rewindDur = c.OutDuration
		}

		// Find target actor if specified
		if c.TargetID != "" {
			if actor, found := ctx.ActorManager.Find(c.TargetID); found {
				c.targetActor = actor
				// Check if target is same as current follow target
				c.sameTarget = c.origFollowing && c.origFollowTarget == actor
				// Only disable following if target is different from current follow target
				if !c.sameTarget {
					c.camera.SetFollowing(false)
					// Store current position for interpolation
					c.startX, c.startY = c.camera.Kamera().Center()
				}
				// If sameTarget is true, keep following enabled so camera naturally tracks the actor
			}
		} else {
			// No target specified, store current position for interpolation
			c.startX, c.startY = c.camera.Kamera().Center()
			c.camera.SetFollowing(false)
		}
	}
}

func (c *CameraZoomCommand) Update() bool {
	if c.camera == nil {
		return true
	}

	// Phase 0: Zoom in
	if c.phase == 0 {
		if c.Duration <= 0 {
			// Instant zoom in
			c.camera.Kamera().ZoomFactor = c.targetZoom
			c.phase = 1 // Move to wait phase
			c.timer = 0 // Reset timer for delay
		} else {
			c.timer++
			progress := float64(c.timer) / float64(c.Duration)
			if progress >= 1.0 {
				c.camera.Kamera().ZoomFactor = c.targetZoom
				c.phase = 1 // Move to wait phase
				c.timer = 0 // Reset timer for delay
			} else {
				// Linear interpolation for smooth zoom in
				currentZoom := c.currentZoom + (c.targetZoom-c.currentZoom)*progress
				c.camera.Kamera().ZoomFactor = currentZoom
			}
		}
	}

	// Phase 1: Wait at peak zoom
	if c.phase == 1 {
		c.timer++
		if c.timer >= c.Delay {
			c.phase = 2 // Move to zoom out phase
			c.timer = 0 // Reset timer for zoom out
		}
	}

	// Phase 2: Zoom out (rewind to original position and zoom)
	if c.phase == 2 {
		if c.rewindDur <= 0 {
			// Instant zoom out
			c.camera.Kamera().ZoomFactor = c.currentZoom
			// Only restore position if target was different from current follow target
			if !c.sameTarget {
				c.camera.SetCenter(c.startX, c.startY)
				// Restore original following state and target
				c.camera.SetFollowTarget(c.origFollowTarget)
				c.camera.SetFollowing(c.origFollowing)
			}
			return true
		}
		c.timer++
		progress := float64(c.timer) / float64(c.rewindDur)
		if progress >= 1.0 {
			c.camera.Kamera().ZoomFactor = c.currentZoom
			// Only restore position if target was different from current follow target
			if !c.sameTarget {
				c.camera.SetCenter(c.startX, c.startY)
				// Restore original following state and target
				c.camera.SetFollowTarget(c.origFollowTarget)
				c.camera.SetFollowing(c.origFollowing)
			}
			return true
		}
		// Zoom back out and move to original position (reverse interpolation)
		currentZoom := c.targetZoom + (c.currentZoom-c.targetZoom)*progress
		c.camera.Kamera().ZoomFactor = currentZoom
		// Only interpolate position if target was different from current follow target
		if !c.sameTarget {
			currX, currY := c.camera.Kamera().Center()
			currentX := currX + (c.startX-currX)*progress
			currentY := currY + (c.startY-currY)*progress
			c.camera.SetCenter(currentX, currentY)
		}
	}

	return false
}

// CameraMoveCommand moves the camera to a specified position
type CameraMoveCommand struct {
	X        float64
	Y        float64
	Duration int  // Duration in frames (0 for instant)
	Smooth   bool // Use smoothing if true
	startX   float64
	startY   float64
	timer    int
	camera   *camera.Controller
}

func (c *CameraMoveCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	currentScene := ctx.SceneManager.CurrentScene()

	// Try to get camera from TilemapScene or scenes that embed it
	var cam *camera.Controller

	if tilemapScene, ok := currentScene.(*scene.TilemapScene); ok {
		cam = tilemapScene.Camera()
	} else if phasesScene, ok := currentScene.(interface{ Camera() *camera.Controller }); ok {
		// For scenes that embed TilemapScene like PhasesScene
		cam = phasesScene.Camera()
	}

	c.camera = cam
	if c.camera != nil {
		// Store current position for interpolation
		c.startX, c.startY = c.camera.Kamera().Center()
		c.timer = 0
	}
}

func (c *CameraMoveCommand) Update() bool {
	if c.camera == nil {
		return true
	}

	if c.Duration <= 0 {
		// Instant move
		c.camera.SetCenter(c.X, c.Y)
		return true
	}

	c.timer++
	progress := float64(c.timer) / float64(c.Duration)
	if progress >= 1.0 {
		c.camera.SetCenter(c.X, c.Y)
		return true
	}

	// Linear interpolation for smooth movement
	currentX := c.startX + (c.X-c.startX)*progress
	currentY := c.startY + (c.Y-c.startY)*progress
	c.camera.SetCenter(currentX, currentY)
	return false
}

// CameraResetCommand resets the camera to default settings
type CameraResetCommand struct {
	DefaultZoom float64
	Duration    int // Duration in frames (0 for instant)
	camera      *camera.Controller
	startZoom   float64
	startX      float64
	startY      float64
	timer       int
}

func (c *CameraResetCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	currentScene := ctx.SceneManager.CurrentScene()

	// Try to get camera from TilemapScene or scenes that embed it
	var cam *camera.Controller

	if tilemapScene, ok := currentScene.(*scene.TilemapScene); ok {
		cam = tilemapScene.Camera()
	} else if phasesScene, ok := currentScene.(interface{ Camera() *camera.Controller }); ok {
		// For scenes that embed TilemapScene like PhasesScene
		cam = phasesScene.Camera()
	}

	c.camera = cam
	if c.camera != nil {
		c.startZoom = c.camera.Kamera().ZoomFactor
		c.startX, c.startY = c.camera.Kamera().Center()
		c.timer = 0

		// Set default zoom if not specified
		if c.DefaultZoom <= 0 {
			c.DefaultZoom = 1.0
		}
	}
}

func (c *CameraResetCommand) Update() bool {
	if c.camera == nil {
		return true
	}

	if c.Duration <= 0 {
		// Instant reset
		c.camera.Kamera().ZoomFactor = c.DefaultZoom
		c.camera.SetFollowing(true) // Re-enable following
		return true
	}

	c.timer++
	progress := float64(c.timer) / float64(c.Duration)
	if progress >= 1.0 {
		c.camera.Kamera().ZoomFactor = c.DefaultZoom
		c.camera.SetFollowing(true) // Re-enable following
		return true
	}

	// Smooth transition
	currentZoom := c.startZoom + (c.DefaultZoom-c.startZoom)*progress
	c.camera.Kamera().ZoomFactor = currentZoom
	return false
}

// CameraSetTargetCommand sets the camera target and can transition smoothly
type CameraSetTargetCommand struct {
	camera   *camera.Controller
	TargetID string
	Duration int // Duration in frames (0 for instant)

	target body.Body
	startX float64
	startY float64
	timer  int
}

func (c *CameraSetTargetCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	currentScene := ctx.SceneManager.CurrentScene()

	// Try to get camera from TilemapScene or scenes that embed it
	var cam *camera.Controller

	if tilemapScene, ok := currentScene.(*scene.TilemapScene); ok {
		cam = tilemapScene.Camera()
	} else if phasesScene, ok := currentScene.(interface{ Camera() *camera.Controller }); ok {
		// For scenes that embed TilemapScene like PhasesScene
		cam = phasesScene.Camera()
	}

	c.camera = cam
	if c.camera == nil {
		return
	}

	collidable := ctx.Space.Find(c.TargetID)
	if collidable == nil {
		return
	}
	c.target = collidable

	if c.Duration <= 0 {
		// Instant transition
		c.camera.SetFollowTarget(c.target)
		c.camera.SetFollowing(true)
		return
	}

	// Prepare smooth transition
	c.startX, c.startY = c.camera.Kamera().Center()
	c.camera.SetFollowing(false)
	c.timer = 0
}

func (c *CameraSetTargetCommand) Update() bool {
	if c.camera == nil || c.target == nil {
		return true
	}

	if c.Duration <= 0 {
		return true
	}

	c.timer++
	progress := float64(c.timer) / float64(c.Duration)

	// Ease-out quadratic interpolation: progress * (2 - progress)
	easedProgress := progress * (2 - progress)

	// Current target position (it might be moving!)
	x, y := c.target.GetPositionMin()
	w, h := c.target.GetShape().Width(), c.target.GetShape().Height()
	targetX := float64(x) + float64(w)/2
	targetY := float64(y) + float64(h)/2

	if progress >= 1.0 {
		c.camera.SetFollowTarget(c.target)
		c.camera.SetFollowing(true)
		return true
	}

	// Interpolate
	currentX := c.startX + (targetX-c.startX)*easedProgress
	currentY := c.startY + (targetY-c.startY)*easedProgress
	c.camera.SetCenter(currentX, currentY)

	return false
}

// CameraShakeCommand triggers a screen shake using trauma
type CameraShakeCommand struct {
	Trauma float64
	camera *camera.Controller
}

func (c *CameraShakeCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	currentScene := ctx.SceneManager.CurrentScene()

	var cam *camera.Controller
	if tilemapScene, ok := currentScene.(*scene.TilemapScene); ok {
		cam = tilemapScene.Camera()
	} else if phasesScene, ok := currentScene.(interface{ Camera() *camera.Controller }); ok {
		cam = phasesScene.Camera()
	}

	c.camera = cam
	if c.camera != nil {
		c.camera.AddTrauma(c.Trauma)
	}
}

func (c *CameraShakeCommand) Update() bool {
	return true
}
