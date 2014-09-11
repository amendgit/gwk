package views

import (
	"gwk/vango"
)

var g_graphic_context *vango.Context

func init_graphic_context() {
	g_graphic_context = vango.NewContext()
}

func GraphicContext() *vango.Context {
	return g_graphic_context
}
