package views

import (
	"gwk/vango"
)

var g_draw_context *vango.Context

func init_draw_context() {
	g_draw_context = vango.NewContext()
}

func GlobalDrawContext() *vango.Context {
	return g_draw_context
}
