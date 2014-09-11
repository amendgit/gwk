package views

import (
	"gwk/views/resc"
)

func InitViews() {
	resc.InitResc()
	init_mockup()
	init_layout()
	init_graphic_context()
}
