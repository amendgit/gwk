package views

import (
	"gwk/views/resc"
)

func InitViews() {
	resc.InitResc()
	init_mockup()
	init_layout()
}
