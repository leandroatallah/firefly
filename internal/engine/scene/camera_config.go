package scene

type CameraMode string

const (
	CameraModeFixed   CameraMode = "fixed"
	CameraModeFollow  CameraMode = "follow"
)

type CameraConfig struct {
	Mode CameraMode
}
