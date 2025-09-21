package screenutil

import "github.com/leandroatallah/firefly/internal/config"

func GetCenterOfScreenPosition(width, height int) (int, int) {
	x := config.ScreenWidth/2 - width/2
	y := config.ScreenHeight/2 - height/2
	return x, y
}
